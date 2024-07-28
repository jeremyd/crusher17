package crusher17

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"os"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip44"
)

var logger *log.Logger

func init() {
	// Initialize the logger to write to stderr
	logger = log.New(os.Stderr, "crusher17: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func SendMessage(senderSk string, receiverPub string, message string, senderRelay string, receiverRelay string, subject string) (string, error) {
	senderPub, err := nostr.GetPublicKey(senderSk)
	if err != nil {
		logger.Printf("error getting public key 1: %v", err)
		return "", err
	}
	tags := nostr.Tags{
		nostr.Tag{"p", senderPub, senderRelay},
		nostr.Tag{"p", receiverPub, receiverRelay},
	}

	if subject != "" {
		tags = append(tags, nostr.Tag{"subject", subject})
	}

	salt := make([]byte, 32)
	rand.Read(salt)
	// Create a chat message
	ev := nostr.Event{
		PubKey:    senderPub,
		CreatedAt: nostr.Now(),
		Kind:      14,
		Tags:      tags,
		Content:   message,
	}
	ev.ID = ev.GetID()
	// Encrypt the message using NIP-44
	conversationKey, err := nip44.GenerateConversationKey(receiverPub, senderSk)
	if err != nil {
		logger.Printf("error generating convo key 1: %v", err)
		return "", err
	}
	encryptedMsg, err := nip44.Encrypt(ev.String(), conversationKey, nip44.WithCustomSalt(salt))
	if err != nil {
		logger.Printf("error encrypting nip44 1: %v", err)
		return "", err
	}
	// Create a seal (kind 13) using the encrypted message
	// make a random key
	randoSk := nostr.GeneratePrivateKey()
	randoPk, _ := nostr.GetPublicKey(randoSk)
	seal := nostr.Event{
		PubKey:    senderPub,
		CreatedAt: nostr.Now(),
		Kind:      13,
		Tags:      nostr.Tags{},
		Content:   encryptedMsg,
	}
	// sign the seal with the sender key
	seal.Sign(senderSk)
	// Encrypt the seal using NIP-44 for sending to receiver (1)
	sealConvoKey, err := nip44.GenerateConversationKey(receiverPub, randoSk)
	if err != nil {
		logger.Printf("error generating convo key for seal: %v", err)
		return "", err
	}
	encryptedSeal, err := nip44.Encrypt(seal.String(), sealConvoKey, nip44.WithCustomSalt(salt))
	if err != nil {
		logger.Printf("error encrypting seal: %v", err)
		return "", err
	}

	// TODO: should we create two giftwraps, one for sender and one for receiver?
	// Create a gift wrap (kind 1059) using the encrypted seal
	giftWrap := nostr.Event{
		PubKey:    randoPk,
		CreatedAt: nostr.Now(),
		Kind:      1059,
		Tags:      nostr.Tags{{"p", receiverPub, receiverRelay}},
		Content:   encryptedSeal,
	}
	giftWrap.Sign(randoSk)
	return giftWrap.String(), nil
}

func ReceiveMessage(receiverSk string, giftWrap string) (string, error) {
	var ev nostr.Event
	err := json.Unmarshal([]byte(giftWrap), &ev)
	if err != nil {
		logger.Printf("error unmarshalling giftwrap: %v", err)
		return "", err
	}
	unSealConvoKey, err := nip44.GenerateConversationKey(ev.PubKey, receiverSk)
	if err != nil {
		logger.Printf("error generating convo key 3: %v", err)
		return "", err
	}
	decryptedSeal, err := nip44.Decrypt(ev.Content, unSealConvoKey)
	if err != nil {
		logger.Printf("error decrypting seal: %v", err)
		return "", err
	}
	var k13 nostr.Event
	newerr := json.Unmarshal([]byte(decryptedSeal), &k13)
	if newerr != nil {
		logger.Printf("error unmarshalling decrypted seal: %v", newerr)
		return "", newerr
	}
	isOk, err := k13.CheckSignature()
	if err != nil || !isOk {
		logger.Printf("error checking signature: %v, isOk: %v", err, isOk)
		return "", err
	}
	k14ConvoKey, err := nip44.GenerateConversationKey(k13.PubKey, receiverSk)
	if err != nil {
		logger.Printf("error generating convo key 4: %v", err)
		return "", err
	}
	decryptedK14, err := nip44.Decrypt(k13.Content, k14ConvoKey)
	if err != nil {
		logger.Printf("error decrypting k14: %v", err)
		return "", err
	}
	var k14 nostr.Event
	k14err := json.Unmarshal([]byte(decryptedK14), &k14)
	if k14err != nil {
		logger.Printf("error unmarshaling kind 14 event: %v", k14err)
		return "", k14err
	}
	if k13.PubKey != k14.PubKey {
		logger.Println("impersonation risk: public keys don't match on kind13 and kind14")
		return "", err
	}
	return decryptedK14, nil
}

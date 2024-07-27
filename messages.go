package crusher17

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip44"
)


func SendMessage(senderSk string, receiverPub string, message string) (string, error) {

	senderPub, err := nostr.GetPublicKey(senderSk)

	if err != nil {
		fmt.Println("error getting public key 1: ", err)
		return "", err
	}

	tags := nostr.Tags{
		nostr.Tag{"p", senderPub, "wss://auth.nostr1.com"},
		nostr.Tag{"p", receiverPub, "wss://auth.nostr1.com"},
		nostr.Tag{"subject", "conversation-title"},
	}

	salt := make([]byte, 32)
	rand.Read(salt)

	// Create a chat message
	ev := nostr.Event{
		PubKey:    senderPub,
		CreatedAt: nostr.Now(),
		Kind:      14,
		Tags: tags,
		Content: message,
	}

	ev.ID = ev.GetID()

	// Encrypt the message using NIP-44

	conversationKey, err := nip44.GenerateConversationKey(receiverPub, senderSk)
	if err != nil {
		fmt.Println("error generating convo key 1:", err)
		return "", err
	}

	encryptedMsg, err := nip44.Encrypt(ev.String(),conversationKey, nip44.WithCustomSalt(salt))
	if err != nil {
		fmt.Println("error encrypting nip44 1: ", err)
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
		fmt.Println("error generating convo key for seal: ", err)
		return "", err
	}

	encryptedSeal, err := nip44.Encrypt(seal.String(),sealConvoKey, nip44.WithCustomSalt(salt))
	if err != nil {
		fmt.Println("error encrypting seal: ", err)
		return "", err
	}

	// Create a gift wrap (kind 1059) using the encrypted seal
	giftWrap := nostr.Event{
		PubKey:    randoPk,
		CreatedAt: nostr.Now(),
		Kind:      1059,
		Tags:      nostr.Tags{{"p", receiverPub, "wss://auth.nostr1.com"}},
		Content:   encryptedSeal,
	}

	giftWrap.Sign(randoSk)

	return giftWrap.String(), nil
}

func ReceiveMessage(receiverSk string, giftWrap string) (string, error) {
	var ev nostr.Event

	err := json.Unmarshal([]byte(giftWrap), &ev)
	if err != nil {
		fmt.Println("error unmarshalling giftwrap: ", err)
		return "", err
	}

	unSealConvoKey, err := nip44.GenerateConversationKey(ev.PubKey, receiverSk)

	if err != nil {
		fmt.Println("error generating convo key 3: ", err)
		return "", err
	}

	decryptedSeal, err := nip44.Decrypt(ev.Content, unSealConvoKey)
	if err != nil {
		fmt.Println("error decrypting seal: ", err)
		return "", err
	}

	var k13 nostr.Event
	newerr := json.Unmarshal([]byte(decryptedSeal), &k13)
	if newerr != nil {
		fmt.Println("error unmarshalling decrypted seal: ", newerr)
		return "", newerr
	}

	isOk, err := k13.CheckSignature()
	if err != nil || isOk == false {
		fmt.Println("error checking signature: ", err, isOk)
		return "", err
	}

	k14ConvoKey, err := nip44.GenerateConversationKey(k13.PubKey, receiverSk)

	if err != nil {
		fmt.Println("error generating convo key 4: ", err)
		return "", err
	}

	decryptedK14, err := nip44.Decrypt(k13.Content, k14ConvoKey)
	if err != nil {
		fmt.Println("error decrypting k14: ", err)
		return "", err
	}

	return decryptedK14, nil
}
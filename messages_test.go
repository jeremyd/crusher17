package crusher17

import (
	"fmt"
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
)

func TestSendAndReceiveMessage(t *testing.T) {
	// Generate sender and receiver key pairs
	senderSk := nostr.GeneratePrivateKey()
	senderPk, _ := nostr.GetPublicKey(senderSk)

	receiverSk := nostr.GeneratePrivateKey()
	receiverPk, _ := nostr.GetPublicKey(receiverSk)

	// Test message
	message := "Hello, this is a test message!"

	// Send the message
	giftWrap, err := SendMessage(senderSk, receiverPk, message, "wss://somerelay.com", "wss://somerelay.com", "conversation-title")
	assert.NoError(t, err, "SendMessage should not return an error")
	assert.NotEmpty(t, giftWrap, "Gift wrap should not be empty")

	// Receive the message
	decryptedMessage, err := ReceiveMessage(receiverSk, giftWrap)
	assert.NoError(t, err, "ReceiveMessage should not return an error")

	// Parse the decrypted message
	var receivedEvent nostr.Event
	err = receivedEvent.UnmarshalJSON([]byte(decryptedMessage))
	assert.NoError(t, err, "Unmarshalling decrypted message should not return an error")

	// Check the received message
	assert.Equal(t, message, receivedEvent.Content, "Received message content should match sent message")
	assert.Equal(t, senderPk, receivedEvent.PubKey, "Sender public key should match")
	assert.Equal(t, 14, receivedEvent.Kind, "Event kind should be 14")

	// Check tags
	assert.Len(t, receivedEvent.Tags, 3, "Event should have 3 tags")
	assert.Equal(t, "p", receivedEvent.Tags[0][0], "First tag should be 'p'")
	assert.Equal(t, senderPk, receivedEvent.Tags[0][1], "First 'p' tag should contain sender's public key")
	assert.Equal(t, "p", receivedEvent.Tags[1][0], "Second tag should be 'p'")
	assert.Equal(t, receiverPk, receivedEvent.Tags[1][1], "Second 'p' tag should contain receiver's public key")
	assert.Equal(t, "subject", receivedEvent.Tags[2][0], "Third tag should be 'subject'")

	fmt.Println("Test completed successfully!")
}

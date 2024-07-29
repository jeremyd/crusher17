# Nostr Secure Messaging Module

This Go module provides a secure messaging implementation of NIP-17 private messages using the Nostr protocol. 

[NIP-17](https://github.com/nostr-protocol/nips/blob/master/17.md)
[NIP-44](https://github.com/nostr-protocol/nips/blob/master/44.md)

## Features

- Create encrypted nostr events 
- Decrypt nostr events
- Implements NIP-44 encryption for secure communication
- Uses gift wrapping (Kind 1059) for additional security

## Installation

To use this module in your Go project, run:

```
go get github.com/jeremyd/crusher17
```

## Usage

### Encrypt/Wrap a message and return the json string result

```go
package main

import (
	"fmt"
	"github.com/jeremyd/crusher17"
)

func main() {
	// Create a new GiftWrapEvent
	gwe := &crusher17.GiftWrapEvent{
		SenderSecretKey: "your_sender_secret_key",
		SenderRelay:     "wss://your-relay.com",
		ReceiverPubkeys: map[string]string{
			"receiver_pubkey1": "wss://receiver-relay1.com",
			"receiver_pubkey2": "wss://receiver-relay2.com",
		},
		Content: "Hello, this is a secure message!",
		Subject: "Greetings",
		GiftWraps: make(map[string]string),
	}

	// Wrap the message
	err := gwe.Wrap()
	if err != nil {
		fmt.Printf("Error wrapping message: %v\n", err)
		return
	}

	// Print the gift wraps for each receiver
	for pubkey, giftWrap := range gwe.GiftWraps {
		fmt.Printf("Gift wrap for %s:\n%s\n\n", pubkey, giftWrap)
	}

	// Example of receiving a message
	receiverSecretKey := "receiver_secret_key"
	receivedGiftWrap := gwe.GiftWraps["receiver_pubkey1"] // Assuming this is the gift wrap for the receiver

	decryptedMessage, err := crusher17.ReceiveMessage(receiverSecretKey, receivedGiftWrap)
	if err != nil {
		fmt.Printf("Error receiving message: %v\n", err)
		return
	}

	fmt.Printf("Decrypted message: %s\n", decryptedMessage)
}

```

todo:
- [ ] randomly set timestamp on the giftwrap to "up to 2 days in the past" according to the NIP-17 spec.
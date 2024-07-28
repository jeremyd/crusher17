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
import "github.com/jeremyd/crusher17"

senderSk := "sender_private_key_hex"
receiverPk := "receiver_public_key_hex"
message := "Hello, this is a NIP-17 private message!"

eventJson, err := crusher17.SendMessage(senderSk, receiverPk, message)
if err != nil {
    // Handle error
}
```

### Decrypt/de-Wrap an event (kind 1059) and return the json string result
```go
import "github.com/jeremyd/crusher17"

receiverSk := "receiver_private_key_hex"
eventJson := "received_event_json_to_parse_and_decrypt"

decryptedMessage, err := crusher17.ReceiveMessage(receiverSk, eventJson)
if err != nil {
    // Handle error
}
```

todo:
- [ ] randomly set timestamp on the giftwrap to "up to 2 days in the past" according to the NIP-17 spec.
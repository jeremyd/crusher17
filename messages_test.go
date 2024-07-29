package crusher17

import (
	"encoding/json"
	"testing"

	"github.com/nbd-wtf/go-nostr"
	"github.com/stretchr/testify/assert"
)

func TestGiftWrapEvent(t *testing.T) {
	// Generate sender and receiver key pairs
	senderSk := nostr.GeneratePrivateKey()
	senderPk, _ := nostr.GetPublicKey(senderSk)

	receiverSk := nostr.GeneratePrivateKey()
	receiverPk, _ := nostr.GetPublicKey(receiverSk)

	// Create a GiftWrapEvent
	gwe := &GiftWrapEvent{
		SenderSecretKey: senderSk,
		SenderRelay:     "wss://sender.relay.com",
		CreatedAt:       nostr.Now(),
		ReceiverPubkeys: map[string]string{receiverPk: "wss://receiver.relay.com"},
		Content:         "Test message content",
		Subject:         "Test subject",
		GiftWraps:       make(map[string]string),
		ReplyTo:         "",
	}

	// Wrap the message
	err := gwe.Wrap()
	assert.NoError(t, err, "Wrapping message should not return an error")

	// Check if gift wraps were created
	assert.Len(t, gwe.GiftWraps, 2, "Should have 2 gift wraps (one for receiver, one for sender)")
	assert.Contains(t, gwe.GiftWraps, receiverPk, "Should contain gift wrap for receiver")
	assert.Contains(t, gwe.GiftWraps, senderPk, "Should contain gift wrap for sender")

	// Decrypt and verify the receiver's gift wrap
	decryptedMessage, err := ReceiveMessage(receiverSk, gwe.GiftWraps[receiverPk])
	assert.NoError(t, err, "Decrypting message should not return an error")

	var receivedEvent nostr.Event
	err = json.Unmarshal([]byte(decryptedMessage), &receivedEvent)
	assert.NoError(t, err, "Unmarshalling decrypted message should not return an error")

	// Verify the decrypted message
	assert.Equal(t, gwe.Content, receivedEvent.Content, "Decrypted content should match original content")
	assert.Equal(t, senderPk, receivedEvent.PubKey, "Sender public key should match")
	assert.Equal(t, 14, receivedEvent.Kind, "Event kind should be 14")

	// Check tags
	assert.GreaterOrEqual(t, len(receivedEvent.Tags), 2, "Event should have at least 2 tags")
	assert.Equal(t, "p", receivedEvent.Tags[0][0], "First tag should be 'p'")
	assert.Equal(t, senderPk, receivedEvent.Tags[0][1], "First 'p' tag should contain sender's public key")
	assert.Equal(t, "p", receivedEvent.Tags[1][0], "Second tag should be 'p'")
	assert.Equal(t, receiverPk, receivedEvent.Tags[1][1], "Second 'p' tag should contain receiver's public key")

	if gwe.Subject != "" {
		assert.Equal(t, "subject", receivedEvent.Tags[2][0], "Third tag should be 'subject'")
		assert.Equal(t, gwe.Subject, receivedEvent.Tags[2][1], "Subject should match")
	}
}

func TestNip17Example(t *testing.T) {
	// Generate sender and receiver key pairs
	// these values were taken from the NIP17 documentation example

	sk1 := "71f8de50a46c9996a21123280c6217c48f67d1378ff4fb14d4f7612181a1ebde"
	sk2 := "511cbb07ec2028bd2dcd039c447581a7f754df9d9a0e5c16b19a5422ab391563"

	giftToReceiver := `{"id":"2886780f7349afc1344047524540ee716f7bdc1b64191699855662330bf235d8", "pubkey":"8f8a7ec43b77d25799281207e1a47f7a654755055788f7482653f9c9661c6d51", "created_at":1703128320, "kind":1059, "tags":[ [ "p", "918e2da906df4ccd12c8ac672d8335add131a4cf9d27ce42b3bb3625755f0788"] ], "content":"AsqzdlMsG304G8h08bE67dhAR1gFTzTckUUyuvndZ8LrGCvwI4pgC3d6hyAK0Wo9gtkLqSr2rT2RyHlE5wRqbCOlQ8WvJEKwqwIJwT5PO3l2RxvGCHDbd1b1o40ZgIVwwLCfOWJ86I5upXe8K5AgpxYTOM1BD+SbgI5jOMA8tgpRoitJedVSvBZsmwAxXM7o7sbOON4MXHzOqOZpALpS2zgBDXSAaYAsTdEM4qqFeik+zTk3+L6NYuftGidqVluicwSGS2viYWr5OiJ1zrj1ERhYSGLpQnPKrqDaDi7R1KrHGFGyLgkJveY/45y0rv9aVIw9IWF11u53cf2CP7akACel2WvZdl1htEwFu/v9cFXD06fNVZjfx3OssKM/uHPE9XvZttQboAvP5UoK6lv9o3d+0GM4/3zP+yO3C0NExz1ZgFmbGFz703YJzM+zpKCOXaZyzPjADXp8qBBeVc5lmJqiCL4solZpxA1865yPigPAZcc9acSUlg23J1dptFK4n3Tl5HfSHP+oZ/QS/SHWbVFCtq7ZMQSRxLgEitfglTNz9P1CnpMwmW/Y4Gm5zdkv0JrdUVrn2UO9ARdHlPsW5ARgDmzaxnJypkfoHXNfxGGXWRk0sKLbz/ipnaQP/eFJv/ibNuSfqL6E4BnN/tHJSHYEaTQ/PdrA2i9laG3vJti3kAl5Ih87ct0w/tzYfp4SRPhEF1zzue9G/16eJEMzwmhQ5Ec7jJVcVGa4RltqnuF8unUu3iSRTQ+/MNNUkK6Mk+YuaJJs6Fjw6tRHuWi57SdKKv7GGkr0zlBUU2Dyo1MwpAqzsCcCTeQSv+8qt4wLf4uhU9Br7F/L0ZY9bFgh6iLDCdB+4iABXyZwT7Ufn762195hrSHcU4Okt0Zns9EeiBOFxnmpXEslYkYBpXw70GmymQfJlFOfoEp93QKCMS2DAEVeI51dJV1e+6t3pCSsQN69Vg6jUCsm1TMxSs2VX4BRbq562+VffchvW2BB4gMjsvHVUSRl8i5/ZSDlfzSPXcSGALLHBRzy+gn0oXXJ/447VHYZJDL3Ig8+QW5oFMgnWYhuwI5QSLEyflUrfSz+Pdwn/5eyjybXKJftePBD9Q+8NQ8zulU5sqvsMeIx/bBUx0fmOXsS3vjqCXW5IjkmSUV7q54GewZqTQBlcx+90xh/LSUxXex7UwZwRnifvyCbZ+zwNTHNb12chYeNjMV7kAIr3cGQv8vlOMM8ajyaZ5KVy7HpSXQjz4PGT2/nXbL5jKt8Lx0erGXsSsazkdoYDG3U", "sig":"a3c6ce632b145c0869423c1afaff4a6d764a9b64dedaf15f170b944ead67227518a72e455567ca1c2a0d187832cecbde7ed478395ec4c95dd3e71749ed66c480"}`
	decryptedGift1, err := ReceiveMessage(sk2, giftToReceiver)
	assert.NoError(t, err, "ReceiveMessage should not return an error")
	assert.NotEmpty(t, decryptedGift1, "Gift wrap should not be empty")

	giftToSender := ` { "id":"162b0611a1911cfcb30f8a5502792b346e535a45658b3a31ae5c178465509721", "pubkey":"626be2af274b29ea4816ad672ee452b7cf96bbb4836815a55699ae402183f512", "created_at":1702711587, "kind":1059, "tags":[ [ "p", "44900586091b284416a0c001f677f9c49f7639a55c3f1e2ec130a8e1a7998e1b"] ], "content":"AsTClTzr0gzXXji7uye5UB6LYrx3HDjWGdkNaBS6BAX9CpHa+Vvtt5oI2xJrmWLen+Fo2NBOFazvl285Gb3HSM82gVycrzx1HUAaQDUG6HI7XBEGqBhQMUNwNMiN2dnilBMFC3Yc8ehCJT/gkbiNKOpwd2rFibMFRMDKai2mq2lBtPJF18oszKOjA+XlOJV8JRbmcAanTbEK5nA/GnG3eGUiUzhiYBoHomj3vztYYxc0QYHOx0WxiHY8dsC6jPsXC7f6k4P+Hv5ZiyTfzvjkSJOckel1lZuE5SfeZ0nduqTlxREGeBJ8amOykgEIKdH2VZBZB+qtOMc7ez9dz4wffGwBDA7912NFS2dPBr6txHNxBUkDZKFbuD5wijvonZDvfWq43tZspO4NutSokZB99uEiRH8NAUdGTiNb25m9JcDhVfdmABqTg5fIwwTwlem5aXIy8b66lmqqz2LBzJtnJDu36bDwkILph3kmvaKPD8qJXmPQ4yGpxIbYSTCohgt2/I0TKJNmqNvSN+IVoUuC7ZOfUV9lOV8Ri0AMfSr2YsdZ9ofV5o82ClZWlWiSWZwy6ypa7CuT1PEGHzywB4CZ5ucpO60Z7hnBQxHLiAQIO/QhiBp1rmrdQZFN6PUEjFDloykoeHe345Yqy9Ke95HIKUCS9yJurD+nZjjgOxZjoFCsB1hQAwINTIS3FbYOibZnQwv8PXvcSOqVZxC9U0+WuagK7IwxzhGZY3vLRrX01oujiRrevB4xbW7Oxi/Agp7CQGlJXCgmRE8Rhm+Vj2s+wc/4VLNZRHDcwtfejogjrjdi8p6nfUyqoQRRPARzRGUnnCbh+LqhigT6gQf3sVilnydMRScEc0/YYNLWnaw9nbyBa7wFBAiGbJwO40k39wj+xT6HTSbSUgFZzopxroO3f/o4+ubx2+IL3fkev22mEN38+dFmYF3zE+hpE7jVxrJpC3EP9PLoFgFPKCuctMnjXmeHoiGs756N5r1Mm1ffZu4H19MSuALJlxQR7VXE/LzxRXDuaB2u9days/6muP6gbGX1ASxbJd/ou8+viHmSC/ioHzNjItVCPaJjDyc6bv+gs1NPCt0qZ69G+JmgHW/PsMMeL4n5bh74g0fJSHqiI9ewEmOG/8bedSREv2XXtKV39STxPweceIOh0k23s3N6+wvuSUAJE7u1LkDo14cobtZ/MCw/QhimYPd1u5HnEJvRhPxz0nVPz0QqL/YQeOkAYk7uzgeb2yPzJ6DBtnTnGDkglekhVzQBFRJdk740LEj6swkJ", "sig":"c94e74533b482aa8eeeb54ae72a5303e0b21f62909ca43c8ef06b0357412d6f8a92f96e1a205102753777fd25321a58fba3fb384eee114bd53ce6c06a1c22bab" }`

	decryptedGift2, err := ReceiveMessage(sk1, giftToSender)
	assert.NoError(t, err, "ReceiveMessage 2 should not return an error")
	assert.NotEmpty(t, decryptedGift2, "Gift wrap2 should not be empty")

	//fmt.Println(decryptedGift1)
	//fmt.Println(decryptedGift2)

}

package types

import (
	"encoding/hex"
	"testing"

	"github.com/scroll-tech/go-ethereum/common"
	"github.com/scroll-tech/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestAuthMessageSignAndVerify(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	publicKeyHex := common.Bytes2Hex(crypto.CompressPubkey(&privateKey.PublicKey))

	var authMsg LoginParameter
	t.Run("sign", func(t *testing.T) {
		authMsg = LoginParameter{
			Message: Message{
				ProverName:    "test1",
				ProverVersion: "v0.0.1",
				Challenge:     "abcdef",
				ProverTypes:   []ProverType{ProverTypeBatch},
				VKs:           []string{"vk1", "vk2"},
			},
			PublicKey: publicKeyHex,
		}

		err = authMsg.SignWithKey(privateKey)
		assert.NoError(t, err)
	})

	t.Run("valid verify", func(t *testing.T) {
		ok, verifyErr := authMsg.Verify()
		assert.True(t, ok)
		assert.NoError(t, verifyErr)
	})

	t.Run("invalid verify", func(t *testing.T) {
		authMsg.Message.Challenge = "abcdefgh"
		ok, verifyErr := authMsg.Verify()
		assert.False(t, ok)
		assert.NoError(t, verifyErr)
	})
}

// TestGenerateSignature this unit test isn't for test, just generate the signature for manually test.
func TestGenerateSignature(t *testing.T) {
	privateKeyHex := "8b8df68fddf7ee2724b79ccbd07799909d59b4dd4f4df3f6ecdc4fb8d56bdf4c"
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	assert.Nil(t, err)
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	assert.NoError(t, err)
	assert.NoError(t, err)
	publicKeyHex := common.Bytes2Hex(crypto.CompressPubkey(&privateKey.PublicKey))

	t.Log("publicKey: ", publicKeyHex)

	authMsg := LoginParameter{
		Message: Message{
			ProverName:    "test",
			ProverVersion: "v4.1.115-4dd11c6-000000-000000",
			Challenge:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTk1NjkyNDAsIm9yaWdfaWF0IjoxNzE5NTY1NjQwLCJyYW5kb20iOiJPRExnNEZtUW1MOEwzTDRvZ3BMcnl6c09EN1ZXd0FoNmd3bVpzVURJV3M0PSJ9.3Oq7fDtFnKGbPyjc8fslzfftyzreQbi-lAr0_HFy54w",
			ProverTypes:   []ProverType{ProverTypeChunk},
			VKs:           []string{"mock_chunk_vk"},
		},
		PublicKey: publicKeyHex,
	}
	err = authMsg.SignWithKey(privateKey)
	assert.NoError(t, err)
	t.Log("signature: ", authMsg.Signature)

	verify, err := authMsg.Verify()
	assert.NoError(t, err)
	assert.True(t, verify)
}
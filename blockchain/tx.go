package blockchain

import (
	"blockchain-go/wallet"
	"bytes"
	"encoding/gob"
)

type TXInput struct {
	ID        []byte // Reference to the transaction containing the output
	Out       int    // Index of the output in the referenced transaction
	Signature []byte // Script/signature to unlock the output
	PubKey    []byte // Public key of the sender
}

type TXOutput struct {
	Value      int    // Amount of coins
	PubKeyHash []byte // Script to lock the output
}

type TXOutputs struct {
	Outputs []*TXOutput // List of transaction outputs
}

func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey) // Hash of the public key
	// Compare the locking hash with the provided public key hash
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) isLockedWithKey(pubKeyHash []byte) bool {
	// Compare the output's public key hash with the provided public key hash
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func (outs TXOutputs) SerializeOutputs() []byte {
	var buffer bytes.Buffer
	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	Handle(err)

	return buffer.Bytes()
}

func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&outputs)
	Handle(err)

	return outputs
}

package blockchain

import (
	"blockchain-go/wallet"
	"bytes"
	"encoding/gob"
)

// TXInput représente une entrée de transaction
type TXInput struct {
	ID        []byte // Référence à la transaction contenant la sortie
	Out       int    // Index de la sortie dans la transaction référencée
	Signature []byte // Signature pour débloquer la sortie
	PubKey    []byte // Clé publique de l'expéditeur
}

// TXOutput représente une sortie de transaction
type TXOutput struct {
	Value      int    // Montant de coins
	PubKeyHash []byte // Script pour verrouiller la sortie
}

// TXOutputs représente une liste de sorties de transaction
type TXOutputs struct {
	Outputs []*TXOutput // Liste des sorties de transaction
}

// NewTXOutput crée une nouvelle sortie de transaction
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

// UsesKey vérifie si l'entrée utilise la clé publique donnée
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey) // Hash of the public key
	// Compare the locking hash with the provided public key hash
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// Lock verrouille la sortie avec une adresse donnée
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

// isLockedWithKey vérifie si la sortie est verrouillée avec la clé publique donnée
func (out *TXOutput) isLockedWithKey(pubKeyHash []byte) bool {
	// Compare the output's public key hash with the provided public key hash
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// SerializeOutputs sérialise une liste de sorties de transaction
func (outs TXOutputs) SerializeOutputs() []byte {
	var buffer bytes.Buffer
	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	Handle(err)

	return buffer.Bytes()
}

// DeserializeOutputs désérialise des données en une liste de sorties de transaction
func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&outputs)
	Handle(err)

	return outputs
}

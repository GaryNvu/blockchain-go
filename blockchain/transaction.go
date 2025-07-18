package blockchain

import (
	"blockchain-go/wallet"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

// Transaction représente une transaction dans la blockchain
type Transaction struct {
	ID      []byte     // Unique identifier of the transaction
	Inputs  []TXInput  // List of transaction inputs
	Outputs []TXOutput // List of transaction outputs
}

// Serialize sérialise une transaction en bytes
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	Handle(err)

	return encoded.Bytes()
}

// DeserializeTransaction désérialise des bytes en transaction
func DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	Handle(err)
	return transaction
}

// Hash calcule le hash SHA256 d'une transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{} // Set ID to nil to avoid hashing it again

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// CoinbaseTx crée une transaction coinbase (récompense de minage)
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 24)
		_, err := rand.Read(randData)
		Handle(err)
		data = fmt.Sprintf("%x", randData)
	}

	txIn := TXInput{[]byte{}, -1, nil, []byte(data)}
	txOut := NewTXOutput(20, to)
	tx := Transaction{nil, []TXInput{txIn}, []TXOutput{*txOut}}
	tx.ID = tx.Hash()

	return &tx
}

// IsCoinbase vérifie si une transaction est une transaction coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// Sign signe les entrées d'une transaction avec la clé privée donnée
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("ERROR: Previous transaction not found")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inId, in := range txCopy.Inputs {
		prevTx := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTx.Outputs[in.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		Handle(err)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[inId].Signature = signature
	}
}

// TrimmedCopy crée une copie de la transaction sans les signatures pour la signature
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TXInput{in.ID, in.Out, nil, nil})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TXOutput{out.Value, out.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// Verify vérifie les signatures d'une transaction
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("ERROR: Previous transaction not found")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, in := range tx.Inputs {
		prevTx := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTx.Outputs[in.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen / 2)])
		y.SetBytes(in.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) {
			return false
		}
	}

	return true
}

// NewTransaction crée une nouvelle transaction normale
func NewTransaction(w *wallet.Wallet, to string, amount int, UTXO *UTXOSet) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)
	fmt.Printf("Finding spendable outputs for address %s, amount needed: %d\n", w.Address(), amount)
	acc, validOutputs := UTXO.FindSpendableOutputs(pubKeyHash, amount)
	fmt.Printf("Found %d coins in spendable outputs\n", acc)

	if acc < amount {
		log.Panic("Error: not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {
			input := TXInput{txID, out, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}

	from := fmt.Sprintf("%s", w.Address())

	outputs = append(outputs, *NewTXOutput(amount, to))

	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXO.Blockchain.SignTransaction(&tx, w.PrivateKey)

	return &tx
}

// String retourne une représentation string de la transaction
func (tx *Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("-- Transaction %x:", tx.ID))
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("		Input %d:", i))
		lines = append(lines, fmt.Sprintf("			TXID: %x", input.ID))
		lines = append(lines, fmt.Sprintf("			Out: %d", input.Out))
		lines = append(lines, fmt.Sprintf("			Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("			PubKey: %x", input.PubKey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("		Output %d:", i))
		lines = append(lines, fmt.Sprintf("			Value: %d", output.Value))
		lines = append(lines, fmt.Sprintf("			PubKeyHash: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

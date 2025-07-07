package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// Represents a block in the blockchain
type Block struct {
	Timestamp    int64          // Timestamp of when the block was created
	Hash         []byte         // Hash of the current block
	Transactions []*Transaction // List of transactions contained in the block
	PrevHash     []byte         // Hash of the previous block in the chain
	Nonce        int            // Nonce used for the proof of work algorithm
	Height       int            // Height of the block in the blockchain
}

// Creates a Merkle root of all the block's transactions
// Returns a byte slice containing the hash of all transactions
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}
	tree := NewMerkleTree(txHashes)
	return tree.RootNode.Data
}

// Creates a new block with the given transactions and previous block hash
// It performs proof of work and returns the newly created block
func CreateBlock(txs []*Transaction, prevHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), []byte{}, txs, prevHash, 0, height}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Creates the first block of the blockchain with a given coinbase transaction
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{}, 0)
}

// Converts the block into a byte slice using gob encoding
// Returns the serialized block data
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

// Converts a byte slice back into a Block structure
// Returns a pointer to the deserialized block
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

// Utility function for error handling
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

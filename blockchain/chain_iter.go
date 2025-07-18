package blockchain

import "github.com/dgraph-io/badger"

// BlockChainIterator permet de parcourir la blockchain depuis le dernier bloc vers le premier
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// Iterator crée un nouvel itérateur pour parcourir la blockchain
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

// Next récupère le bloc suivant dans l'itération (vers les blocs plus anciens)
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err := item.ValueCopy(nil)
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}

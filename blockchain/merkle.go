package blockchain

import "crypto/sha256"

// MerkleTree représente un arbre de Merkle pour vérifier l'intégrité des transactions
type MerkleTree struct {
	RootNode *MerkleNode
}

// MerkleNode représente un nœud dans l'arbre de Merkle
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// NewMerkleNode crée un nouveau nœud de l'arbre de Merkle
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}

	node.Left = left
	node.Right = right
	return &node
}

// NewMerkleTree construit un arbre de Merkle à partir d'un ensemble de données
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1]) // Duplicate last element if odd number
	}

	for _, dat := range data {
		nodes = append(nodes, *NewMerkleNode(nil, nil, dat))
	}

	// Build the tree
	for range len(data) / 2 {
		var level []MerkleNode

		for i := 0; i < len(nodes); i += 2 {
			node := NewMerkleNode(&nodes[i], &nodes[i+1], nil)
			level = append(level, *node)
		}

		nodes = level
	}

	return &MerkleTree{&nodes[0]}
}

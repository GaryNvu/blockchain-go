package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

// Base58Encode encode des données en utilisant l'encodage Base58
func Base58Encode(input []byte) []byte {
	encoded := base58.Encode(input)
	return []byte(encoded)
}

// Base58Decode décode des données encodées en Base58
func Base58Decode(input []byte) []byte {
	decoded, err := base58.Decode(string(input[:]))
	if err != nil {
		log.Panic(err)
	}
	return decoded
}

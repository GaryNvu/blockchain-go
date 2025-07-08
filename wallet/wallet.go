package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

const (
	checkSumLen = 4          // Length of the checksum in bytes
	versionLen  = byte(0x00) // Length of the version byte
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey // Private key for signing transactions
	PublicKey  []byte           // Public key for verifying signatures
}

// Structure pour la sérialisation
type SerializableWallet struct {
	PrivateKeyD []byte `gob:"private_key_d"`
	PublicKeyX  []byte `gob:"public_key_x"`
	PublicKeyY  []byte `gob:"public_key_y"`
	PublicKey   []byte `gob:"public_key"`
}

// Convertit un Wallet en structure sérialisable
func (w *Wallet) ToSerializable() SerializableWallet {
	return SerializableWallet{
		PrivateKeyD: w.PrivateKey.D.Bytes(),
		PublicKeyX:  w.PrivateKey.PublicKey.X.Bytes(),
		PublicKeyY:  w.PrivateKey.PublicKey.Y.Bytes(),
		PublicKey:   w.PublicKey,
	}
}

// Crée un Wallet à partir d'une structure sérialisable
func FromSerializable(sw SerializableWallet) *Wallet {
	curve := elliptic.P256()

	d := big.NewInt(0)
	d.SetBytes(sw.PrivateKeyD)

	x := big.NewInt(0)
	x.SetBytes(sw.PublicKeyX)

	y := big.NewInt(0)
	y.SetBytes(sw.PublicKeyY)

	privateKey := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: d,
	}

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  sw.PublicKey,
	}
}

func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{versionLen}, pubHash...) // Prepend version byte to the hash
	checksum := checksum(versionedHash)                     // Calculate checksum
	fullHash := append(versionedHash, checksum...)          // Append checksum to the payload
	address := Base58Encode(fullHash)                       // Encode the full hash using base58

	return address
}

func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))                        // Decode the address from base58
	actualChecksum := pubKeyHash[len(pubKeyHash)-checkSumLen:]         // Extract the checksum from the address
	version := pubKeyHash[0]                                           // Extract the version byte from the address
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checkSumLen]           // Remove the version byte and checksum from the hash
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...)) // Calculate the target checksum

	return bytes.Compare(actualChecksum, targetChecksum) == 0 // Compare the actual checksum with the target checksum
}

// NewKeyPair generates a new ECDSA key pair and returns the private key and public key
// The public key is represented as a byte slice containing the X and Y coordinates of the elliptic curve point
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...) // Concatenate X and Y coordinates of the public key
	return *private, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

// Hash the public key using SHA-256 followed by RIPEMD-160
func PublicKeyHash(publicKey []byte) []byte {
	pubHash := sha256.Sum256(publicKey)
	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}

func checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:checkSumLen]
}

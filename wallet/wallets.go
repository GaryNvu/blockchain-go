package wallet

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

const walletFile = "./tmp/wallets_%s.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

// CreateWallets crée ou charge une collection de wallets pour un nœud donné
func CreateWallets(nodeId string) (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile(nodeId)

	return &wallets, err
}

// GetWallet récupère un wallet spécifique par son adresse
func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// GetAllAddresses retourne toutes les adresses des wallets
func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// AddWallet crée un nouveau wallet et l'ajoute à la collection
func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet
	return address
}

// LoadFile charge les wallets depuis un fichier
func (ws *Wallets) LoadFile(nodeId string) error {
	walletFile := fmt.Sprintf(walletFile, nodeId)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		return err
	}

	// Charger les wallets sérialisables
	var serializableWallets map[string]SerializableWallet
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&serializableWallets)
	if err != nil {
		return err
	}

	// Convertir en wallets normaux
	ws.Wallets = make(map[string]*Wallet)
	for address, sw := range serializableWallets {
		ws.Wallets[address] = FromSerializable(sw)
	}

	return nil
}

// SaveFile sauvegarde les wallets dans un fichier
func (ws *Wallets) SaveFile(nodeId string) {
	var content bytes.Buffer
	walletFile := fmt.Sprintf(walletFile, nodeId)

	// Convertir les wallets en structure sérialisable
	serializableWallets := make(map[string]SerializableWallet)
	for address, wallet := range ws.Wallets {
		serializableWallets[address] = wallet.ToSerializable()
	}

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(serializableWallets)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

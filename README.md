# blockchain-go

Creation of a blockchain with Golang

## Description

Ce projet implémente une blockchain complète en Go avec support multi-nœuds, propagation de transactions, mining et persistence des données. La blockchain utilise un algorithme de Proof of Work et permet la communication entre plusieurs nœuds via TCP.

## Fonctionnalités

- ✅ Création et gestion de wallets
- ✅ Transactions avec signatures cryptographiques
- ✅ Mining avec Proof of Work
- ✅ Réseau multi-nœuds avec propagation de transactions
- ✅ Persistence des données avec BadgerDB
- ✅ UTXO set pour optimiser les performances
- ✅ Interface CLI complète

## Structure du projet

```
blockchain-go/
├── blockchain/          # Core blockchain logic
├── cli/                # Interface en ligne de commande
├── network/            # Logique réseau et propagation
├── wallet/             # Gestion des wallets et cryptographie
├── tmp/                # Données temporaires (wallets et blocks)
├── setup_blockchain.ps1 # Script de configuration (Windows)
├── setup_blockchain.sh  # Script de configuration (Linux/macOS)
└── main.go             # Point d'entrée principal
```

## Prérequis

- Go 1.19 ou plus récent
- Git

## Installation

1. Clonez le repository :
```bash
git clone <repository-url>
cd blockchain-go
```

2. Installez les dépendances :
```bash
go mod tidy
```

## Configuration rapide

### Windows (PowerShell)

Utilisez le script de configuration automatique :

```powershell
.\setup_blockchain.ps1
```

### Linux/macOS (Bash)

```bash
chmod +x setup_blockchain.sh
./setup_blockchain.sh
```

Ce script va automatiquement :
- Nettoyer les données existantes
- Créer des wallets pour les nœuds 3000, 3001 et 3002
- Créer la blockchain avec l'adresse du nœud 3000
- Copier les données blockchain vers tous les nœuds
- Vérifier les soldes initiaux
- Créer des scripts de démarrage pour chaque nœud
- Mettre à jour le fichier `wallets.json`

## Démarrage des nœuds

Commandes manuelles - Remplacer le X par la version du noeud désirée (0,1,2)
Le noeud 3000 est central

**Terminal X (Nœud 300X):**
```bash
export NODE_ID=300X  # Linux/macOS
# ou
$env:NODE_ID = '300X'  # Windows PowerShell

go run main.go startnode
```

## Test de propagation de transactions

Pour envoyez des transactions, veillez à ce que le noeud sur lequel vous travaillez soit éteint.

```bash
export NODE_ID=3000
go run main.go send -from [ADRESSE_3000] -to [ADRESSE_3001] -amount 10 -mine
```

### Windows :
```powershell
$env:NODE_ID = '3000'
go run main.go send -from [ADRESSE_3000] -to [ADRESSE_3001] -amount 10 -mine
```

### Linux/macOS :
```bash
export NODE_ID=3000
go run main.go send -from [ADRESSE_3000] -to [ADRESSE_3001] -amount 10 -mine
```

### Logs de propagation attendus :

Vous devriez voir des logs similaires à sur les autres noeuds :
```
Recevied a new block!
Added block [hash]
Received block command
Recevied a new block!
Added block [hash]
```

## Vérification des données

### Consulter les soldes

```bash
# Nœud 300X - Remplacer le X par la version de Node désirée
export NODE_ID=300X  # Linux/macOS ou $env:NODE_ID = '3000' pour Windows
go run main.go getbalance -address [ADRESSE_300X]

### Afficher la blockchain

```bash
export NODE_ID=300X
go run main.go printchain
```

### Lister les adresses

```bash
export NODE_ID=300X
go run main.go listaddresses
```

## Arrêt des nœuds

Pour arrêter proprement les nœuds, utilisez `Ctrl+C` dans chaque terminal. Les nœuds sauvegarderont automatiquement leurs données.

## Commandes CLI disponibles

- `createwallet` - Créer un nouveau wallet
- `listaddresses` - Lister toutes les adresses
- `createblockchain -address ADDRESS` - Créer une nouvelle blockchain
- `getbalance -address ADDRESS` - Obtenir le solde d'une adresse
- `send -from FROM -to TO -amount AMOUNT [-mine]` - Envoyer des tokens
- `printchain` - Afficher tous les blocs
- `reindexutxo` - Reconstruire l'UTXO set
- `startnode [-miner ADDRESS]` - Démarrer un nœud réseau

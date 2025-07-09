#!/bin/bash
# Script pour configurer la blockchain multi-noeuds
# Auteur: Assistant
# Date: 2025-07-10

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
GRAY='\033[0;37m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Configuration de la blockchain multi-noeuds ===${NC}"

# Fonction pour attendre l'entree utilisateur
wait_user_input() {
    echo -e "${YELLOW}$1${NC}"
    read -p "Appuyez sur Entree pour continuer..."
}

# Nettoyer tout le contenu du dossier tmp
echo -e "${CYAN}Nettoyage complet du dossier tmp...${NC}"
if [ -d "tmp" ]; then
    rm -rf tmp/*
    echo -e "${RED}Suppression de tout le contenu de tmp/${NC}"
else
    echo -e "${GREEN}Creation du dossier tmp${NC}"
    mkdir -p tmp
fi

echo -e "\n${GREEN}Etape 1: Creation des wallets${NC}"

# Creer le wallet pour le noeud 3000
echo -e "${CYAN}Creation du wallet pour le noeud 3000...${NC}"
export NODE_ID=3000
output3000=$(go run main.go createwallet)
echo -e "${WHITE}$output3000${NC}"
address3000=$(echo "$output3000" | grep "New address is:" | sed 's/New address is: //')
echo -e "${YELLOW}Adresse du noeud 3000: $address3000${NC}"

# Creer le wallet pour le noeud 3001
echo -e "\n${CYAN}Creation du wallet pour le noeud 3001...${NC}"
export NODE_ID=3001
output3001=$(go run main.go createwallet)
echo -e "${WHITE}$output3001${NC}"
address3001=$(echo "$output3001" | grep "New address is:" | sed 's/New address is: //')
echo -e "${YELLOW}Adresse du noeud 3001: $address3001${NC}"

# Creer le wallet pour le noeud 3002
echo -e "\n${CYAN}Creation du wallet pour le noeud 3002...${NC}"
export NODE_ID=3002
output3002=$(go run main.go createwallet)
echo -e "${WHITE}$output3002${NC}"
address3002=$(echo "$output3002" | grep "New address is:" | sed 's/New address is: //')
echo -e "${YELLOW}Adresse du noeud 3002: $address3002${NC}"

echo -e "\n${GREEN}Etape 2: Creation de la blockchain${NC}"

# Creer la blockchain avec l'adresse du noeud 3000
echo -e "${CYAN}Creation de la blockchain avec l'adresse du noeud 3000...${NC}"
export NODE_ID=3000
blockchainOutput=$(go run main.go createblockchain -address $address3000)
echo -e "${WHITE}$blockchainOutput${NC}"

echo -e "\n${GREEN}Etape 3: Copie des donnees blockchain${NC}"

# Verifier que le dossier blocks_3000 existe
if [ ! -d "tmp/blocks_3000" ]; then
    echo -e "${RED}ERREUR: Le dossier tmp/blocks_3000 n'existe pas!${NC}"
    exit 1
fi

# Copier blocks_3000 vers blocks_3001
echo -e "${CYAN}Copie de tmp/blocks_3000 vers tmp/blocks_3001...${NC}"
cp -r tmp/blocks_3000 tmp/blocks_3001
echo -e "${GREEN}Copie terminee: tmp/blocks_3001${NC}"

# Copier blocks_3000 vers blocks_3002
echo -e "${CYAN}Copie de tmp/blocks_3000 vers tmp/blocks_3002...${NC}"
cp -r tmp/blocks_3000 tmp/blocks_3002
echo -e "${GREEN}Copie terminee: tmp/blocks_3002${NC}"

# Copier blocks_3000 vers blocks_gen
echo -e "${CYAN}Copie de tmp/blocks_3000 vers tmp/blocks_gen...${NC}"
cp -r tmp/blocks_3000 tmp/blocks_gen
echo -e "${GREEN}Copie terminee: tmp/blocks_gen${NC}"

echo -e "\n${GREEN}Etape 4: Verification des soldes${NC}"

# Verifier le solde du noeud 3000
echo -e "${CYAN}Verification du solde du noeud 3000...${NC}"
export NODE_ID=3000
balance3000=$(go run main.go getbalance -address $address3000)
echo -e "${WHITE}$balance3000${NC}"

# Verifier le solde du noeud 3001
echo -e "\n${CYAN}Verification du solde du noeud 3001...${NC}"
export NODE_ID=3001
balance3001=$(go run main.go getbalance -address $address3001)
echo -e "${WHITE}$balance3001${NC}"

# Verifier le solde du noeud 3002
echo -e "\n${CYAN}Verification du solde du noeud 3002...${NC}"
export NODE_ID=3002
balance3002=$(go run main.go getbalance -address $address3002)
echo -e "${WHITE}$balance3002${NC}"

echo -e "\n${GREEN}Etape 5: Mise a jour du fichier wallets.json${NC}"

# Creer le fichier wallets.json avec les nouvelles adresses
cat > wallets.json << EOF
{
    "3000": "$address3000",
    "3001": "$address3001",
    "3002": "$address3002"
}
EOF

echo -e "${GREEN}Fichier wallets.json mis a jour avec les nouvelles adresses${NC}"

echo -e "\n${GREEN}Etape 6: Creation des scripts de demarrage${NC}"

# Creer les scripts bash pour demarrer les noeuds
cat > start_node_3000.sh << EOF
#!/bin/bash
echo "Demarrage du noeud 3000 (noeud central)"
echo "Adresse: $address3000"
echo
export NODE_ID=3000
go run main.go startnode
EOF

cat > start_node_3001.sh << EOF
#!/bin/bash
echo "Demarrage du noeud 3001"
echo "Adresse: $address3001"
echo
export NODE_ID=3001
go run main.go startnode
EOF

cat > start_node_3002.sh << EOF
#!/bin/bash
echo "Demarrage du noeud 3002"
echo "Adresse: $address3002"
echo
export NODE_ID=3002
go run main.go startnode
EOF

# Rendre les scripts executables
chmod +x start_node_3000.sh
chmod +x start_node_3001.sh
chmod +x start_node_3002.sh

echo -e "${GREEN}Scripts de demarrage crees:${NC}"
echo -e "${CYAN}- start_node_3000.sh (noeud central)${NC}"
echo -e "${CYAN}- start_node_3001.sh${NC}"
echo -e "${CYAN}- start_node_3002.sh${NC}"

# Creer un script de test de transaction
cat > test_transaction.sh << EOF
#!/bin/bash
echo "Test de transaction"
echo "De: $address3000"
echo "Vers: $address3001"
echo "Montant: 10"
echo
export NODE_ID=3000
go run main.go send -from $address3000 -to $address3001 -amount 10
EOF

chmod +x test_transaction.sh
echo -e "${CYAN}- test_transaction.sh (script de test)${NC}"

echo -e "\n${GREEN}=== CONFIGURATION TERMINEE ===${NC}"
echo -e "${YELLOW}Resume des adresses creees:${NC}"
echo -e "${WHITE}Noeud 3000: $address3000${NC}"
echo -e "${WHITE}Noeud 3001: $address3001${NC}"
echo -e "${WHITE}Noeud 3002: $address3002${NC}"

echo -e "\n${YELLOW}Pour demarrer les noeuds facilement:${NC}"
echo -e "${RED}IMPORTANT: Demarrez d'abord le noeud central (3000)!${NC}"
echo -e "${CYAN}1. Dans un terminal: ./start_node_3000.sh${NC}"
echo -e "${CYAN}2. Dans un autre terminal: ./start_node_3001.sh${NC}"
echo -e "${CYAN}3. Dans un troisieme terminal: ./start_node_3002.sh${NC}"
echo -e "${CYAN}4. Pour tester: ./test_transaction.sh${NC}"

echo -e "\n${YELLOW}Ou utilisez les commandes suivantes dans des terminaux separes:${NC}"
echo -e "${GRAY}Terminal 1: export NODE_ID=3000; go run main.go startnode${NC}"
echo -e "${GRAY}Terminal 2: export NODE_ID=3001; go run main.go startnode${NC}"
echo -e "${GRAY}Terminal 3: export NODE_ID=3002; go run main.go startnode${NC}"

echo -e "\n${YELLOW}Pour envoyer une transaction:${NC}"
echo -e "${RED}IMPORTANT: Assurez-vous d'abord que les noeuds sont demarres!${NC}"
echo -e "${CYAN}export NODE_ID=3000; go run main.go send -from $address3000 -to $address3001 -amount 10${NC}"
echo -e "${CYAN}export NODE_ID=3000; go run main.go send -from $address3000 -to $address3002 -amount 10${NC}"

echo -e "\n${YELLOW}Pour tester sans le reseau (minage local):${NC}"
echo -e "${CYAN}export NODE_ID=3000; go run main.go send -from $address3000 -to $address3001 -amount 10 -mine${NC}"

echo -e "\n${YELLOW}Commandes utiles:${NC}"
echo -e "${CYAN}Verifier le solde: export NODE_ID=3000; go run main.go getbalance -address $address3000${NC}"
echo -e "${CYAN}Lister les adresses: export NODE_ID=3000; go run main.go listaddresses${NC}"
echo -e "${CYAN}Afficher la blockchain: export NODE_ID=3000; go run main.go printchain${NC}"

echo -e "\n${GREEN}Script termine avec succes!${NC}"

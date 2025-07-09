# Script pour configurer la blockchain multi-noeuds
# Auteur: Assistant
# Date: 2025-07-10

Write-Host "=== Configuration de la blockchain multi-noeuds ===" -ForegroundColor Green

# Fonction pour attendre l'entrée utilisateur
function Wait-UserInput {
    param([string]$message)
    Write-Host $message -ForegroundColor Yellow
    Read-Host "Appuyez sur Entrée pour continuer..."
}

# Nettoyer tout le contenu du dossier tmp
Write-Host "Nettoyage complet du dossier tmp..." -ForegroundColor Cyan
if (Test-Path "tmp") {
    Remove-Item -Path "tmp\*" -Recurse -Force
    Write-Host "Suppression de tout le contenu de tmp\" -ForegroundColor Red
} else {
    Write-Host "Création du dossier tmp" -ForegroundColor Green
    New-Item -Path "tmp" -ItemType Directory -Force
}

Write-Host "`nEtape 1: Création des wallets" -ForegroundColor Green

# Créer le wallet pour le noeud 3000
Write-Host "Création du wallet pour le noeud 3000..." -ForegroundColor Cyan
$env:NODE_ID = '3000'
$output3000 = go run main.go createwallet
Write-Host $output3000 -ForegroundColor White
$address3000 = ($output3000 | Select-String "New address is: (.+)").Matches[0].Groups[1].Value
Write-Host "Adresse du noeud 3000: $address3000" -ForegroundColor Yellow

# Créer le wallet pour le noeud 3001
Write-Host "`nCréation du wallet pour le noeud 3001..." -ForegroundColor Cyan
$env:NODE_ID = '3001'
$output3001 = go run main.go createwallet
Write-Host $output3001 -ForegroundColor White
$address3001 = ($output3001 | Select-String "New address is: (.+)").Matches[0].Groups[1].Value
Write-Host "Adresse du noeud 3001: $address3001" -ForegroundColor Yellow

# Créer le wallet pour le noeud 3002
Write-Host "`nCréation du wallet pour le noeud 3002..." -ForegroundColor Cyan
$env:NODE_ID = '3002'
$output3002 = go run main.go createwallet
Write-Host $output3002 -ForegroundColor White
$address3002 = ($output3002 | Select-String "New address is: (.+)").Matches[0].Groups[1].Value
Write-Host "Adresse du noeud 3002: $address3002" -ForegroundColor Yellow

Write-Host "`nEtape 2: Création de la blockchain" -ForegroundColor Green

# Créer la blockchain avec l'adresse du noeud 3000
Write-Host "Création de la blockchain avec l'adresse du noeud 3000..." -ForegroundColor Cyan
$env:NODE_ID = '3000'
$blockchainOutput = go run main.go createblockchain -address $address3000
Write-Host $blockchainOutput -ForegroundColor White

Write-Host "`nEtape 3: Copie des données blockchain" -ForegroundColor Green

# Vérifier que le dossier blocks_3000 existe
if (!(Test-Path "tmp\blocks_3000")) {
    Write-Host "ERREUR: Le dossier tmp\blocks_3000 n'existe pas!" -ForegroundColor Red
    exit 1
}

# Copier blocks_3000 vers blocks_3001
Write-Host "Copie de tmp\blocks_3000 vers tmp\blocks_3001..." -ForegroundColor Cyan
Copy-Item -Path "tmp\blocks_3000" -Destination "tmp\blocks_3001" -Recurse
Write-Host "Copie terminée: tmp\blocks_3001" -ForegroundColor Green

# Copier blocks_3000 vers blocks_3002
Write-Host "Copie de tmp\blocks_3000 vers tmp\blocks_3002..." -ForegroundColor Cyan
Copy-Item -Path "tmp\blocks_3000" -Destination "tmp\blocks_3002" -Recurse
Write-Host "Copie terminée: tmp\blocks_3002" -ForegroundColor Green

# Copier blocks_3000 vers blocks_gen
Write-Host "Copie de tmp\blocks_3000 vers tmp\blocks_gen..." -ForegroundColor Cyan
Copy-Item -Path "tmp\blocks_3000" -Destination "tmp\blocks_gen" -Recurse
Write-Host "Copie terminée: tmp\blocks_gen" -ForegroundColor Green

Write-Host "`nEtape 4: Vérification des soldes" -ForegroundColor Green

# Vérifier le solde du noeud 3000
Write-Host "Vérification du solde du noeud 3000..." -ForegroundColor Cyan
$env:NODE_ID = '3000'
$balance3000 = go run main.go getbalance -address $address3000
Write-Host $balance3000 -ForegroundColor White

# Vérifier le solde du noeud 3001
Write-Host "`nVérification du solde du noeud 3001..." -ForegroundColor Cyan
$env:NODE_ID = '3001'
$balance3001 = go run main.go getbalance -address $address3001
Write-Host $balance3001 -ForegroundColor White

# Vérifier le solde du noeud 3002
Write-Host "`nVérification du solde du noeud 3002..." -ForegroundColor Cyan
$env:NODE_ID = '3002'
$balance3002 = go run main.go getbalance -address $address3002
Write-Host $balance3002 -ForegroundColor White

Write-Host "`nEtape 5: Mise à jour du fichier wallets.json" -ForegroundColor Green

# Créer le fichier wallets.json avec les nouvelles adresses
$walletsJson = @{
    "3000" = $address3000
    "3001" = $address3001
    "3002" = $address3002
} | ConvertTo-Json -Depth 2

$walletsJson | Out-File -FilePath "wallets.json" -Encoding utf8
Write-Host "Fichier wallets.json mis à jour avec les nouvelles adresses" -ForegroundColor Green

Write-Host "`n=== CONFIGURATION TERMINÉE ===" -ForegroundColor Green
Write-Host "Résumé des adresses créées:" -ForegroundColor Yellow
Write-Host "Noeud 3000: $address3000" -ForegroundColor White
Write-Host "Noeud 3001: $address3001" -ForegroundColor White
Write-Host "Noeud 3002: $address3002" -ForegroundColor White

Write-Host "`nPour démarrer les noeuds, utilisez les commandes suivantes dans des terminaux séparés:" -ForegroundColor Yellow
Write-Host "Terminal 1: `$env:NODE_ID = '3000'; go run main.go startnode" -ForegroundColor Cyan
Write-Host "Terminal 2: `$env:NODE_ID = '3001'; go run main.go startnode" -ForegroundColor Cyan
Write-Host "Terminal 3: `$env:NODE_ID = '3002'; go run main.go startnode" -ForegroundColor Cyan

Write-Host "`nPour envoyer une transaction:" -ForegroundColor Yellow
Write-Host "`$env:NODE_ID = '3000'; go run main.go send -from $address3000 -to $address3001 -amount 10" -ForegroundColor Cyan
Write-Host "`$env:NODE_ID = '3000'; go run main.go send -from $address3000 -to $address3002 -amount 10" -ForegroundColor Cyan
Write-Host "`nScript terminé avec succès!" -ForegroundColor Green

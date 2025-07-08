# Script d'initialisation automatique pour la blockchain
param(
    [string]$NodeId = "3000",
    [switch]$Clean = $false
)

Write-Host "🚀 Initialisation de la blockchain (Node $NodeId)" -ForegroundColor Green

# Nettoyer les données existantes si demandé
if ($Clean) {
    Write-Host "🧹 Nettoyage des données existantes..." -ForegroundColor Yellow
    Remove-Item -Path "tmp\blocks_$NodeId" -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item -Path "tmp\wallets_$NodeId.data" -Force -ErrorAction SilentlyContinue
}

# Définir NODE_ID
$env:NODE_ID = $NodeId
Write-Host "📝 NODE_ID défini sur: $NodeId" -ForegroundColor Cyan

# Créer le répertoire tmp si nécessaire
if (!(Test-Path "tmp")) {
    New-Item -ItemType Directory -Path "tmp" | Out-Null
    Write-Host "📁 Répertoire tmp créé" -ForegroundColor Cyan
}

# Créer un wallet principal
Write-Host "🔑 Création du wallet principal..." -ForegroundColor Cyan
$walletOutput = go run main.go createwallet 2>&1
$mainAddress = ($walletOutput | Select-String "New address is: (.+)" | ForEach-Object { $_.Matches[0].Groups[1].Value })

if ($mainAddress) {
    Write-Host "✅ Wallet créé: $mainAddress" -ForegroundColor Green
    
    # Créer la blockchain
    Write-Host "⛓️  Création de la blockchain..." -ForegroundColor Cyan
    go run main.go createblockchain -address $mainAddress | Out-Null
    Write-Host "✅ Blockchain créée avec succès" -ForegroundColor Green
    
    # Créer un deuxième wallet pour les tests
    Write-Host "🔑 Création d'un wallet secondaire..." -ForegroundColor Cyan
    $secondWalletOutput = go run main.go createwallet 2>&1
    $secondAddress = ($secondWalletOutput | Select-String "New address is: (.+)" | ForEach-Object { $_.Matches[0].Groups[1].Value })
    
    if ($secondAddress) {
        Write-Host "✅ Wallet secondaire créé: $secondAddress" -ForegroundColor Green
    }
    
    # Afficher le résumé
    Write-Host "`n🎉 Configuration terminée!" -ForegroundColor Green
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
    Write-Host "📊 Node ID: $NodeId" -ForegroundColor White
    Write-Host "🏦 Wallet principal: $mainAddress" -ForegroundColor White
    Write-Host "👤 Wallet secondaire: $secondAddress" -ForegroundColor White
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray
    Write-Host "`n💡 Commandes utiles:" -ForegroundColor Yellow
    Write-Host "   Voir la balance: go run main.go balance -address $mainAddress" -ForegroundColor Cyan
    Write-Host "   Envoyer des tokens: go run main.go send -from $mainAddress -to $secondAddress -amount 10 -mine" -ForegroundColor Cyan
    Write-Host "   Démarrer le nœud: go run main.go startnode" -ForegroundColor Cyan
    Write-Host "   Voir la blockchain: go run main.go printchain" -ForegroundColor Cyan
    
} else {
    Write-Host "❌ Erreur lors de la création du wallet" -ForegroundColor Red
    exit 1
}

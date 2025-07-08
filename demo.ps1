# Script de demonstration simple pour blockchain
param(
    [string]$Action = "demo",
    [string]$NodeId = "3000"
)

$env:NODE_ID = $NodeId

Write-Host "Blockchain Demo Script" -ForegroundColor Green
Write-Host "==============================================" -ForegroundColor Gray
Write-Host "Node ID: $NodeId" -ForegroundColor Cyan

switch ($Action) {
    "demo" {
        Write-Host "Demarrage de la demonstration complete..." -ForegroundColor Yellow
        
        # Nettoyer les donnees existantes
        Write-Host "Nettoyage des donnees..." -ForegroundColor Cyan
        Remove-Item -Path "tmp\blocks_$NodeId" -Recurse -Force -ErrorAction SilentlyContinue
        Remove-Item -Path "tmp\wallets_$NodeId.data" -Force -ErrorAction SilentlyContinue
        
        # Creer le premier wallet
        Write-Host "Creation du wallet principal..." -ForegroundColor Cyan
        $output1 = go run main.go createwallet 2>&1
        $wallet1 = ($output1 | Select-String "New address is: (.+)" | ForEach-Object { $_.Matches[0].Groups[1].Value })
        Write-Host "Wallet 1: $wallet1" -ForegroundColor Green
        
        # Creer la blockchain
        Write-Host "Creation de la blockchain..." -ForegroundColor Cyan
        go run main.go createblockchain -address $wallet1 | Out-Null
        Write-Host "Blockchain creee" -ForegroundColor Green
        
        # Creer un second wallet
        Write-Host "Creation du wallet secondaire..." -ForegroundColor Cyan
        $output2 = go run main.go createwallet 2>&1
        $wallet2 = ($output2 | Select-String "New address is: (.+)" | ForEach-Object { $_.Matches[0].Groups[1].Value })
        Write-Host "Wallet 2: $wallet2" -ForegroundColor Green
        
        # Verifier les balances
        Write-Host "`nBalances initiales:" -ForegroundColor Yellow
        $balance1 = go run main.go getbalance -address $wallet1
        $balance2 = go run main.go getbalance -address $wallet2
        Write-Host "   Wallet 1: $balance1" -ForegroundColor White
        Write-Host "   Wallet 2: $balance2" -ForegroundColor White
        
        # Envoyer une transaction
        Write-Host "`nEnvoi de 10 tokens du wallet 1 vers le wallet 2..." -ForegroundColor Cyan
        go run main.go send -from $wallet1 -to $wallet2 -amount 10 -mine | Out-Null
        Write-Host "Transaction envoyee" -ForegroundColor Green
        
        # Verifier les nouvelles balances
        Write-Host "`nBalances apres transaction:" -ForegroundColor Yellow
        $balance1 = go run main.go getbalance -address $wallet1
        $balance2 = go run main.go getbalance -address $wallet2
        Write-Host "   Wallet 1: $balance1" -ForegroundColor White
        Write-Host "   Wallet 2: $balance2" -ForegroundColor White
        
        Write-Host "`nDemonstration terminee!" -ForegroundColor Green
        Write-Host "Vous pouvez maintenant utiliser:" -ForegroundColor Yellow
        Write-Host "   - go run main.go startnode (Demarrer le noeud)" -ForegroundColor Cyan
    }
    
    "clean" {
        Write-Host "Nettoyage de toutes les donnees..." -ForegroundColor Yellow
        Remove-Item -Path "tmp\*" -Recurse -Force -ErrorAction SilentlyContinue
        Write-Host "Nettoyage termine" -ForegroundColor Green
    }
    
    default {
        Write-Host "Usage:" -ForegroundColor Yellow
        Write-Host "  .\demo.ps1 demo     - Demonstration complete" -ForegroundColor Cyan
        Write-Host "  .\demo.ps1 clean    - Nettoyer les donnees" -ForegroundColor Cyan
    }
}

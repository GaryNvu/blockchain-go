package main

import (
	cli "blockchain-go/cli"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type WebServer struct {
	nodeID string
}

type WalletInfo struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

type TransactionRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}

func NewWebServer() *WebServer {
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		nodeID = "3000"
		os.Setenv("NODE_ID", nodeID)
	}

	return &WebServer{
		nodeID: nodeID,
	}
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>🔗 Blockchain Interface</title>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        .container {
            background: white;
            border-radius: 10px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.3);
            padding: 30px;
            margin-bottom: 20px;
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
            font-size: 2.5em;
        }
        h2 {
            color: #555;
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
            color: #555;
        }
        input, select, button {
            width: 100%;
            padding: 12px;
            border: 2px solid #ddd;
            border-radius: 6px;
            font-size: 16px;
            box-sizing: border-box;
        }
        input:focus, select:focus {
            border-color: #667eea;
            outline: none;
        }
        button {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            cursor: pointer;
            font-weight: bold;
            margin-top: 10px;
            transition: transform 0.2s;
        }
        button:hover {
            transform: translateY(-2px);
        }
        .wallet-list {
            background: #f8f9fa;
            border-radius: 6px;
            padding: 15px;
            margin-top: 20px;
        }
        .wallet-item {
            background: white;
            padding: 15px;
            margin-bottom: 10px;
            border-radius: 6px;
            border-left: 4px solid #667eea;
        }
        .address {
            font-family: monospace;
            background: #e9ecef;
            padding: 5px;
            border-radius: 4px;
            font-size: 14px;
            word-break: break-all;
        }
        .balance {
            color: #28a745;
            font-weight: bold;
            font-size: 18px;
        }
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
            gap: 20px;
        }
        .alert {
            padding: 15px;
            border-radius: 6px;
            margin-bottom: 20px;
        }
        .alert-success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        .alert-error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        .status {
            text-align: center;
            font-size: 18px;
            margin: 20px 0;
        }
        .mining-status {
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            color: #856404;
            padding: 15px;
            border-radius: 6px;
            margin: 20px 0;
            text-align: center;
        }
    </style>
    <script>
        function refreshWallets() {
            location.reload();
        }

        function submitForm(formId, endpoint) {
            const form = document.getElementById(formId);
            const formData = new FormData(form);
            const data = {};
            formData.forEach((value, key) => {
                data[key] = value;
            });

            fetch(endpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    alert('✅ ' + data.message);
                    location.reload();
                } else {
                    alert('❌ ' + data.error);
                }
            })
            .catch(error => {
                alert('❌ Error: ' + error);
            });
        }

        function mineTokens() {
            const address = document.getElementById('mineAddress').value;
            if (!address) {
                alert('❌ Please select a wallet address');
                return;
            }

            document.getElementById('miningStatus').style.display = 'block';
            
            fetch('/mine', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({address: address})
            })
            .then(response => response.json())
            .then(data => {
                document.getElementById('miningStatus').style.display = 'none';
                if (data.success) {
                    alert('🎉 ' + data.message);
                    location.reload();
                } else {
                    alert('❌ ' + data.error);
                }
            })
            .catch(error => {
                document.getElementById('miningStatus').style.display = 'none';
                alert('❌ Error: ' + error);
            });
        }
    </script>
</head>
<body>
    <div class="container">
        <h1>🔗 Blockchain Interface</h1>
        
        <div class="grid">
            <!-- Create Wallet Section -->
            <div>
                <h2>💰 Create New Wallet</h2>
                <form id="createWalletForm" onsubmit="event.preventDefault(); submitForm('createWalletForm', '/create-wallet');">
                    <button type="submit">Create New Wallet</button>
                </form>
            </div>

            <!-- Mining Section -->
            <div>
                <h2>⛏️ Mine Tokens</h2>
                <div class="form-group">
                    <label for="mineAddress">Select Wallet to Receive Mining Reward:</label>
                    <select id="mineAddress">
                        <option value="">Choose a wallet...</option>
                        {{range .Wallets}}
                        <option value="{{.Address}}">{{.Address}} ({{.Balance}} tokens)</option>
                        {{end}}
                    </select>
                </div>
                <button onclick="mineTokens()">⛏️ Mine Block</button>
                <div id="miningStatus" class="mining-status" style="display: none;">
                    ⛏️ Mining in progress... This may take a moment.
                </div>
            </div>
        </div>

        <!-- Send Transaction Section -->
        <div>
            <h2>💸 Send Transaction</h2>
            <form id="sendForm" onsubmit="event.preventDefault(); submitForm('sendForm', '/send');">
                <div class="grid">
                    <div class="form-group">
                        <label for="from">From Address:</label>
                        <select name="from" id="from" required>
                            <option value="">Choose sender...</option>
                            {{range .Wallets}}
                            <option value="{{.Address}}">{{.Address}} ({{.Balance}} tokens)</option>
                            {{end}}
                        </select>
                    </div>
                    <div class="form-group">
                        <label for="to">To Address:</label>
                        <select name="to" id="to" required>
                            <option value="">Choose recipient...</option>
                            {{range .Wallets}}
                            <option value="{{.Address}}">{{.Address}}</option>
                            {{end}}
                        </select>
                    </div>
                </div>
                <div class="form-group">
                    <label for="amount">Amount:</label>
                    <input type="number" name="amount" id="amount" min="1" step="1" required>
                </div>
                <button type="submit">💸 Send Transaction</button>
            </form>
        </div>

        <!-- Wallets List -->
        <div>
            <h2>👛 Your Wallets</h2>
            <button onclick="refreshWallets()" style="margin-bottom: 20px;">🔄 Refresh Balances</button>
            {{if .Wallets}}
                <div class="wallet-list">
                    {{range .Wallets}}
                    <div class="wallet-item">
                        <div><strong>Address:</strong></div>
                        <div class="address">{{.Address}}</div>
                        <div style="margin-top: 10px;"><strong>Balance:</strong> <span class="balance">{{.Balance}} tokens</span></div>
                    </div>
                    {{end}}
                </div>
            {{else}}
                <div class="alert alert-error">
                    No wallets found. Create a wallet to get started!
                </div>
            {{end}}
        </div>

        <!-- Instructions -->
        <div>
            <h2>📋 Quick Start Guide</h2>
            <ol>
                <li><strong>Create Wallet:</strong> Click "Create New Wallet" to generate a new wallet address</li>
                <li><strong>Mine Tokens:</strong> Select a wallet and click "Mine Block" to earn tokens (mining reward: 100 tokens)</li>
                <li><strong>Send Tokens:</strong> Choose sender, recipient, and amount to transfer tokens between wallets</li>
                <li><strong>Check Balances:</strong> Use "Refresh Balances" to see updated token amounts</li>
            </ol>
            <p><strong>💡 Tip:</strong> Start by creating a few wallets, then mine some tokens to have funds for transactions!</p>
        </div>
    </div>
</body>
</html>
`

func (ws *WebServer) getWallets() []WalletInfo {
	var wallets []WalletInfo

	// Get wallets using CLI
	cmd := exec.Command(".\\blockchain.exe", "-cli", "listaddresses")
	cmd.Env = append(os.Environ(), fmt.Sprintf("NODE_ID=%s", ws.nodeID))
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error getting wallets: %v", err)
		return wallets
	}

	addresses := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, address := range addresses {
		address = strings.TrimSpace(address)
		if address == "" {
			continue
		}

		// Get balance for this address
		cmd := exec.Command(".\\blockchain.exe", "-cli", "getbalance", "-address", address)
		cmd.Env = append(os.Environ(), fmt.Sprintf("NODE_ID=%s", ws.nodeID))
		balanceOutput, err := cmd.Output()
		if err != nil {
			log.Printf("Error getting balance for %s: %v", address, err)
			continue
		}

		balanceStr := strings.TrimSpace(string(balanceOutput))
		re := regexp.MustCompile(`Balance of .*: (\d+)`)
		matches := re.FindStringSubmatch(balanceStr)

		balance := "0"
		if len(matches) > 1 {
			balance = matches[1]
		}

		wallets = append(wallets, WalletInfo{
			Address: address,
			Balance: balance,
		})
	}

	return wallets
}

func (ws *WebServer) handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("home").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Wallets []WalletInfo
	}{
		Wallets: ws.getWallets(),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ws *WebServer) handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create wallet using CLI
	cmd := exec.Command(".\\blockchain.exe", "-cli", "createwallet")
	cmd.Env = append(os.Environ(), fmt.Sprintf("NODE_ID=%s", ws.nodeID))
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error creating wallet: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Failed to create wallet: %v", err),
		})
		return
	}

	outputStr := string(output)
	re := regexp.MustCompile(`Your new address: ([a-zA-Z0-9]+)`)
	matches := re.FindStringSubmatch(outputStr)

	if len(matches) < 2 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to extract wallet address from output",
		})
		return
	}

	address := matches[1]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Wallet created successfully! Address: %s", address),
		"address": address,
	})
}

func (ws *WebServer) handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	// Validate inputs
	if req.From == "" || req.To == "" || req.Amount == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "All fields are required",
		})
		return
	}

	if req.From == req.To {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Cannot send to the same address",
		})
		return
	}

	// Check sender balance first
	cmd := exec.Command(".\\blockchain.exe", "-cli", "getbalance", "-address", req.From)
	cmd.Env = append(os.Environ(), fmt.Sprintf("NODE_ID=%s", ws.nodeID))
	balanceOutput, err := cmd.Output()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Failed to check balance: %v", err),
		})
		return
	}

	balanceStr := strings.TrimSpace(string(balanceOutput))
	re := regexp.MustCompile(`Balance of .*: (\d+)`)
	matches := re.FindStringSubmatch(balanceStr)

	if len(matches) < 2 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to parse balance",
		})
		return
	}

	currentBalance, err := strconv.Atoi(matches[1])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid balance format",
		})
		return
	}

	requestedAmount, err := strconv.Atoi(req.Amount)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid amount format",
		})
		return
	}

	if requestedAmount <= 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Amount must be greater than 0",
		})
		return
	}

	if currentBalance < requestedAmount {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Insufficient funds. Current balance: %d, Requested: %d", currentBalance, requestedAmount),
		})
		return
	}

	// Send transaction using CLI
	cmd = exec.Command(".\\blockchain.exe", "-cli", "send", "-from", req.From, "-to", req.To, "-amount", req.Amount, "-mine")
	cmd.Env = append(os.Environ(), fmt.Sprintf("NODE_ID=%s", ws.nodeID))
	_, err = cmd.Output()
	if err != nil {
		// Try to get stderr for more detailed error
		cmd = exec.Command(".\\blockchain.exe", "-cli", "send", "-from", req.From, "-to", req.To, "-amount", req.Amount, "-mine")
		cmd.Env = append(os.Environ(), fmt.Sprintf("NODE_ID=%s", ws.nodeID))
		combinedOutput, _ := cmd.CombinedOutput()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Transaction failed: %s", string(combinedOutput)),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Transaction successful! Sent %s tokens from %s to %s", req.Amount, req.From, req.To),
	})
}

func (ws *WebServer) handleMine(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	if req.Address == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Address is required",
		})
		return
	}

	// Check if blockchain exists first
	cmd := exec.Command(".\\blockchain.exe", "-cli", "printchain")
	cmd.Env = append(os.Environ(), fmt.Sprintf("NODE_ID=%s", ws.nodeID))
	_, err := cmd.Output()
	if err != nil {
		// Blockchain doesn't exist, create it
		cmd = exec.Command(".\\blockchain.exe", "-cli", "createblockchain", "-address", req.Address)
		cmd.Env = append(os.Environ(), fmt.Sprintf("NODE_ID=%s", ws.nodeID))
		output, err := cmd.CombinedOutput()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to create blockchain: %s", string(output)),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Blockchain created and mining reward received! You got 20 tokens.",
		})
		return
	}

	// Blockchain exists, we need to mine a new block
	// The best way is to send a transaction from this address to itself with mining
	// This will create both the transaction and the coinbase reward
	cmd = exec.Command(".\\blockchain.exe", "-cli", "send", "-from", req.Address, "-to", req.Address, "-amount", "1", "-mine")
	cmd.Env = append(os.Environ(), fmt.Sprintf("NODE_ID=%s", ws.nodeID))
	output, err := cmd.CombinedOutput()

	if err != nil {
		// If the address doesn't have enough balance, try with a different approach
		// Let's check if the error is about insufficient funds
		outputStr := string(output)
		if strings.Contains(outputStr, "not enough funds") || strings.Contains(outputStr, "insufficient") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Cannot mine: wallet has insufficient funds for a transaction. Please send some tokens to this wallet first or mine using a wallet that has funds.",
			})
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Mining failed: %s", outputStr),
			})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Mining completed! You received mining rewards (20 tokens) plus kept your transaction.",
	})
}

func (ws *WebServer) start() {
	fmt.Printf("🌐 Starting Blockchain Web Interface on http://localhost:8080\n")
	fmt.Printf("📊 Node ID: %s\n", ws.nodeID)
	fmt.Printf("💡 Tip: Create wallets, mine tokens, and send transactions through the web interface!\n\n")

	http.HandleFunc("/", ws.handleHome)
	http.HandleFunc("/create-wallet", ws.handleCreateWallet)
	http.HandleFunc("/send", ws.handleSend)
	http.HandleFunc("/mine", ws.handleMine)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	var useCLI bool
	flag.BoolVar(&useCLI, "cli", false, "Use CLI interface instead of web interface")
	flag.Parse()

	if useCLI {
		args := flag.Args()
		if len(args) > 0 {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = os.Args[0]
			copy(newArgs[1:], args)
			os.Args = newArgs
		}

		defer os.Exit(0)
		cli := cli.CommandLine{}
		cli.Run()
	} else {
		// Use web interface by default
		server := NewWebServer()
		server.start()
	}
}

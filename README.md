# 🧠 Wallet Watcher — Ethereum API in Go

A lightweight REST API built in **Go** that connects to the **Ethereum blockchain** (mainnet or testnets like Sepolia) to:

- 🔍 Fetch wallet balances
- 📬 Send ETH transactions (planned)
- 🧾 Get transaction details (planned)
- 📡 Subscribe to new blocks and smart contract events (planned)

## 🚀 Features

| Feature                      | Status    |
|-----------------------------|-----------||
| `GET /health`               | ✅ Done    |
| `GET /wallet/{address}/balance` | ✅ Done    |
| `GET /tx/{hash}`            | ⬜️ Planned |
| `POST /transfer`            | ⬜️ Planned |
| `GET /contract/{address}`   | ⬜️ Planned |
| `POST /watch`               | ⬜️ Planned |
| Subscribe to new blocks     | ⬜️ Planned |
| Watch for token transfers   | ⬜️ Planned |

## 🛠️ Tech Stack

- Language: [Go](https://golang.org/)
- Ethereum Client: [go-ethereum](https://pkg.go.dev/github.com/ethereum/go-ethereum)
- HTTP Routing: [chi](https://github.com/go-chi/chi)
- Environment: [godotenv](https://github.com/joho/godotenv)
- Ethereum Node Provider: [Infura](https://infura.io/)

## 📦 Setup

### 1. Clone and Install Dependencies

```bash
git clone https://github.com/yourusername/wallet-watcher.git
cd wallet-watcher
go mod tidy
```

### 2. Configure Environment

Create a `.env` file in the project root:

```env
# Ethereum Network Configuration
INFURA_PROJECT_ID=your-project-id-here
ETHEREUM_NETWORK=sepolia  # Options: mainnet, sepolia, etc
```

### 3. Run the Server

```bash
go run main.go
```

The server will start on port 3000.

## 🔍 API Endpoints

### Get Wallet Balance

```bash
curl http://localhost:8080/wallet/0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b/balance
```

Response:
```json
{
  "address": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b",
  "balance": 1.23456789
}
```

### Watch Wallet Transactions

Monitor wallet transactions in real-time from the current block forward and save them to a file. The endpoint monitors new blocks for 5 minutes, capturing any transactions involving the specified address.

**Endpoint:** `POST /wallet/{address}/watch`

**Behavior:**
- Starts monitoring from the current blockchain block
- Watches for new blocks every 5 seconds
- Captures transactions where the address is either sender or receiver
- Monitors for 5 minutes total duration
- Saves all found transactions to a timestamped JSON file

**Example:**
```bash
curl -X POST http://localhost:8080/wallet/0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b/watch
```

**Response:**
```json
{
  "message": "Transactions saved successfully",
  "file_path": "transactions/0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b_20240101_120000.json",
  "address": "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b"
}
```

**Note:** This endpoint will run for 5 minutes before returning results. During active periods on the blockchain, you may capture multiple transactions.

### Health Check

```bash
curl http://localhost:3000/health
```

Response: `OK`

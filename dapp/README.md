# NFT Contract Event Listener

A Go-based backend service for monitoring Web3 smart contract events and storing them in a MySQL database.

## Project Structure

```
dapp/
├── cmd/
│   └── server/           # Main application entry point
├── internal/
│   ├── config/          # Configuration structures
│   ├── models/          # Data models
│   ├── repository/      # Database layer
│   ├── service/         # Business logic
│   ├── handler/         # HTTP handlers
│   └── middleware/      # HTTP middleware
├── pkg/
│   ├── logger/          # Logging utilities
│   └── utils/           # Utility functions
├── scripts/             # Database scripts
├── config.json          # Application configuration
├── go.mod               # Go module definition
└── README.md            # Documentation
```

## Features

- 🎯 Real-time blockchain event listening
- 💾 MySQL database storage with GORM ORM
- 🔄 Automatic block height tracking
- 🌐 RESTful API endpoints
- 🔒 CORS support
- 📝 Request logging middleware
- 🚀 Concurrent-safe event processing

## Prerequisites

- Go 1.25.2 or higher
- MySQL 8.0 or higher
- Ethereum node access (Infura, Alchemy, or local node)

## Installation

### 1. Clone and navigate to the project

```bash
cd dapp
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Configure the application

Edit `config.json` with your settings:

```json
{
  "server": {
    "port": ":8080",
    "env": "development"
  },
  "database": {
    "host": "localhost",
    "port": "3306",
    "user": "root",
    "password": "your_password",
    "dbname": "nft_events",
    "max_idle": 10,
    "max_open": 100
  },
  "web3": {
    "rpc_url": "https://mainnet.infura.io/v3/YOUR_PROJECT_ID",
    "contract_addr": "0xYourContractAddress",
    "start_block": 0
  }
}
```

### 4. Initialize the database

Option 1: Let the application auto-migrate tables (recommended)

Option 2: Run the SQL script manually

```bash
mysql -u root -p < scripts/init.sql
```

## Running the Application

### Start the server

```bash
go run cmd/server/main.go
```

The server will start on port 8080 (or as configured).

## API Endpoints

### 1. Get Service Status

```bash
GET /api/status
```

Response:
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "is_running": true,
    "last_block": 12345678,
    "current_block": 12345680,
    "contract_addr": "0x...",
    "rpc_url": "https://..."
  }
}
```

### 2. Get Events

```bash
GET /api/events?block_number=12345678
GET /api/events?contract_address=0x...
GET /api/events?tx_hash=0x...
```

### 3. Start Event Listener

```bash
POST /api/start
```

### 4. Stop Event Listener

```bash
POST /api/stop
```

## Development

### Build the application

```bash
go build -o server cmd/server/main.go
```

### Run tests

```bash
go test ./...
```

### Code structure

- **Models**: Define data structures in `internal/models/`
- **Repository**: Database operations in `internal/repository/`
- **Service**: Business logic in `internal/service/`
- **Handler**: HTTP request handling in `internal/handler/`
- **Middleware**: HTTP middleware in `internal/middleware/`

## Configuration

### Server Configuration

- `port`: HTTP server port (default: :8080)
- `env`: Environment (development/production)

### Database Configuration

- `host`: MySQL host
- `port`: MySQL port
- `user`: Database username
- `password`: Database password
- `dbname`: Database name
- `max_idle`: Maximum idle connections
- `max_open`: Maximum open connections

### Web3 Configuration

- `rpc_url`: Ethereum RPC endpoint URL
- `contract_addr`: Smart contract address to monitor
- `start_block`: Block number to start listening from

## Database Schema

### contract_events table

- `id`: Primary key
- `event_name`: Name of the event
- `contract_address`: Contract address that emitted the event
- `tx_hash`: Transaction hash (unique)
- `block_number`: Block number where event occurred
- `block_hash`: Block hash
- `log_index`: Log index within the block
- `from_address`: Source address
- `to_address`: Destination address
- `data`: Event data (hex encoded)
- `topics`: Event topics (JSON array)
- `timestamp`: Event timestamp
- `created_at`: Record creation time
- `updated_at`: Record update time

### block_heights table

- `id`: Primary key
- `height`: Last processed block height
- `updated_at`: Update timestamp

## Monitoring and Logging

The application provides comprehensive logging:

- INFO: General operational messages
- WARN: Warning messages
- ERROR: Error messages
- DEBUG: Detailed debugging information

Logs include timestamp, file location, and message.

## Troubleshooting

### Database connection issues

- Verify MySQL is running
- Check credentials in config.json
- Ensure database exists

### Web3 connection issues

- Verify RPC URL is accessible
- Check contract address format
- Ensure you have sufficient API credits (if using Infura/Alchemy)

### Event listener not starting

- Check if already running via `/api/status`
- Review logs for error messages
- Verify contract address has events

## License

MIT

## Support

For issues and questions, please create an issue in the repository.

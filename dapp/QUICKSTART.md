# Quick Start Guide

## Prerequisites

Before you begin, ensure you have:
- Go 1.25.2 or higher installed
- MySQL 8.0 or higher running
- Access to an Ethereum RPC endpoint (Infura, Alchemy, or local node)

## Step-by-Step Setup

### 1. Configure the Application

Copy and edit the configuration file:

```bash
# On Windows
notepad config.json

# On Linux/Mac
vim config.json
```

Update the following fields in `config.json`:

```json
{
  "database": {
    "user": "your_mysql_user",
    "password": "your_mysql_password"
  },
  "web3": {
    "rpc_url": "https://mainnet.infura.io/v3/YOUR_INFURA_ID",
    "contract_addr": "0xYourNFTContractAddress"
  }
}
```

### 2. Create Database

The application will auto-create tables on first run. Alternatively, you can manually initialize:

```bash
mysql -u root -p < scripts/init.sql
```

### 3. Install Dependencies

```bash
make deps
# or
go mod tidy
```

### 4. Run the Application

Using Make:
```bash
make run
```

Or directly:
```bash
go run cmd/server/main.go
```

### 5. Verify It's Working

Check the service status:

```bash
curl http://localhost:8080/api/status
```

Expected response:
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "is_running": false,
    "last_block": 0,
    "contract_addr": "0x..."
  }
}
```

### 6. Start Event Listening

Start the event listener:

```bash
curl -X POST http://localhost:8080/api/start
```

Monitor the logs to see events being processed.

### 7. Query Events

Get events by block number:
```bash
curl "http://localhost:8080/api/events?block_number=12345678"
```

Get events by contract address:
```bash
curl "http://localhost:8080/api/events?contract_address=0x..."
```

Get events by transaction hash:
```bash
curl "http://localhost:8080/api/events?tx_hash=0x..."
```

## Docker Deployment (Optional)

If you prefer using Docker:

```bash
# Update config.json for Docker environment
# Set DB_HOST=mysql in your environment

# Start all services
docker-compose up -d

# Check logs
docker-compose logs -f app

# Stop services
docker-compose down
```

## Common Issues

### Database Connection Failed
- Verify MySQL is running: `mysqladmin status`
- Check credentials in config.json
- Ensure database exists

### Web3 Connection Failed
- Verify RPC URL is accessible
- Check if you need API keys (Infura/Alchemy)
- Test with curl: `curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}' YOUR_RPC_URL`

### Port Already in Use
- Change port in config.json
- Or stop the process using port 8080

## Next Steps

1. **Customize Event Parsing**: Modify `internal/service/event_listener.go` to parse specific event types
2. **Add Authentication**: Implement JWT or API key authentication in `internal/middleware/`
3. **Add Monitoring**: Integrate Prometheus metrics or health check endpoints
4. **Scale Up**: Consider adding Redis for caching frequently accessed events

## Development Tips

- Use `make test` to run tests
- Use `make fmt` to format code
- Use `make lint` to check code quality
- Use `make build` to create a production binary

## Support

For detailed documentation, see [README.md](README.md)

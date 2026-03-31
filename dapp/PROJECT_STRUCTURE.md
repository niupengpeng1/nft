# Project Structure - NFT Contract Event Listener

## Directory Layout

```
dapp/
├── cmd/
│   └── server/              # Application entry point
│       └── main.go          # Main server initialization and startup
│
├── internal/                # Private application code
│   ├── config/
│   │   └── config.go        # Configuration structures and loader
│   │
│   ├── models/
│   │   └── event.go         # Data models (ContractEvent, BlockHeight)
│   │
│   ├── repository/
│   │   ├── database.go      # Database connection and migration
│   │   └── event_repository.go  # CRUD operations for events
│   │
│   ├── service/
│   │   └── event_listener.go    # Web3 event listening business logic
│   │
│   ├── handler/
│   │   └── event_handler.go     # HTTP request handlers
│   │
│   └── middleware/
│       └── middleware.go        # HTTP middleware (logging, CORS)
│
├── pkg/                     # Public library code
│   ├── logger/
│   │   └── logger.go        # Logging utilities
│   │
│   └── utils/
│       └── utils.go         # General utility functions
│
├── scripts/                 # Database and deployment scripts
│   └── init.sql             # Database initialization script
│
├── config.json              # Application configuration
├── go.mod                   # Go module dependencies
├── go.sum                   # Dependency checksums
├── Makefile                 # Build automation
├── Dockerfile               # Docker container definition
├── docker-compose.yml       # Multi-container orchestration
├── .gitignore               # Git ignore rules
├── .env.example             # Environment variables template
├── README.md                # Main documentation
└── QUICKSTART.md            # Quick start guide
```

## Module Dependencies

### Core Dependencies

1. **github.com/ethereum/go-ethereum** (v1.14.0)
   - Ethereum client for blockchain interaction
   - Used for: RPC connection, event filtering, log parsing

2. **gorm.io/gorm** (v1.31.1)
   - Go ORM library
   - Used for: Database operations, model mapping, migrations

3. **gorm.io/driver/mysql** (v1.6.0)
   - MySQL driver for GORM
   - Used for: MySQL database connectivity

### Standard Library Packages

- `context` - Context management for cancellation
- `encoding/json` - JSON encoding/decoding
- `fmt` - Formatted I/O
- `log` - Logging package
- `net/http` - HTTP server and client
- `os` - Operating system functionality
- `os/signal` - Signal handling
- `strconv` - String conversion
- `sync` - Synchronization primitives
- `syscall` - System calls
- `time` - Time functionality

## Architecture Layers

### 1. Entry Point (`cmd/server/main.go`)
- Application bootstrap
- Configuration loading
- Dependency injection
- Server startup

### 2. Configuration Layer (`internal/config/`)
- JSON configuration parsing
- Environment-specific settings
- Global configuration access

### 3. Model Layer (`internal/models/`)
- Data structures
- Database schema definitions
- GORM model interfaces

### 4. Repository Layer (`internal/repository/`)
- Database connection management
- CRUD operations
- Transaction handling
- Data persistence abstraction

### 5. Service Layer (`internal/service/`)
- Business logic
- Web3 integration
- Event processing
- External API communication

### 6. Handler Layer (`internal/handler/`)
- HTTP request handling
- Request/response mapping
- Input validation
- API endpoint implementation

### 7. Middleware Layer (`internal/middleware/`)
- Cross-cutting concerns
- Logging
- CORS handling
- Authentication (future)

### 8. Utilities (`pkg/`)
- Shared utilities
- Logger wrapper
- Helper functions

## Data Flow

```
HTTP Request
    ↓
Middleware (CORS, Logging)
    ↓
Handler (Request Processing)
    ↓
Service (Business Logic)
    ↓
Repository (Data Access)
    ↓
Database / Web3 Client
```

## Event Listening Flow

```
1. Start EventListenerService
    ↓
2. Get last processed block from DB
    ↓
3. Poll for new blocks (every 12 seconds)
    ↓
4. Filter logs from contract address
    ↓
5. Convert logs to ContractEvent models
    ↓
6. Save events to database (batch)
    ↓
7. Update block height tracker
    ↓
8. Repeat from step 3
```

## Database Schema

### contract_events Table

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| event_name | VARCHAR(100) | Event signature name |
| contract_address | VARCHAR(42) | Contract that emitted event |
| tx_hash | VARCHAR(66) | Transaction hash (unique) |
| block_number | BIGINT | Block number |
| block_hash | VARCHAR(66) | Block hash |
| log_index | UINT | Log position in block |
| from_address | VARCHAR(42) | Source address |
| to_address | VARCHAR(42) | Destination address |
| data | TEXT | Event data (hex) |
| topics | TEXT | Event topics (JSON) |
| timestamp | DATETIME | Event timestamp |
| created_at | DATETIME | Record creation time |
| updated_at | DATETIME | Record update time |

### block_heights Table

| Column | Type | Description |
|--------|------|-------------|
| id | BIGINT | Primary key |
| height | BIGINT | Last processed block |
| updated_at | DATETIME | Update timestamp |

## API Endpoints

### GET /api/status
Returns current service status and block heights.

### GET /api/events
Query events with filters:
- `block_number` - Filter by block
- `contract_address` - Filter by contract
- `tx_hash` - Filter by transaction

### POST /api/start
Start the event listener service.

### POST /api/stop
Stop the event listener service.

## Configuration Options

### Server
- `port` - HTTP server port
- `env` - Environment mode

### Database
- `host` - MySQL hostname
- `port` - MySQL port
- `user` - Username
- `password` - Password
- `dbname` - Database name
- `max_idle` - Max idle connections
- `max_open` - Max open connections

### Web3
- `rpc_url` - Ethereum RPC endpoint
- `contract_addr` - Contract to monitor
- `start_block` - Initial block number

## Build Targets

```bash
make deps          # Install dependencies
make build         # Build binary
make run           # Run application
make test          # Run tests
make clean         # Clean build artifacts
make init-db       # Initialize database
make fmt           # Format code
make lint          # Lint code
```

## Docker Deployment

### Services
- `mysql` - MySQL 8.0 database
- `app` - Go application container

### Ports
- `3306` - MySQL
- `8080` - HTTP API

### Volumes
- `mysql_data` - Persistent database storage
- `./config.json` - Application configuration

## Security Considerations

1. **Configuration**: Never commit `config.json` with real credentials
2. **Database**: Use strong passwords and restrict network access
3. **API**: Implement authentication middleware for production
4. **Web3**: Use API keys with rate limits for RPC endpoints
5. **Logs**: Sanitize sensitive data in logs

## Extensibility

### Adding New Event Types
1. Define event ABI in service layer
2. Parse specific topics in `logToContractEvent()`
3. Add event-specific fields to model if needed

### Adding New APIs
1. Add handler method in `event_handler.go`
2. Register route in `main.go`
3. Add middleware if required

### Adding New Databases
1. Import new GORM driver
2. Update `InitDB()` function
3. Modify DSN construction

## Performance Optimization

- Batch event saving (100 records per batch)
- Configurable connection pooling
- Block range processing (100 blocks per batch)
- Concurrent-safe mutex protection
- Efficient indexing on database tables

## Monitoring Points

- Last processed block height
- Current blockchain height
- Events processed per block
- Database query performance
- HTTP request latency
- Web3 RPC response times

## Testing Strategy

- Unit tests for repository layer
- Integration tests for service layer
- API endpoint tests
- Mock Web3 client for testing
- Database transaction tests

## Future Enhancements

1. WebSocket subscription for real-time events
2. Multiple contract monitoring
3. Event filtering by signature
4. Redis caching layer
5. Prometheus metrics
6. GraphQL API
7. Event replay functionality
8. Dead letter queue for failed saves
9. Retry mechanism for RPC failures
10. Horizontal scaling support

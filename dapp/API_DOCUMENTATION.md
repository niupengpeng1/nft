# API Documentation

## Base URL
```
http://localhost:8080/api
```

## Content Type
All requests and responses use `application/json`

## Authentication
Currently not implemented (planned for production)

---

## Endpoints

### 1. Get Service Status

Returns the current status of the event listener service.

**Endpoint:** `GET /api/status`

**Parameters:** None

**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "is_running": true,
    "last_block": 18500123,
    "current_block": 18500145,
    "contract_addr": "0x1234567890abcdef1234567890abcdef12345678",
    "rpc_url": "https://mainnet.infura.io/v3/..."
  }
}
```

**Response Fields:**
- `is_running` - Boolean indicating if listener is active
- `last_block` - Last processed block number
- `current_block` - Current blockchain height
- `contract_addr` - Monitored contract address
- `rpc_url` - RPC endpoint URL

**Example:**
```bash
curl http://localhost:8080/api/status
```

---

### 2. Get Events

Query contract events with optional filters.

**Endpoint:** `GET /api/events`

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| block_number | uint64 | No | Filter by specific block number |
| contract_address | string | No | Filter by contract address |
| tx_hash | string | No | Filter by transaction hash |

**Note:** At least one filter parameter is recommended. Without filters, returns all events (may be large).

**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "event_name": "Transfer",
      "contract_address": "0x1234567890abcdef1234567890abcdef12345678",
      "tx_hash": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
      "block_number": 18500123,
      "block_hash": "0x9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba",
      "log_index": 5,
      "from_address": "0xaaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffff",
      "to_address": "0xffffffffeeeeeedddddddcccccccbbbbbbbbaaaaaaaa",
      "data": "0x0000000000000000000000000000000000000000000000000000000000000001",
      "topics": "[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x...\",\"0x...\"]",
      "timestamp": "2024-01-15T10:30:00Z",
      "created_at": "2024-01-15T10:30:05Z",
      "updated_at": "2024-01-15T10:30:05Z"
    }
  ]
}
```

**Event Object Fields:**
- `id` - Database primary key
- `event_name` - Name/signature of the event
- `contract_address` - Contract that emitted the event
- `tx_hash` - Transaction hash (unique identifier)
- `block_number` - Block where event occurred
- `block_hash` - Hash of the block
- `log_index` - Position of log in block
- `from_address` - Source address (if available)
- `to_address` - Destination address (if available)
- `data` - Event data (hex-encoded)
- `topics` - Event topics (JSON array)
- `timestamp` - Event timestamp
- `created_at` - Record creation time
- `updated_at` - Record update time

**Examples:**

Get events by block number:
```bash
curl "http://localhost:8080/api/events?block_number=18500123"
```

Get events by contract address:
```bash
curl "http://localhost:8080/api/events?contract_address=0x1234567890abcdef1234567890abcdef12345678"
```

Get single event by transaction hash:
```bash
curl "http://localhost:8080/api/events?tx_hash=0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
```

**Error Responses:**

400 Bad Request:
```json
{
  "code": 400,
  "message": "Invalid block number"
}
```

500 Internal Server Error:
```json
{
  "code": 500,
  "message": "Failed to fetch events: <error details>"
}
```

---

### 3. Start Event Listener

Starts the blockchain event listening service.

**Endpoint:** `POST /api/start`

**Parameters:** None

**Request:**
```bash
curl -X POST http://localhost:8080/api/start
```

**Response:**
```json
{
  "code": 200,
  "message": "Event listener starting"
}
```

**Behavior:**
- Loads last processed block from database
- Connects to Ethereum RPC node
- Begins polling for new blocks every 12 seconds
- Filters logs from configured contract address
- Saves events to database in batches
- Updates block height tracker

**Notes:**
- Idempotent: Can be called multiple times safely
- Runs asynchronously in background
- Check `/api/status` to monitor progress

**Error Responses:**

405 Method Not Allowed:
```json
{
  "code": 405,
  "message": "Method not allowed"
}
```

---

### 4. Stop Event Listener

Stops the blockchain event listening service.

**Endpoint:** `POST /api/stop`

**Parameters:** None

**Request:**
```bash
curl -X POST http://localhost:8080/api/stop
```

**Response:**
```json
{
  "code": 200,
  "message": "Event listener stopped"
}
```

**Behavior:**
- Gracefully stops the polling loop
- Releases Web3 client resources
- Preserves last processed block in database

**Notes:**
- Safe to call when already stopped
- Does not delete stored events
- Can be restarted with `/api/start`

**Error Responses:**

405 Method Not Allowed:
```json
{
  "code": 405,
  "message": "Method not allowed"
}
```

---

## Response Format Standard

All API responses follow this structure:

```json
{
  "code": <number>,
  "message": "<string>",
  "data": <any>
}
```

**HTTP Status Codes:**

| Code | Meaning | Usage |
|------|---------|-------|
| 200 | OK | Successful request |
| 400 | Bad Request | Invalid parameters |
| 404 | Not Found | Endpoint doesn't exist |
| 405 | Method Not Allowed | Wrong HTTP method |
| 500 | Internal Server Error | Server error |

**Code Field Values:**

- `200` - Success
- `400` - Client error (bad request)
- `405` - Client error (wrong method)
- `500` - Server error

---

## Rate Limiting

Currently not implemented (planned for production)

Recommendations:
- Maximum 100 requests per minute per IP
- Implement token bucket algorithm
- Return 429 Too Many Requests when exceeded

---

## CORS

Cross-Origin Resource Sharing is enabled by default.

**Headers:**
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Accept, Content-Type, Content-Length, Accept-Encoding, Authorization
```

**Preflight Requests:**
OPTIONS requests return 204 No Content immediately.

---

## Logging

All HTTP requests are logged with:
- HTTP method
- Request path
- Response status code
- Request duration

**Example Log Entry:**
```
INFO	2024/01/15 10:30:00 main.go:45	Starting HTTP server on :8080
INFO	2024/01/15 10:30:05 middleware.go:23	GET /api/status 200 1.5ms
```

---

## Health Check

For load balancer health checks:

**Endpoint:** `GET /api/status`

**Expected Response:** HTTP 200 with JSON body

**Timeout Recommendation:** 5 seconds

**Check Interval:** 30 seconds

---

## WebSocket Support (Future)

Planned enhancement for real-time event streaming:

**Endpoint:** `ws://localhost:8080/ws/events`

**Features:**
- Subscribe to specific contract addresses
- Filter by event types
- Real-time event delivery
- Heartbeat/ping-pong mechanism

---

## GraphQL Support (Future)

Planned enhancement for flexible querying:

**Endpoint:** `POST /api/graphql`

**Example Query:**
```graphql
query {
  events(blockNumber: 18500123) {
    id
    eventName
    txHash
    blockNumber
    data
  }
}
```

---

## Best Practices

### 1. Query Optimization
- Always use filters when querying events
- Avoid fetching all events without pagination
- Use specific queries (tx_hash is most specific)

### 2. Error Handling
- Always check response code field
- Handle network timeouts gracefully
- Implement exponential backoff for retries

### 3. Monitoring
- Poll `/api/status` periodically
- Monitor `last_block` progression
- Alert if `is_running` is false unexpectedly

### 4. Production Deployment
- Enable HTTPS/TLS
- Implement API authentication
- Add rate limiting
- Configure proper CORS origins
- Use connection pooling
- Set up log aggregation

---

## Example Client Implementation

### JavaScript/Node.js

```javascript
const axios = require('axios');

const API_BASE = 'http://localhost:8080/api';

async function getStatus() {
  const response = await axios.get(`${API_BASE}/status`);
  return response.data;
}

async function getEventsByBlock(blockNumber) {
  const response = await axios.get(`${API_BASE}/events`, {
    params: { block_number: blockNumber }
  });
  return response.data;
}

async function startListener() {
  const response = await axios.post(`${API_BASE}/start`);
  return response.data;
}

async function stopListener() {
  const response = await axios.post(`${API_BASE}/stop`);
  return response.data;
}

// Usage example
(async () => {
  try {
    console.log('Starting listener...');
    await startListener();
    
    // Wait and check status
    setTimeout(async () => {
      const status = await getStatus();
      console.log('Status:', status);
      
      // Get recent events
      const events = await getEventsByBlock(status.data.last_block);
      console.log('Events:', events);
    }, 5000);
  } catch (error) {
    console.error('Error:', error.message);
  }
})();
```

### Python

```python
import requests
import time

API_BASE = 'http://localhost:8080/api'

def get_status():
    response = requests.get(f'{API_BASE}/status')
    return response.json()

def get_events_by_block(block_number):
    params = {'block_number': block_number}
    response = requests.get(f'{API_BASE}/events', params=params)
    return response.json()

def start_listener():
    response = requests.post(f'{API_BASE}/start')
    return response.json()

def stop_listener():
    response = requests.post(f'{API_BASE}/stop')
    return response.json()

# Usage example
if __name__ == '__main__':
    print('Starting listener...')
    start_listener()
    
    time.sleep(5)
    
    status = get_status()
    print(f'Status: {status}')
    
    events = get_events_by_block(status['data']['last_block'])
    print(f'Found {len(events["data"])} events')
```

### Go

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func getStatus() (*Response, error) {
    resp, err := http.Get("http://localhost:8080/api/status")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result Response
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}

func main() {
    status, err := getStatus()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Status: %+v\n", status)
}
```

---

## Troubleshooting

### Issue: Connection refused
**Solution:** Verify server is running and port 8080 is not blocked

### Issue: Empty events returned
**Solution:** 
- Check if listener is running (`/api/status`)
- Verify contract address has events
- Check block number range

### Issue: Slow response times
**Solution:**
- Add database indexes
- Optimize MySQL configuration
- Use connection pooling
- Consider caching layer

### Issue: Listener not processing events
**Solution:**
- Check RPC URL accessibility
- Verify contract address format
- Review application logs for errors
- Ensure sufficient API credits (Infura/Alchemy)

---

## Support

For issues or questions:
1. Check application logs
2. Review this documentation
3. See README.md for general information
4. Create issue in repository

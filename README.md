# RediFire

A simple service in Go which pops JSON records from a Redis list and saves them to Firestore.

## Features

- Written in Go 1.24
- Transfers JSON records from Redis lists to Firestore collections
- Supports multiple source-to-target mappings
- Configurable via YAML
- Docker support with scratch-based runtime
- External Redis connection with password authentication
- Automatic insertion timestamps
- Dead Letter Queue (DLQ) for failed message processing

## Configuration

Create a `config.yaml` file based on `config.example.yaml`:

```yaml
# Redis configuration
redis:
  host: "localhost:6379"
  password: ""  # Or use REDIS_PASSWORD env var
  db: 0

# Firestore configuration
firestore:
  projectID: "your-gcp-project-id"
  credentialsFile: ""  # Or use GOOGLE_APPLICATION_CREDENTIALS env var

# Mapping of Redis lists to Firestore collections
mappings:
  - source: "events_queue"
    target: "events"
  - source: "users_queue"
    target: "users"
```

## Environment Variables

- `CONFIG_PATH`: Path to configuration file (default: `config.yaml`)
- `REDIS_PASSWORD`: Redis password (overrides config file)
- `GOOGLE_APPLICATION_CREDENTIALS`: Path to Google Cloud credentials JSON file

## Usage

### Running Locally

```bash
# Build
go build -o redifire .

# Run
export REDIS_PASSWORD="your-redis-password"
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/credentials.json"
./redifire
```

### Running with Docker Compose

1. Create your `config.yaml`
2. Set environment variables in `.env` file or export them
3. Run:

```bash
docker-compose up -d
```

### Building Docker Image

```bash
docker build -t redifire:latest .
```

## Data Format

Messages are popped from Redis lists and stored in Firestore with the following structure:

```json
{
  "payload": {
    // Original JSON data from Redis
  },
  "timestamp": "2026-02-01T20:15:00Z"
}
```

## Dead Letter Queue (DLQ)

RediFire includes built-in support for a Dead Letter Queue to handle failed message processing:

### Behavior

- When a message fails to process (e.g., Firestore write error or invalid JSON), it is automatically pushed to a DLQ
- The DLQ is a Redis list with the same name as the source list, with `-dlq` appended
- Example: Messages from `events_queue` that fail will be pushed to `events_queue-dlq`

### Use Cases

The DLQ captures failed messages in these scenarios:
- **Invalid JSON**: Messages that cannot be parsed as JSON
- **Firestore write errors**: Messages that fail to save to Firestore (network issues, permissions, etc.)

### Monitoring

All DLQ events are logged for monitoring and debugging:
```
[events_queue] Message pushed to DLQ events_queue-dlq
[users_queue] Invalid JSON message pushed to DLQ users_queue-dlq
```

### Recovery

Failed messages in the DLQ can be:
- Inspected for debugging using Redis commands: `LRANGE events_queue-dlq 0 -1`
- Reprocessed by pushing them back to the source queue: `RPOPLPUSH events_queue-dlq events_queue`
- Cleared if no longer needed: `DEL events_queue-dlq`

## Requirements

- Go 1.24 or later
- Access to Redis server
- Google Cloud Firestore project and credentials

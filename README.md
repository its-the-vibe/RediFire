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

## Requirements

- Go 1.24 or later
- Access to Redis server
- Google Cloud Firestore project and credentials

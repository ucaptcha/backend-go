# μCaptcha Backend

A backend implementation of the μCaptcha system written in Go.

## Installation

1. Clone the repository:
```bash
git clone https://github.com/ucaptcha/backend-go.git
cd backend-go
```

2. Install dependencies:
```bash
go mod download
```

3. Copy the example config:
```bash
cp config.example.yaml config.yaml
```

4. Edit `config.yaml` with your settings.

## Configuration

Example configuration (`config.yaml`):
```yaml
mode: "memory"  # or "redis"
redis:
  addr: "localhost:6379"
  password: ""
  db: 0
key_length: 1536
key_rotation_interval: "24h"
port: 8080
host: "0.0.0.0"
```

- `mode`: Storage backend ("memory" or "redis")
- `redis`: Redis connection settings (only used when mode is "redis")
- `key_length`: RSA key length in bits (recommended minimum 1536)
- `key_rotation_interval`: Key rotation interval (e.g. "24h", "1h30m")
- `port`: Server port
- `host`: Server host

## Usage

Run the server:
```bash
go run main.go
```

## API Documentation

### Generate Challenge
`POST /challenge`

Generates a new captcha challenge.

**Example Response:**
```json
{
    "id": "8756d5cc-d3fc-35a7-940f-1e388c9f0df8",
    "g": "77642905874398787632272558597266110899559428963277249926018544312322752",
    "n": "950959592177192295820512556855602325965163906380534204976341268132239",
    "t": 1000000
}
```

In which:

- `id`: Challenge ID
- `g`: The input `g` of the VDF function
- `n`: The public key `N` of the RSA key
- `t`: Challenge difficulty

### Verify Solution
`POST /challenge/:id/validation`

Verifies a captcha solution.

**Route Parameters:**
- `id`: Challenge ID

**Request:**
```json
{
  "y": "the answer calculated by client",
}
```

**Response:**
```json
{
  "success": true
}
```

or:

```json
{
    "error": "Challenge not found"
}
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

MIT - See [LICENSE](LICENSE).
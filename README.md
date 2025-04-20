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
challenge_storage: "redis"
key_storage: "memory"
redis:
    addr: "localhost:6379"
    password: ""
    db: 0
key_length: 1536
key_rotation_interval: "24h"
port: 8080
host: "0.0.0.0"
difficulty: 100000
```

- `challenge_storage`: Storage mode for challenges ("memory" or "redis")
- `key_storage`: Storage mode for keys ("memory" or "redis")
- `redis`: Redis connection settings (only used when mode is "redis")
- `key_length`: RSA key length in bits (recommended minimum 1536)
- `key_rotation_interval`: Key rotation interval (e.g. "24h", "1h30m")
- `port`: Server port
- `host`: Server host
- `difficulty`: The initial difficulty of the challenge

We recommended to use `redis` for challenge, since it will automatically clean up expired challenges,
and `memory` for key, since the performance of choosing random key for a generating a challenge
is bad in the current Redis mode implementation.

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
    "y": "the answer calculated by client"
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

### Update Difficulty

`PUT /difficulty`
Updates the difficulty configuration.

```json
{
    "difficulty": 1000000
}
```

**Response:**

````json
{
    "difficulty": 1000000,
    "success": true
}
```

## Performance

Tested on M2 MacBook Air (16GB),
with config:

```yaml
challenge_storage: "redis"
key_storage: "memory"
redis:
  addr: "localhost:6379"
  password: ""
  db: 0
key_length: 1024
key_rotation_interval: "12m"
key_pool_size: 20
port: 8080
host: "0.0.0.0"
````

### `POST /challenge/`

Command: `ab -n 100000 -c 128 -m POST http://127.0.0.1:8080/challenge`

```text
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /challenge
Document Length:        688 bytes

Concurrency Level:      128
Time taken for tests:   4.971 seconds
Complete requests:      100000
Failed requests:        34445
   (Connect: 0, Receive: 0, Length: 34445, Exceptions: 0)
Total transferred:      81219051 bytes
HTML transferred:       68819051 bytes
Requests per second:    20116.26 [#/sec] (mean)
Time per request:       6.363 [ms] (mean)
Time per request:       0.050 [ms] (mean, across all concurrent requests)
Transfer rate:          15955.30 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.6      0      94
Processing:     1    6   4.2      6     101
Waiting:        1    6   4.2      6     101
Total:          1    6   4.3      6     102

Percentage of the requests served within a certain time (ms)
  50%      6
  66%      6
  75%      7
  80%      7
  90%      7
  95%      8
  98%      9
  99%     13
 100%    102 (longest request)
```

### `POST /challenge/:id/validation`

Command: `ab -n 100000 -c 128 -m POST -p post.txt http://127.0.0.1:8080/challenge/4082e12d-7e60-8a4c-36d8-7bacd1a94500/validation`

> Temporiarily disabled deletion for verified challenges, and mocked an incorrect answer for the challenge. (So the test will always fail)

```text
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /challenge/4082e12d-7e60-8a4c-36d8-7bacd1a94500/validation
Document Length:        17 bytes

Concurrency Level:      128
Time taken for tests:   8.069 seconds
Complete requests:      100000
Failed requests:        0
Non-2xx responses:      100000
Total transferred:      15000000 bytes
Total body sent:        50600000
HTML transferred:       1700000 bytes
Requests per second:    12392.66 [#/sec] (mean)
Time per request:       10.329 [ms] (mean)
Time per request:       0.081 [ms] (mean, across all concurrent requests)
Transfer rate:          1815.33 [Kbytes/sec] received
                        6123.72 kb/s sent
                        7939.05 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.2      0      26
Processing:     1   10   3.8     10      52
Waiting:        1   10   3.8     10      52
Total:          1   10   3.8     10      54

Percentage of the requests served within a certain time (ms)
  50%     10
  66%     11
  75%     12
  80%     12
  90%     14
  95%     15
  98%     18
  99%     21
 100%     54 (longest request)
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

MIT - See [LICENSE](LICENSE).

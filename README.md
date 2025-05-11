# μCaptcha Backend

This is a backend implementation of the μCaptcha system, written in Go.

## Installation

### Using Docker

1. Pull the Docker image:

    ```bash
    docker pull alikia2x/ucaptcha-backend
    ```

2. Prepare the configuration file: Copy `config.example.yaml` to `/path/to/your/config.yaml` and modify it according to your settings.

3. Run the container:

    ```bash
    docker run -d -p 8080:8080 -v /path/to/your/config.yaml:/app/config.yaml alikia2x/ucaptcha-backend
    ```

### Manual Build

1. Clone the repository:

    ```bash
    git clone https://github.com/ucaptcha/backend-go.git
    cd backend-go
    ```

2. Install dependencies:

    ```bash
    go mod download
    ```

3. Copy the example configuration:

    ```bash
    cp config.example.yaml config.yaml
    ```

4. Edit `config.yaml` to suit your needs.

## Configuration

Here is an example configuration (`config.yaml`):

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

### Configuration Options

- `challenge_storage`: Mode for storing challenges ("memory" or "redis").
- `key_storage`: Mode for storing keys ("memory" or "redis").
- `redis`: Redis connection settings (applicable only when using "redis").
- `key_length`: RSA key length in bits (recommended minimum is 1536).
- `key_rotation_interval`: Interval for key rotation (e.g., "24h", "1h30m").
- `port`: Port for the server.
- `host`: Host for the server.
- `difficulty`: Initial difficulty level of the challenge.

We recommend using `redis` for challenge storage, as it automatically cleans up expired challenges, and `memory` for key storage, since the current Redis implementation has performance issues when selecting random keys for challenge generation.

## Usage

To run the server, execute:

```bash
go run main.go
```

## API Documentation

Note: An [OpenAPI specification](api-doc.yaml) is also available.

**IMPORTANT:** This API **should not** be directly exposed to the public. You must integrate it into your own backend code and implement additional features as needed (e.g., dynamic difficulty, rate limiting, etc.).

### 1. Creating a Challenge

`POST` `/challenge`

To obtain a new captcha challenge, send a `POST` request to the `/challenge` endpoint. You can optionally specify the `difficulty` in the JSON body.

**Example Request (with difficulty):**

```json
{
    "difficulty": 100000
}
```

**Successful Response (HTTP 201):**

```json
{
    "success": true,
    "id": "dqfUjQbmpT",
    "g": "6806008247175178...254",
    "n": "1087355592116148...087",
    "t": 100000
}
```

You can pass this response to the client. The client will need the `g`, `n`, and `t` values to solve the challenge. Remember to store the `id` for later validation.

### 2. Verifying the Answer

`POST` `/challenge/{id}/validation`

After the client solves the challenge, send their answer (`y`) in a `POST` request to `/challenge/{id}/validation`, replacing `{id}` with the challenge ID you received earlier.

**Example Request:**

`POST /challenge/dqfUjQbmpT/validation`

```json
{
  "y": "32341712...9832"
}
```

**Successful Response (HTTP 200 - Correct Answer):**

```json
{
  "success": true
}
```

**Incorrect Answer Response (HTTP 401):**

```json
{
  "success": false
}
```

**Other Possible Responses:**

- `400`: Invalid format in your request.
- `404`: The provided `id` does not exist.
- `500`: An error occurred on the server.

### 3. Changing Default Difficulty

`PUT` `/difficulty`

You can change the default difficulty for new challenges by sending a `PUT` request to `/difficulty` with the desired `difficulty` in the body.

**Example Request:**

```json
{
  "difficulty": 200000
}
```

**Successful Response (HTTP 200):**

```json
{
  "success": true,
  "difficulty": 200000
}
```

## Performance

Tested on an M2 MacBook Air (16GB) with the following configuration:

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
```

### Performance Testing Results

#### `POST /challenge/`

Command used for testing:

```bash
ab -n 100000 -c 128 -m POST http://127.0.0.1:8080/challenge
```

**Results:**

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

#### `POST /challenge/{id}/validation`

Command used for testing:

```bash
ab -n 100000 -c 128 -m POST -p post.txt http://127.0.0.1:8080/challenge/4082e12d-7e60-8a4c-36d8-7bacd1a94500/validation
```

> Note: Deletion for verified challenges is temporarily disabled, and an incorrect answer is mocked for the challenge, ensuring that the test will always fail.

**Results:**

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

We welcome contributions! If you have suggestions, bug reports, or would like to submit a pull request, please feel free to do so. Your input helps improve the μCaptcha system.

## License

This project is licensed under the MIT License. For more details, please refer to the [LICENSE](LICENSE) file.

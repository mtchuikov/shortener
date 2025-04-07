# shortener
This service can generate short URLs for you.

TODO:

- [ ] Improve error handling 
- [ ] Add more tests

## Downloads

* Download the latest version of the binary from the releases page. You can also use the following command to download it:
```
TODO
```
* Or pull the Docker image from the repository:
```
docker pull ghcr.io/mtchuikov/shortener:latest
```
* Or download the source code and compile it:
```
git clone https://github.com/mtchuikov/shortener
go build -ldflags="-s -w" -o ./build/shortener ./cmd/main.go
```

## Config

The program's behavior can be configured using command-line flags. Here is an overview of supported flags:

| Flag | Shorthand | Default value | Description |
|------|-----------|---------------|-------------|
| `--server-addr`| `-a` | `127.0.0.1:8080` | Specify the IP address and port for the server to listen on |
| `--base-url` | `-b` | `http://127.0.0.1:8080/` | Define the base URL used to generate shortened links |
| `--dsn` | `-d` | |  Define the database DSN (Data Source Name) used to connect to the database |
| `--verbose` | `-v` | `false` | Enable verbose logging at the debug level. Overrides the log-level flag to 'debug' |
| `--help` | `-h` | | Print information about supported flags |

## Endpoints

TODO
# Forward-Proxy

A feature-rich, terminal-based HTTP/HTTPS forward proxy server with a real-time interactive management console.  
Built in Go, this proxy allows you to monitor, block, and manage web traffic with a powerful TUI (terminal user interface).

---

## Features

- **HTTP & HTTPS Proxying**: Handles both HTTP and HTTPS (via CONNECT tunneling) traffic.
- **Real-Time Traffic Monitoring**: View all incoming requests in a live-updating terminal UI.
- **URL Blocking/Unblocking**: Instantly block or unblock any URL from the console.
- **Request Logging**: All requests are timestamped and logged for review.
- **In-Memory Caching**: HTTP responses are cached for performance (HTTP only).
- **Interactive TUI**: Navigate requests and blocked URLs, and manage them with keyboard shortcuts.
- **Configurable via .env**: Easily set proxy host and port.
- **Graceful Shutdown**: Handles SIGINT/SIGTERM for safe shutdown.

---

## Getting Started

### Prerequisites

- Go 1.22 or later

### Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/patnaikankit/Forward-Proxy.git
   cd Forward-Proxy
   ```

2. **Install dependencies:**
   ```sh
   go mod tidy
   ```

3. **Create a `.env` file:**
   ```
   PROXY_HOST=127.0.0.1
   PROXY_PORT=8080
   ```

---

## Usage

### Running the Proxy

```sh
go run main.go
```

By default, the proxy listens on the host and port specified in your `.env` file.

### Enabling the Management Console

- The interactive console UI is enabled if `DebugMode` is set to `true` in `utils/debug.go`.
- To run the proxy in headless mode (no UI), set `DebugMode = false`.

### Configuring Your Browser or System

- Set your HTTP/HTTPS proxy to `127.0.0.1:8080` (or your configured host/port).

---

## Terminal UI Guide

The TUI provides a split-screen interface:

- **Top Half:** Management Console (live request log)
- **Bottom Left:** Keymap (keyboard shortcuts)
- **Bottom Right:** Blocked URLs

### Key Bindings

| Key                | Action                                      |
|--------------------|---------------------------------------------|
| Q                  | Quit                                        |
| S                  | Switch Requests/Blocked pane                |
| R                  | Refresh requests list                       |
| B                  | Block URL (of selected packet)              |
| U                  | Unblock URL (of selected blocked)           |
| Up/Down Arrow      | Select packet or URL                        |

- Use `S` to toggle focus between the requests list and blocked URLs.
- Use `B` to block the selected request's URL.
- Use `U` to unblock the selected blocked URL.

---

## How It Works

- **Proxy Server:** Listens for TCP connections, parses HTTP/HTTPS requests, and forwards them.
- **Request Parsing:** Extracts method, URL, version, host, and port.
- **Blocking:** If a URL is blocked, the proxy returns a 403 Forbidden response.
- **Caching:** HTTP responses are cached in memory for repeated requests.
- **Logging:** All requests are logged with timestamp, protocol, method, and URL.
- **Console UI:** Uses [termbox-go](https://github.com/nsf/termbox-go) for a fast, cross-platform TUI.

---

## Project Structure

```
Forward-Proxy/
  ├── main.go                # Entry point
  ├── proxy/                 # Proxy server logic
  │   ├── proxy.go           # Initialization and listener
  │   ├── handler.go         # Connection handler
  │   ├── request.go         # Request parsing
  │   ├── http_handler.go    # HTTP request handling
  │   ├── http_tunnel.go     # HTTPS tunneling
  ├── utils/                 # Utilities (logging, blocking, caching)
  │   ├── helper.go
  │   ├── debug.go
  │   ├── http_request_model.go
  ├── console/               # Terminal UI
  │   └── ui.go
  ├── go.mod, go.sum
```

---

## Configuration

- **.env**: Set `PROXY_HOST` and `PROXY_PORT` for the proxy server.
- **Debug Mode**: Toggle `DebugMode` in `utils/debug.go` to enable/disable the console UI.

---

## Dependencies

- [termbox-go](https://github.com/nsf/termbox-go) - Terminal UI library
- [joho/godotenv](https://github.com/joho/godotenv) - .env file loader

---

## Contributing

Pull requests and issues are welcome!  
If you have feature requests or bug reports, please open an issue.

---

## License

MIT License

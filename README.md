# qBittorrent Auto-Delete Manager
![Version](https://img.shields.io/badge/Version-1.0.0-brightgreen.svg)  
![Go](https://img.shields.io/badge/Go-1.24-brightgreen.svg)  
![License](https://img.shields.io/badge/License-MIT-blue.svg)  

qBittorrent Auto-Delete Manager is a lightweight Go service that monitors your qBittorrent downloads and automatically deletes completed torrents (and optionally their files) after a configurable amount of time. Designed to prevent unwanted seeding and save disk space, it can run continuously or in a container.  

## Table of Contents
- [Features](#features)
- [Requirements](#requirements)
- [Configuration](#configuration)
- [Development Setup](#development-setup)
- [Docker](#docker)
- [Continuous Integration](#continuous-integration)
- [Contributing](#contributing)
- [License](#license)

## Features
- Connects to qBittorrent via its Web API.
- Automatically deletes torrents a set number of minutes after completion.
- Optionally removes downloaded files from disk.
- Configurable polling interval.
- Environment variable–based configuration.
- Supports `.env` files via `godotenv`.
- Lightweight Docker image for easy deployment.

## Requirements
- Go 1.24 or later.
- qBittorrent with **Web UI enabled**.
- Docker (optional, for containerized deployment).

## Configuration
Create a `.env` file in the project root with the following variables:

```dotenv
QBITTORRENT_URL="http://localhost:8080"
QBITTORRENT_USERNAME="admin"
QBITTORRENT_PASSWORD="adminadmin"
DELETE_FILES=true
DELETE_AFTER_MINUTES=5
POLL_INTERVAL_SECONDS=60
```

**Environment Variable Details**:  
- **QBITTORRENT_URL** – Base URL of your qBittorrent Web UI.  
- **QBITTORRENT_USERNAME** – qBittorrent Web UI username.  
- **QBITTORRENT_PASSWORD** – qBittorrent Web UI password.  
- **DELETE_FILES** – Whether to delete the downloaded files (true/false).  
  - If set to `true`, files will be deleted after the torrent is removed.
  - If set to `false`, only the torrent will be removed, leaving files intact.
- **DELETE_AFTER_MINUTES** – How many minutes after completion a torrent should be deleted.  
- **POLL_INTERVAL_SECONDS** – How often to check torrents (in seconds).  

## Development Setup
1. Install Go 1.24 or later.
2. Clone the repository:
   ```bash
   git clone https://github.com/dmarts05/qbit-autodelete.git
   cd qbit-autodelete
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Build:
   ```bash
   make build
   ```
5. Run:
   ```bash
   make run
   ```

## Docker
Build and run the Docker image:

```bash
docker build -t qbit-autodelete .
docker run --rm --env-file .env qbit-autodelete
```

## Continuous Integration
CI can be configured with GitHub Actions to:
- Lint code with `golangci-lint`.
- Check formatting with `gofmt`.
- Run tests and build binaries.
- Publish Docker images.

See `.github/workflows` for pipeline details.

## Contributing
Contributions are welcome. Please open an issue or submit a pull request.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

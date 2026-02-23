# SaveAny-Bot Minimal Docker Build

This folder contains the artifacts and instructions for building the most minimal version of SaveAny-Bot as a Docker image.

## What's included:

- `Dockerfile.pico` - The minimal Dockerfile with the smallest possible build containing:
  - no_jsparser (disables JavaScript parser support)
  - no_minio (disables MinIO storage backend)  
  - sqlite_glebarez (uses minimal SQLite driver)
  - no_bubbletea (disables bubbletea TUI components)
- `entrypoint.sh` - The container entrypoint script
- `build-minimal-docker.sh` - Script to build the minimal Docker image

## How to build:

When Docker is available, run:

```bash
chmod +x build-minimal-docker.sh
./build-minimal-docker.sh
```

Or build directly with:

```bash 
docker build -f dist/Dockerfile.pico -t saveany-bot:pico . --no-cache
```

## Why this minimizes issues:

Compared to deployments with Telegram interaction issues:

1. Reduced attack surface by disabling unnecessary features
2. Simplified dependency chain minimizing possible conflicts 
3. Minimal runtime components reducing potential failure points
4. Smaller image reduces load time and resource usage
5. Uses scratch base image for the smallest possible footprint
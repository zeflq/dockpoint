# 📘 Dockpoint

**Dockpoint** is a savepoint-aware Docker build CLI tool that enables efficient and incremental Docker builds using inline savepoints inside your Dockerfile.

It’s designed for CI/CD pipelines, DevOps engineers, and platform teams managing Dockerized monorepos.
[![CI](https://github.com/zeflq/dockpoint/actions/workflows/ci.yml/badge.svg)](https://github.com/zeflq/dockpoint/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/zeflq/dockpoint/branch/main/graph/badge.svg)](https://codecov.io/gh/zeflq/dockpoint)
---

## 🎯 Purpose

- Build Dockerfiles incrementally using savepoints
- Reuse previously built layers using savepoints and registry caching
- Speed up CI/CD pipelines dramatically
- Generate smarter cache tags based on build context

---

## 🛍️ Target Users

- CI/CD Pipelines
- DevOps Engineers
- Platform Engineering Teams
- Backend Developers with complex Dockerfiles

---

## 🔑 Core Concepts

### Savepoints

Savepoints are inline markers in your Dockerfile:

```dockerfile
# savepoint: base
FROM node:20

# savepoint: deps
COPY package*.json ./
RUN npm ci
```

Each savepoint becomes a buildable/taggable image layer, enabling incremental builds.

### Full Target Images (`-t repo/image:tag`)

You must specify a full `repo/image:tag` via the `-t` flag:

Example:

```bash
# Build and tag savepoints
 dockpoint build-savepoint -t docker.io/myuser/myapp:1.1.2
```

Inner savepoints will be tagged like `docker.io/myuser/myapp:base`, `:deps`, etc.

The final savepoint will be tagged `docker.io/myuser/myapp:1.1.2`.

---

## 🧪 Features & Commands

### `build-savepoint`

Build all savepoints or only one:

```bash
# Build and tag all savepoints
dockpoint build-savepoint -t docker.io/myuser/myapp:1.1.2

# Build and push
dockpoint build-savepoint -t docker.io/myuser/myapp:1.1.2 --push

# Build only a specific savepoint
dockpoint build-savepoint deps -t docker.io/myuser/myapp:1.1.2

# Dry-run (simulate builds, no Docker actions)
dockpoint build-savepoint -t docker.io/myuser/myapp:1.1.2 --dry-run
```

---

## 🧪 Quick Install

Run this to install latest:

```bash
curl -sSL https://raw.githubusercontent.com/zeflq/dockpoint/main/install.sh | bash
```
## 📥 Install

Download the latest Dockpoint binary from [Releases](https://github.com/zeflq/dockpoint/releases).

Example for Linux/macOS:

```bash
# Download
curl -L https://github.com/zeflq/dockpoint/releases/latest/download/dockpoint_1.0.0_linux_amd64.tar.gz -o dockpoint.tar.gz

# Extract and move
tar -xzf dockpoint.tar.gz
chmod +x dockpoint
sudo mv dockpoint /usr/local/bin/dockpoint

# Test it
dockpoint --help
```
---

## 🔐 Smarter Caching (Optional)

For full context-based caching, hash your full build context before building:

```bash
HASH=$(find . -type f -not -path './.git/*' -exec sha256sum {} + | sort | sha256sum | awk '{print $1}')

docker build -t my-repo/service/my-app:1.0.0.$HASH .
docker push my-repo/service/my-app:1.0.0.$HASH
```

- ✅ Ensures that any change in any file will regenerate a new image tag.
- ✅ 100% reliable for CI/CD systems.

---

## 🛠️ Project Architecture

This project follows **Clean Architecture** adapted for Go CLI applications.

### Folder Layout

```bash
src/
├── application/         # Use cases
├── domain/              # Domain entities
├── core/                # Shared interfaces and DTOs
├── infrastructure/      # Docker, registry, filesystem adapters
├── exposition/cli/      # Cobra commands
└── bootstrap/           # CLI entrypoint (main.go)
```

### ✅ Design Principles

- Separate Domain, UseCases, Infrastructure, and CLI wiring
- All I/O done in infrastructure
- Business rules only depend on interfaces
- DTOs between CLI and UseCases
- Clean testable architecture

See [CONTRIBUTING.md](CONTRIBUTING.md) for full contribution guide.

---

## 🚀 Roadmap

| Epic | Feature | Status |
|:---|:---|:---|
| Savepoint builds | Build and tag savepoints | ✅ Done |
| Smarter Caching | Full context hash caching |✅ Done |
| Verbosity Control | Add `--verbose` mode | 📋 Planned |
| Structured Logging | Zap or Logrus integration | 📋 Planned |
| Pre-commit: Coverage to 80-90% |  | 📋 Planned |
| Pre-commit: golangci-lint |  | 📋 Planned |


---

## 🔗 Technologies

- Go 1.22+
- Cobra CLI
- Docker Engine API (CLI-level)
- Clean Architecture Design

---

## 👥 Contributing

Want to help improve Dockpoint? 🚀

Check out the [CONTRIBUTING.md](CONTRIBUTING.md) for setup instructions, coding conventions, and architecture.

- PRs welcome!
- Please respect Clean Architecture guidelines.

---

## 📄 License

MIT © 2025 — **Build fast, build smart, with Dockpoint!** 🚀


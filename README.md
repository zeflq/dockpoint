# 📘 dockpoint

dockpoint is a **savepoint-aware Docker build CLI** tool that enables efficient and incremental Docker builds using inline savepoints in your Dockerfile. It's designed for CI/CD pipelines, DevOps engineers, and platform teams who manage Dockerized monorepos.

---

## 🎯 Purpose
- Build Dockerfiles incrementally using savepoints
- Reuse previously built image layers via registry cache
- Speed up CI/CD pipelines and reduce build time

---

## 🧭 Target Users
- CI/CD Pipelines
- DevOps Engineers
- Backend Developers
- Platform Teams with Docker Monorepos

---

## 🔑 Core Concepts

### 🔹 Savepoints
Savepoints are named inline comments in the Dockerfile:
```Dockerfile
# savepoint: base
FROM node:20

# savepoint: deps
COPY package*.json ./
RUN npm ci
```
Each becomes a tag like:
```
ghcr.io/org/service:deps
```

### 🔹 .dockpointrc.json
Project-level configuration for base repo:
```json
{
  "repo": "ghcr.io/org/service"
}
```

---

## 🧪 Features & Commands

### ✅ `build-savepoint`
Build a specific or all savepoints:
```bash
dockpoint build-savepoint [savepoint?] --repo <repo> [--push] [--force]
```
- Skips builds if tag exists (unless `--force`)
- Tags image as `<repo>:<savepoint>`
- Pushes image if `--push`

### ✅ `from-savepoint`
Continue a Docker build from a previous savepoint:
```bash
dockpoint from-savepoint <savepoint> --tag <output-tag> --repo <repo>
```

### ✅ `validate`
Validate Dockerfile and `.dockpointrc.json`:
```bash
dockpoint validate
```

### ✅ `list`
List all savepoints and resolved image tags:
```bash
dockpoint list --repo <repo>
```

---

## 🛠 Project Architecture

This project uses **Clean Architecture** adapted for Go CLI tooling with Cobra:

### 📁 Folder Layout
```
src/
├── application/         # Use cases
├── domain/              # Core business logic
├── core/                # Shared types/interfaces/config
├── infrastructure/      # External services: Docker, FS, registry
├── exposition/cli/      # CLI interface (Cobra commands)
└── bootstrap/           # Entry point and CLI wiring
```

### ✅ Contribution Rules
- Use interfaces for all infrastructure (e.g., Docker, registry)
- All I/O in `infrastructure/`, never in `application/`
- Use DTOs between CLI and use cases

See [CONTRIBUTING.md](./CONTRIBUTING.md) for full architectural instructions.

---

## 🔐 Registry Caching Logic
| Condition             | Behavior        |
|----------------------|-----------------|
| Tag exists           | Skip build      |
| Tag doesn't exist    | Build & tag     |
| `--force` is set     | Always rebuild  |

---

## 🚀 Roadmap
| Epic              | Feature          | CLI Command       |
|------------------|------------------|-------------------|
| Savepoint builds | Build to step    | `build-savepoint` |
| Incremental build| Resume from step | `from-savepoint`  |
| Validation        | Validate config  | `validate`        |
| Visibility        | List savepoints  | `list`            |

---

## 🧱 Technologies
- Go 1.22+
- Cobra CLI
- Docker Engine

---

## 👥 Contributing
Want to help improve `dockpoint`? Check out the [CONTRIBUTING.md](./CONTRIBUTING.md) for setup, architecture, and best practices.

---

## 📄 License
MIT © 2025 — Happy building!


# Contributing to dockpoint

Welcome to the dockpoint CLI project! This document outlines the standards and folder structure we follow based on Clean Architecture principles, adapted to a Go CLI tool using Cobra.

---

## 📁 Folder Structure (Clean Architecture)

```
src/
├── application/         # Use cases (orchestrate domain logic via interfaces)
├── domain/              # Core domain entities and business rules
├── core/                # Shared interfaces, config, error types
├── infrastructure/      # External services: Docker, filesystem, registry, etc.
├── exposition/cli/      # CLI interface layer using Cobra
└── bootstrap/           # Main entry point and app wiring

src/
├── application/
│   ├── build_savepoint/
│   │   ├── usecase.go
│   │   ├── dto.go        <-- BuildSavepointRequest & Result
│   │   └── usecase_test.go

```

---

## 🔄 Responsibilities by Layer

| Layer              | Description                                                  | Can Import From                |
|--------------------|--------------------------------------------------------------|--------------------------------|
| `domain/`          | Domain models, value objects, business logic                 | Nothing                        |
| `application/`     | Use case logic, coordinates domain and interfaces            | `domain/`, `core/`             |
| `core/`            | Interfaces and shared utilities (e.g., config, error types)  | Nothing                        |
| `infrastructure/`  | Implements `core` interfaces (e.g., Docker, file parsing)    | `core/`                        |
| `exposition/cli/`  | CLI layer: cobra commands and CLI input handling            | `application/`, `core/`        |
| `bootstrap/`       | Entry point, CLI init, dependency injection                 | All except `infrastructure/`   |

---

## 🧩 Guidelines for Contributions

### ✅ General Rules
- Keep each layer pure and focused.
- Never let `application` or `domain` import from `infrastructure`.
- All I/O must go through interfaces in `core/` and be implemented in `infrastructure/`.
- CLI commands should only convert CLI inputs into DTOs and call `application/` logic.

### ✅ When Adding a Command (e.g., `build-savepoint`)
- **CLI binding**: `exposition/cli/build_savepoint.go`
- **Use case**: `application/build_savepoint.go`
- **Entities**: If needed, `domain/savepoint.go`
- **Interfaces**: Define in `core/`, like `DockerBuilder`
- **Implementation**: In `infrastructure/`, like `docker/docker_builder.go`

### ✅ Testing
- Use dependency injection to allow mocking in tests.
- Unit test `application/` logic in isolation.
- No `os.Exit`, `fmt.Println`, or I/O in `application/` or `domain/`.

---

## 🧪 Sample Prompt for AI Tools

```
You are working inside a clean architecture Go project with the following folders: application/, domain/, core/, infrastructure/, exposition/cli/, bootstrap/.

Implement the command `build-savepoint`:
- Put Cobra CLI command in `exposition/cli/build_savepoint.go`
- Use case orchestration in `application/build_savepoint.go`
- Define entities in `domain/` if needed
- Use interfaces defined in `core/`
- Implement Docker logic in `infrastructure/docker/`
```
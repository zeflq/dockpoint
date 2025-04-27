# 📑 Contributing to Dockpoint

Thank you for your interest in improving **Dockpoint**! 🚀  
We welcome contributions to make Dockpoint the best savepoint-aware Docker builder.

---

## 📚 Code of Conduct

Be respectful, professional, and positive in discussions.  
We value clear communication and clean collaboration.

---

## 🏗 Project Architecture

Dockpoint follows **Clean Architecture** with clear separation of responsibilities:

| Layer              | Responsibility                                 |
|--------------------|-------------------------------------------------|
| `domain/`          | Entities and value objects                     |
| `application/`     | Use cases (business logic)                      |
| `core/`            | Shared interfaces, DTOs, types                 |
| `infrastructure/`  | Implementations: Docker, filesystem, registry   |
| `exposition/cli/`  | Cobra CLI commands                              |
| `bootstrap/`       | Entry point (`main.go`)                         |

---

## 🧰 How to Set Up Locally

1. Install [Go 1.22+](https://golang.org/dl/)
2. Install [Docker](https://docs.docker.com/get-docker/)
3. Clone the repository:

    ```bash
    git clone https://github.com/your-org/dockpoint.git
    cd dockpoint
    ```

4. Build and test:

    ```bash
    go build ./...
    go test ./...
    ```

---

## 🛠 Contribution Guidelines

### ✅ General Rules

- Small, focused PRs — solve one problem at a time.
- Use interfaces inside `application/`.
- Never tie business logic directly to Docker SDK, filesystem, or CLI input.
- No I/O inside `application/` — move it to `infrastructure/`.
- Use DTOs only between CLI and UseCases (`core/dto`).
- No third-party libraries unless necessary (e.g., logging).
- Write clean, readable code following Go best practices.

### 📋 Before You Open a Pull Request

- Update or add unit tests if you touch `application/` or `infrastructure/` logic.
- Ensure all tests pass:

    ```bash
    go test ./...
    ```

- Format your code:

    ```bash
    go fmt ./...
    ```

- Update `README.md`, `CHANGELOG.md`, or `CONTRIBUTING.md` if necessary.
- For major changes, open an Issue first to discuss design.

---

## 🧪 Tests

Write unit tests for:

- Parsers
- Slicers
- Builders
- Registry checkers

We aim for high confidence in core behavior (especially `build-savepoint`).

Use Go’s built-in `testing` package:

```go
import "testing"
```
# 📦 Feature Branch Naming

Use clear branch names:

```bash
feature/build-multi-savepoint
fix/registry-tag-exists
refactor/cleanup-cli-parsing
doc/update-readme
```

# 📦 How to Submit a Contribution
- Fork the repo.
- Create your feature branch.
- Commit your changes.
- Push to your branch.
- Open a Pull Request (PR) into main.
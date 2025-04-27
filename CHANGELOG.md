# 📄 Changelog

All notable changes to this project will be documented in this file.

---

## [1.0.0] - 2025-05-01

### 🚀 Added
- Core command: `build-savepoint`
  - Build all savepoints or a specific one
  - Automatic tagging for each savepoint
  - Final savepoint uses user-provided tag via `-t repo/image:tag`
- Full dry-run mode (`--dry-run`) to preview generated Dockerfile(s)
- Incremental savepoint support:
  - Each savepoint builds as a separate layer
  - Subsequent builds can reuse cached layers
- Registry check before building:
  - Skip building if the image already exists (unless `--force`)

### 🛠️ Changed
- **Require** `-t repo/image:tag` for all builds (no fallback to config anymore)
- **Enforce** full repo/image naming (error if missing `/` or `:tag`)
- **Auto-clean** temporary Dockerfiles after each build (no `--cleanup` flag anymore)

### 🔥 Removed
- Removed `dockpoint init`
- Removed `dockpoint config show`
- Removed `.dockpointrc.json` dependency
- Removed `list`, `validate`, `from-savepoint` commands

### 📦 Internal
- Refactored `Execute()` to strictly parse and validate the `-t` parameter
- Restructured `handleBuildAll()` and `buildOneWithTag()` for better dry-run output
- Full Clean Architecture with separate Domain, Application, Infrastructure, CLI exposition layers
- First CI/CD ready version following Docker CLI practices

---

## 🚀 Next Planned (future versions)

| Planned Feature          | Status |
|---------------------------|--------|
| `--verbose` logging       | 🔜 Planned |
| Structured logging (zap)  | 🔜 Planned |
| Smarter remote cache hints| 🔜 Planned |

---

## 📜 Notes
- This is the **first stable production release** of Dockpoint!
- Designed for modern CI/CD pipelines and mono-repo Docker workflows.
- 100% Go 1.22+ standard.

---

# CLAUDE.md - Project Guidelines for AI-Assisted Development

## Project Overview

GizmoSQL Grafana Data Source Plugin - connects Grafana to GizmoSQL via Apache Arrow Flight SQL.

- **Plugin ID**: `gizmodata-gizmosql-datasource`
- **Backend**: Go (`pkg/` directory)
- **Frontend**: TypeScript/React (`src/` directory)
- **Build system**: Official Grafana plugin tooling (@grafana/create-plugin scaffold)

## Build & Test Commands

```bash
# Frontend
npm install          # Install dependencies
npm run build        # Production build
npm run dev          # Watch mode
npm run typecheck    # Type checking
npm run lint         # Linting
npm test             # Frontend unit tests

# Backend
go build ./pkg/...   # Verify compilation
go test ./...        # Backend unit tests
mage build           # Build current platform via Mage
mage buildall        # Build all platforms

# Full dev environment
docker compose up -d    # Start GizmoSQL + Grafana
docker compose down     # Stop
```

## Best Practices

### Keep a Changelog
- All notable changes MUST be documented in `CHANGELOG.md`
- Follow [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format
- Use [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
- Update CHANGELOG.md with every release or significant change

### Testing
- Write unit tests for all new backend functionality (`go test ./...`)
- Write unit tests for frontend components and logic (`npm test`)
- Integration tests should cover Arrow Flight SQL connectivity
- Verify builds compile cleanly before committing (`go build ./pkg/...` and `npm run build`)
- Run type checking (`npm run typecheck`) and linting (`npm run lint`) before committing

### Grafana Plugin Compliance
- Keep `grafana-plugin-sdk-go` up to date (must not be older than 5 months)
- Keep `@grafana/data` and `@grafana/ui` dependencies current
- Ensure `plugin.json` version matches `package.json` version
- CI workflow Go/Node versions must match `go.mod` and `.nvmrc`
- All referenced assets (logos, screenshots) must exist in `img/`
- LICENSE, README.md, and CHANGELOG.md must be present

### Code Style
- Go: follow standard Go conventions (`go fmt`, `go vet`)
- TypeScript: follow project ESLint config (`.config/.eslintrc`)
- Keep frontend and backend changes in sync when modifying data models

## Key Files

- `src/plugin.json` - Plugin metadata and Grafana dependencies
- `pkg/plugin/datasource.go` - Backend datasource implementation
- `src/datasource.ts` - Frontend datasource implementation
- `src/components/QueryEditor.tsx` - Query editor UI
- `src/components/ConfigEditor.tsx` - Configuration editor UI
- `.github/workflows/` - CI/CD pipelines (ci.yml, build.yml, release.yml)

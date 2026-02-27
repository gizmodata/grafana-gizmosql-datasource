# GizmoSQL Data Source Plugin for Grafana

A Grafana data source plugin that connects to [GizmoSQL](https://github.com/gizmodata/gizmosql) via Apache Arrow Flight SQL.

## Features

- Native Arrow Flight SQL connectivity for high-performance data transfer
- TLS/SSL support with optional certificate verification skip
- Username/password and token-based authentication
- SQL query editor with syntax highlighting
- Time series and table format support
- Grafana template variable support
- Time range macros (`$__timeFrom`, `$__timeTo`, `$__timeFilter`)

### Supported Data Types

| Arrow Type | Grafana Type |
|------------|--------------|
| INT8/16/32/64 | number |
| UINT8/16/32/64 | number |
| FLOAT32/64 | number |
| DECIMAL128/256 | number (float64) |
| STRING | string |
| BOOL | boolean |
| DATE32/64 | time |
| TIMESTAMP | time |

## Requirements

- Grafana >= 10.0.0
- GizmoSQL server running with Flight SQL enabled

## Quick Start (Development)

The easiest way to test the plugin is using the provided Docker Compose setup, which includes:
- GizmoSQL server with TPC-H sample data
- Grafana with the plugin pre-installed
- Pre-configured datasource and sample dashboards

```bash
# Clone the repository
git clone https://github.com/gizmodata/grafana-gizmosql-source.git
cd grafana-gizmosql-source

# Install dependencies and build
npm install
npm run build

# Build the backend (requires Go 1.25+)
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o dist/gpx_gizmosql_datasource_linux_arm64 ./pkg
# Or for Linux amd64:
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/gpx_gizmosql_datasource_linux_amd64 ./pkg

# Start the development environment
docker compose up -d

# Open Grafana
open http://localhost:3000
```

No login required - anonymous access is enabled for testing. You'll find:
- **GizmoSQL** datasource pre-configured
- **GizmoSQL TPC-H Overview** dashboard with sample visualizations
- **GizmoSQL Data Types Demo** dashboard showing all supported types

## Installation (Production)

### From Release

1. Download the latest release from the [releases page](https://github.com/gizmodata/grafana-gizmosql-source/releases)
2. Extract to your Grafana plugins directory (usually `/var/lib/grafana/plugins/`)
3. Add to `grafana.ini`:
   ```ini
   [plugins]
   allow_loading_unsigned_plugins = gizmodata-gizmosql-datasource
   ```
4. Restart Grafana

### Manual Build

```bash
# Install dependencies
npm install

# Build frontend
npm run build

# Build backend for your platform
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/gpx_gizmosql_datasource_linux_amd64 ./pkg

# Or use Mage for all platforms
mage buildall
```

## Configuration

1. In Grafana, go to **Configuration > Data Sources**
2. Click **Add data source**
3. Search for "GizmoSQL"
4. Configure the connection:

| Setting | Description | Default |
|---------|-------------|---------|
| Host | GizmoSQL server hostname | `localhost` |
| Port | Flight SQL port | `31337` |
| Use TLS | Enable TLS encryption | `false` |
| Skip TLS Verify | Skip certificate verification | `false` |
| Username | Authentication username | - |
| Password | Authentication password | - |
| Token | Bearer token (alternative to password) | - |

## Usage

### Basic Query

```sql
SELECT * FROM my_table LIMIT 100
```

### Time Series Query

For time series visualizations, ensure your query returns:
- A time column named `time`, `timestamp`, `ts`, or `datetime`
- One or more value columns

```sql
SELECT
  order_date AS time,
  SUM(total_price) AS revenue
FROM orders
WHERE $__timeFilter
GROUP BY order_date
ORDER BY time
```

### Macros

| Macro | Description | Example Output |
|-------|-------------|----------------|
| `$__timeFrom` | Start of time range | `'2024-01-01T00:00:00Z'` |
| `$__timeTo` | End of time range | `'2024-01-02T00:00:00Z'` |
| `$__timeFilter` | Time range filter | `time >= '...' AND time <= '...'` |

## Sample Dashboards

The plugin includes two sample dashboards that demonstrate its capabilities:

### TPC-H Overview Dashboard
- Total orders, customers, revenue, and parts statistics
- Customers and revenue by region (pie charts, bar charts)
- Orders by priority distribution
- Top 10 customers and suppliers tables
- Orders over time visualization

### Data Types Demo Dashboard
- All supported data types in a single query
- Line items table with multiple column types
- Daily revenue time series chart
- Region and nation reference tables

## Project Structure

```
├── dashboards/              # Sample Grafana dashboards
│   ├── tpch-overview.json
│   └── data-types-demo.json
├── dist/                    # Build output
├── pkg/                     # Go backend source
│   ├── main.go
│   └── plugin/
│       └── datasource.go
├── provisioning/            # Grafana provisioning configs
│   ├── dashboards/
│   └── datasources/
├── src/                     # TypeScript frontend source
│   ├── components/
│   ├── datasource.ts
│   ├── module.ts
│   ├── plugin.json
│   └── types.ts
├── docker-compose.yml       # Development environment
├── go.mod
├── package.json
└── webpack.config.js
```

## Development

### Prerequisites

- Node.js >= 22
- Go >= 1.25
- Docker & Docker Compose (for testing)

### Build Commands

```bash
# Frontend
npm install          # Install dependencies
npm run build        # Production build
npm run dev          # Watch mode
npm run typecheck    # Type checking
npm run lint         # Linting

# Backend
go build -o dist/gpx_gizmosql_datasource ./pkg  # Current platform
mage build           # Current platform with Mage
mage buildall        # All platforms

# Testing
npm test             # Frontend tests
go test ./...        # Backend tests

# Development environment
docker compose up -d    # Start GizmoSQL + Grafana
docker compose down     # Stop environment
docker compose logs -f  # View logs
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## License

Apache License 2.0

## Links

- [GizmoSQL](https://github.com/gizmodata/gizmosql)
- [Apache Arrow Flight SQL](https://arrow.apache.org/docs/format/FlightSql.html)
- [Grafana Plugin Development](https://grafana.com/developers/plugin-tools/)

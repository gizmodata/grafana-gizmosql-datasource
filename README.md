# GizmoSQL Data Source Plugin for Grafana

A Grafana data source plugin that connects to [GizmoSQL](https://github.com/gizmodata/gizmosql) via Apache Arrow Flight SQL.

## Features

- Native Arrow Flight SQL connectivity
- TLS/SSL support
- Username/password and token-based authentication
- SQL query editor with syntax highlighting
- Time series and table format support
- Grafana template variable support
- Time range macros (`$__timeFrom`, `$__timeTo`, `$__timeFilter`)

## Requirements

- Grafana >= 10.0.0
- GizmoSQL server running with Flight SQL enabled

## Installation

### From Release

1. Download the latest release from the [releases page](https://github.com/gizmodata/grafana-gizmosql-source/releases)
2. Extract to your Grafana plugins directory (usually `/var/lib/grafana/plugins/`)
3. Restart Grafana
4. Enable unsigned plugins in `grafana.ini`:
   ```ini
   [plugins]
   allow_loading_unsigned_plugins = gizmodata-gizmosql-datasource
   ```

### Development Setup

```bash
# Install dependencies
npm install

# Build frontend
npm run build

# Build backend (requires Go 1.22+)
go build -o dist/gpx_gizmosql_datasource ./pkg

# Or use Mage
mage build

# Start Grafana with Docker
docker-compose up
```

## Configuration

1. In Grafana, go to **Configuration > Data Sources**
2. Click **Add data source**
3. Search for "GizmoSQL"
4. Configure the connection:
   - **Host**: GizmoSQL server address (e.g., `localhost`)
   - **Port**: Flight SQL port (default: `31337`)
   - **Use TLS**: Enable for encrypted connections
   - **Username/Password**: Optional authentication
   - **Token**: Bearer token authentication (alternative to password)

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
  timestamp as time,
  value,
  metric_name
FROM metrics
WHERE $__timeFilter
ORDER BY time
```

### Macros

| Macro | Description | Example Output |
|-------|-------------|----------------|
| `$__timeFrom` | Start of time range | `'2024-01-01T00:00:00Z'` |
| `$__timeTo` | End of time range | `'2024-01-02T00:00:00Z'` |
| `$__timeFilter` | Time range filter | `time >= '...' AND time <= '...'` |

## Building

### Prerequisites

- Node.js >= 20
- Go >= 1.22
- Mage (optional, for build automation)

### Build Commands

```bash
# Frontend only
npm run build

# Backend only
go build -o dist/gpx_gizmosql_datasource ./pkg

# Backend for all platforms
mage buildall

# Run tests
npm test
mage test
```

## License

Apache License 2.0

## Links

- [GizmoSQL](https://github.com/gizmodata/gizmosql)
- [Apache Arrow Flight SQL](https://arrow.apache.org/docs/format/FlightSql.html)
- [Grafana Plugin Development](https://grafana.com/developers/plugin-tools/)

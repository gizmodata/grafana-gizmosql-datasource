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
- Alerting and annotation support

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

The plugin includes sample dashboards demonstrating its capabilities with the TPC-H dataset, including revenue analytics, customer breakdowns, and data type examples.

## License

Apache License 2.0

## Links

- [GizmoSQL](https://github.com/gizmodata/gizmosql)
- [Apache Arrow Flight SQL](https://arrow.apache.org/docs/format/FlightSql.html)

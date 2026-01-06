package plugin

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/flight"
	"github.com/apache/arrow-go/v18/arrow/flight/flightsql"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Make sure Datasource implements required interfaces
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

// DatasourceSettings holds the parsed data source settings
type DatasourceSettings struct {
	Host          string `json:"host"`
	Port          int    `json:"port"`
	Username      string `json:"username"`
	UseTLS        bool   `json:"useTLS"`
	SkipTLSVerify bool   `json:"skipTLSVerify"`
}

// QueryModel holds the parsed query from the frontend
type QueryModel struct {
	RawSQL string `json:"rawSql"`
	Format string `json:"format"`
}

// Datasource is the GizmoSQL data source implementation
type Datasource struct {
	settings DatasourceSettings
	password string
	token    string
}

// NewDatasource creates a new datasource instance
func NewDatasource(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var dsSettings DatasourceSettings
	if err := json.Unmarshal(settings.JSONData, &dsSettings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	// Set default port
	if dsSettings.Port == 0 {
		dsSettings.Port = 31337
	}

	// Get secure settings
	password := settings.DecryptedSecureJSONData["password"]
	token := settings.DecryptedSecureJSONData["token"]

	return &Datasource{
		settings: dsSettings,
		password: password,
		token:    token,
	}, nil
}

// Dispose cleans up the datasource instance
func (d *Datasource) Dispose() {
	// Clean up any resources
}

// QueryData handles multiple queries and returns multiple responses
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)
		response.Responses[q.RefID] = res
	}

	return response, nil
}

// query executes a single query
func (d *Datasource) query(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var qm QueryModel
	if err := json.Unmarshal(query.JSON, &qm); err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to unmarshal query: %v", err))
	}

	if qm.RawSQL == "" {
		return backend.ErrDataResponse(backend.StatusBadRequest, "query is empty")
	}

	// Replace time macros
	sql := d.replaceMacros(qm.RawSQL, query.TimeRange)

	// Execute the query via Flight SQL
	frame, err := d.executeFlightSQL(ctx, sql)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("query execution failed: %v", err))
	}

	// Set frame name
	frame.Name = query.RefID

	// Convert to time series if requested
	if qm.Format == "time_series" {
		frame = d.convertToTimeSeries(frame)
	}

	return backend.DataResponse{Frames: data.Frames{frame}}
}

// replaceMacros replaces Grafana macros in the SQL query
func (d *Datasource) replaceMacros(sql string, timeRange backend.TimeRange) string {
	sql = strings.ReplaceAll(sql, "$__timeFrom", fmt.Sprintf("'%s'", timeRange.From.UTC().Format(time.RFC3339)))
	sql = strings.ReplaceAll(sql, "$__timeTo", fmt.Sprintf("'%s'", timeRange.To.UTC().Format(time.RFC3339)))
	sql = strings.ReplaceAll(sql, "$__timeFilter", fmt.Sprintf("time >= '%s' AND time <= '%s'",
		timeRange.From.UTC().Format(time.RFC3339),
		timeRange.To.UTC().Format(time.RFC3339)))
	return sql
}

// executeFlightSQL connects to GizmoSQL and executes the query
func (d *Datasource) executeFlightSQL(ctx context.Context, sql string) (*data.Frame, error) {
	addr := fmt.Sprintf("%s:%d", d.settings.Host, d.settings.Port)

	// Set up transport credentials
	var creds credentials.TransportCredentials
	if d.settings.UseTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: d.settings.SkipTLSVerify,
		}
		creds = credentials.NewTLS(tlsConfig)
	} else {
		creds = insecure.NewCredentials()
	}

	// Create Flight SQL client
	client, err := flightsql.NewClient(addr, nil, nil, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to create Flight SQL client: %w", err)
	}
	defer client.Close()

	// Authenticate using Basic Auth if credentials provided
	if d.password != "" {
		username := d.settings.Username
		if username == "" {
			username = "gizmosql" // default username
		}
		authCtx, err := client.Client.AuthenticateBasicToken(ctx, username, d.password)
		if err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
		ctx = authCtx
	} else if d.token != "" {
		// For bearer token auth, add to metadata
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+d.token)
	}

	// Execute query
	info, err := client.Execute(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Get result data
	if len(info.Endpoint) == 0 {
		// No data returned, create empty frame
		return data.NewFrame("result"), nil
	}

	// Read from the first endpoint
	reader, err := client.DoGet(ctx, info.Endpoint[0].Ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to get query results: %w", err)
	}
	defer reader.Release()

	// Convert Arrow records to Grafana data frame
	return d.arrowReaderToFrame(reader)
}

// arrowReaderToFrame converts an Arrow record reader to a Grafana data frame
func (d *Datasource) arrowReaderToFrame(reader *flight.Reader) (*data.Frame, error) {
	alloc := memory.NewGoAllocator()
	frame := data.NewFrame("result")

	schema := reader.Schema()
	fields := schema.Fields()

	// Initialize columns based on schema
	columns := make([]interface{}, len(fields))
	for i, field := range fields {
		columns[i] = d.createColumn(field.Type, alloc)
	}

	// Read all records
	for reader.Next() {
		record := reader.Record()
		for colIdx := 0; colIdx < int(record.NumCols()); colIdx++ {
			col := record.Column(colIdx)
			d.appendColumnData(columns[colIdx], col)
		}
	}

	if err := reader.Err(); err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	// Convert columns to Grafana fields
	for i, field := range fields {
		grafanaField := d.columnToField(field.Name, columns[i])
		if grafanaField != nil {
			frame.Fields = append(frame.Fields, grafanaField)
		}
	}

	return frame, nil
}

// createColumn creates a slice for storing column data based on Arrow type
func (d *Datasource) createColumn(dt arrow.DataType, alloc memory.Allocator) interface{} {
	switch dt.ID() {
	case arrow.INT8:
		return &[]int8{}
	case arrow.INT16:
		return &[]int16{}
	case arrow.INT32:
		return &[]int32{}
	case arrow.INT64:
		return &[]int64{}
	case arrow.UINT8:
		return &[]uint8{}
	case arrow.UINT16:
		return &[]uint16{}
	case arrow.UINT32:
		return &[]uint32{}
	case arrow.UINT64:
		return &[]uint64{}
	case arrow.FLOAT32:
		return &[]float32{}
	case arrow.FLOAT64:
		return &[]float64{}
	case arrow.STRING, arrow.LARGE_STRING:
		return &[]string{}
	case arrow.BOOL:
		return &[]bool{}
	case arrow.TIMESTAMP:
		return &[]time.Time{}
	case arrow.DATE32, arrow.DATE64:
		return &[]time.Time{}
	case arrow.DECIMAL128, arrow.DECIMAL256:
		return &[]float64{}
	default:
		// Default to string for unknown types
		return &[]string{}
	}
}

// appendColumnData appends data from an Arrow array to a column slice
func (d *Datasource) appendColumnData(col interface{}, arr arrow.Array) {
	for i := 0; i < arr.Len(); i++ {
		if arr.IsNull(i) {
			d.appendNull(col)
			continue
		}

		switch a := arr.(type) {
		case *array.Int8:
			if c, ok := col.(*[]int8); ok {
				*c = append(*c, a.Value(i))
			}
		case *array.Int16:
			if c, ok := col.(*[]int16); ok {
				*c = append(*c, a.Value(i))
			}
		case *array.Int32:
			if c, ok := col.(*[]int32); ok {
				*c = append(*c, a.Value(i))
			}
		case *array.Int64:
			if c, ok := col.(*[]int64); ok {
				*c = append(*c, a.Value(i))
			}
		case *array.Uint8:
			if c, ok := col.(*[]uint8); ok {
				*c = append(*c, a.Value(i))
			}
		case *array.Uint16:
			if c, ok := col.(*[]uint16); ok {
				*c = append(*c, a.Value(i))
			}
		case *array.Uint32:
			if c, ok := col.(*[]uint32); ok {
				*c = append(*c, a.Value(i))
			}
		case *array.Uint64:
			if c, ok := col.(*[]uint64); ok {
				*c = append(*c, a.Value(i))
			}
		case *array.Float32:
			if c, ok := col.(*[]float32); ok {
				*c = append(*c, a.Value(i))
			}
		case *array.Float64:
			if c, ok := col.(*[]float64); ok {
				*c = append(*c, a.Value(i))
			}
		case *array.String:
			if c, ok := col.(*[]string); ok {
				*c = append(*c, a.Value(i))
			}
		case *array.LargeString:
			if c, ok := col.(*[]string); ok {
				*c = append(*c, a.Value(i))
			}
		case *array.Boolean:
			if c, ok := col.(*[]bool); ok {
				*c = append(*c, a.Value(i))
			}
		case *array.Timestamp:
			if c, ok := col.(*[]time.Time); ok {
				*c = append(*c, a.Value(i).ToTime(a.DataType().(*arrow.TimestampType).Unit))
			}
		case *array.Date32:
			if c, ok := col.(*[]time.Time); ok {
				*c = append(*c, a.Value(i).ToTime())
			}
		case *array.Date64:
			if c, ok := col.(*[]time.Time); ok {
				*c = append(*c, a.Value(i).ToTime())
			}
		case *array.Decimal128:
			if c, ok := col.(*[]float64); ok {
				// Convert Decimal128 to float64
				val := a.Value(i)
				scale := a.DataType().(*arrow.Decimal128Type).Scale
				*c = append(*c, val.ToFloat64(scale))
			}
		case *array.Decimal256:
			if c, ok := col.(*[]float64); ok {
				// Convert Decimal256 to float64
				val := a.Value(i)
				scale := a.DataType().(*arrow.Decimal256Type).Scale
				*c = append(*c, val.ToFloat64(scale))
			}
		default:
			// Convert unknown types to string
			if c, ok := col.(*[]string); ok {
				*c = append(*c, fmt.Sprintf("%v", arr.ValueStr(i)))
			}
		}
	}
}

// appendNull appends a zero value for null handling (simplified - Grafana handles nulls differently)
func (d *Datasource) appendNull(col interface{}) {
	switch c := col.(type) {
	case *[]int8:
		*c = append(*c, 0)
	case *[]int16:
		*c = append(*c, 0)
	case *[]int32:
		*c = append(*c, 0)
	case *[]int64:
		*c = append(*c, 0)
	case *[]uint8:
		*c = append(*c, 0)
	case *[]uint16:
		*c = append(*c, 0)
	case *[]uint32:
		*c = append(*c, 0)
	case *[]uint64:
		*c = append(*c, 0)
	case *[]float32:
		*c = append(*c, 0)
	case *[]float64:
		*c = append(*c, 0)
	case *[]string:
		*c = append(*c, "")
	case *[]bool:
		*c = append(*c, false)
	case *[]time.Time:
		*c = append(*c, time.Time{})
	}
}

// columnToField converts a column slice to a Grafana data field
func (d *Datasource) columnToField(name string, col interface{}) *data.Field {
	switch c := col.(type) {
	case *[]int8:
		return data.NewField(name, nil, *c)
	case *[]int16:
		return data.NewField(name, nil, *c)
	case *[]int32:
		return data.NewField(name, nil, *c)
	case *[]int64:
		return data.NewField(name, nil, *c)
	case *[]uint8:
		return data.NewField(name, nil, *c)
	case *[]uint16:
		return data.NewField(name, nil, *c)
	case *[]uint32:
		return data.NewField(name, nil, *c)
	case *[]uint64:
		return data.NewField(name, nil, *c)
	case *[]float32:
		return data.NewField(name, nil, *c)
	case *[]float64:
		return data.NewField(name, nil, *c)
	case *[]string:
		return data.NewField(name, nil, *c)
	case *[]bool:
		return data.NewField(name, nil, *c)
	case *[]time.Time:
		return data.NewField(name, nil, *c)
	default:
		return nil
	}
}

// convertToTimeSeries converts a table frame to time series format
func (d *Datasource) convertToTimeSeries(frame *data.Frame) *data.Frame {
	// Look for a time column
	var timeFieldIdx = -1
	for i, field := range frame.Fields {
		name := strings.ToLower(field.Name)
		if name == "time" || name == "timestamp" || name == "ts" || name == "datetime" {
			timeFieldIdx = i
			break
		}
		// Also check if it's a time type
		if field.Type() == data.FieldTypeTime {
			timeFieldIdx = i
			break
		}
	}

	if timeFieldIdx == -1 {
		// No time column found, return as-is
		return frame
	}

	// Reorder so time is first (Grafana convention for time series)
	if timeFieldIdx != 0 {
		fields := make([]*data.Field, len(frame.Fields))
		fields[0] = frame.Fields[timeFieldIdx]
		j := 1
		for i, f := range frame.Fields {
			if i != timeFieldIdx {
				fields[j] = f
				j++
			}
		}
		frame.Fields = fields
	}

	return frame
}

// CheckHealth handles health checks for the data source
func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	// Try to execute a simple query to verify connectivity
	frame, err := d.executeFlightSQL(ctx, "SELECT 1")
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("Failed to connect to GizmoSQL: %v", err),
		}, nil
	}

	if frame == nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Connection succeeded but no data returned",
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Successfully connected to GizmoSQL",
	}, nil
}

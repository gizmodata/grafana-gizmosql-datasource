import React from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { InlineField, Select, CodeEditor } from '@grafana/ui';
import { DataSource } from '../datasource';
import { GizmoSQLDataSourceOptions, GizmoSQLQuery, defaultQuery } from '../types';

type Props = QueryEditorProps<DataSource, GizmoSQLQuery, GizmoSQLDataSourceOptions>;

const formatOptions: Array<SelectableValue<string>> = [
  { label: 'Table', value: 'table', description: 'Return results as a table' },
  { label: 'Time Series', value: 'time_series', description: 'Return results as time series (requires a time column)' },
];

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const { rawSql, format } = query;

  const onSqlChange = (sql: string) => {
    onChange({ ...query, rawSql: sql });
  };

  const onFormatChange = (option: SelectableValue<string>) => {
    onChange({ ...query, format: option.value as 'table' | 'time_series' });
    onRunQuery();
  };

  const onBlur = () => {
    onRunQuery();
  };

  return (
    <div>
      <InlineField label="Format" labelWidth={10} tooltip="Choose how to format the query results">
        <Select
          width={20}
          options={formatOptions}
          value={format || defaultQuery.format}
          onChange={onFormatChange}
        />
      </InlineField>
      <div style={{ marginTop: '8px' }}>
        <CodeEditor
          height="200px"
          language="sql"
          value={rawSql || defaultQuery.rawSql || ''}
          onBlur={onBlur}
          onSave={onBlur}
          onChange={onSqlChange}
          showMiniMap={false}
          showLineNumbers={true}
        />
      </div>
      <div style={{ marginTop: '8px', color: '#8e8e8e', fontSize: '12px' }}>
        <strong>Macros:</strong> Use <code>$__timeFrom</code>, <code>$__timeTo</code> for time range filters.
        For time series, include a <code>time</code> column and a <code>value</code> column.
      </div>
    </div>
  );
}


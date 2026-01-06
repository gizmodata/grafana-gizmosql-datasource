import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { GizmoSQLQuery, GizmoSQLDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, GizmoSQLQuery, GizmoSQLDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);


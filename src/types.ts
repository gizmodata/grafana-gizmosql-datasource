import { DataQuery, DataSourceJsonData } from '@grafana/data';

/**
 * Query model for GizmoSQL
 */
export interface GizmoSQLQuery extends DataQuery {
  rawSql: string;
  format?: 'table' | 'time_series';
}

/**
 * Default query values
 */
export const defaultQuery: Partial<GizmoSQLQuery> = {
  rawSql: 'SELECT 1',
  format: 'table',
};

/**
 * Configuration options saved in the data source settings
 */
export interface GizmoSQLDataSourceOptions extends DataSourceJsonData {
  host?: string;
  port?: number;
  username?: string;
  useTLS?: boolean;
  skipTLSVerify?: boolean;
}

/**
 * Secure configuration (stored encrypted)
 */
export interface GizmoSQLSecureJsonData {
  password?: string;
  token?: string;
}

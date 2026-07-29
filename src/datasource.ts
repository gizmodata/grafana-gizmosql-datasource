import {
  DataSourceInstanceSettings,
  DataSourceVariableSupport,
  CoreApp,
  ScopedVars,
} from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { GizmoSQLDataSourceOptions, GizmoSQLQuery, defaultQuery } from './types';

/**
 * Enables Query variables using the plugin's own query editor: the variable
 * query runs through the normal backend query path and Grafana maps the
 * resulting frame to variable options (first field, or __text/__value fields).
 */
class GizmoSQLVariableSupport extends DataSourceVariableSupport<DataSource, GizmoSQLQuery, GizmoSQLDataSourceOptions> {
  getDefaultQuery(): Partial<GizmoSQLQuery> {
    return defaultQuery;
  }
}

export class DataSource extends DataSourceWithBackend<GizmoSQLQuery, GizmoSQLDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<GizmoSQLDataSourceOptions>) {
    super(instanceSettings);
    this.variables = new GizmoSQLVariableSupport();
  }

  /**
   * Get default query for new panels
   */
  getDefaultQuery(_: CoreApp): Partial<GizmoSQLQuery> {
    return defaultQuery;
  }

  /**
   * Apply Grafana template variables to the query
   */
  applyTemplateVariables(query: GizmoSQLQuery, scopedVars: ScopedVars): GizmoSQLQuery {
    const templateSrv = getTemplateSrv();

    return {
      ...query,
      rawSql: templateSrv.replace(query.rawSql, scopedVars, this.interpolateVariable),
    };
  }

  /**
   * Custom variable interpolation for SQL queries
   */
  interpolateVariable = (value: string | string[], variable: { multi?: boolean; includeAll?: boolean }): string => {
    if (typeof value === 'string') {
      return this.quoteLiteral(value);
    }

    if (Array.isArray(value)) {
      return value.map((v) => this.quoteLiteral(v)).join(',');
    }

    return String(value);
  };

  /**
   * Quote a string literal for SQL
   */
  private quoteLiteral(value: string): string {
    // Check if it's a number
    if (!isNaN(Number(value))) {
      return value;
    }
    // Escape single quotes and wrap in quotes
    return `'${value.replace(/'/g, "''")}'`;
  }

  /**
   * Filter query - only run queries that have SQL
   */
  filterQuery(query: GizmoSQLQuery): boolean {
    return !!query.rawSql && query.rawSql.trim().length > 0;
  }
}


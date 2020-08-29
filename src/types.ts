import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface SQLiteQuery extends DataQuery {
  queryText?: string;
  format: 'time_series';
}

export const defaultQuery: Partial<SQLiteQuery> = {
  queryText: 'SELECT 1',
  format: 'time_series',
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  path?: string;
}

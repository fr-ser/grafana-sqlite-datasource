import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface SQLiteQuery extends DataQuery {
  queryText?: string;
}

export const defaultQuery: Partial<SQLiteQuery> = {
  queryText: 'SELECT 1 as time, 4 as value',
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  path?: string;
}

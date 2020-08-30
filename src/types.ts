import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface SQLiteQuery extends DataQuery {
  rawQueryText: string;
  queryText: string;
}

export const defaultQuery: Partial<SQLiteQuery> = {
  rawQueryText: 'SELECT 1 as time, 4 as value where time >= $__from and time < $__to',
  queryText: 'SELECT 1 as time, 4 as value where time >= 1234 and time < 134567',
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  path?: string;
}

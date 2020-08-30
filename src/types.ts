import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface SQLiteQuery extends DataQuery {
  rawQueryText: string;
  queryText: string;
  timeColumns: string[];
}

export const defaultQuery: Partial<SQLiteQuery> = {
  rawQueryText: 'SELECT 1 as time, 4 as value where time >= $__from / 1000 and time < $__to / 1000',
  queryText: 'SELECT 1 as time, 4 as value where time >= 1234 and time < 134567',
  timeColumns: ['time', 'ts'],
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  path?: string;
}

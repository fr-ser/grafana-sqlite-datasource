import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface SQLiteQuery extends DataQuery {
  rawQueryText: string;
  queryText: string;
  timeColumns: string[];
}

export const defaultQuery: Partial<SQLiteQuery> = {
  rawQueryText:
    "SELECT CAST(strftime('%s', 'now', '-1 minute') as INTEGER) as time, 4 as value \n" +
    'WHERE time >= $__from / 1000 and time < $__to / 1000',
  queryText: `
    SELECT CAST(strftime('%s', 'now', '-1 minute') as INTEGER) as time, 4 as value
    WHERE time >= 1234 and time < 134567
  `,
  timeColumns: ['time', 'ts'],
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  path?: string;
}

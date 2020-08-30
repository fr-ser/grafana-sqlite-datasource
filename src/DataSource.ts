import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { MyDataSourceOptions, SQLiteQuery } from './types';

export class DataSource extends DataSourceWithBackend<SQLiteQuery, MyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: SQLiteQuery): SQLiteQuery {
    query.queryText = getTemplateSrv().replace(query.rawQueryText);
    return query;
  }
}

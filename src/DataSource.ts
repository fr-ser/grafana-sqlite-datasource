import { DataSourceInstanceSettings, DataFrame, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { MyDataSourceOptions, SQLiteQuery } from './types';

export class DataSource extends DataSourceWithBackend<SQLiteQuery, MyDataSourceOptions> {
  templateSrv;

  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);

    this.templateSrv = getTemplateSrv();
    this.annotations = {};
  }

  applyTemplateVariables(query: SQLiteQuery, scopedVars: ScopedVars): SQLiteQuery {
    query.queryText = this.templateSrv.replace(query.rawQueryText, scopedVars);
    return query;
  }

  async metricFindQuery(query: string, options?: any) {
    if (!query) {
      return [];
    }
    const response = await this.query({
      targets: [
        {
          refId: 'metricFindQuery',
          rawQueryText: query,
          queryText: query,
          timeColumns: [],
        },
      ],
    } as any).toPromise();

    if (response === undefined) {
      throw new Error('Received no response at all');
    } else if (response.error) {
      throw new Error(response.error.message);
    }

    const data = response.data[0] as DataFrame;
    if (data.fields.length === 1) {
      return data.fields[0].values.toArray().map((text) => ({ text }));
    } else if (data.fields.length === 2) {
      const textIndex = data.fields.findIndex((x) => x.name === '__text');
      const valueIndex = data.fields.findIndex((x) => x.name === '__value');
      if (textIndex === -1 || valueIndex === -1) {
        throw new Error(
          `No columns named "__text" and "__value" were found. Columns: ${data.fields.map((x) => x.name).join(',')}`
        );
      }

      const valueArray = data.fields[valueIndex].values.toArray();
      return data.fields[textIndex].values.toArray().map((text, index) => ({ text, value: valueArray[index] }));
    } else {
      throw new Error(
        `Received more than two (${data.fields.length}) fields: ${data.fields.map((x) => x.name).join(',')}`
      );
    }
  }
}

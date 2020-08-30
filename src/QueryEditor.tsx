import { TextArea } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import defaults from 'lodash/defaults';
import React, { ChangeEvent, PureComponent } from 'react';
import { getTemplateSrv } from '@grafana/runtime';

import { DataSource } from './DataSource';
import { defaultQuery, MyDataSourceOptions, SQLiteQuery } from './types';

type Props = QueryEditorProps<DataSource, SQLiteQuery, MyDataSourceOptions>;

const templateSrv = getTemplateSrv();

export class QueryEditor extends PureComponent<Props> {
  onQueryTextChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { onChange, query } = this.props;
    onChange({
      ...query,
      rawQueryText: event.target.value,
      queryText: templateSrv.replace(event.target.value),
    });
  };

  sendQuery = () => {
    this.props.onRunQuery();
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { rawQueryText } = query;

    return (
      <div className="gf-form">
        <TextArea
          role="query-editor-input"
          value={rawQueryText || ''}
          onBlur={this.sendQuery}
          onChange={this.onQueryTextChange}
        />
      </div>
    );
  }
}

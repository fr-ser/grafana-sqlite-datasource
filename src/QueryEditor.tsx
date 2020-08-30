import { TextArea, TagsInput, InlineFormLabel } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import defaults from 'lodash/defaults';
import React, { ChangeEvent, PureComponent } from 'react';

import { DataSource } from './DataSource';
import { defaultQuery, MyDataSourceOptions, SQLiteQuery } from './types';

type Props = QueryEditorProps<DataSource, SQLiteQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onQueryTextChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { onChange, query } = this.props;
    onChange({
      ...query,
      rawQueryText: event.target.value,
    });
  };

  sendQuery = () => this.props.onRunQuery();

  onUpdateColumnTypes = (columnKey: string, columns: string[]) => {
    const { onChange, query } = this.props;
    onChange({
      ...query,
      [columnKey]: columns,
    });

    this.sendQuery();
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { rawQueryText, timeColumns } = query;

    return (
      <>
        <div className="gf-form">
          <TextArea
            role="query-editor-input"
            value={rawQueryText}
            onBlur={this.sendQuery}
            onChange={this.onQueryTextChange}
          />
        </div>
        <div className="gf-form">
          <div style={{ display: 'flex', flexDirection: 'column', marginRight: 15 }} role="time-column-selector">
            <InlineFormLabel>Time formatted columns</InlineFormLabel>
            <TagsInput
              onChange={(tags: string[]) => this.onUpdateColumnTypes('timeColumns', tags)}
              tags={timeColumns}
            />
          </div>
        </div>
      </>
    );
  }
}

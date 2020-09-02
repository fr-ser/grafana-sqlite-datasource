import { TextArea, TagsInput, Icon, Alert, InlineFormLabel } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import defaults from 'lodash/defaults';
import React, { ChangeEvent, useState } from 'react';

import { DataSource } from './DataSource';
import { defaultQuery, MyDataSourceOptions, SQLiteQuery } from './types';

type Props = QueryEditorProps<DataSource, SQLiteQuery, MyDataSourceOptions>;

export function QueryEditor(props: Props) {
  function onQueryTextChange(event: ChangeEvent<HTMLTextAreaElement>) {
    const { onChange, query } = props;
    onChange({
      ...query,
      rawQueryText: event.target.value,
    });
  }

  function onUpdateColumnTypes(columnKey: string, columns: string[]) {
    const { onChange, query } = props;
    onChange({
      ...query,
      [columnKey]: columns,
    });

    props.onRunQuery();
  }

  const query = defaults(props.query, defaultQuery);
  const { rawQueryText, timeColumns } = query;
  const [showHelp, setShowHelp] = useState(false);

  return (
    <>
      <div className="gf-form">
        <TextArea
          style={{ height: 100 }}
          role="query-editor-input"
          value={rawQueryText}
          onBlur={() => props.onRunQuery()}
          onChange={onQueryTextChange}
        />
      </div>
      <div className="gf-form">
        <div style={{ display: 'flex', flexDirection: 'column', marginRight: 15 }} role="time-column-selector">
          <InlineFormLabel>
            <div style={{ whiteSpace: 'nowrap' }} onClick={() => setShowHelp(!showHelp)}>
              Time formatted columns <Icon name={showHelp ? 'angle-down' : 'angle-right'} />
            </div>
          </InlineFormLabel>
          <TagsInput onChange={(tags: string[]) => onUpdateColumnTypes('timeColumns', tags)} tags={timeColumns} />
        </div>
      </div>
      {showHelp && (
        <Alert title="Time formatted columns" severity="info">
          Columns with these names, will be formatted as time. This is required as SQLite has no native "time" format,
          but mostly strings and numbers. See:{' '}
          <a href="https://www.sqlite.org/datatype3.html" target="_blank">
            SQLite3 Data Types Documentation
          </a>
        </Alert>
      )}
    </>
  );
}

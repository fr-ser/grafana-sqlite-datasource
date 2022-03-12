import { TextArea, TagsInput, Icon, Alert, InlineFormLabel, Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
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
      queryText: event.target.value,
    });
  }

  function onQueryTypeChange(value: SelectableValue<string>) {
    const { onChange, query } = props;
    onChange({
      ...query,
      queryType: value.value || 'table',
    });

    props.onRunQuery();
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

  const options: Array<SelectableValue<string>> = [
    { label: 'Table', value: 'table' },
    { label: 'Time series', value: 'time series' },
  ];
  const selectedOption = options.find((options) => options.value === query.queryType) || options[0];
  return (
    <>
      <div className="gf-form max-width-8" role="query-type-container">
        <Select
          allowCustomValue={false}
          isSearchable={false}
          onChange={onQueryTypeChange}
          options={options}
          value={selectedOption}
        />
      </div>
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
          Columns with these names, will be formatted as time. This is required as SQLite has no native &quot;time&quot;
          format, but mostly strings and numbers. See:{' '}
          <a href="https://www.sqlite.org/datatype3.html" target="_blank" rel="noreferrer">
            SQLite3 Data Types Documentation
          </a>
          <br />
          For more information (like supported formats) see:{' '}
          <a
            href="https://github.com/fr-ser/grafana-sqlite-datasource#support-for-time-formatted-columns"
            target="_blank"
            rel="noreferrer"
          >
            Plugin documentation
          </a>
        </Alert>
      )}
    </>
  );
}

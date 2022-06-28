import { TextArea, TagsInput, Icon, Alert, InlineFormLabel, Select, CodeEditor, Switch } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import defaults from 'lodash/defaults';
import React, { ChangeEvent, useState } from 'react';

import { DataSource } from './DataSource';
import { defaultQuery, MyDataSourceOptions, SQLiteQuery } from './types';

type Props = QueryEditorProps<DataSource, SQLiteQuery, MyDataSourceOptions>;

function calculateHeight(queryText: string): number {
  const minHeight = 200;
  const maxHeight = 500;

  // assume 20 px per row
  let desiredHeight = queryText.split('\n').length * 20;

  // return the value in a range between the min and max height
  return Math.min(maxHeight, Math.max(minHeight, desiredHeight));
}

export function QueryEditor(props: Props) {
  function onQueryTextChange(value: string) {
    const { onChange, query } = props;
    onChange({
      ...query,
      rawQueryText: value,
      queryText: value,
    });

    props.onRunQuery();
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
  const [useLegacyEditor, setUseLegacyEditor] = useState(false);

  const options: Array<SelectableValue<string>> = [
    { label: 'Table', value: 'table' },
    { label: 'Time series', value: 'time series' },
  ];
  const selectedOption = options.find((options) => options.value === query.queryType) || options[0];

  return (
    <>
      {useLegacyEditor ? (
        <div className="gf-form">
          <TextArea
            style={{ height: 100 }}
            role="query-editor-input"
            value={rawQueryText}
            onBlur={() => props.onRunQuery()}
            onChange={(event: ChangeEvent<HTMLTextAreaElement>) => onQueryTextChange(event.target.value)}
          />
        </div>
      ) : (
        <CodeEditor
          height={calculateHeight(rawQueryText)}
          value={rawQueryText}
          onBlur={onQueryTextChange}
          onSave={onQueryTextChange}
          language="sql"
          showMiniMap={false}
        />
      )}
      <div className="gf-form-inline">
        <div className="gf-form" role="query-type-container" style={{ marginRight: 15 }}>
          <InlineFormLabel>
            <div style={{ whiteSpace: 'nowrap' }}>Format as:</div>
          </InlineFormLabel>
          <Select
            allowCustomValue={false}
            isSearchable={false}
            onChange={onQueryTypeChange}
            options={options}
            value={selectedOption}
          />
        </div>
        <div className="gf-form">
          <div style={{ display: 'flex', flexDirection: 'row', marginRight: 15 }} role="time-column-selector">
            <InlineFormLabel>
              <div style={{ whiteSpace: 'nowrap' }} onClick={() => setShowHelp(!showHelp)}>
                Time formatted columns <Icon name={showHelp ? 'angle-down' : 'angle-right'} />
              </div>
            </InlineFormLabel>
            <TagsInput onChange={(tags: string[]) => onUpdateColumnTypes('timeColumns', tags)} tags={timeColumns} />
          </div>
          <div className="gf-form" style={{ alignItems: 'center' }}>
            <InlineFormLabel>
              <div style={{ whiteSpace: 'nowrap' }}>Use legacy code editor:</div>
            </InlineFormLabel>
            <Switch
              role="use-legacy-editor-switch"
              value={useLegacyEditor}
              onChange={() => setUseLegacyEditor(!useLegacyEditor)}
            />
          </div>
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

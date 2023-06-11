import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms, Alert } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';

import { MyDataSourceOptions, MySecureJsonData } from './types';

const { FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions, MySecureJsonData> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onPathPrefixChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      pathPrefix: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onPathOptionsChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      pathOptions: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onSecurePathOptionsChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        securePathOptions: event.target.value,
      },
    });
  };

  onPathChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      path: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onAttachLimitChange = (event: ChangeEvent<HTMLInputElement>) => {
    let value: number | undefined = undefined;

    if (event.target.value !== '') {
      value = parseInt(event.target.value, 10);
      if (Number.isNaN(value)) {
        return;
      }
    }

    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      attachLimit: value,
    };

    onOptionsChange({ ...options, jsonData });
  };

  render() {
    const { options, onOptionsChange } = this.props;
    const { jsonData, secureJsonFields, secureJsonData } = options;
    if (jsonData.pathPrefix === undefined) {
      onOptionsChange({
        ...options,
        jsonData: { ...options.jsonData, pathPrefix: 'file:' },
      });
    }
    if (jsonData.attachLimit === undefined) {
      onOptionsChange({
        ...options,
        jsonData: { ...options.jsonData, attachLimit: 0 },
      });
    }

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            label="Path"
            tooltip="(absolute) path to the SQLite database file"
            labelWidth={10}
            inputWidth={20}
            onChange={this.onPathChange}
            value={jsonData.path}
            placeholder="/path/to/the/database.db"
          />
        </div>
        <div className="gf-form">
          <FormField
            label="Path Prefix"
            tooltip={
              'This string is prefixed before the path in the connection string. ' +
              'Unless you know what you are doing this should be "file:" (without the quotes). ' +
              'Not using "file:" can cause the Path Options to not take effect.'
            }
            labelWidth={10}
            inputWidth={20}
            onChange={this.onPathPrefixChange}
            value={jsonData.pathPrefix}
          />
        </div>
        <div className="gf-form">
          <FormField
            label="Path Options"
            tooltip={
              'This string is appended to the path (after adding a "?") when opening the ' +
              'database. A typical example is "mode=ro" (without the quotes) for readonly mode.'
            }
            labelWidth={10}
            inputWidth={20}
            onChange={this.onPathOptionsChange}
            value={jsonData.pathOptions}
            placeholder="mode=ro&_ignore_check_constraints=1"
          />
        </div>
        <div className="gf-form">
          <FormField
            label="Secure Path Options"
            tooltip={
              'This is combined with the regular path options. Typical for the secure options ' +
              'are credentials (options starting with _auth).'
            }
            labelWidth={10}
            inputWidth={20}
            placeholder={secureJsonFields?.securePathOptions ? 'configured' : ''}
            value={secureJsonData?.securePathOptions}
            onChange={this.onSecurePathOptionsChange}
          />
        </div>
        <div className="gf-form">
          <FormField
            label="Attach limit"
            tooltip="The runtime limit for attached databases (see: https://www.sqlite.org/limits.html)."
            labelWidth={10}
            inputWidth={20}
            value={jsonData.attachLimit}
            onChange={this.onAttachLimitChange}
          />
        </div>
        <div className="gf-form">
          <Alert title="File System Permissions" severity="info">
            <div>
              The plugin runs with the same permissions as the Grafana user. Any file that can be opened with the
              Grafana user can be opened with the SQLite plugin.
            </div>
            <div>
              Beware that by enabling attaching databases (setting an &quot;attach limit&quot; above 0) you enable any
              user with access to the plugin to attach any database that the Grafana user has access to.
            </div>
            <div>It is the most secure (and recommended) approach to set the &quot;attach limit&quot; to 0.</div>
          </Alert>
        </div>
      </div>
    );
  }
}

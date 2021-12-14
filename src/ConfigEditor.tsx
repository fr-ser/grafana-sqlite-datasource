import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
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

  render() {
    const { options } = this.props;
    const { jsonData, secureJsonFields, secureJsonData } = options;
    if (jsonData.pathPrefix === undefined) {
      jsonData.pathPrefix = 'file:';
    }

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            label="Path"
            tooltip="(absolute) path to the SQLite database"
            labelWidth={10}
            inputWidth={20}
            onChange={this.onPathChange}
            value={jsonData.path || ''}
            placeholder="/path/to/the/database.db"
          />
        </div>
        <div className="gf-form">
          <FormField
            label="Path Prefix"
            tooltip={
              'This string is prefixed before the path in the connection string. </br>' +
              'Unless you know what you are doing this should be "file:" (without the quotes). ' +
              'Not using "file:" can cause the Path Options to not take effect.'
            }
            labelWidth={10}
            inputWidth={20}
            onChange={this.onPathPrefixChange}
            value={jsonData.pathPrefix || ''}
            placeholder="Connection String Prefix"
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
            value={jsonData.pathOptions || ''}
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
            value={secureJsonData?.securePathOptions ?? ''}
            onChange={this.onSecurePathOptionsChange}
          />
        </div>
      </div>
    );
  }
}

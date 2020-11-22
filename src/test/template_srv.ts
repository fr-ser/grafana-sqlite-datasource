import * as _ from 'lodash';
import { TimeRange, VariableModel } from '@grafana/data';

export class TemplateSrv {
  variables: any[] = [];
  timeRange: any = {};

  private regex = /\$(\w+)|\[\[([\s\S]+?)(?::(\w+))?\]\]|\${(\w+)(?::(\w+))?}/g;
  private index: any = {};
  private grafanaVariables: any = {};

  constructor() {}

  init(variables: any, timeRange?: TimeRange) {
    this.variables = variables;
    this.updateTemplateData();
    this.timeRange = timeRange;
  }

  highlightVariablesAsHtml() {}
  updateTemplateData() {
    this.index = {};

    for (var i = 0; i < this.variables.length; i++) {
      var variable = this.variables[i];

      if (!variable.current || (!variable.current.isNone && !variable.current.value)) {
        continue;
      }

      this.index[variable.name] = variable;
    }
  }

  formatValue(value: any, format: any, variable: any) {
    if (typeof format === 'function') {
      return format(value, variable, this.formatValue);
    }

    if (_.isString(value)) {
      return value;
    }
    return value.join(',');
  }

  replace(target: string, scopedVars?: any, format?: string | Function): any {
    if (!target) {
      return target;
    }

    let variable, systemValue, value, fmt;
    this.regex.lastIndex = 0;

    return target.replace(this.regex, (match, var1, var2, fmt2, var3, fmt3) => {
      variable = this.index[var1 || var2 || var3];
      fmt = fmt2 || fmt3 || format;
      if (scopedVars) {
        value = scopedVars[var1 || var2 || var3];
        if (value) {
          return this.formatValue(value.value, fmt, variable);
        }
      }

      if (!variable) {
        return match;
      }

      systemValue = this.grafanaVariables[variable.current.value];
      if (systemValue) {
        return this.formatValue(systemValue, fmt, variable);
      }

      value = variable.current.value;
      if (this.isAllValue(value)) {
        value = this.getAllValue(variable);
        // skip formatting of custom all values
        if (variable.allValue) {
          return this.replace(value as any);
        }
      }

      const res = this.formatValue(value, fmt, variable);
      return res;
    });
  }

  isAllValue(value: any) {
    return false;
  }

  getAllValue(variable: any) {
    return null;
  }

  getVariables(): VariableModel[] {
    return this.variables;
  }
}

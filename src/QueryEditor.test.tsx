import * as ui from '@grafana/ui';
import { fireEvent, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';

import { QueryEditor } from './QueryEditor';

// based on examples from the grafana repo
// https://github.com/grafana/grafana/blob/467e375fe6c3de0309a69664b32301e22c0f5f7e/public/app/plugins/datasource/cloudwatch/components/MetricsQueryEditor/MetricsQueryEditor.test.tsx#L18
jest.mock('@grafana/ui', () => ({
  ...jest.requireActual<typeof ui>('@grafana/ui'),
  CodeEditor: function CodeEditor({ value }: { value: string }) {
    return <pre>{value}</pre>;
  },
}));

describe('QueryEditor', () => {
  let onChangeMock: jest.Mock;
  let onRunQueryMock: jest.Mock;
  let queryEditor: React.ReactElement;

  beforeEach(() => {
    onChangeMock = jest.fn();
    onRunQueryMock = jest.fn();
    queryEditor = (
      <QueryEditor onChange={onChangeMock} onRunQuery={onRunQueryMock} query={null as any} datasource={null as any} />
    );
  });

  it.only('allows editing the queryType', async () => {
    const { findByRole } = render(queryEditor);
    const queryTypeContainer = await findByRole('query-type-container');

    fireEvent.focus(queryTypeContainer.querySelector('input') as HTMLInputElement);
    fireEvent.keyDown(queryTypeContainer.querySelector('input') as HTMLInputElement, { key: 'Down', code: 'Down' });
    fireEvent.click(await screen.findByText('Time series'));

    expect(onRunQueryMock).toHaveBeenCalled();
    expect(onChangeMock).toHaveBeenLastCalledWith({
      queryType: 'time series',
    });
  });

  it('allows adding time columns', async () => {
    const { findByRole } = render(queryEditor);

    const selector = await findByRole('time-column-selector');
    const selectorInput = selector.querySelector('input') as HTMLInputElement;

    await userEvent.type(selectorInput, 'test_column', { delay: 1 });
    await userEvent.keyboard('{enter}');

    expect(onRunQueryMock).toHaveBeenCalled();
    expect(onChangeMock).toHaveBeenLastCalledWith({
      timeColumns: ['time', 'ts', 'test_column'],
    });
  });

  it('allows removing time columns', async () => {
    const { findByText } = render(queryEditor);

    const timeTag = await findByText('time', { selector: "div[role='time-column-selector'] span" });
    await userEvent.click(timeTag.parentElement!.querySelector('svg') as SVGElement);

    expect(onRunQueryMock).toHaveBeenCalled();
    expect(onChangeMock).toHaveBeenLastCalledWith({
      timeColumns: ['ts'],
    });
  });
});

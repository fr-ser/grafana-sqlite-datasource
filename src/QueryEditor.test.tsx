import React from 'react';
import { render, fireEvent, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { QueryEditor } from './QueryEditor';

describe('QueryEditor', () => {
  let onChangeMock: jest.Mock;
  let onRunQueryMock: jest.Mock;
  let queryEditor: JSX.Element;

  beforeEach(() => {
    onChangeMock = jest.fn();
    onRunQueryMock = jest.fn();
    queryEditor = (
      <QueryEditor onChange={onChangeMock} onRunQuery={onRunQueryMock} query={null as any} datasource={null as any} />
    );
  });

  it('allows editing the queryType', async () => {
    const { findByRole, findByText } = render(queryEditor);
    const queryTypeContainer = await findByRole('query-type-container');

    await act(async () => {
      fireEvent.focus(queryTypeContainer.querySelector('input') as HTMLInputElement);
      fireEvent.keyDown(queryTypeContainer.querySelector('input') as HTMLInputElement, { key: 'Down', code: 'Down' });
      fireEvent.click(await findByText('Time series'));
    });

    expect(onRunQueryMock).toHaveBeenCalled();
    expect(onChangeMock).toHaveBeenLastCalledWith({
      queryType: 'time series',
    });
  });

  it('allows editing the rawQuery', async () => {
    const { findByRole } = render(queryEditor);

    const queryInput = await findByRole('query-editor-input');

    await act(async () => {
      fireEvent.click(queryInput);

      fireEvent.change(queryInput, {
        target: { value: 'Some Input' },
      });

      fireEvent.blur(queryInput);
    });

    expect(onRunQueryMock).toHaveBeenCalled();
    expect(onChangeMock).toHaveBeenLastCalledWith({
      queryText: 'Some Input',
      rawQueryText: 'Some Input',
    });
  });

  it('allows adding time columns', async () => {
    const { findByRole } = render(queryEditor);

    const selector = await findByRole('time-column-selector');
    const selectorInput = selector.querySelector('input') as HTMLInputElement;

    await act(async () => {
      await userEvent.type(selectorInput, 'test_column', { delay: 1 });
      userEvent.keyboard('{enter}');
    });

    expect(onRunQueryMock).toHaveBeenCalled();
    expect(onChangeMock).toHaveBeenLastCalledWith({
      timeColumns: ['time', 'ts', 'test_column'],
    });
  });

  it('allows removing time columns', async () => {
    const { findByText } = render(queryEditor);

    await act(async () => {
      const timeTag = await findByText('time', { selector: 'div>span' });
      userEvent.click(timeTag.parentElement!.querySelector('svg') as SVGElement);
    });

    expect(onRunQueryMock).toHaveBeenCalled();
    expect(onChangeMock).toHaveBeenLastCalledWith({
      timeColumns: ['ts'],
    });
  });
});

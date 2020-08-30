import React from 'react';
import { render, fireEvent, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { QueryEditor } from './QueryEditor';

describe('QueryEditor', () => {
  it('allows editing the rawQuery', async () => {
    const onChangeMock = jest.fn();
    const onRunQueryMock = jest.fn();
    const { findByRole } = render(<QueryEditor onChange={onChangeMock} onRunQuery={onRunQueryMock} />);

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
      rawQueryText: 'Some Input',
    });
  });

  it('allows setting time columns', async () => {
    const onChangeMock = jest.fn();
    const onRunQueryMock = jest.fn();
    const { findByRole, findByText } = render(<QueryEditor onChange={onChangeMock} onRunQuery={onRunQueryMock} />);

    const selector = await findByRole('time-column-selector');
    const selectorInput = selector.querySelector('input') as HTMLInputElement;
    const addButton = await findByText('Add', { selector: 'button span' });

    await act(async () => {
      // add a column
      await userEvent.type(selectorInput, 'test_column', { delay: 1 });
      userEvent.click(addButton);

      // remove a default column
      const timeTag = await findByText('time', { selector: 'div>span' });
      userEvent.click(timeTag.parentElement!.querySelector('svg') as SVGElement);
    });

    expect(onRunQueryMock).toHaveBeenCalled();
    expect(onChangeMock).toHaveBeenLastCalledWith({
      timeColumns: ['ts', 'test_column'],
    });
  });
});

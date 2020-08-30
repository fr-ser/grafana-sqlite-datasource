import React from 'react';
import { render, fireEvent, act } from '@testing-library/react';

import { QueryEditor } from './QueryEditor';

jest.mock('@grafana/runtime', () => ({
  getTemplateSrv: () => {
    return {
      replace: jest.fn(input => `${input}-replaced`),
    };
  },
}));
describe('QueryEditor', () => {
  it('converts query variables on change', async () => {
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

    expect(queryInput).toBeTruthy();
    expect(onRunQueryMock).toHaveBeenCalled();
    expect(onChangeMock).toHaveBeenLastCalledWith({
      rawQueryText: 'Some Input',
      queryText: 'Some Input-replaced',
    });
  });
});

import { setTemplateSrv } from '@grafana/runtime';

import { TemplateSrv } from './test/template_srv';
import { DataSource } from './DataSource';
import { FieldType, MutableDataFrame } from '@grafana/data';

describe('DataSource', () => {
  beforeEach(() => {
    setTemplateSrv(new TemplateSrv());
  });
  describe('variable replacing', () => {
    it('used the template service for replacing variables', () => {
      const ds = new DataSource({} as any);
      const mockReplace = jest.fn((input: string) => 'mock response');
      ds.templateSrv.replace = mockReplace;

      const result = ds.applyTemplateVariables(
        {
          rawQueryText: 'SELECT 1',
          queryText: '',
        } as any,
        {
          __interval: { text: '2m', value: '2m' },
          __interval_ms: { text: '120000', value: 120000 },
        }
      );

      expect(mockReplace.mock.calls[0][0]).toBe('SELECT 1');
      expect(result.queryText).toBe('mock response');
    });
  });

  describe('query variables', () => {
    it('returns the result of a query', async () => {
      const ds = new DataSource({} as any);
      const mockResponse = {
        data: [
          new MutableDataFrame({
            fields: [{ name: 'value', type: FieldType.number, values: [1, 2] }],
          }),
        ],
      };
      ds.query = jest.fn(() => ({ toPromise: async () => mockResponse })) as any;

      const result = await ds.metricFindQuery("SELECT 'my-query'", { variable: { datasource: 'sqlite' } });

      expect(result).toStrictEqual([{ text: 1 }, { text: 2 }]);
    });

    it('throws for multiple columns', async () => {
      const ds = new DataSource({} as any);
      const mockResponse = {
        data: [
          new MutableDataFrame({
            fields: [
              { name: 'value', type: FieldType.number, values: [1, 2] },
              { name: 'name', type: FieldType.number, values: ['a', 'b'] },
            ],
          }),
        ],
      };
      ds.query = jest.fn(() => ({ toPromise: async () => mockResponse })) as any;

      try {
        await ds.metricFindQuery("SELECT 'my-query'", { variable: { datasource: 'sqlite' } });
        fail('did not receive an error');
      } catch (error) {
        const errorMessage = (error as Error).toString();
        expect(errorMessage).toContain('Received more than one (2) fields');
      }
    });

    it('throws for a server error', async () => {
      const ds = new DataSource({} as any);
      const mockResponse = {
        error: { message: 'test error' },
      };
      ds.query = jest.fn(() => ({ toPromise: async () => mockResponse })) as any;

      try {
        await ds.metricFindQuery("SELECT 'my-query'", { variable: { datasource: 'sqlite' } });
        fail('did not receive an error');
      } catch (error) {
        const errorMessage = (error as Error).toString();
        expect(errorMessage).toContain('test error');
      }
    });
  });
});

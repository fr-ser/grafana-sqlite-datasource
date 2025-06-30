// Jest setup provided by Grafana scaffolding
import './.config/jest-setup';

// mock the intersection observer and just say everything is in view
// copied from https://github.com/grafana/grafana/blob/0a90b7b5e9c6fc691743389b78f8116427ec5210/public/test/jest-setup.ts#L43C1-L53C56
const mockIntersectionObserver = jest.fn().mockImplementation((callback) => ({
  observe: jest.fn().mockImplementation((elem) => {
    callback([{ target: elem, isIntersecting: true }]);
  }),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));
global.IntersectionObserver = mockIntersectionObserver;

// @vitest-environment jsdom
import { describe, it, expect, vi } from 'vitest';
import { render } from '@testing-library/react';
import App from './App';

Object.defineProperty(window, 'localStorage', {
  value: {
    getItem: vi.fn().mockReturnValue(null),
    setItem: vi.fn(),
    removeItem: vi.fn(),
    clear: vi.fn(),
  },
  writable: true,
});

describe('App Component', () => {
  it('renders heading', () => {
    const { getByText } = render(<App />);
    expect(getByText('BooksLib')).toBeDefined();
  });
});

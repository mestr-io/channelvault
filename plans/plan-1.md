# Testing Plan: Go + Wails + Svelte FE Application

## Overview

This document outlines the testing strategy for the ChannelVault (Slack Export Viewer) application, which uses:
- **Backend**: Go with Wails v2
- **Frontend**: Svelte 5 + TypeScript + Vite

---

## Current State

- **Backend (Go)**: Test infrastructure exists
  - Uses `testing` package + `testify/assert`
  - Test files: `*_test.go` throughout backend
  - Build tags: `headless`, `webkit2_41`, `fts5`

- **Frontend (Svelte)**: Now has full testing setup
  - Vitest + @testing-library/svelte + jsdom installed
  - Configured for Svelte 5 component testing (Runes support + Browser environment)
  - Test files:
    - `frontend/src/lib/components/DayDivider.test.ts` (Visual rendering test)
    - `frontend/src/lib/components/Prompt.test.ts` (Interaction test)

---

## Implemented

### Frontend Setup

```bash
cd frontend
npm install -D vitest @testing-library/svelte jsdom @testing-library/jest-dom
```

Added to `frontend/package.json`:
```json
{
  "scripts": {
    "test": "vitest",
    "test:run": "vitest run"
  }
}
```

Created `frontend/vitest.config.ts`:
```typescript
import { defineConfig } from 'vitest/config'
import { svelte } from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  plugins: [
    svelte({
      hot: !process.env.VITEST,
      compilerOptions: {
        runes: true,
      },
    }),
  ],
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/test/setup.ts'],
    include: ['src/**/*.{test,spec}.{js,ts}'],
    pool: 'forks',
  },
  resolve: {
    conditions: ['browser'], // Critical for Svelte 5 component testing in Vitest
  },
})
```

Created `frontend/src/test/setup.ts`:
```typescript
import '@testing-library/jest-dom'
```

### Example Tests

1. `frontend/src/lib/components/Prompt.test.ts` - Integration test that opens a modal, types into an input, and clicks a button to verify callback execution.
2. `frontend/src/lib/components/DayDivider.test.ts` - Component test verifying visual output based on props.

---

## Test Commands

### Backend
```bash
# Unit tests
make test

# Headless tests (with xvfb)
make testh
```

### Frontend
```bash
# Run tests (watch mode)
cd frontend && npm run test

# Run tests once
cd frontend && npm run test:run

# Type check
cd frontend && npm run check
```

### All
```bash
# Run all tests
make test && make test-fe
```

---

## CI/CD Integration

The Makefile has been updated:
- `make test` - Backend tests (uses headless tag)
- `make test-fe` - Frontend tests

---

## Test Organization

```
backend/
  */*_test.go           # Unit tests

frontend/
  src/
    lib/
      components/
        DayDivider.test.ts       # Visual rendering tests
        Prompt.test.ts           # Interaction tests
    test/
      setup.ts          # Test setup
  vitest.config.ts      # Vitest config
```

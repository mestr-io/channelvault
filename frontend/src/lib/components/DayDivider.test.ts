import { render, screen } from '@testing-library/svelte'
import { describe, it, expect } from 'vitest'
import DayDivider from './DayDivider.svelte'
import { models } from '../../../wailsjs/go/models'

describe('DayDivider', () => {
  it('shows day pill for first message', () => {
    const current = new models.Message({
      ts: '1704067200',
      user_name: 'Test User',
      ImportID: 1,
    })

    render(DayDivider, { current, previous: undefined })

    expect(screen.getByText(/January/i)).toBeInTheDocument()
  })

  it('does not show day pill when same day as previous', () => {
    const current = new models.Message({
      ts: '1704067200',
      user_name: 'Test User',
      ImportID: 1,
    })
    const previous = new models.Message({
      ts: '1704067200',
      user_name: 'Test User',
      ImportID: 1,
    })

    const { container } = render(DayDivider, { current, previous })
    expect(container.querySelector('.day-divider')).toBeNull()
  })

  it('shows day pill when different day from previous', () => {
    const current = new models.Message({
      ts: '1704153600',
      user_name: 'Test User',
      ImportID: 1,
    })
    const previous = new models.Message({
      ts: '1704067200',
      user_name: 'Test User',
      ImportID: 1,
    })

    render(DayDivider, { current, previous })
    expect(screen.getByText(/January/i)).toBeInTheDocument()
  })
})

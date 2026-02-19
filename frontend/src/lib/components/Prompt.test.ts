import { render, screen, fireEvent } from '@testing-library/svelte'
import { describe, it, expect, vi } from 'vitest'
import Prompt from './Prompt.svelte'

describe('Prompt', () => {
  it('renders correctly when open', () => {
    render(Prompt, {
      open: true,
      title: 'Test Title',
      message: 'Test Message',
    })

    expect(screen.getByText('Test Title')).toBeInTheDocument()
    expect(screen.getByText('Test Message')).toBeInTheDocument()
    expect(screen.getByRole('textbox')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'OK' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument()
  })

  it('does not render when closed', () => {
    render(Prompt, {
      open: false,
      title: 'Test Title',
    })

    expect(screen.queryByText('Test Title')).not.toBeInTheDocument()
  })

  it('calls onConfirm with input value when OK is clicked', async () => {
    const onConfirm = vi.fn()
    const { component } = render(Prompt, {
      open: true,
      onConfirm,
    })

    const input = screen.getByRole('textbox')
    await fireEvent.input(input, { target: { value: 'New Value' } })
    
    const okButton = screen.getByRole('button', { name: 'OK' })
    await fireEvent.click(okButton)

    expect(onConfirm).toHaveBeenCalledWith('New Value')
  })

  it('calls onCancel when Cancel is clicked', async () => {
    const onCancel = vi.fn()
    render(Prompt, {
      open: true,
      onCancel,
    })

    const cancelButton = screen.getByRole('button', { name: 'Cancel' })
    await fireEvent.click(cancelButton)

    expect(onCancel).toHaveBeenCalled()
  })
})

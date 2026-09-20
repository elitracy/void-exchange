import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import StatTile from './StatTile'

describe('StatTile', () => {
    it('renders the label and value', () => {
        render(<StatTile label="Tick" value={42} />)
        expect(screen.getByText('Tick')).toBeInTheDocument()
        expect(screen.getByText('42')).toBeInTheDocument()
    })

    it('renders string values as-is', () => {
        render(<StatTile label="Territories Claimed" value="3/10" />)
        expect(screen.getByText('3/10')).toBeInTheDocument()
    })

    it('defaults to the neutral accent when none is given', () => {
        render(<StatTile label="Factions" value={0} />)
        expect(screen.getByText('0')).toHaveClass('text-ink-200')
    })

    it('applies the requested accent class', () => {
        render(<StatTile label="Tick" value={1} accent="violet" />)
        expect(screen.getByText('1')).toHaveClass('text-violet-glow')
    })
})

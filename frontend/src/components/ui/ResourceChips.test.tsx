import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import ResourceChips from './ResourceChips'

describe('ResourceChips', () => {
    it('shows the default empty label when there are no types', () => {
        render(<ResourceChips types={[]} />)
        expect(screen.getByText('None')).toBeInTheDocument()
    })

    it('shows a custom empty label', () => {
        render(<ResourceChips types={[]} emptyLabel="No holdings" />)
        expect(screen.getByText('No holdings')).toBeInTheDocument()
    })

    it('treats a missing types array the same as empty (undefined guard)', () => {
        // @ts-expect-error exercising the runtime guard for callers that skip the type check
        render(<ResourceChips types={undefined} />)
        expect(screen.getByText('None')).toBeInTheDocument()
    })

    it('renders a chip per resource type with its formatted label', () => {
        render(<ResourceChips types={['resource_mineral', 'resource_gas']} />)
        expect(screen.getByText('MINERAL')).toBeInTheDocument()
        expect(screen.getByText('GAS')).toBeInTheDocument()
    })

    it('renders the amount next to the label when amounts are provided', () => {
        render(
            <ResourceChips
                types={['resource_mineral']}
                amounts={{ resource_mineral: 7 }}
            />,
        )
        expect(screen.getByText('7')).toBeInTheDocument()
    })

    it('does not render an amount span when amounts is omitted', () => {
        const { container } = render(<ResourceChips types={['resource_mineral']} />)
        expect(container.querySelector('.font-mono')).toBeNull()
    })
})

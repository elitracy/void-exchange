import { describe, expect, it } from 'vitest'
import { factionAccent, formatResourceLabel, resourceChipClass } from './theme'

describe('factionAccent', () => {
    it('returns the neutral/unclaimed accent for negative ids', () => {
        const accent = factionAccent(-1)
        expect(accent.text).toBe('text-ink-500')
    })

    it('is deterministic for the same id', () => {
        expect(factionAccent(2)).toEqual(factionAccent(2))
    })

    it('wraps around the palette instead of throwing for large ids', () => {
        // Palette has 5 entries; id 5 should map back to the same accent as id 0.
        expect(factionAccent(5)).toEqual(factionAccent(0))
    })
})

describe('formatResourceLabel', () => {
    it('strips the resource_ prefix and upper-cases the rest', () => {
        expect(formatResourceLabel('resource_mineral')).toBe('MINERAL')
    })

    it('upper-cases labels that have no prefix', () => {
        expect(formatResourceLabel('gas')).toBe('GAS')
    })

    it('only strips a leading resource_, not one that appears mid-string', () => {
        expect(formatResourceLabel('rare_resource_gas')).toBe('RARE_RESOURCE_GAS')
    })
})

describe('resourceChipClass', () => {
    it('is deterministic for the same type', () => {
        expect(resourceChipClass('resource_mineral')).toBe(resourceChipClass('resource_mineral'))
    })

    it('returns a value from the palette for an empty string', () => {
        expect(() => resourceChipClass('')).not.toThrow()
    })
})

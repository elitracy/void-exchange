// Deterministic color assignment so the same faction/resource always
// renders with the same accent, without needing a name->color config.

const FACTION_PALETTE = [
    { dot: "bg-violet-glow", text: "text-violet-glow", ring: "ring-violet-core/40", bg: "bg-violet-dim" },
    { dot: "bg-signal-up", text: "text-signal-up", ring: "ring-signal-up/40", bg: "bg-signal-up-dim" },
    { dot: "bg-signal-warn", text: "text-signal-warn", ring: "ring-signal-warn/40", bg: "bg-space-700" },
    { dot: "bg-signal-down", text: "text-signal-down", ring: "ring-signal-down/40", bg: "bg-signal-down-dim" },
    { dot: "bg-sky-400", text: "text-sky-400", ring: "ring-sky-400/40", bg: "bg-space-700" },
]

export function factionAccent(id: number) {
    if (id < 0) {
        return { dot: "bg-ink-600", text: "text-ink-500", ring: "ring-ink-600/40", bg: "bg-space-700" }
    }
    return FACTION_PALETTE[id % FACTION_PALETTE.length]
}

// Strips the "resource_" prefix scenarios use (resource_mineral -> MINERAL)
// for compact HUD-style chip labels.
export function formatResourceLabel(type: string): string {
    return type.replace(/^resource_/, "").toUpperCase()
}

const RESOURCE_PALETTE = [
    "text-violet-glow border-violet-core/40 bg-violet-dim/60",
    "text-signal-up border-signal-up/40 bg-signal-up-dim/60",
    "text-signal-warn border-signal-warn/40 bg-space-700",
    "text-sky-400 border-sky-400/40 bg-space-700",
    "text-signal-down border-signal-down/40 bg-signal-down-dim/60",
]

export function resourceChipClass(type: string): string {
    let hash = 0
    for (let i = 0; i < type.length; i++) {
        hash = (hash * 31 + type.charCodeAt(i)) >>> 0
    }
    return RESOURCE_PALETTE[hash % RESOURCE_PALETTE.length]
}

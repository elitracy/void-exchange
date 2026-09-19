type StatTileProps = {
    label: string
    value: string | number
    accent?: "violet" | "up" | "down" | "neutral"
}

const ACCENTS: Record<NonNullable<StatTileProps["accent"]>, string> = {
    violet: "text-violet-glow",
    up: "text-signal-up",
    down: "text-signal-down",
    neutral: "text-ink-200",
}

function StatTile({ label, value, accent = "neutral" }: StatTileProps) {
    return (
        <div className="flex flex-col gap-1 rounded-xl border border-space-600 bg-space-850/80 px-4 py-3 min-w-[7.5rem]">
            <span className="text-[0.65rem] font-semibold uppercase tracking-widest text-ink-500">
                {label}
            </span>
            <span className={`font-mono text-2xl font-semibold tabular-nums ${ACCENTS[accent]}`}>
                {value}
            </span>
        </div>
    )
}

export default StatTile

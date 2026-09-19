import { formatResourceLabel, resourceChipClass } from "../../lib/theme"

type ResourceChipsProps = {
    types: string[]
    amounts?: Record<string, number>
    emptyLabel?: string
}

function ResourceChips({ types, amounts, emptyLabel = "None" }: ResourceChipsProps) {
    if (!types || types.length === 0) {
        return <span className="text-xs text-ink-500">{emptyLabel}</span>
    }

    return (
        <div className="flex flex-wrap justify-center gap-1.5">
            {types.map((type) => (
                <span
                    key={type}
                    className={`rounded-full border px-2 py-0.5 text-[0.65rem] font-semibold tracking-wide ${resourceChipClass(type)}`}
                >
                    {formatResourceLabel(type)}
                    {amounts && (
                        <span className="ml-1 font-mono tabular-nums opacity-80">{amounts[type]}</span>
                    )}
                </span>
            ))}
        </div>
    )
}

export default ResourceChips

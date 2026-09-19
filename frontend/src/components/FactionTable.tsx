import { tableFeatures, useTable } from '@tanstack/react-table';
import type { ColumnDef } from '@tanstack/react-table';
import { useEffect, useState } from "react";
import type { api } from "../../wailsjs/go/models";
import { useFactionsQuery } from "../hooks/useGameData";
import { factionAccent } from "../lib/theme";
import ResourceChips from "./ui/ResourceChips";

type Faction = {
    id: number
    name: string
    territories: number
    drones: number
    resourceTypes: string[]
    resourceAmounts: Record<string, number>
    totalResources: number
}

const features = tableFeatures({})

function FactionBadge({ id, name }: { id: number, name: string }) {
    const accent = factionAccent(id)
    return (
        <span className={`inline-flex items-center gap-2 rounded-full border border-space-600 px-3 py-1 text-xs font-semibold ${accent.text}`}>
            <span className={`h-1.5 w-1.5 rounded-full ${accent.dot}`} />
            {name}
        </span>
    )
}

const columns: Array<ColumnDef<typeof features, Faction>> = [
    {
        id: 'name',
        header: "Faction",
        cell: ({ row }) => <FactionBadge id={row.original.id} name={row.original.name} />,
    },
    {
        accessorKey: 'territories',
        header: "Territories",
        cell: (info) => <span className="font-mono tabular-nums">{info.getValue() as number}</span>,
    },
    {
        accessorKey: 'drones',
        header: "Fleet",
        cell: (info) => <span className="font-mono tabular-nums">{info.getValue() as number}</span>,
    },
    {
        accessorKey: 'totalResources',
        header: "Stockpile",
        cell: (info) => <span className="font-mono tabular-nums text-signal-up">{info.getValue() as number}</span>,
    },
    {
        id: 'resources',
        header: "Resources",
        cell: ({ row }) => (
            <ResourceChips
                types={row.original.resourceTypes}
                amounts={row.original.resourceAmounts}
                emptyLabel="No holdings"
            />
        ),
    },
]

type TableProps = {
    tick: number
    scenario: string
}

function transformFactions(data: api.FactionView[]): Faction[] {
    return (data ?? []).map(f => {
        const resources = f?.resources ?? {}
        const resourceTypes = Object.keys(resources).sort()
        const totalResources = resourceTypes.reduce((sum, key) => sum + (resources[key] ?? 0), 0)

        return {
            id: f.id,
            name: f.name,
            territories: f?.territories?.length ?? 0,
            drones: f?.drones?.length ?? 0,
            resourceTypes,
            resourceAmounts: resources,
            totalResources,
        }
    })
}

function FactionTable({ tick, scenario }: TableProps) {
    const [factions, setFactions] = useState<Array<Faction>>([])

    const { data, isLoading } = useFactionsQuery(tick, scenario)

    useEffect(() => {
        setFactions(transformFactions(data ?? []))
    }, [scenario, tick, data])

    const table = useTable({
        features,
        columns,
        data: factions,
    })

    if (!scenario) {
        return (
            <div className="flex flex-1 items-center justify-center text-sm text-ink-500">
                Select a scenario to view factions.
            </div>
        )
    }

    if (isLoading && factions.length === 0) {
        return (
            <div className="flex flex-1 items-center justify-center text-sm text-ink-500">
                Loading factions…
            </div>
        )
    }

    return (
        <div className="overflow-auto rounded-xl border border-space-600 bg-space-850/60">
            <table className="w-full min-w-[640px] border-collapse text-sm">
                <thead>
                    {table.getHeaderGroups().map((headerGroup) => (
                        <tr key={headerGroup.id} className="sticky top-0 z-10 bg-space-800/95 backdrop-blur">
                            {headerGroup.headers.map((header) => (
                                <th
                                    key={header.id}
                                    className="border-b border-space-600 px-4 py-3 text-center text-[0.7rem] font-bold uppercase tracking-widest text-ink-500"
                                >
                                    {header.isPlaceholder ? null : (
                                        <table.FlexRender header={header} />
                                    )}
                                </th>
                            ))}
                        </tr>
                    ))}
                </thead>
                <tbody>
                    {table.getRowModel().rows.map((row) => (
                        <tr
                            key={row.id}
                            className="border-b border-space-700/60 transition-colors last:border-b-0 hover:bg-violet-core/5"
                        >
                            {row.getAllCells().map(cell => (
                                <td key={cell.id} className="px-4 py-3 text-center text-ink-200">
                                    <table.FlexRender cell={cell} />
                                </td>
                            ))}
                        </tr>
                    ))}
                    {table.getRowModel().rows.length === 0 && (
                        <tr>
                            <td colSpan={columns.length} className="px-4 py-8 text-center text-sm text-ink-500">
                                No factions yet.
                            </td>
                        </tr>
                    )}
                </tbody>
            </table>
        </div>
    )
}

export default FactionTable

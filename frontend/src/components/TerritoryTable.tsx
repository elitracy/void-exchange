import { tableFeatures, useTable } from '@tanstack/react-table';
import type { ColumnDef } from '@tanstack/react-table';
import { useEffect, useMemo, useState } from "react";
import type { api } from "../../wailsjs/go/models";
import { useFactionsQuery, useTerritoriesQuery } from "../hooks/useGameData";
import { factionAccent } from "../lib/theme";
import ResourceChips from "./ui/ResourceChips";

type Territory = {
    id: number
    ownerId: number
    ownerName: string
    factions: number
    deposits: number
    resources: string[]
}

const features = tableFeatures({})

function OwnerBadge({ ownerId, ownerName }: { ownerId: number, ownerName: string }) {
    const accent = factionAccent(ownerId)
    return (
        <span className={`inline-flex items-center gap-2 rounded-full border border-space-600 px-3 py-1 text-xs font-semibold ${accent.text}`}>
            <span className={`h-1.5 w-1.5 rounded-full ${accent.dot}`} />
            {ownerName}
        </span>
    )
}

const columns: Array<ColumnDef<typeof features, Territory>> = [
    {
        accessorKey: 'id',
        header: "Territory",
        cell: (info) => (
            <span className="font-mono text-ink-400">#{String(info.getValue()).padStart(3, "0")}</span>
        ),
    },
    {
        id: 'owner',
        header: "Owner",
        cell: ({ row }) => <OwnerBadge ownerId={row.original.ownerId} ownerName={row.original.ownerName} />,
    },
    {
        accessorKey: 'factions',
        header: "Factions",
        cell: (info) => <span className="font-mono tabular-nums">{info.getValue() as number}</span>,
    },
    {
        accessorKey: 'deposits',
        header: "Deposits",
        cell: (info) => <span className="font-mono tabular-nums">{info.getValue() as number}</span>,
    },
    {
        accessorKey: 'resources',
        header: "Resources",
        cell: (info) => <ResourceChips types={info.getValue() as string[]} />,
    },
]

type TableProps = {
    tick: number
    scenario: string
}

function transformTerritories(data: api.TerritoryView[], factionsById: Map<number, string>): Territory[] {
    return (data ?? []).map(t => {
        const ownerId = t?.owner ?? -1
        return {
            id: t.id,
            ownerId,
            ownerName: ownerId < 0 ? "Unclaimed" : factionsById.get(ownerId) ?? `Faction #${ownerId}`,
            factions: t?.factions?.length ?? 0,
            deposits: t?.deposits?.length ?? 0,
            resources: t?.deposit_types ?? [],
        }
    })
}

function TerritoryTable({ tick, scenario }: TableProps) {
    const [territories, setTerritories] = useState<Array<Territory>>([])

    const { data, isLoading } = useTerritoriesQuery(tick, scenario)
    const { data: factionsData } = useFactionsQuery(tick, scenario)

    const factionsById = useMemo(
        () => new Map((factionsData ?? []).map(f => [f.id, f.name])),
        [factionsData],
    )

    useEffect(() => {
        setTerritories(transformTerritories(data ?? [], factionsById))
    }, [scenario, tick, data, factionsById])

    const table = useTable({
        features,
        columns,
        data: territories,
    })

    if (!scenario) {
        return (
            <div className="flex flex-1 items-center justify-center text-sm text-ink-500">
                Select a scenario to view territories.
            </div>
        )
    }

    if (isLoading && territories.length === 0) {
        return (
            <div className="flex flex-1 items-center justify-center text-sm text-ink-500">
                Loading territories…
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
                                No territories yet.
                            </td>
                        </tr>
                    )}
                </tbody>
            </table>
        </div>
    )
}

export default TerritoryTable

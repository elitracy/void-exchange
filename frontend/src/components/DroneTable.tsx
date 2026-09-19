import { tableFeatures, useTable } from '@tanstack/react-table';
import type { ColumnDef } from '@tanstack/react-table';
import { useEffect, useMemo, useState } from "react";
import type { api } from "../../wailsjs/go/models";
import { DispatchDrones, RecallDrones } from "../../wailsjs/go/main/App";
import { useDronesQuery, useFactionsQuery, useInvalidateGameData, useTerritoriesQuery } from "../hooks/useGameData";
import {
    DRONE_FIGHTING,
    DRONE_IDLE,
    DRONE_MINING,
    DRONE_TYPE_FIGHTER,
    DRONE_TYPE_MINER,
    activityAccent,
    factionAccent,
    formatDroneLabel,
} from "../lib/theme";

type Drone = {
    id: number
    factionId: number
    factionName: string
    name: string
    level: number
    hp: number
    attack: number
    type: string
    activity: string
    target: number
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

function ActivityBadge({ activity, target }: { activity: string, target: number }) {
    const accent = activityAccent(activity)
    return (
        <span className={`inline-flex items-center gap-2 rounded-full border border-space-600 px-3 py-1 text-xs font-semibold ${accent.text}`}>
            <span className={`h-1.5 w-1.5 rounded-full ${accent.dot}`} />
            {formatDroneLabel(activity)}
            {activity !== DRONE_IDLE && <span className="font-mono text-ink-500">#{String(target).padStart(3, "0")}</span>}
        </span>
    )
}

type DroneActionsProps = {
    drone: Drone
    territoryOptions: { id: number, label: string }[]
    onChanged: () => void
}

function DroneActions({ drone, territoryOptions, onChanged }: DroneActionsProps) {
    const [targetId, setTargetId] = useState<number | "">(territoryOptions[0]?.id ?? "")
    const [error, setError] = useState<string>()
    const [busy, setBusy] = useState(false)

    async function dispatch(activity: string) {
        if (targetId === "") return
        setBusy(true)
        setError(undefined)
        try {
            await DispatchDrones(drone.factionId, Number(targetId), [drone.id], activity)
            onChanged()
        } catch (e) {
            setError(String(e))
        } finally {
            setBusy(false)
        }
    }

    async function recall() {
        setBusy(true)
        setError(undefined)
        try {
            await RecallDrones(drone.factionId, [drone.id])
            onChanged()
        } catch (e) {
            setError(String(e))
        } finally {
            setBusy(false)
        }
    }

    if (drone.activity !== DRONE_IDLE) {
        return (
            <div className="flex flex-col items-center gap-1">
                <button
                    disabled={busy}
                    onClick={recall}
                    className="rounded-md border border-space-600 bg-space-800 px-3 py-1 text-xs font-bold uppercase tracking-wide text-ink-200 transition-colors hover:border-signal-warn hover:text-signal-warn disabled:opacity-40"
                >
                    Recall
                </button>
                {error && <span className="text-[0.65rem] text-signal-down">{error}</span>}
            </div>
        )
    }

    if (drone.type !== DRONE_TYPE_FIGHTER && drone.type !== DRONE_TYPE_MINER) {
        return <span className="text-xs text-ink-500">No actions</span>
    }

    return (
        <div className="flex flex-col items-center gap-1.5">
            <div className="flex items-center gap-1.5">
                <select
                    value={targetId}
                    onChange={(e) => setTargetId(e.target.value === "" ? "" : Number(e.target.value))}
                    className="rounded-md border border-space-600 bg-space-800 px-2 py-1 text-xs text-ink-200 focus:border-violet-core focus:outline-none"
                >
                    {territoryOptions.length === 0 && <option value="">No territories</option>}
                    {territoryOptions.map(t => (
                        <option key={t.id} value={t.id}>{t.label}</option>
                    ))}
                </select>
                {drone.type === DRONE_TYPE_FIGHTER && (
                    <button
                        disabled={busy || targetId === ""}
                        onClick={() => dispatch(DRONE_FIGHTING)}
                        className="rounded-md bg-signal-down/20 border border-signal-down/40 px-3 py-1 text-xs font-bold uppercase tracking-wide text-signal-down transition-opacity hover:bg-signal-down/30 disabled:opacity-40"
                    >
                        Fight
                    </button>
                )}
                {drone.type === DRONE_TYPE_MINER && (
                    <button
                        disabled={busy || targetId === ""}
                        onClick={() => dispatch(DRONE_MINING)}
                        className="rounded-md bg-signal-up/20 border border-signal-up/40 px-3 py-1 text-xs font-bold uppercase tracking-wide text-signal-up transition-opacity hover:bg-signal-up/30 disabled:opacity-40"
                    >
                        Mine
                    </button>
                )}
            </div>
            {error && <span className="text-[0.65rem] text-signal-down">{error}</span>}
        </div>
    )
}

type TableProps = {
    tick: number
    scenario: string
}

function transformDrones(data: api.DroneView[], factionsById: Map<number, string>): Drone[] {
    return (data ?? []).map(d => ({
        id: d.id,
        factionId: d.faction_id,
        factionName: factionsById.get(d.faction_id) ?? `Faction #${d.faction_id}`,
        name: d.name,
        level: d.level,
        hp: d.hp,
        attack: d.attack,
        type: d.type,
        activity: d.activity,
        target: d.target,
    }))
}

function DroneTable({ tick, scenario }: TableProps) {
    const [drones, setDrones] = useState<Array<Drone>>([])

    const { data, isLoading } = useDronesQuery(tick, scenario)
    const { data: factionsData } = useFactionsQuery(tick, scenario)
    const { data: territoriesData } = useTerritoriesQuery(tick, scenario)
    const invalidate = useInvalidateGameData(scenario)

    const factionsById = useMemo(
        () => new Map((factionsData ?? []).map(f => [f.id, f.name])),
        [factionsData],
    )

    const territoryOptions = useMemo(
        () => (territoriesData ?? []).map(t => ({
            id: t.id,
            label: `#${String(t.id).padStart(3, "0")}${t.owner < 0 ? " (unclaimed)" : ""}`,
        })),
        [territoriesData],
    )

    useEffect(() => {
        setDrones(transformDrones(data ?? [], factionsById))
    }, [scenario, tick, data, factionsById])

    const columns: Array<ColumnDef<typeof features, Drone>> = useMemo(() => [
        {
            id: 'faction',
            header: "Faction",
            cell: ({ row }) => <FactionBadge id={row.original.factionId} name={row.original.factionName} />,
        },
        {
            accessorKey: 'name',
            header: "Name",
            cell: (info) => <span className="font-semibold text-ink-200">{info.getValue() as string}</span>,
        },
        {
            accessorKey: 'level',
            header: "Level",
            cell: (info) => <span className="font-mono tabular-nums">{info.getValue() as number}</span>,
        },
        {
            id: 'health',
            header: "Health",
            cell: ({ row }) => (
                <span className={`font-mono tabular-nums ${row.original.hp <= 0 ? "text-signal-down" : "text-ink-200"}`}>
                    {row.original.hp}
                </span>
            ),
        },
        {
            accessorKey: 'type',
            header: "Type",
            cell: (info) => <span className="text-xs uppercase tracking-wide text-ink-400">{formatDroneLabel(info.getValue() as string)}</span>,
        },
        {
            id: 'activity',
            header: "Activity",
            cell: ({ row }) => <ActivityBadge activity={row.original.activity} target={row.original.target} />,
        },
        {
            id: 'actions',
            header: "Actions",
            cell: ({ row }) => (
                <DroneActions drone={row.original} territoryOptions={territoryOptions} onChanged={invalidate} />
            ),
        },
    ], [territoryOptions, invalidate])

    const table = useTable({
        features,
        columns,
        data: drones,
    })

    if (!scenario) {
        return (
            <div className="flex flex-1 items-center justify-center text-sm text-ink-500">
                Select a scenario to view drones.
            </div>
        )
    }

    if (isLoading && drones.length === 0) {
        return (
            <div className="flex flex-1 items-center justify-center text-sm text-ink-500">
                Loading drones…
            </div>
        )
    }

    return (
        <div className="overflow-auto rounded-xl border border-space-600 bg-space-850/60">
            <table className="w-full min-w-[760px] border-collapse text-sm">
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
                                No drones yet.
                            </td>
                        </tr>
                    )}
                </tbody>
            </table>
        </div>
    )
}

export default DroneTable

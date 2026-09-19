const MET_SECONDS_PER_TICK = 1

function pad(n: number): string {
    return n.toString().padStart(2, "0")
}

export function formatMissionClock(tick: number): string {
    const totalSeconds = Math.max(0, tick) * MET_SECONDS_PER_TICK
    const days = Math.floor(totalSeconds / 86400)
    const hours = Math.floor((totalSeconds % 86400) / 3600)
    const minutes = Math.floor((totalSeconds % 3600) / 60)
    const seconds = totalSeconds % 60

    const clock = `${pad(hours)}:${pad(minutes)}:${pad(seconds)}`
    return days > 0 ? `${days}d ${clock}` : clock
}

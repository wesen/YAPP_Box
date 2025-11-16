import time


def _clamp(v: int, lo: int, hi: int) -> int:
    return lo if v < lo else hi if v > hi else v


def parse_mmss(s: str) -> int:
    """Parse MM:SS into total seconds; returns 0 on invalid."""
    try:
        mm, ss = s.split(":")
        return int(mm) * 60 + int(ss)
    except Exception:
        return 0


def format_mmss(total_seconds: int) -> str:
    if total_seconds < 0:
        total_seconds = -total_seconds
    mm = total_seconds // 60
    ss = total_seconds % 60
    return "{:02d}:{:02d}".format(mm, ss)


class TimerEngine:
    """
    Minimal non-blocking multi-stage timer.
    - Stages: list of dicts with {'name': str, 'planned_sec': int}
    - States: READY, RUNNING, PAUSED, COMPLETE
    """
    READY = "READY"
    RUNNING = "RUNNING"
    PAUSED = "PAUSED"
    COMPLETE = "COMPLETE"

    def __init__(self, stages: list[dict[str, object]]) -> None:
        self.stages = stages
        self.stage_index = 0
        self.state = self.READY
        self._planned_ms = 0
        self._remaining_ms = 0
        self._overtime_ms = 0
        self._last_ms = 0
        self._load_current_stage()

    # Internal helpers
    def _load_current_stage(self) -> None:
        st = self.stages[self.stage_index]
        self._planned_ms = int(st["planned_sec"]) * 1000
        self._remaining_ms = self._planned_ms
        self._overtime_ms = 0
        self._last_ms = time.ticks_ms()
        self.state = self.READY

    # Public controls
    def start(self) -> None:
        if self.state in (self.READY, self.PAUSED):
            self._last_ms = time.ticks_ms()
            self.state = self.RUNNING

    def pause(self) -> None:
        if self.state == self.RUNNING:
            self.update()
            self.state = self.PAUSED

    def resume(self) -> None:
        if self.state == self.PAUSED:
            self._last_ms = time.ticks_ms()
            self.state = self.RUNNING

    def next_stage(self) -> bool:
        """
        Advance to next stage. Returns True if advanced, False if complete.
        """
        if self.stage_index < len(self.stages) - 1:
            self.stage_index += 1
            self._load_current_stage()
            self.start()
            return True
        self.state = self.COMPLETE
        return False

    # Update and state derived values
    def update(self) -> None:
        if self.state != self.RUNNING:
            return
        now = time.ticks_ms()
        dt = time.ticks_diff(now, self._last_ms)
        self._last_ms = now
        if self._remaining_ms > 0:
            self._remaining_ms -= dt
            if self._remaining_ms < 0:
                self._overtime_ms += -self._remaining_ms
                self._remaining_ms = 0
        else:
            self._overtime_ms += dt

    # Getters for UI
    def get_stage_name(self) -> str:
        return str(self.stages[self.stage_index]["name"])

    def get_remaining_sec(self) -> int:
        self.update()
        return max(0, self._remaining_ms // 1000)

    def get_remaining_str(self) -> str:
        return format_mmss(self.get_remaining_sec())

    def get_planned_sec(self) -> int:
        return self._planned_ms // 1000

    def get_overtime_sec(self) -> int:
        self.update()
        return self._overtime_ms // 1000

    def get_progress_ratio(self) -> float:
        planned = self._planned_ms
        if planned <= 0:
            return 1.0
        elapsed = planned - self._remaining_ms
        return _clamp(elapsed / planned, 0.0, 1.0)



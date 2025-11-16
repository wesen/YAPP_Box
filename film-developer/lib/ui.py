import time
from .display import Display
from .input import Input
from .temp import Temp
from .timer import TimerEngine, format_mmss
from .presets import load_presets, get_preset_by_id


class UI:
    """
    UI state machine for navigation and core screens on a 16x8 text grid.
    States (subset; more can be added to match full spec):
    - splash → main → menu → (presets/system-info/back)
    - timer (developer) ↔ paused → next stage (stop/fixer/wash) → done
    """
 
    # Simple string constants (MicroPython-friendly)
    S_SPLASH = "splash"
    S_MAIN = "main"
    S_MENU = "menu"
    S_PRESETS = "presets"
    S_SYSTEM = "system_info"
    S_TIMER = "timer"   # running
    S_PAUSED = "paused"
    S_DONE = "done"
 
    # Development stages (timer)
    STAGES = ("developer", "stop_bath", "fixer", "wash")
 
    def __init__(self, display: Display, input_dev: Input, temp: Temp, enable_temp: bool = True) -> None:
        self.display = display
        self.input = input_dev
        self.temp = temp
        self.enable_temp = enable_temp
 
        self.state = self.S_SPLASH
        self.splash_start_ms = time.ticks_ms()
 
        # Menu context
        self.menu_items = ("Timer Presets", "System Info", "Back")
        self.menu_index = 0
 
        # Timer context
        default_stages = [
            {"name": "developer", "planned_sec": 14 * 60 + 30},
            {"name": "stop_bath", "planned_sec": 30},
            {"name": "fixer", "planned_sec": 10 * 60},
            {"name": "wash", "planned_sec": 15 * 60},
        ]
        self.timer = TimerEngine(default_stages)
        # Presets
        self.presets = load_presets()
        # Default to TMax 400 + Xtol if available
        self.current_preset = get_preset_by_id(self.presets, "kodak_tmax_400_xtol_1_1_20c") or (self.presets[0] if self.presets else None)
        self.presets_index = 0
 
    # ----------------------------
    # Rendering helpers
    # ----------------------------
    def _draw_buttons(self, left: str, mid: str, right: str) -> None:
        """
        Button row helper targeting 16 columns total.
        Format: [XXX][YYY][ZZZZ] -> 5 + 5 + 6 = 16 chars.
        Use 3-char left/mid labels and 4-char right label to fit.
        """
        l = (left or "   ")[:3]
        m = (mid or "   ")[:3]
        r = (right or "    ")[:4]
        line = "[{}][{}][{}]".format(l, m, r)
        self.display.text_at(6, 0, line)
 
    def _render_header(self, title: str) -> None:
        self.display.text_at(0, 0, title[:16])
 
    def _render_temp_line(self, row: int) -> None:
        if not self.enable_temp:
            txt = "Temp: --.-C"
        else:
            # Prefer non-blocking getter if available
            try:
                c = self.temp.get_c()
            except Exception:
                c = self.temp.read_c()
            txt = "Temp: --.-C" if c is None else "Temp: {:>4.1f}C".format(c)
        self.display.text_at(row, 0, txt[:16])
 
    # ----------------------------
    # Render screens
    # ----------------------------
    def render(self) -> None:
        d = self.display
        d.clear()
 
        if self.state == self.S_SPLASH:
            self._render_header("PICO TEMP MON")
            d.text_at(2, 0, "Booting...")
            self._draw_buttons("   ", "   ", "    ")
            # Auto-advance after ~1s
            if time.ticks_diff(time.ticks_ms(), self.splash_start_ms) > 1000:
                self.state = self.S_MAIN
 
        elif self.state == self.S_MAIN:
            self._render_header("FILM DEV TIMER")
            # Show current preset (compact label) and planned time
            label = (self.current_preset.get("label") if self.current_preset else "No preset") if self.current_preset else "No preset"
            d.text_at(2, 0, label[:16])
            planned = format_mmss(self.timer.get_planned_sec())
            d.text_at(3, 0, "Timer: {}".format(planned)[:16])
            self._render_temp_line(4)
            self._draw_buttons("STA", "MEN", "    ")  # START / MENU
 
        elif self.state == self.S_MENU:
            self._render_header("MENU")
            for i, name in enumerate(self.menu_items[:3]):
                prefix = "\u25BA " if i == self.menu_index else "  "
                d.text_at(2 + i, 0, (prefix + name)[:16])
            self._draw_buttons("SEL", "NAV", "BACK")
 
        elif self.state == self.S_PRESETS:
            self._render_header("TIMER PRESETS")
            n = len(self.presets)
            if n == 0:
                d.text_at(3, 0, "No presets")
            else:
                # Show up to 3 items around current index
                start = self.presets_index
                # Display three in a looped list for simplicity
                for row_offset in range(3):
                    idx = (start + row_offset) % n
                    name = self.presets[idx].get("label", self.presets[idx].get("id", "preset"))
                    prefix = "\u25BA " if row_offset == 0 else "  "
                    d.text_at(2 + row_offset, 0, (prefix + name)[:16])
            self._draw_buttons("SEL", "NAV", "BACK")

        elif self.state == self.S_SYSTEM:
            self._render_header("SYSTEM INFO")
            d.text_at(2, 0, "FW: v0.0.1")
            d.text_at(3, 0, "Uptime: --:--")
            self._render_temp_line(4)
            self._draw_buttons("   ", "   ", "BACK")
 
        elif self.state in (self.S_TIMER, self.S_PAUSED):
            stage = self.timer.get_stage_name()
            title = {"developer": "DEVELOPING...", "stop_bath": "STOP BATH", "fixer": "FIXER", "wash": "WASH"}.get(stage, "TIMER")
            self._render_header(title)
            # Status lines
            planned = format_mmss(self.timer.get_planned_sec())
            remain = self.timer.get_remaining_str()
            ratio = self.timer.get_progress_ratio()
            bar_len = 10
            filled = int(ratio * bar_len)
            filled = bar_len if filled > bar_len else filled
            bar = ("#" * filled) + (" " * (bar_len - filled))
            d.text_at(2, 0, "TRX+2 D76 22C")
            d.text_at(3, 0, "{} [{}]".format(planned, bar)[:16])
            d.text_at(4, 0, "{} remain".format(remain)[:16])
            # Temperature (non-blocking)
            self._render_temp_line(5)
            if self.state == self.S_TIMER:
                self._draw_buttons("PAU", "   ", "NEXT")
            else:
                self._draw_buttons("RES", "   ", "NEXT")
 
        elif self.state == self.S_DONE:
            self._render_header("DEVELOP DONE!")
            d.text_at(2, 0, "Run: #YYYY-MMDD")
            d.text_at(3, 0, "Stages: 4/4  ")
            self._render_temp_line(4)
            self._draw_buttons("   ", "MEN", "    ")
 
        d.show()
 
    # ----------------------------
    # Event handling
    # ----------------------------
    def _next_stage(self) -> None:
        if self.stage_idx < len(self.STAGES) - 1:
            self.stage_idx += 1
            self.state = self.S_TIMER
        else:
            self.state = self.S_DONE
 
    def handle(self) -> None:
        # Non-blocking temperature polling (if supported)
        if self.enable_temp and hasattr(self.temp, "poll"):
            try:
                self.temp.poll()
            except Exception:
                pass
        for e in self.input.read():
            # BTN1=left, BTN2=middle, BTN3=right
            print("ui: recv event={} in {}".format(e, self.state))
            prev = self.state
            if self.state == self.S_MAIN:
                if e == 1:  # START
                    # Reset and start timer at stage 0
                    self.timer.stage_index = 0
                    self.timer._load_current_stage()
                    self.timer.start()
                    self.state = self.S_TIMER
                elif e == 2:  # MENU
                    self.menu_index = 0
                    self.state = self.S_MENU
 
            elif self.state == self.S_MENU:
                if e == 2:  # NAV
                    self.menu_index = (self.menu_index + 1) % len(self.menu_items)
                elif e == 1:  # SEL
                    sel = self.menu_items[self.menu_index]
                    if sel == "Timer Presets":
                        self.state = self.S_PRESETS
                    elif sel == "System Info":
                        self.state = self.S_SYSTEM
                    else:
                        self.state = self.S_MAIN
                elif e == 3:  # BACK
                    self.state = self.S_MAIN
            elif self.state == self.S_PRESETS:
                if e == 2:  # NAV
                    if self.presets:
                        self.presets_index = (self.presets_index + 1) % len(self.presets)
                elif e == 1:  # SEL
                    if self.presets:
                        # Pick current (top of the list)
                        chosen = self.presets[self.presets_index]
                        self.current_preset = chosen
                        # Apply times to timer
                        dev_time = chosen.get("developer_time", "10:00")
                        try:
                            # Reload current stage with new dev time
                            self.timer.stage_index = 0
                            self.timer.stages[0]["planned_sec"] = (int(dev_time.split(":")[0]) * 60 + int(dev_time.split(":")[1]))
                            self.timer._load_current_stage()
                        except Exception:
                            pass
                        self.state = self.S_MAIN
                elif e == 3:  # BACK
                    self.state = self.S_MENU
 
            elif self.state == self.S_SYSTEM:
                if e == 3:  # BACK
                    self.state = self.S_MENU
 
            elif self.state == self.S_TIMER:
                if e == 1:  # PAUSE
                    self.timer.pause()
                    self.state = self.S_PAUSED
                elif e == 3:  # NEXT
                    advanced = self.timer.next_stage()
                    if not advanced:
                        self.state = self.S_DONE
 
            elif self.state == self.S_PAUSED:
                if e == 1:  # RESUME
                    self.timer.resume()
                    self.state = self.S_TIMER
                elif e == 3:  # NEXT
                    advanced = self.timer.next_stage()
                    if not advanced:
                        self.state = self.S_DONE
 
            elif self.state == self.S_DONE:
                if e == 2:  # MENU
                    self.state = self.S_MENU
            # Serial debug for state transitions
            if prev != self.state:
                print("ui: event={} {} -> {}".format(e, prev, self.state))
 

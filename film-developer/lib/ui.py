import time
from .display import Display
from .input import Input
from .temp import Temp
from .timer import TimerEngine, format_mmss, parse_mmss
from .presets import (
    load_presets,
    get_preset_by_id,
    list_developers,
    list_films_for_developer,
    list_isos_for_film_dev,
    find_best_preset,
)
from .leds import LedSnake


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
    S_PRESETS = "presets"  # legacy single-list (unused after tree)
    S_PRESETS_DEV = "presets_dev"
    S_PRESETS_FILM = "presets_film"
    S_PRESETS_ISO = "presets_iso"
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
        self.presets_dev_index = 0
        self.presets_film_index = 0
        self.presets_iso_index = 0
        self._preset_path = {"dev": None, "film": None, "iso": None}
        # Marquee for long labels (settings/presets)
        self._marquee_offset = 0
        self._marquee_last_ms = time.ticks_ms()
        self._marquee_interval_ms = 200
        # LEDs (optional)
        try:
            self.leds = LedSnake(debug=True)
        except Exception:
            self.leds = None
        self._dbg_last_remain_sec = None
 
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
            d.text_at(2, 0, self._marquee_text(label, 16))
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
 
        elif self.state == self.S_PRESETS_DEV:
            self._render_header("DEVELOPER")
            devs = list_developers(self.presets)
            if not devs:
                d.text_at(3, 0, "No data")
            else:
                visible = 3 if len(devs) > 3 else len(devs)
                for row_offset in range(visible):
                    idx = (self.presets_dev_index + row_offset) % len(devs)
                    name = devs[idx]
                    prefix = "\u25BA " if row_offset == 0 else "  "
                    shown = self._marquee_text(name, 16 - len(prefix)) if row_offset == 0 else name[:16 - len(prefix)]
                    d.text_at(2 + row_offset, 0, (prefix + shown)[:16])
            self._draw_buttons("SEL", "NAV", "BACK")

        elif self.state == self.S_PRESETS_FILM:
            self._render_header("FILM")
            dev = self._preset_path.get("dev")
            films = list_films_for_developer(self.presets, dev) if dev else []
            if not films:
                d.text_at(3, 0, "No films")
            else:
                visible = 3 if len(films) > 3 else len(films)
                for row_offset in range(visible):
                    idx = (self.presets_film_index + row_offset) % len(films)
                    name = films[idx]
                    prefix = "\u25BA " if row_offset == 0 else "  "
                    shown = self._marquee_text(name, 16 - len(prefix)) if row_offset == 0 else name[:16 - len(prefix)]
                    d.text_at(2 + row_offset, 0, (prefix + shown)[:16])
            self._draw_buttons("SEL", "NAV", "BACK")

        elif self.state == self.S_PRESETS_ISO:
            self._render_header("TARGET ISO")
            dev = self._preset_path.get("dev")
            film = self._preset_path.get("film")
            isos = list_isos_for_film_dev(self.presets, dev, film) if (dev and film) else []
            if not isos:
                d.text_at(3, 0, "No ISOs")
            else:
                visible = 3 if len(isos) > 3 else len(isos)
                for row_offset in range(visible):
                    idx = (self.presets_iso_index + row_offset) % len(isos)
                    name = "ISO {}".format(isos[idx])
                    prefix = "\u25BA " if row_offset == 0 else "  "
                    d.text_at(2 + row_offset, 0, (prefix + name)[:16])
            self._draw_buttons("SEL", "NAV", "BACK")
        elif self.state == self.S_PRESETS:
            self._render_header("TIMER PRESETS")
            n = len(self.presets)
            if n == 0:
                d.text_at(3, 0, "No presets")
            else:
                # Show up to 3 items around current index
                start = self.presets_index
                visible = 3 if n > 3 else n
                # Display up to three without duplicating when n < 3
                for row_offset in range(visible):
                    idx = (start + row_offset) % n
                    name = self.presets[idx].get("label", self.presets[idx].get("id", "preset"))
                    prefix = "\u25BA " if row_offset == 0 else "  "
                    # Apply marquee only to selected/top item if it overflows
                    if row_offset == 0:
                        avail = 16 - len(prefix)
                        shown = self._marquee_text(name, avail)
                        d.text_at(2 + row_offset, 0, (prefix + shown)[:16])
                    else:
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
            rem_sec = self.timer.get_remaining_sec()
            remain = format_mmss(rem_sec)
            ratio = self.timer.get_progress_ratio()
            bar_len = 10
            filled = int(ratio * bar_len)
            filled = bar_len if filled > bar_len else filled
            bar = ("#" * filled) + (" " * (bar_len - filled))
            d.text_at(2, 0, "TRX+2 D76 22C")
            d.text_at(3, 0, "{} [{}]".format(planned, bar)[:16])
            d.text_at(4, 0, "{} remain".format(remain)[:16])
            # Debug: log once per second when remaining changes
            if self._dbg_last_remain_sec != rem_sec:
                self._dbg_last_remain_sec = rem_sec
                try:
                    print("timer: stage={}, remain={}, planned={}, ratio={:.2f}".format(
                        stage, remain, planned, ratio
                    ))
                except Exception:
                    pass
            # Temperature (non-blocking)
            self._render_temp_line(5)
            # LED snake pattern (rates vary with remaining seconds)
            if self.leds is not None:
                try:
                    self.leds.update(self.timer.get_remaining_sec(), active=(self.state == self.S_TIMER))
                except Exception:
                    pass
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
                        # Start at developer selection
                        self.presets_dev_index = 0
                        self._preset_path = {"dev": None, "film": None, "iso": None}
                        self._marquee_offset = 0
                        self._marquee_last_ms = time.ticks_ms()
                        self.state = self.S_PRESETS_DEV
                    elif sel == "System Info":
                        self.state = self.S_SYSTEM
                    else:
                        self.state = self.S_MAIN
                elif e == 3:  # BACK
                    self.state = self.S_MAIN
            elif self.state == self.S_PRESETS_DEV:
                devs = list_developers(self.presets)
                if e == 2 and devs:
                    self.presets_dev_index = (self.presets_dev_index + 1) % len(devs)
                    self._marquee_offset = 0
                    self._marquee_last_ms = time.ticks_ms()
                elif e == 1 and devs:
                    self._preset_path["dev"] = devs[self.presets_dev_index]
                    self.presets_film_index = 0
                    self._marquee_offset = 0
                    self._marquee_last_ms = time.ticks_ms()
                    self.state = self.S_PRESETS_FILM
                elif e == 3:
                    self.state = self.S_MENU
            elif self.state == self.S_PRESETS_FILM:
                dev = self._preset_path.get("dev")
                films = list_films_for_developer(self.presets, dev) if dev else []
                if e == 2 and films:
                    self.presets_film_index = (self.presets_film_index + 1) % len(films)
                    self._marquee_offset = 0
                    self._marquee_last_ms = time.ticks_ms()
                elif e == 1 and films:
                    self._preset_path["film"] = films[self.presets_film_index]
                    self.presets_iso_index = 0
                    self.state = self.S_PRESETS_ISO
                elif e == 3:
                    self.state = self.S_PRESETS_DEV
            elif self.state == self.S_PRESETS_ISO:
                dev = self._preset_path.get("dev")
                film = self._preset_path.get("film")
                isos = list_isos_for_film_dev(self.presets, dev, film) if (dev and film) else []
                if e == 2 and isos:
                    self.presets_iso_index = (self.presets_iso_index + 1) % len(isos)
                elif e == 1 and isos:
                    iso = isos[self.presets_iso_index]
                    chosen = find_best_preset(self.presets, dev, film, iso)
                    if chosen:
                        self.current_preset = chosen
                        # Reset marquee for main label after selection
                        self._marquee_offset = 0
                        self._marquee_last_ms = time.ticks_ms()
                        # Apply developer time
                        dev_time = chosen.get("developer_time", "10:00")
                        try:
                            self.timer.stage_index = 0
                            self.timer.stages[0]["planned_sec"] = parse_mmss(dev_time)
                            # Apply stage overrides if provided
                            st = chosen.get("stages") or {}
                            def _apply(name: str, fallback: str) -> int:
                                return parse_mmss(st.get(name, fallback))
                            self.timer.stages[1]["planned_sec"] = _apply("stop_bath", "00:30")
                            self.timer.stages[2]["planned_sec"] = _apply("fixer", "10:00")
                            self.timer.stages[3]["planned_sec"] = _apply("wash", "15:00")
                            self.timer._load_current_stage()
                        except Exception:
                            pass
                    self.state = self.S_MAIN
                elif e == 3:
                    self.state = self.S_PRESETS_FILM
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
    # ----------------------------
    # Marquee helper
    # ----------------------------
    def _marquee_text(self, text: str, width: int) -> str:
        """
        Scroll text horizontally (wrap-around) if it exceeds width.
        Advances ~every _marquee_interval_ms.
        """
        if width <= 0:
            return ""
        if len(text) <= width:
            return text[:width]
        now = time.ticks_ms()
        if time.ticks_diff(now, self._marquee_last_ms) >= self._marquee_interval_ms:
            self._marquee_last_ms = now
            self._marquee_offset = (self._marquee_offset + 1) % (len(text) + 3)
        gap = "   "
        loop = text + gap + text
        start = self._marquee_offset
        end = start + width
        # Ensure slice within bounds
        if end <= len(loop):
            return loop[start:end]
        # Wrap around manually
        part1 = loop[start:]
        part2 = loop[: end - len(loop)]
        return (part1 + part2)[:width]


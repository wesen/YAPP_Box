---
Title: 'Intern Guide: Bring-up Diary, OLED/SH1107, Buttons/Debounce, uasyncio, DS18B20'
Ticket: FILM-DEV-HWTEST-001
Status: active
Topics:
    - hardware
    - oled
    - sh1107
    - buttons
    - leds
    - temperature
    - pico
DocType: reference
Intent: long-term
Owners:
    - manuel
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2025-11-15T20:03:22.872003946-05:00
---


# Intern Guide: Bring-up Diary, OLED/SH1107, Buttons/Debounce, uasyncio, DS18B20

## Goal

Help you (new to embedded and MicroPython) understand the whole bring-up: wiring, OLED driver, async input handling, debouncing strategies, temperature sensor reads, debugging methods, and the rationale behind decisions. This is a practical diary plus a how-to guide.

## Context

Hardware:
- Raspberry Pi Pico W
- 64x128 SH1107 OLED on SPI0 (SCK=GP18, MOSI=GP19, CS=GP17, DC=GP20, RST=GP21)
- 3x buttons (to GND, with internal pull-ups), used initially BTN1 on GP14
- DS18B20 temperature sensor on GP22 with 4.7k pull-up to 3V3

Software:
- MicroPython v1.26.1
- Custom `lib/sh1107.py` minimal driver using `framebuf`
- `lib/display.py` convenience wrapper (text grid)
- `main.py` used to run quick demos (pattern tests, input tests, async handling)

## Quick Reference

This section is your fast lane when you already know what you want to change and just need the exact place to do it. It provides copy-and-paste snippets and precise code references to the parts of the system you will modify most often.

- Display init and usage:

```21:26:/home/manuel/code/others/YAPP_Box/film-developer/lib/display.py
        # SH1107 pins: DC=GP20, CS=GP17, RST=GP21
        # rotation="cw" maps 128x64 logical buffer to 64x128 physical panel.
        # mirror=False to avoid mirrored characters; set True if hardware wiring mirrors X.
        self.dev = SH1107(width=128, height=64, spi=self.spi, dc=Pin(20), cs=Pin(17), rst=Pin(21), rotation="cw", mirror=False)
        self.clear()
```

- Display reset and SPI re-init:

```150:157:/home/manuel/code/others/YAPP_Box/film-developer/lib/sh1107.py
    def reset(self) -> None:
        """
        Public reset: hardware reset, re-init, clear and show.
        """
        self._reset()
        self._init()
        self.fill(0)
        self.show()
```

```52:74:/home/manuel/code/others/YAPP_Box/film-developer/lib/display.py
    def reset_spi(self, spi_baudrate: int | None = None) -> None:
        """
        Fully reset the SPI bus and reattach it to the SH1107 device.
        Optionally change baudrate at the same time.
        """
        try:
            self.spi.deinit()
        except Exception:
            # Some ports may not implement deinit; ignore.
            pass
        if spi_baudrate is not None:
            self.spi_baudrate = spi_baudrate
        # Recreate SPI0
        self.spi = SPI(
            0,
            baudrate=self.spi_baudrate,
            polarity=0,
            phase=0,
            sck=Pin(18),
            mosi=Pin(19),
        )
        # Reattach to device and keep current framebuffer contents
        self.dev.spi = self.spi
```

- Async button ghosting (first IRQ wins) and display loop:

```30:44:/home/manuel/code/others/YAPP_Box/film-developer/main.py
def _btn_irq(pin):
    # Simple ghosting debounce: accept first edge, suppress further for GHOST_MS
    # Keep ISR tiny; schedule userland update via micropython.schedule
    global _irq_total, _suppress_until_ms
    now = time.ticks_ms()
    if time.ticks_diff(now, _suppress_until_ms) < 0:
        _irq_total += 1
        return
    _irq_total += 1
    # Extend suppression window
    _suppress_until_ms = time.ticks_add(now, GHOST_MS) if hasattr(time, "ticks_add") else (now + GHOST_MS)
    try:
        micropython.schedule(_sched_click, 0)
    except Exception:
        pass
```

```62:81:/home/manuel/code/others/YAPP_Box/film-developer/main.py
def _sched_click(_):
    # Runs in soft-IRQ context (scheduled), safe to touch Python state
    s = _state_ref
    if s is None:
        return
    t_now = time.ticks_ms()
    last_press = s["stats"]["last_press_ms"]
    if last_press is not None:
        gap = time.ticks_diff(t_now, last_press)
        st = s["stats"]
        st["gap_sum"] += gap
        st["gap_cnt"] += 1
        if st["gap_min"] is None or gap < st["gap_min"]:
            st["gap_min"] = gap
        if gap > st["gap_max"]:
            st["gap_max"] = gap
    s["stats"]["last_press_ms"] = t_now
    s["presses"] += 1
    s["stats"]["accepted"] += 1
    s["dirty"] = True
```

```142:151:/home/manuel/code/others/YAPP_Box/film-developer/main.py
async def display_task(display, state):
    while True:
        if state.get("dirty"):
            display.clear()
            display.text_at(0, 0, "ASYNC BUTTON DEMO")
            display.text_at(2, 0, "BTN1 on GP14")
            display.text_at(4, 0, "Presses: {}".format(state["presses"]))
            display.show()
            state["dirty"] = False
        await asyncio.sleep_ms(50)
```

## Usage Examples

These examples show typical flows you will run during bring-up and testing. Start with the async button demo to validate input and rendering, then explore display diagnostics and finally the temperature probe once the basics are stable.

1) Run the async button demo (current main.py)
- Deploy project files to the board.
- Observe serial logs every 2 s with stats about IRQs and accepted presses.
- Tap BTN1 (GP14). You should see the counter increase very responsively.

2) Run the OLED pattern test
- Swap `main.py` to call `lib/test_display.run_pattern_demo(display, input_dev)` (or temporarily copy sample from playbook).
- BTN1 cycles diagnostic screens; look for orientation, line spacing, flicker.

3) Run the temp probe demo
- Use `lib/test_temp.run_temp_demo(display, input_dev, temp)`; BTN1 forces a read; BTN3 resets OLED.
- If temperature shows as "--.-C", check the 4.7k pull-up and GP22 wiring.

## Bring-up Diary and Rationale

This narrative explains what we tried, what failed, and why we chose the current design. Reading it will help you internalize the tradeoffs and avoid dead ends when you iterate.

### 1. OLED SH1107 driver and orientation
We built a minimal SH1107 SPI driver (`lib/sh1107.py`) over `framebuf` and a small wrapper `lib/display.py` offering a 16x8 text grid using the builtin 8x8 font. Many 64x128 SH1107 modules mount the panel in portrait. We keep a logical 128x64 buffer and rotate 90 deg CW on `show()` to target 64x128 physical layout.

Rotation and page write:

```74:111:/home/manuel/code/others/YAPP_Box/film-developer/lib/sh1107.py
    def show(self) -> None:
        if self.rotation == "none":
            # SH1107 expects pages (8px tall)
            for page in range(self.height // 8):
                self._write_cmd(_SET_PAGE | page)
                # Column 0..127; many SH1107 modules use column offset 0
                self._write_cmd(_SET_COL_LOW | 0x00)
                self._write_cmd(_SET_COL_HIGH | 0x00)
                start = page * self.width
                end = start + self.width
                self._write_data(self.buffer[start:end])
            return

        if self.rotation == "cw":
            # Logical buffer is 128x64; physical panel is 64x128 (portrait).
            # Rotate 90° clockwise when sending to the panel.
            phys_width = 64
            phys_pages = 128 // 8  # 16 pages
            for page in range(phys_pages):
                self._write_cmd(_SET_PAGE | page)
                self._write_cmd(_SET_COL_LOW | 0x00)
                self._write_cmd(_SET_COL_HIGH | 0x00)
                line = bytearray(phys_width)
                # For each physical column (0..63), pack 8 vertical pixels from the logical buffer
                for xp in range(phys_width):
                    b = 0
                    yl = (self.height - 1) - xp  # map to logical y
                    for bit in range(8):
                        yp = page * 8 + bit
                        xl = yp  # map to logical x
                        # Read pixel from logical framebuffer
                        if self.fb.pixel(xl, yl):
                            b |= (1 << bit)
                    line[xp] = b
                # Optional horizontal mirror at the panel level can be handled by SEG_REMAP,
                # so do not reverse 'line' here. Rely on _init() mirror setting.
                self._write_data(line)
            return
```

### 2. Display corruption and EMI
Symptom: random garbage when hand near wires or after several updates. Likely EMI/signal integrity with breadboard-length wires. Fixes that helped:
- Lower SPI baudrate during bring-up (we tried 1 MHz -> 50 kHz -> 5 kHz; settle based on stability).
- Add decoupling at OLED VCC: 0.1 uF ceramic close to module plus 4.7–10 uF bulk.
- Keep wires short and route close to a ground return (twist with GND if possible).
- Add series resistors (22–100 ohm) in SCK/MOSI near MCU to tame ringing.
- Ensure DC is not floating. We force DC low when idle after data writes to avoid accidental command/data mix:

```165:171:/home/manuel/code/others/YAPP_Box/film-developer/lib/sh1107.py
    def _write_data(self, buf: bytes) -> None:
        self.cs.value(0)
        self.dc.value(1)
        self.spi.write(buf)
        self.cs.value(1)
        # Drive DC low when idle to avoid floating if CS glitches
        self.dc.value(0)
```

### 3. Reliable display reset
Toggle RST with CS high, DC low, and long enough pulses. We extended the reset low/high to 50 ms and re-init afterwards:

```113:123:/home/manuel/code/others/YAPP_Box/film-developer/lib/sh1107.py
    def _reset(self) -> None:
        # Ensure bus is idle and in command mode during reset
        self.cs.value(1)
        self.dc.value(0)
        time.sleep_ms(1)
        # Longer reset pulse to fully reinitialize the controller
        self.rst.value(0)
        time.sleep_ms(50)
        self.rst.value(1)
        time.sleep_ms(50)
```

If corruption persists, power-cycle the OLED module (hard reset beats any software sequence).

### 4. Display test patterns and static rendering
- We created `lib/test_display.py` with pattern grids, lines, checkerboard, borders. Press BTN1 to cycle; we clear between patterns to avoid ghosting.
- The wrapper `lib/display.py` uses an 8x8 font; that gives 16 columns x 8 rows.
- Only call `display.show()` when something changed to avoid flicker.

### 5. Buttons: from polling to IRQ to async
We started with a simple polling `lib/input.py` (read pin values, debounce in software). Then moved to IRQ-based input and a ring buffer with `ThreadSafeFlag`. To avoid blocking or jitter from display/temp work, we:
- Used `uasyncio` tasks: a button consumer, temp sampler, and display refresher.
- Then simplified further: use a tiny ISR that schedules a callback (`micropython.schedule`) and handle a single press immediately with a ghost window (ignore extra edges for a few milliseconds). This captured fast taps while suppressing bounce with minimal complexity.

### 6. Debounce strategies we tried (and why)
- Confirm-low debounce: on interrupt, wait some ms, confirm pin is still low, then accept; wait for release and release-debounce. Pros: robust. Cons: can miss very fast taps; more complexity.
- Time-gated press: ignore any edges for N ms after a press. Pros: simple. Cons: can swallow intended double clicks if window is too large.
- Accept-first-edge + wait-release: accept on first falling edge, then wait until pin goes high and is stable for some ms, flushing bounce edges. Pros: good responsiveness. Cons: still some complexity managing queues.
- Final approach (current): accept on first falling edge, schedule in soft-IRQ, and ghost subsequent edges for GHOST_MS (default 25 ms). Pros: simplest and responsive; easy to tune; no need to track release explicitly for counting taps.

Tuning tips:
- If double-clicks should register, lower GHOST_MS to 10–15 ms.
- If bounce still causes double counts, raise GHOST_MS to 30–40 ms and consider a small cap (10–47 nF) to GND at the button pin and/or a stronger pull-up (10 kOhm external).

### 7. uasyncio usage and patterns
- uasyncio is a cooperative scheduler; tasks must yield (e.g., `await asyncio.sleep_ms(n)`) so others can run.
- We used:
  - A display task that only draws when `state["dirty"]` is set, and throttles to ~20 FPS (50 ms).
  - A stats task that prints debugging metrics every 2 s (helps see if IRQs are being dropped or debounced too aggressively).
  - For buttons, we migrated away from an async consumer to a scheduled callback because it is simpler and avoids managing edge queues.

### 8. DS18B20 temperature reads (non-blocking)
- The one-wire DS18B20 conversion can take ~750 ms at 12-bit. Do not block your loop while waiting.
- Non-blocking pattern: kick `convert_temp()`, after 750 ms `read_temp()`, repeat. In uasyncio this becomes two awaits separated by `await asyncio.sleep_ms(750)`.
- `lib/test_temp.py` demonstrates a periodic sampler and a forced read on BTN1. If reads hang or return None, check wiring (DATA to GP22, 4.7k to 3V3).

### 9. Troubleshooting checklist
- Display mirrored/rotated:
  - Flip `mirror` in `lib/display.py` dev init.
  - If still odd offsets, tweak `_SET_DISP_OFFSET` and `start line` in `_init()`.
- Display corruption:
  - Lower SPI speed, use decoupling, shorten wires, add series resistors, drive DC low when idle, and use the hardened reset. Power-cycle if needed.
- Buttons miss fast taps:
  - Lower GHOST_MS, ensure reliable ground, add hardware RC if needed.
- Temp reads flaky:
  - Confirm 4.7k pull-up, common ground, and stable 3V3 supply.

### 10. What is running now (11/15)
- main.py: single-button async demo with first-IRQ ghosting debounce and a display task.
- Display SPI currently at 1 MHz (adjust in `Display(spi_baudrate=...)` if EMI returns).
- Stats print every 2 s: use these to tune debounce and see if IRQs are being dropped.
- Draws counter on screen (“Draws: N”) increments each frame; fast bursts queue up via `pending_draws` and the display task drains them with a short cadence (~10 ms).

## Debugging Workflow

A consistent debugging routine saves hours. Use this loop to diagnose issues quickly:

- Start at the REPL
  - After flashing files, open the REPL to watch logs. If `main.py` prints nothing, add a one-line print early in `amain()` to confirm it runs.
- Read the stats every 2 seconds
  - The stats task prints IRQ totals, queue depth, accepted presses, and timing:

```153:167:/home/manuel/code/others/YAPP_Box/film-developer/main.py
async def stats_task(state):
    global _irq_total, _irq_dropped, _irq_peak_depth
    while True:
        await asyncio.sleep_ms(2000)
        s = state["stats"]
        gap_avg = (s["gap_sum"] // s["gap_cnt"]) if s["gap_cnt"] else 0
        hold_avg = (s["hold_sum"] // s["hold_cnt"]) if s["hold_cnt"] else 0
        print(
            "stats: irq_total={}, dropped={}, peak_depth={}, accepted={}, ignored={}, "
            "gap_ms(avg/min/max)={}/{}/{} hold_ms(avg/min/max)={}/{}/{}".format(
                _irq_total, _irq_dropped, _irq_peak_depth,
                s["accepted"], s["debounce_ignored"],
                gap_avg, s["gap_min"] if s["gap_min"] is not None else -1, s["gap_max"],
                hold_avg, s["hold_min"] if s["hold_min"] is not None else -1, s["hold_max"],
            )
        )
```

- Adjust one variable at a time
  - Example: tune `GHOST_MS` for button feel, then adjust `spi_baudrate` for OLED stability. Avoid changing multiple variables at once.
- Recover quickly
  - If the OLED corrupts, try BTN3 reset if mapped, or call `display.reset()` from a quick test, or power-cycle the module.
- Add minimal prints
  - When you add prints, include the variable and units (e.g., "gap_ms=...") so you can skim logs fast.
 - Watch the on-screen “Draws:” counter
  - If Draws advances on each tap but less than 1:1 with every press, that’s expected—`show()` is heavier than counting presses. The counter confirms frames are being scheduled and rendered.

### 11. Render scheduling and redraw counter

To make fast bursts feel smoother, the button handler now requests frames and the display task drains a small backlog at a short cadence. This keeps the UI responsive without trying to render on every single micro-event.

- The scheduled click handler increments a small “pending draws” counter:

```62:84:/home/manuel/code/others/YAPP_Box/film-developer/main.py
def _sched_click(_):
    # Runs in soft-IRQ context (scheduled), safe to touch Python state
    s = _state_ref
    if s is None:
        return
    t_now = time.ticks_ms()
    last_press = s["stats"]["last_press_ms"]
    if last_press is not None:
        gap = time.ticks_diff(t_now, last_press)
        st = s["stats"]
        st["gap_sum"] += gap
        st["gap_cnt"] += 1
        if st["gap_min"] is None or gap < st["gap_min"]:
            st["gap_min"] = gap
        if gap > st["gap_max"]:
            st["gap_max"] = gap
    s["stats"]["last_press_ms"] = t_now
    s["presses"] += 1
    s["stats"]["accepted"] += 1
    s["dirty"] = True
    s["pending_draws"] = s.get("pending_draws", 0) + 1
```

- The display task renders while there are pending draws (and shows a clipped redraw counter):

```142:156:/home/manuel/code/others/YAPP_Box/film-developer/main.py
async def display_task(display, state):
    while True:
        if state.get("dirty") or state.get("pending_draws", 0) > 0:
            display.clear()
            display.text_at(0, 0, "ASYNC BUTTON DEMO")
            display.text_at(2, 0, "BTN1 on GP14")
            display.text_at(4, 0, "Presses: {}".format(state["presses"]))
            cnt = (state.get("draw_count", 0) + 1) % 100000
            state["draw_count"] = cnt
            display.text_at(5, 0, "Draws: {}".format(cnt))
            display.show()
            state["dirty"] = False
            if state.get("pending_draws", 0) > 0:
                state["pending_draws"] -= 1
        await asyncio.sleep_ms(10)
```

If you want the UI to feel “snappier,” you can slightly reduce the sleep (e.g., 8–10 ms); if you want to reduce bus load, increase it (e.g., 20–50 ms). The Draws counter gives immediate feedback about how quickly frames are being pushed.

## Exercises

Sharpen skills with focused, hands-on tasks. Each exercise lists goals, hints, and acceptance criteria.

### 1) Tune Ghost Debounce for Double-Clicks
- Goal: Configure `GHOST_MS` so two fast taps count as two presses, but bounce does not double-count.
- Where:

```22:24:/home/manuel/code/others/YAPP_Box/film-developer/main.py
_state_ref = None  # set in amain so scheduled callbacks can update state
GHOST_MS = 25      # ignore further IRQs for this many ms after first edge
```

- Hints:
  - Lower `GHOST_MS` in 5 ms steps (start at 25 → 15 → 10).
  - Watch the stats lines to verify `accepted` increments twice for two quick taps.
- Acceptance:
  - Ten trials of fast double-clicks register ≥ 9/10 as two increments.
  - Single taps never increment more than once.

### 2) Add Long-Press Detection (≥ 600 ms)
- Goal: Count a “long press” separately from a short press.
- Approach options:
  - Simple: track last press time in `_sched_click`; schedule a deferred checker (e.g., 600 ms later) that marks long press if no second click arrived.
  - Alternative: add a rising-edge IRQ path and compute `hold` time between DOWN and UP, then classify.
- Acceptance:
  - Holding BTN1 for ≥ 600 ms increments `long_presses` counter; short taps do not.
- Stretch:
  - Show “LONG” on the display for 300 ms, then restore.

### 3) Implement Double-Click Recognition (Window 300 ms)
- Goal: Distinguish single vs double-click using a time window.
- Hints:
  - Reuse a “pending single” idea: if a second press arrives within 300 ms, count double and cancel the pending single; else, after 300 ms commit as single.
  - Keep the logic in scheduled context (like `_sched_click`) to avoid queue complexity.
- Acceptance:
  - Two quick taps within 300 ms increments `double_clicks`.
  - One tap increments `single_clicks`.
  - No tap produces both counters simultaneously.

### 4) Integrate Non-Blocking DS18B20 Read with uasyncio
- Goal: Sample temperature without blocking button responsiveness.
- Hints:
  - Create a task that does `convert_temp()`, awaits 750 ms, then `read_temp()`, updates a shared value, and marks UI dirty.
  - Keep the display task throttled to 50–100 ms.
- Acceptance:
  - Press responsiveness stays snappy while temperature updates every ~1 s.
  - No REPL stalls; stats still print every 2 s.

### 5) Reduce Display Bus Load
- Goal: Only call `display.show()` when something actually changed.
- Hints:
  - Track a “last_rendered_state” (press count, line(s) of text). Only render when different.
  - Throttle to ≥ 50 ms per frame.
- Acceptance:
  - No flicker; bus activity roughly matches state changes (e.g., burst when clicking, idle otherwise).

### 6) EMI A/B Testing
- Goal: Empirically measure corruption vs SPI speed.
- Steps:
  - Test at `5 kHz`, `50 kHz`, `200 kHz`, `1 MHz`, and record if artifacts appear during pattern demo.
  - Try adding 22–47 Ω in series with SCK/MOSI; note improvements.
  - Vary decoupling (0.1 µF only vs. 0.1 µF + 4.7 µF).
- Acceptance:
  - Document a table of settings vs. observed stability; recommend defaults.

### 7) Visual Feedback: LED Blink on Press
- Goal: Toggle an LED (e.g., GP10) briefly on each press, without blocking.
- Hints:
  - In `_sched_click`, set an LED flag and timestamp; add a tiny async “blinker” task that turns it off after 100 ms.
- Acceptance:
  - LED blips visibly on each accepted press; no missed presses.

### 8) Watchdog OLED Reset
- Goal: Provide a key combo to invoke `display.reset()` when artifacts appear.
- Hints:
  - Map a special sequence (e.g., triple-click within 1 s) to call `display.reset()` and redraw the current UI.
- Acceptance:
  - Trigger resets on demand; panel reliably returns to clean state.

### 9) Stats Overlay
- Goal: Show last `gap_ms` and `hold_ms` on the display for the previous press, for debugging.
- Hints:
  - Update the stats task to store last computed values; display task prints them on a lower row.
- Acceptance:
  - Overlay updates after each press; values make sense vs. tapping cadence.

## Before You Start

This section helps you get a known-good baseline before touching code. It prevents chasing phantom issues caused by wiring or environment.

- Hardware checklist:
  - [ ] Pico W on USB power, common GND across all devices
  - [ ] OLED wired: SCK=GP18, MOSI=GP19, CS=GP17, DC=GP20, RST=GP21, VCC=3V3
  - [ ] BTN1 wired: GP14 -> button -> GND (internal pull-up active)
  - [ ] DS18B20 wired: DATA=GP22, 4.7k to 3V3, common GND
  - [ ] 0.1 uF + 4.7 uF near OLED VCC
- Software checklist:
  - [ ] MicroPython v1.26.x flashed
  - [ ] Copied `film-developer/` to board; files visible from REPL
  - [ ] REPL connected; can see prints

## First Hour Plan (Milestones)

Use this as your roadmap. Check items as you go.

- [ ] Run main.py (async button demo). Verify BTN1 increments on screen.
- [ ] Tune `GHOST_MS` in `main.py` for your button feel (start at 25, try 10-15).
- [ ] Switch to pattern demo in `lib/test_display.py`. Verify orientation and line spacing.
- [ ] If text mirrored, set `mirror=True` in `lib/display.py`.
- [ ] Lower/raise SPI via `Display(spi_baudrate=...)` to balance stability and speed.
- [ ] Try temp demo in `lib/test_temp.py`. Verify temperature shows numeric value.

## Where To Change Things (Fast Map)

Use this section as a jump table. It points you straight to the few lines you will edit most: SPI speed, display orientation, and debounce timing.

- SPI speed:

```161:165:/home/manuel/code/others/YAPP_Box/film-developer/main.py
async def amain():
    # Keep SPI very low for stability during bring-up
    display = Display(spi_baudrate=1_000_000)
```

- OLED mirror/rotation:

```21:26:/home/manuel/code/others/YAPP_Box/film-developer/lib/display.py
        self.dev = SH1107(..., rotation="cw", mirror=False)
```

- Button ghost window:

```22:24:/home/manuel/code/others/YAPP_Box/film-developer/main.py
_state_ref = None  # set in amain so scheduled callbacks can update state
GHOST_MS = 25      # ignore further IRQs for this many ms after first edge
```

## Common Errors and Quick Fixes

When something goes wrong, start here. Each entry lists the likely cause and the fastest remedies that worked for us.

- ImportError: no module named 'typing'
  - Cause: MicroPython firmware lacks CPython typing module.
  - Fix: Remove typing imports/annotations (already done in `lib/temp.py`).
- OLED garbage when hand near wires
  - Cause: EMI. Fix with shorter wires, decoupling, lower SPI, series resistors, hardened reset.
- OLED blank after reset
  - Try `display.reset()`. If still blank, power-cycle the module.
- BTN feels laggy or misses taps
  - Lower `GHOST_MS`. Consider 10-15 ms. Add small cap (10-47 nF) to GND at the pin if needed.
- DS18B20 shows "--.-C"
  - Check the 4.7k pull-up to 3V3 and wiring to GP22. Ensure common ground.

## FAQ

Short answers to questions you will likely have on day 1.

- Q: Why rotate in software instead of configuring the controller?
  - A: The SH1107 modules vary; a software rotation keeps our framebuffer logic simple and consistent across boards.
- Q: Why keep DC low when idle?
  - A: Avoids the controller accidentally sampling data as commands if CS/DC glitches. It reduced corruption in tests.
- Q: Why not handle release in IRQ too?
  - A: We simplified to "first IRQ wins" with a ghost window for bounce. It is minimal, responsive, and easy to tune.
- Q: Should I use uasyncio or a loop?
  - A: uasyncio helps keep tasks cooperative. For buttons, a scheduled callback is even simpler and avoids queue management.

## Glossary

A quick decoder ring for terms you will see in this codebase.

- DC (D/C): Data/Command select pin on OLED controllers.
- CS: Chip Select (active low) for SPI.
- SPI: Serial Peripheral Interface bus (we use SCK/MOSI/CS; no MISO).
- IRQ: Hardware interrupt (we use falling-edge on the button).
- micropyton.schedule: Schedules a function to run soon outside hard IRQ.
- Ghosting debounce: Accept first edge, ignore further edges for a small time window.
- uasyncio: MicroPython’s cooperative task scheduler.

## Related

- Code: `lib/sh1107.py`, `lib/display.py`, `lib/test_display.py`, `lib/test_input.py`, `lib/test_temp.py`, `main.py`
- Playbook: "Initialize and Use SH1107 Display on Pico W" (same ticket)
- MicroPython docs: uasyncio, interrupts, `micropython.schedule`, `time.ticks_ms/diff`

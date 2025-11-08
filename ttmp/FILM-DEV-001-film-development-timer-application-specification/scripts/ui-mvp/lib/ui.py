import time
from .display import Display
from .input import Input
from .temp import Temp


class UI:
    """
    Minimal UI controller for Splash → Main → Menu → Timer stub flows.
    """

    def __init__(self, display: Display, input_dev: Input, temp: Temp) -> None:
        self.display = display
        self.input = input_dev
        self.temp = temp
        self.screen = "splash"
        self.splash_start_ms = time.ticks_ms()
        self.menu_index = 0

    def render(self) -> None:
        d = self.display
        d.clear()
        if self.screen == "splash":
            d.text_at(0, 2, "PICO TEMP MON")
            d.text_at(2, 2, "Booting...")
            d.text_at(6, 0, "[    ][    ][    ]")
            # Auto-advance after ~1s
            if time.ticks_diff(time.ticks_ms(), self.splash_start_ms) > 1000:
                self.screen = "main"
        elif self.screen == "main":
            c = self.temp.read_c()
            temp_line = "Temp: --.-C" if c is None else f"Temp: {c:.1f}C"
            d.text_at(0, 1, "FILM DEV TIMER")
            d.text_at(2, 0, "TRX+2 D76 22C")
            d.text_at(3, 0, "Timer: 14:30")
            d.text_at(4, 0, temp_line)
            d.text_at(6, 0, "[START][MENU][    ]")
        elif self.screen == "menu":
            d.text_at(0, 4, "MENU")
            items = ["Timer Presets", "System Info", "Back"]
            for i, name in enumerate(items):
                prefix = "\u25BA " if i == self.menu_index else "  "
                d.text_at(2 + i, 0, prefix + name[:19])
            d.text_at(6, 0, "[ SEL][ NAV][BACK]")
        elif self.screen == "timer":
            d.text_at(0, 3, "DEVELOPING...")
            d.text_at(2, 0, "TRX+2 D76 22C")
            d.text_at(3, 0, "14:30 [#####     ]")
            d.text_at(4, 0, "12:30 remaining")
            d.text_at(6, 0, "[PAUS][    ][STOP]")

        d.show()

    def handle(self) -> None:
        for e in self.input.read():
            if self.screen == "main":
                if e == 1:  # BTN1: START
                    self.screen = "timer"
                elif e == 2:  # BTN2: MENU
                    self.menu_index = 0
                    self.screen = "menu"
            elif self.screen == "menu":
                if e == 2:  # BTN2: NAV
                    self.menu_index = (self.menu_index + 1) % 3
                elif e == 1:  # BTN1: SEL
                    if self.menu_index == 0:
                        # Timer Presets (stub for MVP)
                        self.screen = "main"
                    elif self.menu_index == 1:
                        # System Info (stub)
                        self.screen = "main"
                    else:
                        self.screen = "main"
                elif e == 3:  # BTN3: BACK
                    self.screen = "main"
            elif self.screen == "timer":
                if e == 3:  # BTN3: STOP
                    self.screen = "main"



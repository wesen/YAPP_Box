import time
from .display import Display
from .input import Input


class MinimalUI:
    """
    Minimal, non-blocking UI for input and display testing.
    - No temperature reads
    - No timer logic
    - Redraws only on change
    """

    def __init__(self, display: Display, input_dev: Input) -> None:
        self.display = display
        self.input = input_dev
        self.counts = [0, 0, 0]
        self.last_label = "-"
        self.dirty = True

    def render(self) -> None:
        if not self.dirty:
            return
        d = self.display
        d.clear()
        d.text_at(0, 0, "UI MINIMAL TEST")
        d.text_at(2, 0, "BTN1: {:>5}".format(self.counts[0]))
        d.text_at(3, 0, "BTN2: {:>5}".format(self.counts[1]))
        d.text_at(4, 0, "BTN3: {:>5}".format(self.counts[2]))
        d.text_at(6, 0, "Last: {}".format(self.last_label)[:16])
        d.show()
        self.dirty = False

    def handle(self) -> None:
        updated = False
        for e in self.input.read():
            idx = e - 1
            if 0 <= idx < 3:
                self.counts[idx] += 1
                self.last_label = "BTN{}".format(e)
                print("minui: event btn{}".format(e))
                updated = True
        if updated:
            self.dirty = True



import micropython
from .async_buttons import GhostDebouncedButton


class AsyncInput:
    """
    IRQ-backed button input using GhostDebouncedButton for BTN1/2/3.
    Captures short presses reliably and exposes a read() API like polling Input.
    """

    def __init__(self, ghost_ms: int = 25, pins=(14, 15, 16), debug: bool = False) -> None:
        self.debug = debug
        self._q = [0] * 32
        self._mask = 31
        self._head = 0
        self._tail = 0
        self._dropped = 0
        # Ensure IRQ exceptions are safe
        try:
            micropython.alloc_emergency_exception_buf(100)
        except Exception:
            pass
        def on_press(button_id: int) -> None:
            nxt = (self._head + 1) & self._mask
            if nxt == self._tail:
                # overflow
                self._dropped += 1
                if self.debug:
                    print("async-input: drop btn{}".format(button_id))
                return
            self._q[self._head] = button_id
            self._head = nxt
            if self.debug:
                print("async-input: enq btn{}".format(button_id))
        # Create three buttons
        self._btns = [
            GhostDebouncedButton(pins[0], ghost_ms=ghost_ms, on_press=lambda _=0: on_press(1), button_id=1),
            GhostDebouncedButton(pins[1], ghost_ms=ghost_ms, on_press=lambda _=0: on_press(2), button_id=2),
            GhostDebouncedButton(pins[2], ghost_ms=ghost_ms, on_press=lambda _=0: on_press(3), button_id=3),
        ]

    def read(self) -> list[int]:
        events = []
        while self._tail != self._head:
            self._tail = (self._tail + 1) & self._mask
            events.append(self._q[(self._tail - 1) & self._mask])
        return events



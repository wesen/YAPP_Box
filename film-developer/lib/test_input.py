import time
from machine import Pin


def run_input_demo(display, input_dev):
    """
    Input-only diagnostics: count BTN1/2/3 press events and show them.
    BTN3 also resets the OLED. LEDs on GP10/11/12 toggle on press.
    """
    # LEDs on GP10/11/12
    LED_PINS = (10, 11, 12)
    leds = [Pin(p, Pin.OUT, value=0) for p in LED_PINS]
    led_state = [0, 0, 0]

    counts = [0, 0, 0]
    last_label = "-"

    def draw():
        display.clear()
        display.text_at(0, 0, "BUTTON TEST")
        display.text_at(2, 0, "BTN1: {:>5}".format(counts[0]))
        display.text_at(3, 0, "BTN2: {:>5}".format(counts[1]))
        display.text_at(4, 0, "BTN3: {:>5}".format(counts[2]))
        display.text_at(6, 0, "Last: {}".format(last_label))
        display.show()

    print("Starting input-only test")
    draw()

    while True:
        updated = False
        for e in input_dev.read():
            print("Button event:", e)
            idx = e - 1
            if 0 <= idx < 3:
                counts[idx] += 1
                led_state[idx] ^= 1
                leds[idx].value(led_state[idx])
                last_label = "BTN{}".format(e)
                updated = True
            if e == 3:
                print("Display reset requested")
                display.reset()
                updated = True
        if updated:
            draw()
        time.sleep_ms(20)



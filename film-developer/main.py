from lib.display import Display
from lib.input import Input
from lib.temp import Temp
import time
from machine import Pin


def format_temp(c):
    if c is None:
        return "--.-C"
    return "{:.1f}C".format(c)


def main() -> None:
    # Temperature probe test + button→LED mapping
    display = Display()
    input_dev = Input()
    temp = Temp()

    # LEDs on GP10/11/12 per README pin map
    LED_PINS = (10, 11, 12)
    leds = [Pin(p, Pin.OUT, value=0) for p in LED_PINS]
    led_state = [0, 0, 0]

    # Initial screen
    print("Starting DS18B20 temperature probe test")
    display.clear()
    display.text_at(0, 0, "TEMP PROBE TEST")
    display.text_at(2, 0, "BTN1: force read")
    display.text_at(3, 0, "BTN2/3: LEDs")
    display.show()
    last_update_ms = 0
    last_temp = None

    while True:
        # Toggle corresponding LED on button press (active-low edge)
        for e in input_dev.read():
            print("Button event:", e)
            idx = e - 1
            if 0 <= idx < 3:
                led_state[idx] ^= 1
                leds[idx].value(led_state[idx])
            # BTN1: force immediate temperature read
            if e == 1:
                c = temp.read_c()
                print("Forced temp read:", c)
                last_temp = c
                display.clear()
                display.text_at(0, 0, "TEMP PROBE TEST")
                display.text_at(2, 0, "BTN1: force read")
                display.text_at(4, 0, "Temp: {}".format(format_temp(c)))
                display.show()
        # Periodic update every ~1s
        now = time.ticks_ms()
        if time.ticks_diff(now, last_update_ms) > 1000:
            c = temp.read_c()
            last_update_ms = now
            if c != last_temp:
                print("Temp update:", c)
            last_temp = c
            display.clear()
            display.text_at(0, 0, "TEMP PROBE TEST")
            display.text_at(2, 0, "BTN1: force read")
            display.text_at(4, 0, "Temp: {}".format(format_temp(c)))
            display.show()
        time.sleep_ms(20)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        pass



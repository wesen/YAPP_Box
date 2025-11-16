import time
from machine import Pin


def format_temp(c):
    if c is None:
        return "--.-C"
    return "{:.1f}C".format(c)


def run_temp_demo(display, input_dev, temp, spi_baudrate_note="200 kHz"):
    """
    Temperature probe smoke test with 1s updates and BTN1 forced read.
    BTN3 resets the OLED. BTN2/BTN3 also toggle LEDs (GP10/11/12).
    """
    # LEDs on GP10/11/12
    LED_PINS = (10, 11, 12)
    leds = [Pin(p, Pin.OUT, value=0) for p in LED_PINS]
    led_state = [0, 0, 0]

    print("Starting DS18B20 temperature probe test (SPI {})".format(spi_baudrate_note))
    last_update_ms = 0
    last_temp = None

    def draw(c):
        display.clear()
        display.text_at(0, 0, "TEMP PROBE TEST")
        display.text_at(2, 0, "BTN1: force read")
        display.text_at(3, 0, "BTN3: reset OLED")
        display.text_at(4, 0, "Temp: {}".format(format_temp(c)))
        display.show()

    draw(last_temp)

    while True:
        for e in input_dev.read():
            print("Button event:", e)
            idx = e - 1
            if 0 <= idx < 3:
                led_state[idx] ^= 1
                leds[idx].value(led_state[idx])
            if e == 1:
                c = temp.read_c()
                print("Forced temp read:", c)
                last_temp = c
                draw(last_temp)
            if e == 3:
                print("Display reset requested")
                display.reset()
                draw(last_temp)

        now = time.ticks_ms()
        if time.ticks_diff(now, last_update_ms) > 1000:
            c = temp.read_c()
            last_update_ms = now
            if c != last_temp:
                print("Temp update:", c)
            last_temp = c
            draw(last_temp)
        time.sleep_ms(20)



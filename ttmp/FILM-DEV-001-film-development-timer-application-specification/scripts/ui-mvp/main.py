from lib.display import Display
from lib.input import Input
from lib.temp import Temp
from lib.ui import UI
import time


def main() -> None:
    display = Display()
    input_dev = Input()
    temp = Temp()
    ui = UI(display, input_dev, temp)

    while True:
        ui.handle()
        ui.render()
        time.sleep_ms(50)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        pass



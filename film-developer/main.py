from lib.display import Display
from lib.input import Input
from lib.test_input import run_input_demo


def main() -> None:
    # Input-only test to avoid DS18B20 hangs
    display = Display(spi_baudrate=1_000_000)
    input_dev = Input()
    run_input_demo(display, input_dev)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        pass



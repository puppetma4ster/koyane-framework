import pyfiglet

def get_banner_lines() -> list[str]:
    ascii_banner = pyfiglet.figlet_format("KYF", font="slant")
    return [
        ascii_banner.rstrip("\n"),
        "Koyane-Framework :: wordlist forge & analysis toolkit",
        "made by Puppetm4ster"
        "\n\n"
    ]
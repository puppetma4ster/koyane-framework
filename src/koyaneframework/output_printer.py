# src/koyaneframework/banner.py

import pyfiglet

def show_banner():
    ascii_banner = pyfiglet.figlet_format("KYF", font="slant")
    print(ascii_banner)
    print("Koyane-Framework :: wordlist forge & analysis toolkit")
    print ("made by Puppetm4ster")
    print()
    print()

# Koyane-Framework :: wordlist forge & analysis toolkit (Early Stage)

```
 _   __ _____ __   __  ___   _   _  _____         ______ ______   ___  ___  ___ _____  _    _  _____ ______  _   __
| | / /|  _  |\ \ / / / _ \ | \ | ||  ___|        |  ___|| ___ \ / _ \ |  \/  ||  ___|| |  | ||  _  || ___ \| | / /
| |/ / | | | | \ V / / /_\ \|  \| || |__   ______ | |_   | |_/ // /_\ \| .  . || |__  | |  | || | | || |_/ /| |/ /
|    \ | | | |  \ /  |  _  || . ` ||  __| |______||  _|  |    / |  _  || |\/| ||  __| | |/\| || | | ||    / |    \
| |\  \\ \_/ /  | |  | | | || |\  || |___         | |    | |\ \ | | | || |  | || |___ \  /\  /\ \_/ /| |\ \ | |\  \
\_| \_/ \___/   \_/  \_| |_/\_| \_/\____/         \_|    \_| \_|\_| |_/\_|  |_/\____/  \/  \/  \___/ \_| \_|\_| \_/   
 Koyane-Framework :: wordlist forge & analysis toolkit
 made by Puppetm4ster
```                   


**Koyane** is a modular framework for generating, editing, finding, and analyzing wordlists, designed for password cracking and ethical security testing.  
The project is named after *Ame-no-Koyane*, a kami (deity) in Japanese mythology associated with structure, ritual, and the power of words.

---
## Installation
### Makefile (recommended) koyane-framework 0.7.0-beta
The installation has only been tested on Linux operating systems so far. 
Due to the program's path handling, it will probably not be able to run on Windows.
Future versions are also planned to be compatible with Windows.

#### requirements 
- [golang](https://go.dev/) >= 1.25.0
- [make](https://www.gnu.org/software/make/)

Debian, Ubuntu:
```bash
sudo apt update
sudo apt install build-essential golang-go
```
- Other used go libs are automatically downloaded from the source code during compilation.

#### Download
Download the last version -> unzip the directory -> go into program directory:
```bash
wget https://gitlab.com/puppetm4ster/koyane-framework/-/archive/0.7.0-beta/koyane-framework-0.7.0-beta.zip?ref_type=tags
unzip koyane-framework-0.7.0-beta.zip\?ref_type\=tags
cd koyane-framework-0.7.0-beta
```
#### Build & Installation
you can use the **auto_install.sh** for installing
```bash
chmod +x auto_install.sh
sudo ./auto_install.sh
```
or install manually with make
```bash
make build
sudo make install
make clean
```
#### Uninstall

Go to the project directory where the **Makefile** is located and run:
```bash
sudo make uninstall koyane-framework
```
---------------------------------------------

### From pypi **(outdated)** koyane-framework 0.2.0-alpha
##### This version is NOT recommended!
##### The pipy version was a very early alpha release.This version is NOT recommended.
You can install the **deprecated** build directly with pip:

```bash
pip install koyaneframework
```

---------------------------------------------
## Status: Beta

The project is in a early stage of development. Functionality is limited and subject to change.
This is my very first bigger coding project so i am grateful for every improvement suggestion at **puppetma4ster@proton.me**
I try to update the project once a week.
---

## Features

The program serves as a helpful framework for developing, editing, analyzing, and finding word lists for
password hash analysis, directory discovery, web fuzzing, and much more.
The following features are listed per module:


### 1. Generation (kyfgen)
+ support for mask-based generation with multiple wildcards using ?d (digits) and
  fixed character segments using ! (e.g. !A for 'A', !abc123 for custom sets).
+ extraction passwords from hashcat potfiles
+ support for string permutation

### 2. Editing  (kyfedit)
+ wordlist sorting
+ remove words by length
+ remove words who match a mask
+ remove words containing given characters
+ filter all words that are not in european chars

### 3. Analyzing (kyfinfo)
+ analyze basic word list info
+ line count
+ word length statistics
+ complexity metrics
+ char statistics and frequency

### 4. Finding outer word lists (kyfdb)
+ search and filter stats from outer word lists
+ download them from terminal

### General Features
+ CLI interface
+ Basic status messages

**Not yet implemented:**

- Deduplication and merging
- TUI or rich CLI frontend
- Contextual generation based on target information (name, location, job, hobbies, birthday...)
- extract and build rules out of wordlists
- console version

---

## Goals

Koyane aims to become a fast, modular, and scriptable framework for:

- Wordlist generation using rules, masks, templates, and personal context
- Large-scale list manipulation and cleanup
- Format conversion and statistical inspection
- Efficient pipelines for password audit and red team workflows

---

## Example Usage

```bash

kyfgen -M !AS?l?l?l?lL?d?d?! -o ~/MyList.txt

kyfinfo --all ~/MyList.txt
kyfedit -r 8:* --european ~/MyList.txt ~/EditList.txt

kyfdb search --biggest-word 8:14 --smallest-word 8:8 --digits-perc 30.0:*

kyfdb show 22

kyfdb get -p ~/OuterList 22

```
---
## License

This project is licensed under the Apache License 2.0.  
You may use, modify, and distribute it freely, provided you include proper attribution and preserve this license.

See the full license text in the [LICENSE](LICENSE) file.

---
## Disclaimer

Koyane is intended for **educational purposes** and **authorized security testing** only.  
The author strictly opposes any misuse of this tool for illegal, unethical, or oppressive purposes.

**Use responsibly. Break rules, not laws.**

> This software may **not** be used by government surveillance agencies, intelligence services, or law enforcement bodies without explicit and provable ethical clearance.  
> Any use for **mass surveillance**, **targeted repression**, or **covert intelligence operations** is explicitly **forbidden** by the author.

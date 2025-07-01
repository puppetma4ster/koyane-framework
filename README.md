# Koyane-Framework :: wordlist forge & analysis toolkit (Early Stage)

```
     __ ____  ________
    / //_/\ \/ / ____/
   / , <    \  / /_    
  / /| |   / / __/    
 /_/ |_|  /_/_/      
 Koyane-Framework :: wordlist forge & analysis toolkit
 made by Puppetm4ster
```                   


**Koyane** is a modular framework for generating, editing, and analyzing wordlists, designed for password cracking and ethical security testing.  
The project is named after *Ame-no-Koyane*, a kami (deity) in Japanese mythology associated with structure, ritual, and the power of words.

---
## Installation

### 📦 From Local Build (Recommended)

You can install the latest build directly from the `dist/` folder:

```bash
pip install koyaneframework
```

## Status: Pre-Alpha

The project is in a very early stage of development. Functionality is limited and subject to change.
This is my very first coding project so i am grateful for every improvement suggestion at **puppetma4ster@proton.me**
I try to update the project once a week.
---

## Features

+ Basic wordlist analysis (line count, word length stats, complexity metrics)
+ Charset-based generation with configurable min/max length
+ Partial support for mask-based generation
+ Basic wordlist sorting
+ CLI interface powered by Typer
+ Basic status messages

**Not yet implemented:**

- Word mutation, filtering, or combinator logic
- Deduplication and merging
- TUI or rich CLI frontend
- Performance optimizations (parallelization, mmap, low-level backend)
- Contextual generation based on target information (name, location, job, hobbies, birthday...)
- extract and build rules out of wordlists
- Mask generation supports wildcards such as ?d for digits,
  and fixed character segments using ! (e.g. !A for 'A', !abc123 for custom sets).
- Search functionality for wordlists (like `search type:ftp` in Metasploit).
- Building a console version of Koyane with **prompt_toolkik**

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
koyane generate --min 4 --max 6 --char-set abc123 output_wordlist.txt
koyane generate --min 1 --max 8 --char-file charset.txt output_wordlist.txt
koyane generate --mask ?L?v?l?l?l?d?d?p output_wordlist.txt
koyane generate --min 5 --mask ?L?v?l?l?l?d?d?p output_wordlist.txt
koyane analyze wordlist.txt
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

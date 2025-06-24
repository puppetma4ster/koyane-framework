import os
import heapq
from pathlib import Path


LOWER_CASE_CHARACTERS: str = "abcdefghijklmnopqrstuvwxyz"       # ?l
UPPER_CASE_CHARACTERS: str = LOWER_CASE_CHARACTERS.upper()      # ?L
LOWER_CASE_VOWELS: str = "aeiou"        # ?v
UPPER_CASE_VOWELS: str = LOWER_CASE_VOWELS.upper()      # ?V
LOWER_CASE_CONSONANTS: str = "bcdfghjklmnpqrstvwxyz"        # ?c
UPPER_CASE_CONSONANTS: str = LOWER_CASE_CONSONANTS.upper()      # ?C

DIGITS: str = "0123456789"      # ?d
SPECIAL_CHARACTERS_MOST_USED: str = "!@#$%^&*()-_+=?"       # ?f
SPECIAL_CHARACTERS_POINTS: str = ".,:;"     # ?p
SPECIAL_CHARACTERS_BRACELET: str = "()[]{}" # ?b

SPECIAL_CHARACTERS: str = "<>|^°!\"§$%&/()=?´{}[]\\¸`+~*#'-_.:,;@€" #?s




temp_dir = Path(__file__).parent / "tmp"
output_file_rem_empty_lines = temp_dir / "removed_empty_lines.kyftmp"

def external_sort(input_file, output_file, chunk_size=1_000_000):


    temp_files = []
    temp_dir = Path(__file__).parent / "tmp" / "chunks"
    temp_dir.mkdir(parents=True, exist_ok=True)

    # Chunks erzeugen
    with open(input_file, 'r', encoding='utf-8', errors='ignore') as f:
        chunk_index = 0
        while True:
            lines = []
            for _ in range(chunk_size):
                line = f.readline()
                if not line:
                    break
                lines.append(line)  # exakt wie eingelesen

            if not lines:
                break

            lines.sort()
            temp_path = temp_dir / f"chunk_{chunk_index}.kyftmp"
            with open(temp_path, 'w', encoding='utf-8') as tf:
                tf.writelines(lines)  # exakt wie eingelesen schreiben

            temp_files.append(temp_path)
            chunk_index += 1

    # Chunks zusammenführen
    files = [open(path, 'r', encoding='utf-8', errors='ignore') for path in temp_files]
    with open(output_file, 'w', encoding='utf-8') as outf:
        iterators = (f for f in files)  # keine Manipulation
        for line in heapq.merge(*iterators):
            outf.write(line)  # direkt schreiben, keine Veränderung

    for f in files:
        f.close()
    # deletes tmp chunk files
    for path in temp_files:
        path.unlink()







def create_new_wordlist(filepath: Path):
    path = os.path.dirname(filepath) or "."

    if os.path.isdir(path):
        with open(filepath, 'w', encoding="utf-8"):
            pass



def add_new_word_to_wordlist(filepath: Path, word: str):
    with open(filepath, "a", encoding="utf-8") as file:
        file.write(f"{word}\n")


def remove_empty_lines(input_path: Path, output_path: Path):
    with input_path.open('r', encoding='utf-8', errors='ignore') as infile, \
         output_path.open('w', encoding='utf-8') as outfile:
        for line in infile:
            if line.strip() != '':
                outfile.write(line)
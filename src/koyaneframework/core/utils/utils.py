import os
import heapq
from pathlib import Path
import tempfile
import random
import  tarfile

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

# Global variables to store paths after calling prepare_temp_dirs
BASE_TEMP_DIR: Path = None
CHUNK_TEMP_DIR: Path = None

LIST_SUFFIX = ".klst"
TEMP_SUFFIX = ".kyftmp"


def external_sort(input_file: Path, output_file: Path, chunk_size=1_000_000):
    """
    sorts a wordlist file by Unicode.
    Large files are split into chunks and later reassembled into a fully sorted wordlist.
    :param input_file: the word list that is sorted
    :param output_file: where to save the sorted list
    :param chunk_size: how big a single chunks should be
    :return: -
    """
    temp_files = []

    # 1. Read and split into sorted chunks
    with open(input_file, 'r', encoding='utf-8', errors='ignore') as f:
        chunk_index = 0
        while True:
            lines = []
            for _ in range(chunk_size):
                line = f.readline()
                if not line:
                    break
                lines.append(line)

            if not lines:
                break

            lines.sort()
            chunk_path = CHUNK_TEMP_DIR / f"chunk_{chunk_index}{TEMP_SUFFIX}"
            with open(chunk_path, 'w', encoding='utf-8') as tf:
                tf.writelines(lines)

            temp_files.append(chunk_path)
            chunk_index += 1

    # 2. Merge sorted chunks
    files = [open(path, 'r', encoding='utf-8', errors='ignore') for path in temp_files]
    with open(output_file, 'w', encoding='utf-8') as outf:
        iterators = (f for f in files)
        for line in heapq.merge(*iterators):
            outf.write(line)

    # 3. Cleanup
    for f in files:
        f.close()
    for path in temp_files:
        path.unlink()



def create_new_wordlist(filepath: Path, temp_wordlist:bool = False) -> Path:
    """
    creates a new empty list file
    :param filepath: where list is saved
    :return: -
    """
    if temp_wordlist:
        filepath = filepath.with_suffix(TEMP_SUFFIX).resolve()
    else:
        filepath = filepath.with_suffix(LIST_SUFFIX).resolve()

    directory = filepath.parent

    if not directory.is_dir():
        raise FileNotFoundError(f"Directory does not exist: {directory}")

    # Create empty wordlist
    with open(filepath, 'w', encoding='utf-8'):
        pass

    return  filepath


def add_new_word_to_wordlist(filepath: Path, word: str):
    """
    adds a sting to an existing list
    :param filepath: where list is saved
    :param word: which word should be added to the list
    :return: -
    """
    with open(filepath, "a", encoding="utf-8") as file:
        file.write(f"{word}\n")


def remove_empty_lines(input_path: Path, output_path: Path):
    with input_path.open('r', encoding='utf-8', errors='ignore') as infile, \
         output_path.open('w', encoding='utf-8') as outfile:
        for line in infile:
            if line.strip() != '':
                outfile.write(line)

def prepare_temp_dirs():
    """
    Creates and prepares necessary temporary directories under the system's temp folder.
    Sets global path variables for later use.
    """
    global BASE_TEMP_DIR, CHUNK_TEMP_DIR

    # Base directory for temporary files related to koyane-framework
    BASE_TEMP_DIR = Path(tempfile.gettempdir()) / "koyane-framework"
    BASE_TEMP_DIR.mkdir(parents=True, exist_ok=True)

    # Subdirectory specifically for chunk files
    CHUNK_TEMP_DIR = BASE_TEMP_DIR / "chunks"
    CHUNK_TEMP_DIR.mkdir(parents=True, exist_ok=True)

def compress_to_tarxz(source: Path, delete_file: bool = True):
    """
    Compresses a list into a tar.xz archive.
    The archive has the name of the file and is stored at the path where the list was saved
    :param source: the path to the file and the storage path of the archive
    :param delete_file: if file is deleted afterward compression
    :raises:
    """
    source = Path(source).with_suffix(LIST_SUFFIX).resolve()
    if not source.is_file():
        raise ValueError()

    # Archive path: same directory, filename.tar.xz
    archive_path = source.parent / (source.name + ".tar.xz")

    with tarfile.open(archive_path, 'w:xz') as tar:
        tar.add(source, arcname=source.name)  # Only include the file name in the archive

    if delete_file:
        source.unlink()

def get_base_temp_dir():
    """
    Method only for returning the temp path    :return:
    """
    if BASE_TEMP_DIR is None:
        raise RuntimeError("BASE_TEMP_DIR not initialized")
    return BASE_TEMP_DIR

def random_temp_number():
    random_num = random.randint(1_000_000, 9_999_999)
    print(random_num)

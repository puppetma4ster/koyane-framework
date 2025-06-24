import os
import sys

def load_words_from_file(path):
    """
    reads words from file line by line
    """
    if not os.path.isfile(path):
        print(f" no existing file to read: {path}", file=sys.stderr)
        sys.exit(1)

    with open(path, 'r', encoding='utf-8') as f:
        return [line.strip() for line in f if line.strip()]

def load_chars_from_input(raw_input: str) -> list[str]:
    """
    Splits a string into individual non-whitespace characters.
    """
    return [char for char in raw_input if not char.isspace()]


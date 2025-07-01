from koyaneframework.utils import external_sort
from pathlib import Path


def sort_wordlist(input_file: Path, output_file: Path):
    external_sort(input_file, output_file)

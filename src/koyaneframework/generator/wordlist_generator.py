import itertools
from pathlib import  Path
from koyaneframework.utils import add_new_word_to_wordlist, create_new_wordlist
from koyaneframework.generator.mask_interpreter import MaskInterpreter
from koyaneframework.load_animation import LoadingSpinner


def generate_wordlist(words, min_len, max_len, outfile: Path):
    load = LoadingSpinner(text="Generating wordlist")
    load.start()
    create_new_wordlist(outfile)
    for i in range(min_len, max_len + 1):
        for combination in itertools.permutations(words, i):
            combined = ''.join(combination)
            add_new_word_to_wordlist(outfile, combined)
    load.stop()

def generate_mask_wordlist(mask_arg: str, outfile: Path, max_len: int = None,min_len: int=1):
    load = LoadingSpinner(text="Generating mask wordlist")
    load.start()

    mask = MaskInterpreter(mask_arg)
    if not max_len:
        max_len = len(mask.mask_segments)

    segments = [segment.permitted_characters for segment in mask.mask_segments]

    create_new_wordlist(outfile)
    for i in range(min_len, max_len + 1):
        for combination in itertools.product(*segments):
            word = ''.join(combination)
            add_new_word_to_wordlist(outfile, word)

    load.stop()


# work in progress
def calculate_mask_storage(mask: MaskInterpreter, min_len: int, max_len: int):
    total_combinations = 0
    for length in range(min_len, max_len + 1):
        combinations = 1
        for seg in mask.mask_segments[:length]:
            combinations *= len(seg.permitted_characters)
        total_combinations += combinations
    return total_combinations

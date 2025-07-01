import itertools
from pathlib import  Path
from koyaneframework.core.utils.utils import add_new_word_to_wordlist, create_new_wordlist
from koyaneframework.core.generator.mask_interpreter import MaskInterpreter
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

def generate_mask_wordlist(mask_arg: str, outfile: Path, min_len: int = None):
    """
    Generates a wordlist based on a given mask.

    If no minimum length is specified, it generates only words of full mask length.
    If a minimum length is provided, it generates all combinations from that length up to the full mask length.

    Args:
        mask_arg (str): The input mask string (e.g. "?d?d?l").
        outfile (Path): The file path where the generated wordlist will be saved.
        min_len (int, optional): The minimum word length to generate. Defaults to full mask length.
    """
    load = LoadingSpinner(text="Generating mask wordlist")
    load.start()

    # Parse the mask into segments of permitted characters
    mask = MaskInterpreter(mask_arg)
    full_segments = [segment.permitted_characters for segment in mask.mask_segments]
    mask_length = len(full_segments)

    # Determine the lengths to generate:
    # If min_len is not provided, only the full mask length is used.
    # If min_len is given, generate all lengths from min_len up to full length.
    if min_len is None:
        lengths = [mask_length]
    else:
        lengths = list(range(min_len, mask_length + 1))

    # Create or overwrite the target wordlist file
    create_new_wordlist(outfile)

    # Generate all possible combinations for each length
    for i in lengths:
        segments = full_segments[:i]  # use only the first i mask segments
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

HELP_TEXTS = {
    "before": {
        "quiet_mode" : "Suppress startup banners and informational output."
    },
    "generate": {
        "min_length": "Specifies the minimum word length.",
        "max_length": "Specifies the maximum word length.",
        "mask_help": (
            "Generate wordlist from a pattern mask. "
            "A mask consists of segments starting with '?' followed by letters that define the character type. "
            "Example: '?ld?d?f' generates a word with a lowercase letter, a digit, and a special character.\n"
            "Available wildcards:\n"
            "  l = lowercase letter\n"
            "  L = uppercase letter\n"
            "  v = lowercase vowel\n"
            "  V = uppercase vowel\n"
            "  c = lowercase consonant\n"
            "  C = uppercase consonant\n"
            "  d = digit\n"
            "  s = any special character\n"
            "  f = common special characters\n"
            "  p = dot special characters\n"
            "  b = bracket special characters"
        ),
        "char_set": "Generate a wordlist from a custom set of characters.",
        "word_file": (
            "Generate a wordlist from entries in a file. Useful for combinations or permutations. "
            "Each line must contain one character or word."
        ),
        "output_file": "Output file path for the generated wordlist."
    },
    "edit": {

    },
    "analyze": {
        "general": "Print only general info (Filename, Path, Size...)",
        "content": "Print only content info (Total words, Smallest word, Biggest word...)",
        "file_path": "Input file which is to be analyzed"

    }

}

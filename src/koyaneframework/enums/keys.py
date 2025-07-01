from enum import Enum


class StatusKeys(str, Enum):
    # Errors

    # Status Analyzer
    ANALYSIS_STARTED = "analysis_started"

    # Status Generator
    CALCULATE_WORDS = "calculate_words"
    CALCULATE_SIZE = "calculate_size"
    BUILDING_CHAR_WORDLIST = "building_char_wordlist"
    BUILDING_FILE_WORDLIST = "building_file_wordlist"
    BUILDING_MASK_WORDLIST = "building_mask_wordlist"

    # Status before
    TEMP_DIRS = "temp_dirs"

    # Success Generator
    WORDLIST_CREATED = "wordlist_created"


class HelpKeys(str, Enum):
    # before
    QUIET_MODE = "quiet_mode"
    # generate category
    MIN_LENGTH = "min_length"
    MAX_LENGTH = "max_length"
    MASK = "mask_help"
    CHAR_SET = "char_set"
    WORD_FILE = "word_file"
    OUTPUT_FILE = "output_file"

    # edit category

    # analyze category
    GENERAL = "general"
    CONTENT = "content"
    FILE_PATH = "file_path"
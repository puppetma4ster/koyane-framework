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
    COMPRESS_WORDLIST = "compress_wordlist"
    WORDLIST_STATS = "wordlist_stats"

    # Status before
    TEMP_DIRS = "temp_dirs"

    # Success Generator
    WORDLIST_CREATED = "wordlist_created"
    ARCHIVE_CREATED = "archive_created"


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
    COMPRESS = "compress"

    # edit category
    SORT = "sort"
    INVERT = "invert"
    REMOVE_FILE = "remove_file"
    REMOVE_MASK = "remove_mask"

    # analyze category
    GENERAL = "general"
    CONTENT = "content"
    FILE_PATH = "file_path"
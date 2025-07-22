import shutil

from koyaneframework.core.generator.mask_interpreter import MaskInterpreter
from koyaneframework.core.utils.utils import external_sort, create_new_wordlist, add_new_word_to_wordlist,get_base_temp_dir, TEMP_SUFFIX, LIST_SUFFIX, random_temp_number
from pathlib import Path
import re

class EditWordList:

    def __init__(self, input_file : Path, output_file: Path):
        self.output_file = output_file

        self.temp_path = _generate_temp_path_file()   #  Path with random Number
        shutil.copy(input_file, self.temp_path)

    def sort_wordlist(self):
        new_temp_file = _generate_temp_path_file()
        external_sort(self.temp_path, new_temp_file)

        self.temp_path.unlink()
        self.temp_path = new_temp_file

    def remove_words_from_mask(self, mask: str):
        """
        Removes all words from the current temporary wordlist file (self.temp_path)
        that match the given mask. Rewrites the wordlist without matching entries.
        """
        new_temp_file: Path = _generate_temp_path_file()
        create_new_wordlist(new_temp_file)

        mask = MaskInterpreter(mask)

        with open(self.temp_path, "r", encoding="utf-8") as old_wl, \
                open(new_temp_file, "w", encoding="utf-8") as new_wl:

            for line in old_wl:
                word = line.rstrip("\n")
                if not mask.matches_word(word):
                    add_new_word_to_wordlist(new_temp_file, word)

        self.temp_path.unlink()
        self.temp_path = new_temp_file

    def invert_wordlist(self, regex: str):
        """
        Inverts the words in a wordlist.

        If no pattern is specified, all words will be inverted.
        If a regular expression is given, only words matching the pattern (via `re.fullmatch`) will be inverted.
        All other words remain unchanged.

        :param regex: Optional regular expression as string
        :return: List of strings with selected words reversed
        """
        new_temp_file = _generate_temp_path_file()
        create_new_wordlist(new_temp_file)

        with open(self.temp_path, "r", encoding="utf-8") as wordlist:
            for word in wordlist:
                word = word.rstrip("\n\r") # Remove line breaks and tabs

                if re.fullmatch(regex, word):
                    out_word = word[::-1]
                    add_new_word_to_wordlist(new_temp_file, out_word)

        self.temp_path.unlink()
        self.temp_path =  new_temp_file

    def remove_words_from_list(self, remove_wordlist: Path):
        """
        subtracts all words of a given wordlist from the instance wordlist

        :param remove_wordlist: words which should be deleted
        :return: -
        """

        # -------------------- SORTING --------------------
        sorted_current_wl_path = _generate_temp_path_file()
        external_sort(self.temp_path, sorted_current_wl_path)# sort wordlist which will be edited

        sorted_remove_wl_path = _generate_temp_path_file()  # path for sorted remove wordlist
        external_sort(remove_wordlist, sorted_remove_wl_path)

        new_temp_file = _generate_temp_path_file()
        # -------------------- SUBTRACTING --------------------
        create_new_wordlist(new_temp_file)

        with open(sorted_current_wl_path, "r", encoding="utf-8") as minuend_wl, \
            open(sorted_remove_wl_path, "r", encoding="utf8") as subtrahend_wl, \
            open(new_temp_file, "w", encoding="utf-8") as new_wl:

            word_minuend = minuend_wl.readline()
            word_subtrahend = subtrahend_wl.readline()

            while word_minuend and word_subtrahend:
                word_minuend = word_minuend.rstrip("\n")
                word_subtrahend = word_subtrahend.rstrip("\n")
                if word_minuend < word_subtrahend:
                    add_new_word_to_wordlist(new_temp_file, word_minuend)
                    word_minuend =  minuend_wl.readline()
                elif word_minuend == word_subtrahend:
                    word_minuend = minuend_wl.readline()
                    word_subtrahend = subtrahend_wl.readline()
                else:   # if minuend word ist bigger than subtrahend word
                    word_subtrahend = subtrahend_wl.readline()

            while word_minuend: # copy remaining lines from a
                add_new_word_to_wordlist(new_temp_file, word_minuend)
                word_minuend =  minuend_wl.readline()

            # -------------------- Clean up --------------------
            self.temp_path.unlink()
            sorted_current_wl_path.unlink()
            sorted_remove_wl_path.unlink()

            self.temp_path = new_temp_file



    def flush_finished_wordlist(self):

        if not self.temp_path.is_file():
            raise FileNotFoundError()

        final_path = self.output_file.with_suffix(LIST_SUFFIX).resolve()
        shutil.copy(self.temp_path.resolve(), final_path)

        self.temp_path.unlink()



    def _copy_file_to_temp(self, source_file: Path):
        shutil.copy(source_file, self.temp_path)



def _generate_temp_path_file() -> Path:

    while True:
        temp_nr = random_temp_number()
        temp_path = Path(get_base_temp_dir()) / f"temp_edit_{temp_nr}{TEMP_SUFFIX}"
        if not temp_path.is_file():     # if the file already exists
            return temp_path

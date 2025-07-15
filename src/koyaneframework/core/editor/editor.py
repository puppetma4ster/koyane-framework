import shutil

from koyaneframework.core.utils.utils import external_sort, create_new_wordlist, add_new_word_to_wordlist,get_base_temp_dir, TEMP_SUFFIX, LIST_SUFFIX, random_temp_number
from pathlib import Path
import re

#TEMP_PATH  = Path(get_base_temp_dir()) / f"temp_{TEMP_SUFFIX}"
class EditWordList:

    def __init__(self, input_file : Path, output_file: Path):
        self.input_file = input_file
        self.output_file = output_file

        self.temp_path = _generate_temp_path_file()   #  Path with random Number
        self._copy_file_to_temp(input_file)

    def sort_wordlist(self, input_file: Path):
        new_temp_file = _generate_temp_path_file()
        external_sort(self.temp_path, new_temp_file)

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

    def remove_words_from_list(self):






    def flush_temp_wordlist(self):

        if not self.temp_path.is_file():
            raise FileNotFoundError()

        final_path = self.output_file.with_suffix(LIST_SUFFIX).resolve()
        shutil.copy(self.temp_path.resolve(), final_path)

        self.temp_path.unlink()



    def _copy_file_to_temp(self, source_file: Path):
        shutil.copy(source_file, self.temp_path)



def _generate_temp_path_file() -> Path:
    temp_nr = random_temp_number()
    temp_path = Path(get_base_temp_dir()) / f"temp_edit_{temp_nr}{TEMP_SUFFIX}"
    return temp_path

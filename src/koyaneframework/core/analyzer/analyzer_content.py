from  pathlib import Path
from koyaneframework.core.utils.utils import external_sort, get_base_temp_dir, SPECIAL_CHARACTERS
class ContentAnalyzer:


    file_path: Path = None
    wordlines = 0

    smallestword = None
    biggestword = None

    av_wordlen = 0

    char_count = {}

    dublicate_words = []

    DIGIT = "Has digit"
    UPPER = "Has upper"
    SPECSIGN = "Has special character"

    DIGITANDUPPER = "Has digit & upper"
    DIGITANDSPECSIGN = "Has digit & special character"
    UPPERANDSPECSIGN = "Has upper & special character"

    DIGITUPPERSPECSIGN = "Has digit & upper & special character"

    pw_type_list = {
        DIGIT: 0,
        UPPER: 0,
        SPECSIGN: 0,

        DIGITANDUPPER: 0,
        DIGITANDSPECSIGN: 0,
        UPPERANDSPECSIGN: 0,

        DIGITUPPERSPECSIGN: 0

    }

    def __init__(self, file_path: Path):
        self.file_path = file_path



    def analyze_file(self, count: bool=False, min_max:bool=False, av_length: bool=False, char_freq: bool=False, av_entropy: bool=False,dupl: bool=False, weak_pw: bool=False, perc_pw: bool=False):
        output_file = self.file_path
        if dupl:
            output_file = Path(get_base_temp_dir()) / "sorted_list.kyftmp"
            external_sort(self.file_path, output_file)
        
        with open(output_file, "r", encoding="utf-8") as file:
            # helper var for calculating $av_wordlen
            totalwordlen = 0
            last_word = None
            for line in file:
                word = line
                # count flag
                if count or av_length:
                    self.wordlines += 1
                if min_max:
                    self._password_min_max_info(word)
                if av_length:
                    totalwordlen += len(word)
                if char_freq:
                    self._char_frequency(word)
                if av_entropy:
                    # displays average password-entropy
                    pass
                if dupl:
                    self._password_duplicate(last_word,word)
                    last_word = word
                if weak_pw:
                    # weak password Identifier ... coming soon
                    pass
                if perc_pw:
                    self._password_type(word)



            if av_length:
                self.av_wordlen = totalwordlen / self.wordlines
            if perc_pw:
                self._password_type_statistics()
            if dupl:
                output_file.unlink()



    def _password_min_max_info(self, word):
        word_length =len(word.strip("\n")) # because line breaks otherwise count as characters
        if self.smallestword is None or self.biggestword is None:
            self.smallestword = word_length
            self.biggestword = word_length

        if word_length < self.smallestword:
            self.smallestword = word_length
        if word_length > self.biggestword:
            self.biggestword = word_length


    def _char_frequency(self, word):
        for char in word:
            if char in self.char_count:
                self.char_count[char] += 1
            else:
                self.char_count[char] = 1

    def _password_duplicate(self, last_word, duplicate_word):
        if last_word == duplicate_word:
            if duplicate_word not in self.dublicate_words:
                self.dublicate_words.append(duplicate_word)

    def _password_type(self, word: str):
        digit = False
        upper = False
        specsign = False

        for char in word:
            if char.isdigit():
                digit = True
            if char.isupper():
                upper = True
            if char in SPECIAL_CHARACTERS:
                specsign = True

        if digit and upper and specsign:
            self.pw_type_list[self.DIGITUPPERSPECSIGN] += 1
        elif digit and upper:
            self.pw_type_list[self.DIGITANDUPPER] += 1
        elif digit and specsign:
            self.pw_type_list[self.DIGITANDSPECSIGN] += 1
        elif upper and specsign:
            self.pw_type_list[self.UPPERANDSPECSIGN] += 1
        elif digit:
            self.pw_type_list[self.DIGIT] += 1
        elif upper:
            self.pw_type_list[self.UPPER] += 1
        elif specsign:
            self.pw_type_list[self.SPECSIGN] += 1


    def _password_type_statistics(self):
        for type in self.pw_type_list:
            PERCENT_MULTIPLIER = 100
            percent = (self.pw_type_list[type] / self.wordlines) * PERCENT_MULTIPLIER
            self.pw_type_list[type] = percent






    #GETTER---------------------------------------------------------------------------------

    def get_wordlines(self):
        return self.wordlines

    def get_smallest_word(self):
        return self.smallestword

    def get_biggest_word(self):
        return self.biggestword

    def get_average_word_length(self):
        return self.av_wordlen

    def get_occurring_characters(self):
        return self.char_count

    def get_duplicate_words(self):
        return self.dublicate_words

    def get_average_entropy(self):
        #comeing soon...
        pass

    def get_weak_passwords(self):
        #comeing soon...
        pass

    def get_type_occur_percent(self):
        return self.pw_type_list


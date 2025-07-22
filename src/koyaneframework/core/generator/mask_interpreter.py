from typing import List

from koyaneframework.core.generator.maskChar import MaskChar


class MaskInterpreter:
    mask: str = None
    mask_segments: list[MaskChar] = []  # Python 3.9+

    def __init__(self, mask: str):
        self.mask = mask.strip()
        self.mask_segments = []
        segment = ""
        inside_mask = False

        for char in self.mask:
            if char == "?" or char == "!":
                if inside_mask: # add built segment
                    if char == "!":
                        if segment.startswith("?"):  # add last segment
                            self.mask_segments.append(MaskChar(segment, True))
                        elif segment.startswith("!"):
                            self.mask_segments.append(MaskChar(segment, False))
                        segment = "!"
                    else:
                        if segment.startswith("?"):  # add last segment
                            self.mask_segments.append(MaskChar(segment, True))
                        elif segment.startswith("!"):
                            self.mask_segments.append(MaskChar(segment, False))
                        segment = "?"
                else:
                    segment = char
                    inside_mask = True
            else:
                segment += char

        if segment.startswith("?"):  # add last segment
            self.mask_segments.append(MaskChar(segment, True))
        elif segment.startswith("!"):
            self.mask_segments.append(MaskChar(segment, False))


    def matches_word(self, word: str) -> bool:
        """
        Checks whether the given word matches the mask defined by self.mask_segments.

        Each segment in self.mask_segments represents a position in the word
        and provides a set of permitted characters via get_permitted_characters().
        The word matches the mask if every character at position i in the word
        is contained in the permitted characters of segment i.

        :param word: The input string to be checked against the mask.
        :return: True if the word matches the mask, False otherwise.
        """
        if len(word) != len(self.mask_segments):    # if word length does not match mask length
            return False
        match: bool = False

        for i, segment in enumerate(self.mask_segments):

            match: bool = False
            for char in segment.get_permitted_characters(): # character matching
                if word[i] == char:
                    match = True
                    break

            if not match:   # mask matches the word
                return False
        return True    # mask matches the word
"""
maskChar.py – Parsing individual mask segments to determine allowed characters

This file defines the `MaskChar` class, which parses individual mask segments
(e.g., "?l", "?d", "?v", etc.) and determines which character groups are permitted
for password or wordlist generation.

Each instance of `MaskChar` interprets a mask segment, sets flags, and compiles
a string of allowed characters to be used later for string generation.

Masks can represent lowercase letters, uppercase letters, vowels, consonants,
digits, and various classes of special characters.

Example:
MaskChar("?lV") allows all lowercase letters plus all uppercase vowels.

Author: puppetm4ster
"""
import koyaneframework.core.utils.utils as msk_c
class MaskChar:

    def __init__(self, mask_character: str, is_wildcard: bool):
        self.mask_character: str = mask_character.strip().lstrip('?')
        self.is_wildcard = is_wildcard

        self.all_upper_case: bool = False
        self.all_lower_case: bool = False
        self.all_vowels_lower_case: bool = False
        self.all_vowels_upper_case: bool = False
        self.all_consonants_lower_case: bool = False
        self.all_consonants_upper_case: bool = False
        self.all_digit: bool = False
        self.all_special: bool = False
        self.special_most_used: bool = False
        self.special_points: bool = False
        self.special_bracelet: bool = False

        self.permitted_characters: str = ""

        self._convert_mask_char()


    def _convert_mask_char(self):
        """
            Converts the mask character(s) into internal flags and a character set.

            This method updates flags like `all_lower_case` or `all_digit`, and fills
            `permitted_characters` based on the given mask (e.g., 'l', 'V', 'd', etc.).
            """
        if self.is_wildcard:    #if wildcard
            for char in self.mask_character:
                if char == "l":
                    self.all_lower_case = True
                    self.permitted_characters += msk_c.LOWER_CASE_CHARACTERS
                elif char == "L":
                    self.all_upper_case = True
                    self.permitted_characters += msk_c.UPPER_CASE_CHARACTERS
                elif char == "v":
                    self.all_vowels_lower_case = True
                    self.permitted_characters += msk_c.LOWER_CASE_VOWELS
                elif char == "V":
                    self.all_vowels_upper_case = True
                    self.permitted_characters += msk_c.UPPER_CASE_VOWELS
                elif char == "c":
                    self.all_consonants_lower_case = True
                    self.permitted_characters += msk_c.LOWER_CASE_CONSONANTS
                elif char == "C":
                    self.all_consonants_upper_case = True
                    self.permitted_characters += msk_c.UPPER_CASE_CONSONANTS
                elif char == "d":
                    self.all_digit = True
                    self.permitted_characters += msk_c.DIGITS
                elif char == "s":
                    self.all_special = True
                    self.permitted_characters += msk_c.SPECIAL_CHARACTERS
                elif char == "f":
                    self.special_most_used = True
                    self.permitted_characters += msk_c.SPECIAL_CHARACTERS_MOST_USED
                elif char == "p":
                    self.special_points = True
                    self.permitted_characters += msk_c.SPECIAL_CHARACTERS_POINTS
                elif char == "b":
                    self.special_bracelet = True
                    self.permitted_characters += msk_c.SPECIAL_CHARACTERS_BRACELET
                else:
                    raise ValueError(f"Unknown character: ? {self.mask_character}")

        elif not self.is_wildcard:  # if literal
            for char in self.mask_character:
                if char == "!":
                    continue
                self.permitted_characters += char


    def get_permitted_characters(self):
        return self.permitted_characters
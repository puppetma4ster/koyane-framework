
from pathlib import Path
from koyaneframework.core.analyzer.analyzer_general import GeneralAnalyzer
from koyaneframework.core.analyzer.analyzer_content import ContentAnalyzer
from koyaneframework.load_animation import LoadingSpinner



def print_all_info_wordlist(filepath: Path):
    print_general_info_wordlist(filepath)
    print_content_info_wordlist(filepath)




# general data about the file itself
def print_general_info_wordlist(filepath: Path):
    general = GeneralAnalyzer(filepath)
    print()
    print("---------- GENERAL WORDLIST INFO ----------")
    print(f"File Name.....{general.name()}")
    print(f"Path.....{general.path()}")
    print(f"Size..... {byte_to_megabyte(general.size())} MB / {general.size()} Byte")
    print(f"Last changes..... {general.last_changes()}")
    print(f"Last access..... {general.last_access()}")
    print(f"Created..... {general.created()}")
    print(f"File extension..... {general.extension()}")
    print(f"Estimated encoding..... {general.encoding()}")
    print(f"File hash..... {general.hash()}")
    print()


def print_content_info_wordlist(filepath: Path):
    load = LoadingSpinner("Loading content")
    load.start()
    content = ContentAnalyzer(filepath)
    content.analyze_file(count=True,min_max=True,av_length=True,char_freq=True,av_entropy=False,dupl=True,weak_pw=False,perc_pw=True)
    load.stop()
    print("---------- WORDLIST CONTENT STATISTICS ----------")
    print(f"Total words..... {content.get_wordlines()}")
    print(f"Smallest word..... {content.get_smallest_word()}")
    print(f"Biggest word..... {content.get_biggest_word()}")
    print(f"Average length..... {round(content.get_average_word_length(), 2)}")
    print(f"chars.....{used_chars_to_string(content.get_occurring_characters())}")
    print(f"duplicates found..... {duplicates_to_string(content.get_duplicate_words())}")
    print()

    print("---------- WORDLIST PERCENT STATISTICS ----------")
    print(f"{type_percentage_to_string(content.get_type_occur_percent())}")
    print()


def byte_to_megabyte(byte):
    mb =  1_048_576
    return  round(byte / mb, 2)


from textwrap import wrap


def used_chars_to_string(chars):
    if not chars:
        return "no characters found"

    sorted_chars = sorted(chars)
    all_chars = ",".join(sorted_chars)

    indent = "chars..... "
    wrapped_lines = wrap(all_chars, width=100 - len(indent))  # zieh Einzugslänge ab
    return ("\n" + " " * len(indent)).join([wrapped_lines[0]] + wrapped_lines[1:])


def duplicates_to_string(dup):
    if not dup:
        return "no duplicates found"

    indent = "dublicates found..... "
    words = [word.strip() for word in sorted(dup)]

    lines = []
    for i in range(0, len(words), 10):
        line = ",".join(words[i:i + 10])
        lines.append(line)

    # Indentation: first line with prefix, all subsequent lines indented accordingly.
    return ("\n" + " " * len(indent)).join([lines[0]] + lines[1:])


def type_percentage_to_string(type_list: dict):
    return "\n".join(f"{k}..... = {v:.2f}%" for k, v in type_list.items())


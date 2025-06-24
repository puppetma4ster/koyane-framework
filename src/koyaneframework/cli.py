import typer
from pathlib import Path
from koyaneframework.output_printer import show_banner
from koyaneframework.analyzer.analyzer import print_all_info_wordlist, print_content_info_wordlist, print_general_info_wordlist
from koyaneframework.generator.wordlist_generator import generate_mask_wordlist, generate_wordlist
from koyaneframework.editor.editor import sort_wordlist
from koyaneframework.word_sources import load_chars_from_input, load_words_from_file

app = typer.Typer()

@app.command(help="generates wordlists from scratch")
def generate(
        min_length: int = typer.Option(
            None,
            "-m",
            "--min-length",
            help="specifies the minimum word length."),
        max_length: int = typer.Option(
            None,
            "-x",
            "--max-length",
            help="specifies the maximum word length."),
        mask: str = typer.Option(
            None,
            "-ms",
            "--mask",
            help=("generate wordlist from mask."
                  "A mask is a string consisting of segments. Each segment begins with a “?” followed by one or more letters that define the character type. Example: “?ld?d?f”."
                  "available letter wildcards are “l” = small letter, “L” = capital letter, “v” = small vowel, “V” = capital vowel, “c” = small consonant, “C” = capital consonant, “d” = number,"
                  " “s” = special character, “f” = most commonly used special character, ‘p’ = dot special character and “b” = bracket special character")),
        char_set: str = typer.Option(
            None,
            "-cs",
            "--char_set",
            help=("generate a wordlist with chars."
                  "Example: -cs abcdef1234 || -cs \'#+:_<>abcd123\'")),
        word_file: Path = typer.Option(
            None,
            "-cf",
            "--char_file",
            exists=True,
            dir_okay=False,
            file_okay=True,
            help=("generate a wordlist with a file."
                  "This is useful for word combinations."
                  "The file must contain one char or word per line")),
        output_file: Path = typer.Argument(
            ...,
            exists=False,
            dir_okay=False,
            file_okay=True,
            help="output file")
):

    if mask:    # simple mask generation
        generate_mask_wordlist(mask, output_file)
    elif mask and min_length:   # maskgeneration with min length
        generate_mask_wordlist(mask, output_file, min_len=min_length)
    elif char_set and min_length and max_length:    #char set
        if min_length > max_length:
            raise typer.BadParameter("min_length cannot be greater than max_length")
        else:
            chars = load_chars_from_input(char_set)
            generate_wordlist(chars, min_length, max_length, output_file)
    elif word_file and min_length and max_length:   # word file
        if min_length > max_length:
            raise typer.BadParameter("min_length cannot be greater than max_length")
        else:
            chars = load_words_from_file(word_file)
            generate_wordlist(chars, min_length, max_length, output_file)

@app.command(help="Edit existing word lists")
def edit(
        sort: bool =typer.Option(
            False,
            "--sort",
            "-s",
            help="sort a wordlist"),
        input_file: Path = typer.Argument(
            ...,
            exists=True,
            dir_okay=False,
            file_okay=True,
            readable=True,
            help="Input file which is to be edited"),
        output_file: Path = typer.Argument(
            None,
            exists=False,
            dir_okay=False,
            file_okay=True,
            help="Output file path - default = output/wl.txt")
):
    if output_file is None:
        output_file = Path(__file__).resolve().parent / "output" / "wl.txt"
    else:
        output_file = output_file.resolve()
    output_file.parent.mkdir(parents=True, exist_ok=True)

    show_banner()
    if sort:

        sort_wordlist(input_file,output_file)


@app.command(help="get detailed properties of a wordlist")
def analyze(
        general: bool =typer.Option(
            False,
            "--general",
            "-g",
            help="Print only general info (Filename, Path, Size...)"),
        content: bool = typer.Option(
            False,
            "--content",
            "-c",
            help="Print only content info (Total words, Smallest word, Biggest word...)"
        ),
        file_path: Path = typer.Argument(
            ...,
            exists=True,
            dir_okay=False,
            file_okay=True,
            readable=True,
            help="Input file which is to be analyzed")
):
    show_banner()
    if content and general or not content and not  general:
        pass
        print_all_info_wordlist(file_path)
    elif general:
         print_general_info_wordlist(file_path)

    elif content:
        print_content_info_wordlist(file_path)





@app.command(help="search online for suitable word lists for a specified application and download them (WPA2 ...)")
def search():
    show_banner()
    print("This function is still in work...")


@app.command(help="configurate this script")
def configurate():
    show_banner()





if __name__ == '__main__':
    app()
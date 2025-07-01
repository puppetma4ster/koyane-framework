import typer
from pathlib import Path

# Imports for text displays and output prints
from koyaneframework.output.banner import get_banner_lines
from koyaneframework.output.output_printer import print_lines, print_status, print_success, print_error, print_warning
from koyaneframework.output.help_cli_texts import HELP_TEXTS
from koyaneframework.enums.keys import StatusKeys, HelpKeys
from koyaneframework.enums.Categories import StatusCategories, HelpCategories

# imports for main logic
from koyaneframework.core.analyzer.analyzer import print_all_info_wordlist, print_content_info_wordlist, print_general_info_wordlist
from koyaneframework.core.generator.wordlist_generator import generate_mask_wordlist, generate_wordlist
from koyaneframework.core.editor.editor import sort_wordlist
from koyaneframework.core.utils.utils import prepare_temp_dirs

# imports for helper functions
from koyaneframework.core.utils.word_sources import load_chars_from_input, load_words_from_file

app = typer.Typer()

@app.command(help="generates wordlists from scratch")
def generate(
        min_length: int = typer.Option(
            None,
            "-m",
            "--min-length",
            help=HELP_TEXTS[HelpCategories.GENERATE][HelpKeys.MIN_LENGTH]),
        max_length: int = typer.Option(
            None,
            "-x",
            "--max-length",
            help=HELP_TEXTS[HelpCategories.GENERATE][HelpKeys.MAX_LENGTH]),
        mask: str = typer.Option(
            None,
            "-ms",
            "--mask",
            help=HELP_TEXTS[HelpCategories.GENERATE][HelpKeys.MASK]),
        char_set: str = typer.Option(
            None,
            "-cs",
            "--char-set",
            help=HELP_TEXTS[HelpCategories.GENERATE][HelpKeys.CHAR_SET]),
        word_file: Path = typer.Option(
            None,
            "-cf",
            "--char-file",
            exists=True,
            dir_okay=False,
            file_okay=True,
            help=HELP_TEXTS[HelpCategories.GENERATE][HelpKeys.WORD_FILE]),
        output_file: Path = typer.Argument(
            ...,
            exists=False,
            dir_okay=False,
            file_okay=True,
            help=HELP_TEXTS[HelpCategories.GENERATE][HelpKeys.OUTPUT_FILE])
):
    """
        Main entry point for generating wordlists using different methods.

        Depending on the selected options, this command can:
        - Generate wordlists from a character set with defined min/max length.
        - Generate wordlists based on a structural mask pattern.
        - Combine or transform entries from a custom word/char file.

        Args:
            min_length (int): Minimum word length (used with --char-set).
            max_length (int): Maximum word length (used with --char-set).
            mask (str): Pattern-based wordlist using character wildcards.
            char_set (str): Custom set of characters to build from.
            word_file (Path): File containing words or characters, one per line.
            output_file (Path): File path to write the generated wordlist to.
        """

    if mask and min_length:   # maskgeneration with min length
        print_status(StatusCategories.STATUS_GENERATOR, StatusKeys.BUILDING_MASK_WORDLIST, mask=mask)
        generate_mask_wordlist(mask, output_file, min_len=min_length)
    elif mask:    # simple mask generation
        print_status(StatusCategories.STATUS_GENERATOR,StatusKeys.BUILDING_MASK_WORDLIST, mask=mask)
        generate_mask_wordlist(mask, output_file)
    elif char_set and min_length and max_length:    #char set
        if min_length > max_length:
            raise typer.BadParameter("min_length cannot be greater than max_length")
        else:
            print_status(StatusCategories.STATUS_GENERATOR,StatusKeys.BUILDING_CHAR_WORDLIST,charset=char_set)
            chars = load_chars_from_input(char_set)
            generate_wordlist(chars, min_length, max_length, output_file)
    elif word_file and min_length and max_length:   # word file
        if min_length > max_length:
            raise typer.BadParameter("min_length cannot be greater than max_length")
        else:
            print_status(StatusCategories.STATUS_GENERATOR, StatusKeys.BUILDING_FILE_WORDLIST, path=word_file)
            chars = load_words_from_file(word_file)
            generate_wordlist(chars, min_length, max_length, output_file)

    print_success(StatusCategories.SUCCESS_GENERATOR, StatusKeys.WORDLIST_CREATED, path=output_file)
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


    if sort:

        sort_wordlist(input_file,output_file)


@app.command(help="get detailed properties of a wordlist")
def analyze(
        general: bool =typer.Option(
            False,
            "--general",
            "-g",
            help=HELP_TEXTS[HelpCategories.ANALYZE][HelpKeys.GENERAL]),
        content: bool = typer.Option(
            False,
            "--content",
            "-c",
            help=HELP_TEXTS[HelpCategories.ANALYZE][HelpKeys.CONTENT]
        ),
        file_path: Path = typer.Argument(
            ...,
            exists=True,
            dir_okay=False,
            file_okay=True,
            readable=True,
            help=HELP_TEXTS[HelpCategories.ANALYZE][HelpKeys.FILE_PATH])
):
    """
        Main entry point for analyzing wordlists.

        This command can:
        - Analyze the entire wordlist.
        - analyze the content or the wordlist.

        Args:
            general (bool): shows only general file infos.
            content (int): Shows only infos about the content of the wordlist.
            file_path (Path): Which wordlist is to be analyzed.
            """

    print_status(StatusCategories.STATUS_ANALYZER, StatusKeys.ANALYSIS_STARTED,path=file_path)
    if content and general or not content and not  general:
        print_all_info_wordlist(file_path)

    elif general:
         print_general_info_wordlist(file_path)

    elif content:
        print_content_info_wordlist(file_path)





@app.command(help="search online for suitable word lists for a specified application and download them (WPA2 ...)")
def search():
    print("This function is still in work...")



@app.callback()
def before(
        quiet_mode: bool= typer.Option(
            False,
            "-q",
            "--quiet-mode",
            help=HELP_TEXTS[HelpCategories.BEFORE][HelpKeys.QUIET_MODE]
        )
):
    """
        Main entry point of the program.
        Everything that should happen before a command is regulated here.

        Here:
        - The banner is printed onto the console.
        Args:
            quiet_mode: prevents banner output
        """
    # generate tmp dirs
    print_status(StatusCategories.STATUS_BEFORE,StatusKeys.TEMP_DIRS)
    prepare_temp_dirs()
    # print ACII banner
    if not quiet_mode:
        print_lines(get_banner_lines(), style="red")



if __name__ == '__main__':
    app()
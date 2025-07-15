import typer

STATUS_MESSAGES = {            #[+] = Success [*] = Status [!] = Warning [-] = Error

    "errors": {
        "prefix": "[-]"

    },
    "status_analyzer": {
        "prefix": "[*]",
        "analysis_started": "Analyzing wordlist: {path}"
    },
    "status_generator": {
        "prefix": "[*]",
        "calculate_words": "Final wordlist has {wordcount} entries.",
        "calculate_size": "Final wordlist size: {size}",
        "building_char_wordlist": "Building wordlist using the following character set: \'{charset}\'",
        "building_file_wordlist": "Building wordlist from file entries at: {path}",
        "building_mask_wordlist": "Building wordlist using the following mask: \'{mask}\'",
        "compress_wordlist": "Compress generated wordlist to: {path}"
    },
    "status_before": {
        "prefix": "[*]",
        "temp_dirs": "Generating temp dictionaries if missing."
    },
    "success_generator": {
        "prefix": "[+]",
        "wordlist_created": "Wordlist successfully created at: {path}",
        "archive_created": "Compressed wordlist successfully created at: {path}"
    }
}



def print_lines(lines: list[str], style: str = "none"):
    """
    Prints a list of text lines to the console using optional color styling.

    Args:
        lines (list[str]): The lines of text to print.
        style (str, optional): Color style for the output (e.g., 'red').
            Use "none" (default) for no styling.

    Returns:
        None
    """
    for line in lines:
        typer.echo(typer.style(line, fg=style if style != "none" else None))

def print_status(category: str, key: str, **kwargs):
    """
    Prints a normal (unstyled) status message.
    """
    _print_message(category, key, style=None, **kwargs)

def print_success(category: str, key: str, **kwargs):
    """
    Prints a success message in green.
    """
    _print_message(category, key, style="green", **kwargs)

def print_error(category: str, key: str, **kwargs):
    """
    Prints an error message in red.
    """
    _print_message(category, key, style="red", **kwargs)

def print_warning(category: str, key: str, **kwargs):
    """
    Prints a warning message in yellow.
    """
    _print_message(category, key, style="yellow", **kwargs)

def _print_message(category: str, key: str, style: str = None, **kwargs):
    """
    Internal helper to format and print a message using the message dictionary.
    """
    prefix = STATUS_MESSAGES.get(category, {}).get("prefix", "")
    message = STATUS_MESSAGES.get(category, {}).get(key, "")
    if kwargs:
        message = message.format(**kwargs)
    typer.echo(typer.style(f"{prefix} {message}", fg=style))


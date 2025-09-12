# Database Documentation: Wordlists & Tags

## Purpose
This documentation describes the database structure for storing wordlists and their associated tags, including metadata about each wordlist and allowed classification tags.

---

## Tables Overview

### `wordlists`
| Column                       | Type       | Description                                                    |
|------------------------------|------------|----------------------------------------------------------------|
| id                           | INTEGER PK | Unique identifier for the wordlist                             |
| name                         | TEXT       | Name of the wordlist file                                      |
| entities                     | INTEGER    | Number of entries in the wordlist                              |
| smallestEntity               | INTEGER    | Length of the smallest entry                                   |
| biggestEntity                | INTEGER    | Length of the biggest entry                                    |
| averageLength                | REAL       | Average length of entries                                      |
| averageEntropy               | REAL       | Average entropy across entries                                 |
| digitsPercent                | REAL       | Percentage of digits in entries                                |
| upperCasePercent             | REAL       | Percentage of uppercase letters                                |
| specialCharPercent           | REAL       | Percentage of special characters                               |
| digitAndUpperCase            | REAL       | Percentage containing digits AND uppercase                     |
| digitAndSpecialChar          | REAL       | Percentage containing digits AND special characters            |
| upperCaseAndSpecialChar      | REAL       | Percentage containing uppercase AND special characters         |
| digitUpperCaseAndSpecialChar | REAL       | Percentage containing digits, uppercase AND special characters |
| encoding                     | TEXT       | Encoding of the wordlist (e.g., UTF-8, ASCII)                  |
| language                     | TEXT       | Language of the wordlist                                       |
| category                     | TEXT       | High-level category (e.g., web, wifi, auth)                    |
| author                       | TEXT       | Creator / contributor of the wordlist                          |
| size                         | INTEGER    | File size in bytes                                             |
| link                         | TEXT       | link to source or reference                                    |
| info                         | TEXT       | detaild extra infos                                            |

### `wordlist_tags`
| Column      | Type    | Description                                           |
|-------------|---------|-------------------------------------------------------|
| wordlist_id | INTEGER | References `wordlists.id`                             |
| tag         | TEXT    | Tag applied to this wordlist (see allowed tags below) |

---

## Allowed Tags

### Platform / Target
`web`, `cms`, `wordpress`, `drupal`, `joomla`, `shop`, `api`, `ftp`, `ssh`, `smtp`, `windows`, `linux`, `android`, `ios`

### Auth / Credentials / Wi-Fi
`wifi`, `wpa`, `wpa2`, `wpa3`, `wep`, `credentials`, `usernames`, `passwords`, `default-creds`

### Attack Type / Use-case
`discovery`, `fuzzing`, `webshell`, `rce`, `lfi`, `sqli`, `xss`, `dir-traversal`, `hashcracking`,
`scanner`, `forensics`

### Content / Pattern / Encoding
`regex`, `string`, `hash`, `sha1`, `sha256`, `base64`, `obfuscated`, `minified`, `binary`, `text`

---

## Allowed Categories

`leaked`, `dictionary`, `generated`, `collection`

---

## Guidelines
- All tags must be lowercase.
- Multi-word tags should use hyphens or underscores (e.g., `double-ext`).
- Each tag in `wordlist_tags.tag` must be valid according to the lists above.
- Wordlist metadata should be kept accurate for automated analysis (entropy, character percentages, encoding).
- Tags are primarily used for classification, discovery workflows, fuzzing, or reporting purposes.
- Authors should ensure consistency and maintain proper attribution in the `author` column.
- Category in `wordlists.category` should reflect the main use-case (platform, attack type, content type).  


# Usage String for koyane_db_importer

`poetry run python src/koyane_db_importer/cli.py --name "cirt-default-usernames.txt" --language "-" --category collection --author "-" --tag usernames --tag web --tag cms --tag ftp --tag ssh --tag windows --tag linux --tag usernames --tag default-creds --tag credentials --tag discovery --tag scanner --tag string --tag text --info "found in github repo https://github.com/danielmiessler/SecLists?tab=readme-ov-file" --link "https://crackstation.net/files/crackstation-human-only.txt.gz"`
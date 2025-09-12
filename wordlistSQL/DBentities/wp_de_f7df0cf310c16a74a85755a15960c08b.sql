PRAGMA foreign_keys = ON;
BEGIN;

-- Insert wordlist if not present
INSERT INTO wordlists (
  name, entities, smallestEntity, biggestEntity,
  averageLength, averageEntropy,
  digitsPercent, upperCasePercent, specialCharPercent,
  digitAndUpperCase, digitAndSpecialChar, upperCaseAndSpecialChar, digitUpperCaseAndSpecialChar,
  encoding, language, category, author,
  size, link
)
SELECT
  'wp_de.txt', 5330677, 8, 24,
  13.448970177709135, 3.080860089572083,
  0.0, 0.0, 0.0,
  0.0, 0.0, 0.0, 0.0,
  'ISO-8859-1', 'german', 'dictionary', 'alex stanev',
  77022793, 'https://sec.stanev.org/dict/wp_de.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='wp_de.txt' AND size=77022793 AND entities=5330677
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='wp_de.txt' AND size=77022793 AND entities=5330677;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='wp_de.txt' AND size=77022793 AND entities=5330677;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='wp_de.txt' AND size=77022793 AND entities=5330677;
COMMIT;

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
  'hashes_2017.txt', 214720062, 8, 24,
  9.756085013611823, 2.8601300819716977,
  49.586113, 1.7327402, 2.2932618,
  8.315037, 3.386834, 0.26127648, 1.0410647,
  'UTF-8', '-', 'leaked', 'alex stanev',
  2309565205, 'https://sec.stanev.org/dict/hashes_2017.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='hashes_2017.txt' AND size=2309565205 AND entities=214720062
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='hashes_2017.txt' AND size=2309565205 AND entities=214720062;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='hashes_2017.txt' AND size=2309565205 AND entities=214720062;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='hashes_2017.txt' AND size=2309565205 AND entities=214720062;
COMMIT;

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
  'wp_fr.txt', 1284151, 8, 24,
  10.597280226390822, 2.878676164079019,
  0.0, 0.0, 0.0,
  0.0, 0.0, 0.0, 0.0,
  'ISO-8859-1', 'french', 'dictionary', 'alex stanev',
  14892659, 'https://sec.stanev.org/dict/wp_fr.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='wp_fr.txt' AND size=14892659 AND entities=1284151
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='wp_fr.txt' AND size=14892659 AND entities=1284151;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='wp_fr.txt' AND size=14892659 AND entities=1284151;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='wp_fr.txt' AND size=14892659 AND entities=1284151;
COMMIT;

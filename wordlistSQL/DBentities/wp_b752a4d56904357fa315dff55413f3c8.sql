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
  'wp.txt', 5843249, 8, 24,
  11.287855437959259, 2.9243065624576823,
  0.0, 0.0, 0.0,
  0.0, 0.0, 0.0, 0.0,
  'Big5', 'english', 'dictionary', 'alex stanev',
  71800999, 'https://sec.stanev.org/dict/wp.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='wp.txt' AND size=71800999 AND entities=5843249
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='wp.txt' AND size=71800999 AND entities=5843249;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='wp.txt' AND size=71800999 AND entities=5843249;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='wp.txt' AND size=71800999 AND entities=5843249;
COMMIT;

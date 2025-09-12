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
  'num8_0.txt', 8905349, 8, 8,
  8.0, 2.379886239587678,
  100.0, 0.0, 0.0,
  0.0, 0.0, 0.0, 0.0,
  'UTF-8', '-', 'generated', 'alex stanev',
  80148141, 'https://sec.stanev.org/dict/num8_0.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='num8_0.txt' AND size=80148141 AND entities=8905349
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='num8_0.txt' AND size=80148141 AND entities=8905349;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='num8_0.txt' AND size=80148141 AND entities=8905349;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='num8_0.txt' AND size=80148141 AND entities=8905349;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'numeric' FROM wordlists WHERE name='num8_0.txt' AND size=80148141 AND entities=8905349;
COMMIT;

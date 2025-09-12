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
  'num8_7.txt', 8982647, 8, 8,
  8.0, 2.3792816189589927,
  100.0, 0.0, 0.0,
  0.0, 0.0, 0.0, 0.0,
  'UTF-8', '-', 'generated', 'alex stanev',
  80843823, 'https://sec.stanev.org/dict/num8_7.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='num8_7.txt' AND size=80843823 AND entities=8982647
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='num8_7.txt' AND size=80843823 AND entities=8982647;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='num8_7.txt' AND size=80843823 AND entities=8982647;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='num8_7.txt' AND size=80843823 AND entities=8982647;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'numeric' FROM wordlists WHERE name='num8_7.txt' AND size=80843823 AND entities=8982647;
COMMIT;

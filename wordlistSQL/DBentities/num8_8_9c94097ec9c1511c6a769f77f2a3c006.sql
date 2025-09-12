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
  'num8_8.txt', 8980138, 8, 8,
  8.0, 2.3793188797097455,
  100.0, 0.0, 0.0,
  0.0, 0.0, 0.0, 0.0,
  'UTF-8', '-', 'generated', 'alex stanev',
  80821242, 'https://sec.stanev.org/dict/num8_8.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='num8_8.txt' AND size=80821242 AND entities=8980138
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='num8_8.txt' AND size=80821242 AND entities=8980138;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='num8_8.txt' AND size=80821242 AND entities=8980138;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='num8_8.txt' AND size=80821242 AND entities=8980138;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'numeric' FROM wordlists WHERE name='num8_8.txt' AND size=80821242 AND entities=8980138;
COMMIT;

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
  'num8_6.txt', 8981079, 8, 8,
  8.0, 2.379306907562021,
  100.0, 0.0, 0.0,
  0.0, 0.0, 0.0, 0.0,
  'UTF-8', '-', 'generated', 'alex stanev',
  80829711, 'https://sec.stanev.org/dict/num8_6.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='num8_6.txt' AND size=80829711 AND entities=8981079
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='num8_6.txt' AND size=80829711 AND entities=8981079;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='num8_6.txt' AND size=80829711 AND entities=8981079;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='num8_6.txt' AND size=80829711 AND entities=8981079;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'numeric' FROM wordlists WHERE name='num8_6.txt' AND size=80829711 AND entities=8981079;
COMMIT;

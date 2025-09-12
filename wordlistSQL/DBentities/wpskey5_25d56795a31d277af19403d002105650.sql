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
  'wpskey5.txt', 1000000, 8, 8,
  8.0, 2.3791421925923752,
  100.0, 0.0, 0.0,
  0.0, 0.0, 0.0, 0.0,
  'UTF-8', '-', 'generated', 'alex stanev',
  9000000, 'https://sec.stanev.org/dict/wpskey5.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='wpskey5.txt' AND size=9000000 AND entities=1000000
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='wpskey5.txt' AND size=9000000 AND entities=1000000;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wps' FROM wordlists WHERE name='wpskey5.txt' AND size=9000000 AND entities=1000000;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'numeric' FROM wordlists WHERE name='wpskey5.txt' AND size=9000000 AND entities=1000000;
COMMIT;

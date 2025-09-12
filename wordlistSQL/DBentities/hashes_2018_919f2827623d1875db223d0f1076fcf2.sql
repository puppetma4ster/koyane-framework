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
  'hashes_2018.txt', 175113966, 8, 24,
  10.477690528692618, 2.8843544772472303,
  66.830734, 1.8490752, 3.2258973,
  8.151176, 6.808151, 0.28532162, 1.4613575,
  'UTF-8', '-', 'leaked', 'alex stanev',
  2009912135, 'https://sec.stanev.org/dict/hashes_2018.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='hashes_2018.txt' AND size=2009912135 AND entities=175113966
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='hashes_2018.txt' AND size=2009912135 AND entities=175113966;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='hashes_2018.txt' AND size=2009912135 AND entities=175113966;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='hashes_2018.txt' AND size=2009912135 AND entities=175113966;
COMMIT;

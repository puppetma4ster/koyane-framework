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
  'old_gold.txt', 1546653, 8, 24,
  12.002617264506002, 2.982381966708667,
  11.690857, 10.688371, 20.114725,
  0.6388634, 0.12931149, 4.4386168, 0.07803948,
  'ISO-8859-1', 'mixed', 'leaked', 'alex stanev',
  20511179, 'https://sec.stanev.org/dict/old_gold.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='old_gold.txt' AND size=20511179 AND entities=1546653
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='old_gold.txt' AND size=20511179 AND entities=1546653;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='old_gold.txt' AND size=20511179 AND entities=1546653;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='old_gold.txt' AND size=20511179 AND entities=1546653;
COMMIT;

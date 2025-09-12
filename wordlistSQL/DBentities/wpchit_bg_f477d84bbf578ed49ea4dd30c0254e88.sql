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
  'wpchit_bg.txt', 1314862, 8, 24,
  10.959560775199222, 2.9231743516477846,
  18.915445, 0.0, 0.0,
  0.0, 0.0, 0.0, 0.0,
  'ISO-8859-1', 'bulgarian', 'dictionary', 'alex stanev',
  15725172, 'https://sec.stanev.org/dict/wpchit_bg.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='wpchit_bg.txt' AND size=15725172 AND entities=1314862
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='wpchit_bg.txt' AND size=15725172 AND entities=1314862;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='wpchit_bg.txt' AND size=15725172 AND entities=1314862;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='wpchit_bg.txt' AND size=15725172 AND entities=1314862;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'leet' FROM wordlists WHERE name='wpchit_bg.txt' AND size=15725172 AND entities=1314862;
COMMIT;

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
  'os.txt', 33879411, 8, 24,
  11.31887112795438, 3.037601942324322,
  26.482943, 16.199068, 7.4201374,
  22.477097, 9.194616, 6.4479337, 0.1169855,
  'ISO-8859-1', '-', 'leaked', 'alex stanev',
  417356323, 'https://sec.stanev.org/dict/os.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='os.txt' AND size=417356323 AND entities=33879411
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='os.txt' AND size=417356323 AND entities=33879411;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='os.txt' AND size=417356323 AND entities=33879411;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='os.txt' AND size=417356323 AND entities=33879411;
COMMIT;

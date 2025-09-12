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
  'insidepro.txt', 7760097, 8, 20,
  9.27515583890253, 2.7777407349633516,
  22.306099, 13.505154, 7.23509,
  3.075026, 1.4137323, 1.271015, 0.17927611,
  'UTF-8', '-', 'leaked', 'alex stanev',
  79736356, 'https://sec.stanev.org/dict/insidepro.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='insidepro.txt' AND size=79736356 AND entities=7760097
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='insidepro.txt' AND size=79736356 AND entities=7760097;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='insidepro.txt' AND size=79736356 AND entities=7760097;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='insidepro.txt' AND size=79736356 AND entities=7760097;
COMMIT;

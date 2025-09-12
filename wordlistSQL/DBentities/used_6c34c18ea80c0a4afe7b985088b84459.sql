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
  'used.txt', 8983644, 7, 24,
  9.87205459165568, 2.833652769027494,
  60.746395, 2.1109364, 2.985982,
  6.0473347, 4.2073574, 0.60654676, 0.7459334,
  'UTF-8', '-', 'leaked', 'alex stanev',
  97696882, 'https://sec.stanev.org/dict/used.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='used.txt' AND size=97696882 AND entities=8983644
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='used.txt' AND size=97696882 AND entities=8983644;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='used.txt' AND size=97696882 AND entities=8983644;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='used.txt' AND size=97696882 AND entities=8983644;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'leet' FROM wordlists WHERE name='used.txt' AND size=97696882 AND entities=8983644;
COMMIT;

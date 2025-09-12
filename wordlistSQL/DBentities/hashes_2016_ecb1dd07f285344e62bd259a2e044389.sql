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
  'hashes_2016.txt', 78829276, 4, 32,
  10.185212915059628, 2.8530238126745826,
  53.01148, 1.8167728, 1.251252,
  10.8808775, 3.406667, 0.43780938, 1.9061737,
  'UTF-8', '-', 'leaked', 'alex stanev',
  882972793, 'https://sec.stanev.org/dict/hashes_2016.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='hashes_2016.txt' AND size=882972793 AND entities=78829276
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='hashes_2016.txt' AND size=882972793 AND entities=78829276;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='hashes_2016.txt' AND size=882972793 AND entities=78829276;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='hashes_2016.txt' AND size=882972793 AND entities=78829276;
COMMIT;

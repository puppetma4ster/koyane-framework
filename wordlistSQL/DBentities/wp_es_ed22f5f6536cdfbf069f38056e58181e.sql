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
  'wp_es.txt', 1521867, 8, 24,
  10.675829096760754, 2.8711850110293944,
  0.0, 0.0, 0.0,
  0.0, 0.0, 0.0, 0.0,
  'ISO-8859-1', 'spanish', 'dictionary', 'alex stanev',
  17769059, 'https://sec.stanev.org/dict/wp_es.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='wp_es.txt' AND size=17769059 AND entities=1521867
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='wp_es.txt' AND size=17769059 AND entities=1521867;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='wp_es.txt' AND size=17769059 AND entities=1521867;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='wp_es.txt' AND size=17769059 AND entities=1521867;
COMMIT;

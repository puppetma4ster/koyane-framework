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
  'wp_ru.txt', 2570308, 8, 24,
  11.513093761525855, 2.991965866824146,
  19.582829, 0.0, 0.0,
  0.0, 0.0, 0.0, 0.0,
  'ISO-8859-1', 'russian', 'dictionary', 'alex stanev',
  32162505, 'https://sec.stanev.org/dict/wp_ru.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='wp_ru.txt' AND size=32162505 AND entities=2570308
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='wp_ru.txt' AND size=32162505 AND entities=2570308;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='wp_ru.txt' AND size=32162505 AND entities=2570308;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='wp_ru.txt' AND size=32162505 AND entities=2570308;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'leet' FROM wordlists WHERE name='wp_ru.txt' AND size=32162505 AND entities=2570308;
COMMIT;

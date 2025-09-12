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
  'cow_ud_pc.txt', 1498112, 8, 24,
  9.399546228853383, 2.771255584835752,
  4.672214, 0.6442109, 3.5489335,
  1.2596521, 0.2150707, 0.20445734, 0.020692712,
  'ISO-8859-1', 'mixed', 'leaked', 'alex stanev',
  15579730, 'https://sec.stanev.org/dict/cow_ud_pc.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='cow_ud_pc.txt' AND size=15579730 AND entities=1498112
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='cow_ud_pc.txt' AND size=15579730 AND entities=1498112;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='cow_ud_pc.txt' AND size=15579730 AND entities=1498112;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='cow_ud_pc.txt' AND size=15579730 AND entities=1498112;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'leet' FROM wordlists WHERE name='cow_ud_pc.txt' AND size=15579730 AND entities=1498112;
COMMIT;

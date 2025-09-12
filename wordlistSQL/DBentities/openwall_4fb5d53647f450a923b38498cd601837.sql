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
  'openwall.txt', 1107279, 8, 24,
  11.878503972350238, 2.9601478125652316,
  0.00614118, 2.5981708, 14.342907,
  0.0069539836, 0.002077164, 5.42212, 0.0026190327,
  'UTF-8', 'mixed', 'leaked', 'alex stanev',
  14373782, 'https://sec.stanev.org/dict/openwall.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='openwall.txt' AND size=14373782 AND entities=1107279
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='openwall.txt' AND size=14373782 AND entities=1107279;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='openwall.txt' AND size=14373782 AND entities=1107279;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='openwall.txt' AND size=14373782 AND entities=1107279;
COMMIT;

PRAGMA foreign_keys = ON;
BEGIN;

-- Insert wordlist if not present
INSERT INTO wordlists (
  name, entities, smallestEntity, biggestEntity,
  averageLength, averageEntropy,
  digitsPercent, upperCasePercent, specialCharPercent,
  digitAndUpperCase, digitAndSpecialChar, upperCaseAndSpecialChar, digitUpperCaseAndSpecialChar,
  encoding, language, category, author,
  size, link, info
)
SELECT
  'realhuman_phill.txt', 63841069, 1, 200,
  10.205027472206012, 2.8027892609038427,
  19.283113, 23.60109, 7.4610486,
  17.081095, 1.0056912, 7.623458, 1.041231,
  'ISO-8859-1', '-', 'leaked', 'taylor hornby (crackstation)',
  716441107, 'https://crackstation.net/files/crackstation-human-only.txt.gz', '-'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='realhuman_phill.txt' AND size=716441107 AND entities=63841069
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'passwords' FROM wordlists WHERE name='realhuman_phill.txt' AND size=716441107 AND entities=63841069;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'hashcracking' FROM wordlists WHERE name='realhuman_phill.txt' AND size=716441107 AND entities=63841069;
COMMIT;

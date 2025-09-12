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
  'hashes_2015.txt', 156622452, 7, 24,
  10.79952019905805, 2.9877566578795074,
  36.535732, 7.4178576, 6.5135336,
  20.727243, 2.9556084, 1.4376515, 2.9328613,
  'UTF-8', '-', 'leaked', 'alex stanev',
  1848109783, 'https://sec.stanev.org/dict/hashes_2015.txt.gz'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='hashes_2015.txt' AND size=1848109783 AND entities=156622452
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wifi' FROM wordlists WHERE name='hashes_2015.txt' AND size=1848109783 AND entities=156622452;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa' FROM wordlists WHERE name='hashes_2015.txt' AND size=1848109783 AND entities=156622452;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'wpa2' FROM wordlists WHERE name='hashes_2015.txt' AND size=1848109783 AND entities=156622452;
COMMIT;

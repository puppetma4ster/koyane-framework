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
  'cirt-default-usernames.txt', 828, 1, 26,
  6.664251207729468, 2.3548963230093083,
  3.6231883, 44.082127, 4.227053,
  6.1594205, 0.24154589, 6.0386477, 1.0869565,
  'ISO-8859-1', '-', 'collection', '-',
  6346, 'https://crackstation.net/files/crackstation-human-only.txt.gz', 'found in github repo https://github.com/danielmiessler/SecLists?tab=readme-ov-file'
WHERE NOT EXISTS (
  SELECT 1 FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828
);

INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'usernames' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'web' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'cms' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'ftp' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'ssh' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'windows' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'linux' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'usernames' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'default-creds' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'credentials' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'discovery' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'scanner' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'string' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
INSERT OR IGNORE INTO wordlist_tags (wordlist_id, tag) SELECT id, 'text' FROM wordlists WHERE name='cirt-default-usernames.txt' AND size=6346 AND entities=828;
COMMIT;

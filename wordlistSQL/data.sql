/*
categories:
    leaked,
    generated,
    Dictionaries,
    frequent,
    rule

tags:
    passwords,
    wifi,
    web fuzzing,
    names,

    small,
    medium,
    big,
    huge,

    md5,
    sha1,
    ntlm,

 */
INSERT INTO wordlists
    (name, entities, smallestEntity, biggestEntity, averageLength, averageEntropy, digitsPercent, upperCasePercent, specialCharPercent,
     digitAndUpperCase, digitAndSpecialChar, upperCaseAndSpecialChar, digitUpperCaseAndSpecialChar,
     encoding, language, category, author,size, link)
VALUES
    ('rockyou.txt',
    14344391,
     1,
     240,
     8.44,
     2.68,
     58.73,
     2.72,
     2.65,
     5.51,
     3.28,
     0.53,
     0.56,
     'Latin1',
     'mixed',
     'leaked',
     'unknown',
     139921506,
     'https://github.com/dw0rsec/rockyou.txt/blob/master/rockyou.txt.zip');

INSERT INTO wordlist_tags(wordlist_id, tag) VALUES (1, 'password');
INSERT INTO wordlist_tags(wordlist_id, tag) VALUES (1, 'medium');


INSERT INTO wordlists
(name, entities, smallestEntity, biggestEntity, averageLength, averageEntropy, digitsPercent, upperCasePercent, specialCharPercent,
 digitAndUpperCase, digitAndSpecialChar, upperCaseAndSpecialChar, digitUpperCaseAndSpecialChar,
 encoding, language, category, author, size, link)
VALUES
    ('cracked.txt',
     7642500,
     1,
     70,
     10.00,
     2.72,
     68.49,
     3.53,
     0.66,
     14.15,
     3.40,
     0.37,
     2.00,
     'Latin-1',
     'mixed',
     'frequent',
     'alex stanev',
     7642500,
     'https://wpa-sec.stanev.org/dict/cracked.txt.gz');

INSERT INTO wordlist_tags(wordlist_id, tag) VALUES (1, 'wifi');
INSERT INTO wordlist_tags(wordlist_id, tag) VALUES (1, 'small');


INSERT INTO wordlists
(name, entities, smallestEntity, biggestEntity, averageLength, averageEntropy, digitsPercent, upperCasePercent, specialCharPercent,
 digitAndUpperCase, digitAndSpecialChar, upperCaseAndSpecialChar, digitUpperCaseAndSpecialChar,
 encoding, language, category, author, size, link)
VALUES
    ('num8_1.txt',
     8822994,
     8,
     8,
     8.00,
     2.38,
     100.00,
     0.00,
     0.00,
     0.00,
     0.00,
     0.00,
     0.00,
     'UTF-8',
     '-',
     'numeric',
     'alex stanev',
     79406946,
     'https://sec.stanev.org/dict/num8_1.txt.gz');

INSERT INTO wordlist_tags(wordlist_id, tag) VALUES (1, 'wifi');
INSERT INTO wordlist_tags(wordlist_id, tag) VALUES (1, 'medium');



SELECT * FROM wordlists;
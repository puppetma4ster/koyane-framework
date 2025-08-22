-- Tbuild table
CREATE TABLE wordlists (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    entities INTEGER NOT NULL,
    smallestEntity INTEGER NOT NULL,
    biggestEntity INTEGER NOT NULL,
    averageLength REAL NOT NULL,
    averageEntropy REAL NOT NULL,
    digitsPercent REAL NOT NULL,
    upperCasePercent REAL NOT NULL,
    specialCharPercent REAL NOT NULL,
    digitAndUpperCase REAL NOT NULL,
    digitAndSpecialChar REAL NOT NULL,
    upperCaseAndSpecialChar REAL NOT NULL,
    digitUpperCaseAndSpecialChar REAL NOT NULL,
    encoding TEXT NOT NULL,
    language TEXT NOT NULL,
    category TEXT NOT NULL,
    author TEXT NOT NULL,
    size INTEGER NOT NULL,
    link TEXT NOT NULL
);

DROP TABLE IF EXISTS wordlist_tags;
CREATE TABLE wordlist_tags (
    wordlist_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
    FOREIGN KEY(wordlist_id) REFERENCES wordlists(id)
);

CREATE TABLE IF NOT EXISTS comments
(id INTEGER PRIMARY KEY, key TEXT, username TEXT, content TEXT, time INTEGER);

CREATE TABLE IF NOT EXISTS url_aliases
(id INTEGER PRIMARY KEY, alias TEXT UNIQUE, value TEXT UNIQUE)
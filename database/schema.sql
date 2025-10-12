CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userName TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS posts(
     id INTEGER PRIMARY KEY AUTOINCREMENT,
     title TEXT NOT NULL,
     post TEXT NOT NULL,
     imageUrl TEXT,
     userId INTEGER,
     creationDate TEXT,
     FOREIGN KEY(userId) REFERENCES users(id)  
);
CREATE TABLE IF NOT EXISTS categories(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category TEXT NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS postCategories (
    postId INTEGER NOT NULL,
    categoryId INTEGER NOT NULL,
    PRIMARY KEY (postId, categoryId),
    FOREIGN KEY (postId) REFERENCES posts(id),
    FOREIGN KEY (categoryId) REFERENCES categories(id)
);
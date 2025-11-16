CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userName TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    session TEXT DEFAULT NULL,
    dateexpired DATETIME DEFAULT NULL
);
CREATE TABLE IF NOT EXISTS posts(
     id INTEGER PRIMARY KEY AUTOINCREMENT,
     title TEXT NOT NULL,
     post TEXT NOT NULL,
     imageUrl TEXT,
     userId INTEGER NOT NULL,
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
CREATE TABLE IF NOT EXISTS comments(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    comment TEXT NOT NULL,
    postId INTEGER NOT NULL,
    userId INTEGER NOT NULL,
    creationDate TEXT,
    FOREIGN KEY(userId) REFERENCES users(id),
    FOREIGN KEY(postId) REFERENCES posts(id)
);
CREATE TABLE IF NOT EXISTS postReactions(
    userId INTEGER NOT NULL,
    postId INTEGER NOT NULL,
    reaction INTEGER NOT NULL,
    PRIMARY KEY(userId,postId),
    FOREIGN KEY(userId) REFERENCES users(id),
    FOREIGN KEY(postId) REFERENCES posts(id)
); 

CREATE TABLE IF NOT EXISTS commentReactions(
    userId INTEGER NOT NULL,
    commentId INTEGER NOT NULL,
    reaction INTEGER NOT NULL,
    PRIMARY KEY(userId,commentId),
    FOREIGN KEY(userId) REFERENCES users(id),
    FOREIGN KEY(commentId) REFERENCES comments(id)
);

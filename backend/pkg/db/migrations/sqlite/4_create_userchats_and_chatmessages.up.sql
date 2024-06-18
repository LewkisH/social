CREATE TABLE IF NOT EXISTS user_chats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user1 INTEGER NOT NULL,
    user2 INTEGER NOT NULL,
    last_message TEXT NOT NULL,
    created_at DATE DEFAULT CURRENT_DATE,
    FOREIGN KEY(user1) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(user2) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chatmessages (
    id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    user_id INTEGER NOT NULL, 
    chat_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    is_group BOOLEAN,
    seen BOOLEAN,
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (chat_id) REFERENCES user_chats(id) ON DELETE CASCADE
);

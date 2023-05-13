		CREATE TABLE IF NOT EXISTS reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			comment_id INTEGER,
			user_id INTEGER NOT NULL,
			type TEXT NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (comment_id) REFERENCES comments(id),
			FOREIGN KEY (user_id) REFERENCES users(id),
		);`
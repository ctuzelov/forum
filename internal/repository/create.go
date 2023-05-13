package repository

import "database/sql"

func CreateTables(DB *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
			token TEXT,
			expires DATETIME
		);
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			parent_id INTEGER DEFAULT 0,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			comment_id INTEGER,
			user_id INTEGER NOT NULL,
			type TEXT NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (comment_id) REFERENCES comments(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category TEXT NOT NULL
			);
		CREATE TABLE IF NOT EXISTS post_cat (
			cat_id INTEGER,
			post_id INTEGER,
			FOREIGN KEY (cat_id) REFERENCES categories(id),			
			FOREIGN KEY (post_id) REFERENCES posts(id)			
		);`
	_, err := DB.Exec(query)
	return err
}

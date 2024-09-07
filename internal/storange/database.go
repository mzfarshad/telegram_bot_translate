package storange

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the database and creates the necessary tables if they don't exist.

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./bot.db")
	if err != nil {
		return fmt.Errorf("failed to create db: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bot_setting(
	user_id INTEGER,
	chat_id INTEGER,
	bot_language TEXT,
	PRIMARY KEY(user_id, chat_id)
	);`)

	if err != nil {
		return fmt.Errorf("failed to create bot_setting table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS menu_state(
	chat_id INTEGER,
	user_id INTEGER,
	state TEXT,
	PRIMARY KEY(user_id, chat_id)
	);`)
	if err != nil {
		return fmt.Errorf("failed to create menu_state table: %v", err)
	}

	// _, err = db.Exec(`DROP TABLE IF EXISTS translation`)
	// if err != nil {
	// 	return fmt.Errorf("failed to drop existing translation table: %v", err)
	// }

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS translation (
    user_id INTEGER,
    sent_message BOOLEAN DEFAULT FALSE,
    sent_source_language TEXT,
    sent_target_language TEXT,
	active_translation BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (user_id)
	);`)
	if err != nil {
		return fmt.Errorf("failed to create translation table: %v", err)
	}

	return nil
}

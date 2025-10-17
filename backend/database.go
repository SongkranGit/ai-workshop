package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT,
			phone TEXT,
			avatar_url TEXT,
			bio TEXT,
			points_balance INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE TABLE IF NOT EXISTS transfers (
			transfer_id INTEGER PRIMARY KEY AUTOINCREMENT,
			from_user_id INTEGER NOT NULL,
			to_user_id INTEGER NOT NULL,
			amount INTEGER NOT NULL CHECK (amount > 0),
			status TEXT NOT NULL CHECK (status IN ('pending','processing','completed','failed','cancelled','reversed')),
			note TEXT,
			idempotency_key TEXT NOT NULL UNIQUE,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			completed_at DATETIME,
			fail_reason TEXT,
			FOREIGN KEY (from_user_id) REFERENCES users(id),
			FOREIGN KEY (to_user_id) REFERENCES users(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_transfers_from ON transfers(from_user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transfers_to ON transfers(to_user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transfers_created ON transfers(created_at)`,
		`CREATE TABLE IF NOT EXISTS point_ledger (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			change INTEGER NOT NULL,
			balance_after INTEGER NOT NULL,
			event_type TEXT NOT NULL CHECK (event_type IN ('transfer_out','transfer_in','adjust','earn','redeem')),
			transfer_id INTEGER,
			reference TEXT,
			metadata TEXT,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (transfer_id) REFERENCES transfers(transfer_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_ledger_user ON point_ledger(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ledger_transfer ON point_ledger(transfer_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ledger_created ON point_ledger(created_at)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return err
		}
	}

	return nil
}

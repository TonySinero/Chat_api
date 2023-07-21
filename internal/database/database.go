package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
)

type PostgresDB struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(database PostgresDB) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		database.Username, database.Password, database.Host, database.Port, database.DBName, database.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %s", err)
	}

	err = db.Ping()
	if err != nil {
		logrus.Errorf("DB ping error: %s", err)
		return nil, err
	}

	err = createDatabaseTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDatabaseTables(db *sql.DB) error {
	schemaQueries := []string{
		USER_SCHEMA,
		GROUP_CHATS_SCHEMA,
		GROUP_CHAT_MEMBERS_SCHEMA,
		GROUP_CHAT_MESSAGES_SCHEMA,
		PRIVATE_MESSAGES_SCHEMA,
		MESSAGE_SCHEMA,
	}

	for _, schemaQuery := range schemaQueries {
		_, err := db.Exec(schemaQuery)
		if err != nil {
			return err
		}
	}

	return nil
}

const (
	GROUP_CHATS_SCHEMA = `
		CREATE TABLE IF NOT EXISTS group_chats (
			id serial not null primary key,
			name varchar(225) NOT NULL,
			description varchar(225) NOT NULL,
			created_at timestamp,
			updated_at timestamp
		);
	`

	GROUP_CHAT_MEMBERS_SCHEMA = `
	    CREATE TABLE IF NOT EXISTS group_chat_members (
			id serial not null primary key,
			group_chat_id int NOT NULL,
			user_id int NOT NULL,
			FOREIGN KEY (group_chat_id) REFERENCES group_chats (id),
			FOREIGN KEY (user_id) REFERENCES users (id)
	    );
    `

	GROUP_CHAT_MESSAGES_SCHEMA = `
		CREATE TABLE IF NOT EXISTS group_chat_messages (
			id serial not null primary key,
			sender_id int NOT NULL,
			group_chat_id int NOT NULL,
			content varchar(225) NOT NULL,
			created_at timestamp,
			FOREIGN KEY (sender_id) REFERENCES users (id),
			FOREIGN KEY (group_chat_id) REFERENCES group_chats (id)
		);
	`

	PRIVATE_MESSAGES_SCHEMA = `
		CREATE TABLE IF NOT EXISTS private_messages (
			id serial not null primary key,
			sender_id int NOT NULL,
			recipient_id int NOT NULL,
			content varchar(225) NOT NULL,
			created_at timestamp,
			FOREIGN KEY (sender_id) REFERENCES users (id),
			FOREIGN KEY (recipient_id) REFERENCES users (id)
		);
	`

	MESSAGE_SCHEMA = `
		CREATE TABLE IF NOT EXISTS messages (
			id serial not null primary key,
			username varchar(225) NOT NULL,
			content varchar(225) NOT NULL
		);
	`
	USER_SCHEMA = `
		CREATE TABLE IF NOT EXISTS users
		(
			id serial not null primary key,
			name varchar(225) not null,
			email varchar(225) not null UNIQUE,
			phone varchar(225) not null UNIQUE,
			password varchar(225) not null,
			role varchar(225) not null default 'USER',
			CONSTRAINT proper_email CHECK (email ~* '^[A-Za-z0-9._+%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$')
		);
	`
)

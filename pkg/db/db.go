package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" //postgres drivers
	"like-service/config"
)

func ConnectToDB(cfg config.Config) (*sql.DB, func(), error) {
	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	connDb, err := sql.Open("postgres", psqlString)
	if err != nil {
		return nil, nil, err
	}

	cleanUpFunc := func() {
		connDb.Close()
	}
	return connDb, cleanUpFunc, nil
}

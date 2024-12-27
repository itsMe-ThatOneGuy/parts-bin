package state

import (
	"database/sql"
	"errors"
	"os"

	"github.com/itsMe-ThatOneGuy/parts-bin/internal/config"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
)

func (s *State) InitConfig() error {
	cfg, err := config.Read()
	if err != nil {
		return err
	}

	s.Config = &cfg

	return nil
}

func (s *State) InitDB() error {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = s.Config.DBUrl
		if dbURL == "" || dbURL == config.DefaultDBurl {
			return errors.New("db_url not set in either .env or ~/.partsbinconfig.json")
		}
	}

	dbCon, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	s.DB = dbCon
	s.DBQueries = database.New(dbCon)
	return nil
}

func (s *State) CloseDB() {
	if s.DB != nil {
		s.DB.Close()
	}
}

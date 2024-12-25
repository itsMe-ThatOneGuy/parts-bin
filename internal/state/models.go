package state

import (
	"database/sql"

	"github.com/itsMe-ThatOneGuy/parts-bin/internal/config"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
)

type State struct {
	Config    *config.Config
	DBQueries *database.Queries
	DB        *sql.DB
}

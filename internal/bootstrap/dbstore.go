package bootstrap

import (
	repo "atm-test/internal/drivers/dbstore"
	"atm-test/internal/pkg/logger"
	"database/sql"
)

type Repos struct {
	repos *repo.Repo
}

func initRepos(db *sql.DB, log logger.Logger) Repos {
	return Repos{
		repos: repo.New(db, log),
	}
}

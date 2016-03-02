package initdb

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/constants"
	"github.com/loopfz/scecret/models"
	_ "github.com/mattn/go-sqlite3"
)

func PopulateDbMap(db *gorp.DbMap) error {

	db.AddTableWithName(models.Photo{}, `photo`).SetKeys(true, "id")
	db.AddTableWithName(models.Album{}, `album`).SetKeys(true, "id")
	
	return db.CreateTablesIfNotExists()
}

func InitPostgres() (*gorp.DbMap, error) {
	sqldb, err := sql.Open("postgres", fmt.Sprintf("dbname=%s", constants.DBName))
	if err != nil {
		return nil, err
	}

	db := &gorp.DbMap{Db: sqldb, Dialect: gorp.PostgresDialect{}}

	err = PopulateDbMap(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitSqlite() (*gorp.DbMap, error) {
	return doInitSqlite(false)
}

func InitSqliteRandom() (*gorp.DbMap, error) {
	return doInitSqlite(true)
}

func doInitSqlite(random bool) (*gorp.DbMap, error) {

	var s string
	if random {
		rand.Seed(time.Now().UTC().UnixNano())
		s = fmt.Sprintf("/tmp/%s%d.db", constants.DBName, rand.Int())
	} else {
		s = fmt.Sprintf("/tmp/%s.db", constants.DBName) // TODO put that somewhere else/constants
	}
	sqldb, err := sql.Open("sqlite3", s)
	if err != nil {
		return nil, err
	}

	db := &gorp.DbMap{Db: sqldb, Dialect: gorp.SqliteDialect{}}

	err = PopulateDbMap(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

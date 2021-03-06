package initdb

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/elektroid/golivia/constants"
	"github.com/elektroid/golivia/models"
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

func InitSqlite(databaseDir string) (*gorp.DbMap, error) {
	return doInitSqlite(databaseDir, false)
}

func InitSqliteRandom(databaseDir string) (*gorp.DbMap, error) {
	return doInitSqlite(databaseDir, true)
}

func doInitSqlite(databaseDir string, random bool) (*gorp.DbMap, error) {

	var s string
	if random {
		rand.Seed(time.Now().UTC().UnixNano())
		s = fmt.Sprintf("%s/%s%d.db", databaseDir, constants.DBName, rand.Int())
	} else {
		s = fmt.Sprintf("%s/%s.db", databaseDir, constants.DBName) // TODO put that somewhere else/constants
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

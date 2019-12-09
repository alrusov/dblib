package dblib

import (
	"fmt"

	"github.com/alrusov/config"
	"github.com/alrusov/misc"
)

//----------------------------------------------------------------------------------------------------------------------------//

type (
	// DBext --
	DBext struct {
		db      *DB
		tuner   queryTuner
		queries misc.StringMap
	}

	queryTuner   func(query string, options tunerOptions) (newQuery string, err error)
	tunerOptions misc.InterfaceMap
)

//----------------------------------------------------------------------------------------------------------------------------//

// DB --
func (me *DBext) DB() *DB {
	return me.db
}

//----------------------------------------------------------------------------------------------------------------------------//

// NewDBext --
func NewDBext(cfg *config.DB, queries misc.StringMap) (dbExt *DBext, err error) {
	dbExt = &DBext{
		db:      &DB{},
		queries: queries,
	}
	err = dbExt.db.Init(cfg.Type, cfg.DSN, cfg.Timeout, cfg.Retry)
	if err != nil {
		return
	}

	dbExt.tuner = emptyTuner

	switch cfg.Type {
	case "mysql":
		dbExt.tuner = mysqlTuner
	case "pgsql":
		dbExt.tuner = pgsqlTuner
	case "mssql":
		dbExt.tuner = mssqlTuner
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

// get --
func (me *DBext) getQuery(name string) (string, error) {
	q, exists := me.queries[name]

	if !exists {
		return "", fmt.Errorf(`SQL query "%s" not found`, name)
	}

	return q, nil
}

//----------------------------------------------------------------------------------------------------------------------------//

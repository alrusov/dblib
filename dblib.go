package dblib

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/alrusov/misc"
)

//----------------------------------------------------------------------------------------------------------------------------//

const (
	// MYSQL --
	MYSQL = "mysql"
	// PGSQL --
	PGSQL = "pgsql"
	// MSSQL --
	MSSQL = "mssql"
)

//----------------------------------------------------------------------------------------------------------------------------//

// DB --
type DB struct {
	driver          string
	dsn             string
	db              *sql.DB
	maxRetry        int
	delayAfterError time.Duration
}

//----------------------------------------------------------------------------------------------------------------------------//

var (
	replacers = map[string]*strings.Replacer{
		MYSQL: strings.NewReplacer(`'`, `\'`, `\`, `\\`),
		PGSQL: strings.NewReplacer(`'`, `\'`, `\`, `\\`),
		MSSQL: strings.NewReplacer(`'`, `''`),
	}
)

// SecureString --
func (me *DB) SecureString(s string) string {
	r, exists := replacers[me.driver]
	if !exists {
		return s
	}
	return r.Replace(s)
}

//----------------------------------------------------------------------------------------------------------------------------//

var (
	errFilterRE      = regexp.MustCompile(`(\sdsn\s*=\s*".*://.*:)(.*)(@.*")`)
	errFilterReplace = `$1*$3`
)

func (me *DB) execSQL(withResult bool, query string, secured bool, prm ...interface{}) (result interface{}, err error) {
	if !secured {
		for i, v := range prm {
			switch v.(type) {
			case string:
				prm[i] = me.SecureString(prm[i].(string))
			}
		}
	}

	preparedQuery := strings.TrimSpace(fmt.Sprintf(query, prm...))

	try := 0
	for {
		if withResult {
			result, err = me.db.Query(preparedQuery)
		} else {
			result, err = me.db.Exec(preparedQuery)
		}

		if err == nil {
			return
		}

		msg := errFilterRE.ReplaceAllString(err.Error(), errFilterReplace)
		if !strings.Contains(err.Error(), " connectex:") {
			// Все ошибки, кроме соединения
			err = fmt.Errorf(`Exec query (%s): "%s"`, preparedQuery, msg)
			return
		}

		// Ошибка соединения
		try++
		err = fmt.Errorf(`Exec query (%s): "%s", was try %d from %d`, preparedQuery, msg, try, me.maxRetry)
		if try == me.maxRetry {
			return
		}

		misc.Sleep(me.delayAfterError)
	}
}

//----------------------------------------------------------------------------------------------------------------------------//

// OpenRecordset --
func (me *DB) OpenRecordset(query string, secured bool, prm ...interface{}) (result *sql.Rows, err error) {
	var r interface{}
	r, err = me.execSQL(true, query, secured, prm...)
	if err == nil {
		var ok bool
		result, ok = r.(*sql.Rows)
		if !ok {
			result = nil
			err = errors.New("Bad result type")
		}
	} else if result != nil {
		result.Close()
		result = nil
	}
	return result, err
}

// Exec --
func (me *DB) Exec(query string, secured bool, prm ...interface{}) (result sql.Result, err error) {
	var r interface{}
	r, err = me.execSQL(false, query, secured, prm...)
	if err == nil {
		var ok bool
		result, ok = r.(sql.Result)
		if !ok {
			result = nil
			err = errors.New("Bad result type")
		}
	}

	return result, err
}

//----------------------------------------------------------------------------------------------------------------------------//

func exit(code int, p interface{}) {
	db, ok := p.(*DB)
	if ok && (db != nil) && (db.db != nil) {
		db.db.Close()
	}
}

// Init --
func (me *DB) Init(driver string, dsn string, maxConn int, maxRetry int) error {
	me.driver = driver
	me.dsn = dsn
	me.maxRetry = maxRetry
	me.delayAfterError = 3 * time.Second

	db, err := sql.Open(me.driver, me.dsn)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(maxConn)
	me.db = db
	misc.AddExitFunc("dblib.exit", exit, me)

	return nil
}

//----------------------------------------------------------------------------------------------------------------------------//

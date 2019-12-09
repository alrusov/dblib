package dblib

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alrusov/misc"
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
		"mysql": strings.NewReplacer(`'`, `\'`, `\`, `\\`),
		"pgsql": strings.NewReplacer(`'`, `\'`, `\`, `\\`),
		"mssql": strings.NewReplacer(`'`, `''`),
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

func (me *DB) execSQL(isOpen bool, query string, secured bool, prm ...interface{}) (result interface{}, err error) {

	if !secured {
		for i, v := range prm {
			switch v.(type) {
			case string:
				prm[i] = me.SecureString(prm[i].(string))
			}
		}
	}

	preparedQuery := strings.TrimSpace(fmt.Sprintf(query, prm...))

	for i := 0; i < me.maxRetry; i++ {
		if isOpen {
			result, err = me.db.Query(preparedQuery)
		} else {
			result, err = me.db.Exec(preparedQuery)
		}

		if err == nil {
			break
		} else {
			//if (err == mysql.ErrInvalidConn) || strings.Contains(err.Error(), " connectex:") {
			if strings.Contains(err.Error(), " connectex:") {
				err = fmt.Errorf("Exec query (%s) error (%s), try %d from %d", preparedQuery, err.Error(), i+1, me.maxRetry)
				if i < me.maxRetry-1 {
					misc.Sleep(me.delayAfterError)
				}
			} else {
				err = fmt.Errorf("Exec query (%s) error (%s)", preparedQuery, err.Error())
				break
			}
		}
	}

	return result, err
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

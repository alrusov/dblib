package dblib

import "database/sql"

//----------------------------------------------------------------------------------------------------------------------------//

func (me *DBext) exec(name string, secured bool, params ...interface{}) (sql.Result, error) {
	query, err := me.getQuery(name)
	if err != nil {
		return nil, err
	}

	result, err := me.db.Exec(query, secured, params...)
	if err != nil {
		return nil, err
	}

	return result, err
}

// Exec --
func (me *DBext) Exec(name string, params ...interface{}) (sql.Result, error) {
	return me.exec(name, false, params...)
}

// ExecWithoutSecuring --
func (me *DBext) ExecWithoutSecuring(name string, params ...interface{}) (sql.Result, error) {
	return me.exec(name, true, params...)
}

//----------------------------------------------------------------------------------------------------------------------------//

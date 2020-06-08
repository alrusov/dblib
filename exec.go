package dblib

import "database/sql"

//----------------------------------------------------------------------------------------------------------------------------//

// Exec --
func (me *DBext) Exec(name string, secured bool, params ...interface{}) (sql.Result, error) {
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

//----------------------------------------------------------------------------------------------------------------------------//

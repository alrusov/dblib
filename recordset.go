package dblib

import (
	"database/sql"
)

//----------------------------------------------------------------------------------------------------------------------------//

func (me *DBext) openRecordset(name string, secured bool, offset uint, count uint, params ...interface{}) (*sql.Rows, error) {
	query, err := me.getQuery(name)
	if err != nil {
		return nil, err
	}

	query, err = me.tuner(query, tunerOptions{"offset": offset, "count": count})
	if err != nil {
		return nil, err
	}

	rs, err := me.db.OpenRecordset(query, secured, params...)
	if err != nil {
		return nil, err
	}

	return rs, err
}

// OpenRecordset -- fetching data
func (me *DBext) OpenRecordset(name string, params ...interface{}) (*sql.Rows, error) {
	return me.openRecordset(name, false, 0, 0, params...)
}

// OpenRecordsetExtended -- fetching data with providing a start position and max count
func (me *DBext) OpenRecordsetExtended(name string, offset uint, count uint, params ...interface{}) (*sql.Rows, error) {
	return me.openRecordset(name, false, offset, count, params...)
}

// OpenRecordsetWithoutSecuring -- fetching data without special symbols escaping
func (me *DBext) OpenRecordsetWithoutSecuring(name string, params ...interface{}) (*sql.Rows, error) {
	return me.openRecordset(name, true, 0, 0, params...)
}

// OpenRecordsetWithoutSecuringExtended -- fetching data with providing a start position, max count and without special symbols escaping
func (me *DBext) OpenRecordsetWithoutSecuringExtended(name string, offset uint, count uint, params ...interface{}) (*sql.Rows, error) {
	return me.openRecordset(name, true, offset, count, params...)
}

//----------------------------------------------------------------------------------------------------------------------------//

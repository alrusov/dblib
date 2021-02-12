package dblib

import (
	"net/http"

	"github.com/alrusov/misc"
)

//----------------------------------------------------------------------------------------------------------------------------//

// Exec -- execute stored procedure or other
func (me *DBext) Exec(id uint64, secured bool, name string, params ...interface{}) (result interface{}, code int, err error) {
	t0 := misc.NowUnixNano()
	defer func() {
		misc.LogProcessingTime(me.db.logFacility.Name(), "", id, "db.call", "", t0)
	}()

	query, err := me.getQuery(name)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	result, err = me.db.Exec(query, secured, params...)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

//----------------------------------------------------------------------------------------------------------------------------//

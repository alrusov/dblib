package dblib

import (
	"net/http"

	"github.com/alrusov/misc"
)

//----------------------------------------------------------------------------------------------------------------------------//

// Call -- execute stored procedure
func (me *DBext) Call(id uint64, name string, params ...interface{}) (_ interface{}, code int, err error) {
	t0 := misc.NowUTC().UnixNano()
	defer misc.LogProcessingTime("", id, "db.call", "", t0)

	_, err = me.Exec(name, params...)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil

}

//----------------------------------------------------------------------------------------------------------------------------//

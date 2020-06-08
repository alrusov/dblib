package dblib

import (
	"net/http"

	"github.com/alrusov/misc"
)

//----------------------------------------------------------------------------------------------------------------------------//

// Call -- execute stored procedure
func (me *DBext) Call(id uint64, secured bool, name string, params ...interface{}) (_ interface{}, code int, err error) {
	t0 := misc.NowUTC().UnixNano()
	defer func() {
		misc.LogProcessingTime(me.db.logFacility.Name(), "", id, "db.call", "", t0)
	}()

	_, err = me.Exec(name, secured, params...)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

//----------------------------------------------------------------------------------------------------------------------------//

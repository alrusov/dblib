package dblib

import (
	"net/http"

	"github.com/alrusov/misc"
)

//----------------------------------------------------------------------------------------------------------------------------//

func (me *DBext) call(withoutSecuring bool, id uint64, name string, params ...interface{}) (_ interface{}, code int, err error) {
	t0 := misc.NowUTC().UnixNano()
	defer func() {
		misc.LogProcessingTime(me.db.logFacility.Name(), "", id, "db.call", "", t0)
	}()

	if withoutSecuring {
		_, err = me.ExecWithoutSecuring(name, params...)
	} else {
		_, err = me.Exec(name, params...)
	}
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

// Call -- execute stored procedure
func (me *DBext) Call(id uint64, name string, params ...interface{}) (_ interface{}, code int, err error) {
	return me.call(false, id, name, params...)
}

// CallWithoutSecuring -- execute stored procedure
func (me *DBext) CallWithoutSecuring(id uint64, name string, params ...interface{}) (_ interface{}, code int, err error) {
	return me.call(true, id, name, params...)
}

//----------------------------------------------------------------------------------------------------------------------------//

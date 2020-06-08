package dblib

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/alrusov/misc"
)

//----------------------------------------------------------------------------------------------------------------------------//

type resultDef struct {
	mutex  *sync.Mutex
	row    interface{}
	fields []interface{}
}

var (
	resultDefCacheMutex = new(sync.RWMutex)
	resultDefCache      = make(map[string]*resultDef)
)

//----------------------------------------------------------------------------------------------------------------------------//

func fieldsFromPattern(rowPattern interface{}) (v []interface{}, err error) {

	// Parse a rowPattern structure

	fields := reflect.ValueOf(rowPattern).Elem()
	n := fields.NumField()

	v = make([]interface{}, n)

	fc := 0
	for i := 0; i < n; i++ {
		a := fields.Field(i).Addr()
		if a.CanInterface() {
			v[fc] = a.Interface()
			fc++
		}
	}

	if fc == 0 {
		err = fmt.Errorf(`No exported fields on the structure "%s"`, reflect.TypeOf(rowPattern))
		return
	}

	if fc < n {
		v = v[0:fc]
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

func makeRowPattern(rs *sql.Rows) (df *resultDef, err error) {

	cols, err := rs.ColumnTypes()
	if err != nil {
		return
	}

	n := len(cols)
	f := make([]reflect.StructField, n)
	names := make(map[string]int, n)

	for i := 0; i < n; i++ {
		name := cols[i].Name()

		n, exists := names[name]
		if !exists {
			n = 1
		} else {
			n++
		}
		names[name] = n

		var tp reflect.Type
		switch cols[i].ScanType() {
		case reflect.TypeOf(bool(false)):
			tp = reflect.TypeOf(NullBool{})

		case reflect.TypeOf(string("")):
			tp = reflect.TypeOf(NullString{})

		case reflect.TypeOf(int64(0)):
			tp = reflect.TypeOf(NullInt64{})

		case reflect.TypeOf(int32(0)):
			tp = reflect.TypeOf(NullInt32{})

		case reflect.TypeOf(float64(0)):
			tp = reflect.TypeOf(NullFloat64{})

		case reflect.TypeOf(time.Time{}):
			tp = reflect.TypeOf(NullTime{})

		case reflect.TypeOf([]uint8{}):
			tp = reflect.TypeOf(NullString{})

		}

		if n > 1 {
			name += strconv.Itoa(n)
		}
		f[i] = reflect.StructField{
			Name: "X" + name,
			Type: tp,
			Tag:  `json:"` + reflect.StructTag(name) + `"`,
		}
	}

	tp := reflect.StructOf(f)
	obj := reflect.New(tp)

	df = &resultDef{}
	df.row = obj.Interface()
	df.fields = make([]interface{}, n)

	fields := obj.Elem()
	for i := 0; i < n; i++ {
		df.fields[i] = fields.Field(i).Addr().Interface()
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

// Select -- Get data from table or view
func (me *DBext) Select(id uint64, secured bool, rowPattern interface{}, name string, cacheName string, offset uint, count uint, params ...interface{}) (result []interface{}, code int, err error) {
	t0 := misc.NowUTC().UnixNano()
	defer func() {
		misc.LogProcessingTime(me.db.logFacility.Name(), "", id, `db.select "`+name+`(`+cacheName+`)"`, "", t0)
	}()

	if cacheName == "" {
		cacheName = name
	}

	// Make a query

	rs, err := me.OpenRecordsetExtended(name, secured, offset, count, params...)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rs.Close()

	// Get result definition

	var df *resultDef

	if rowPattern != nil {
		// Prepare a list of fields from the pattern
		df = &resultDef{
			row: rowPattern,
		}
		df.fields, err = fieldsFromPattern(rowPattern)
	} else {
		exists := false
		resultDefCacheMutex.RLock()
		df, exists = resultDefCache[cacheName]
		resultDefCacheMutex.RUnlock()

		if !exists {
			// Make row Pattern and fields
			df, err = makeRowPattern(rs)
			df.mutex = new(sync.Mutex)

			// Save to cache for future use
			resultDefCacheMutex.Lock()
			resultDefCache[cacheName] = df
			resultDefCacheMutex.Unlock()
		}
	}

	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Read data and prepare answer

	answer := []interface{}{}

	if df.mutex != nil {
		df.mutex.Lock()
	}

	for rs.Next() {
		err = rs.Scan(df.fields...)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		row := reflect.Indirect(reflect.ValueOf(df.row)).Interface()
		answer = append(answer, row)
	}

	if df.mutex != nil {
		df.mutex.Unlock()
	}

	return answer, http.StatusOK, nil
}

//----------------------------------------------------------------------------------------------------------------------------//

// GetVals --
func GetVals(data interface{}, fieldNames []string) (values []interface{}, err error) {
	values = make([]interface{}, len(fieldNames))

	fields := reflect.ValueOf(data)

	for i, name := range fieldNames {
		f := fields.FieldByName("X" + name)
		if !f.IsValid() {
			return nil, fmt.Errorf(`Unknown field "%s"`, name)
		}

		v := f.Interface()
		var vv interface{}

		switch v.(type) {
		case NullBool:
			vv = v.(NullBool).Bool
		case NullString:
			vv = v.(NullString).String
		case NullInt64:
			vv = v.(NullInt64).Int64
		case NullFloat64:
			vv = v.(NullFloat64).Float64
		case NullTime:
			vv = v.(NullTime).Time
		default:
			return nil, fmt.Errorf(`Bad type of the field "%s"`, name)
		}

		values[i] = vv
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

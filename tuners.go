package dblib

import (
	"fmt"
)

//----------------------------------------------------------------------------------------------------------------------------//

func (o tunerOptions) getUintOption(name string) (uint, error) {
	i, exists := o[name]
	if !exists {
		return 0, nil
	}

	v, ok := i.(uint)
	if !ok {
		return 0, fmt.Errorf(`Option "%s" has type %T, expected %T`, name, i, v)
	}

	return v, nil
}

func (o tunerOptions) getOffsetAndCount() (offset uint, count uint, err error) {
	offset, err = o.getUintOption("offset")
	if err != nil {
		return
	}

	count, err = o.getUintOption("count")
	if err != nil {
		return
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

func emptyTuner(query string, o tunerOptions) (string, error) {
	return "", fmt.Errorf("tuner not defined")
}

//----------------------------------------------------------------------------------------------------------------------------//

func mysqlTuner(query string, o tunerOptions) (string, error) {
	offset, count, err := o.getOffsetAndCount()
	if err != nil {
		return "", err
	}

	if offset > 0 {
		query += fmt.Sprintf(` LIMIT %d OFFSET %d`, count, offset)
	} else if count > 0 {
		query += fmt.Sprintf(` LIMIT %d`, count)
	}

	return query, nil
}

//----------------------------------------------------------------------------------------------------------------------------//

func pgsqlTuner(query string, o tunerOptions) (string, error) {
	return mysqlTuner(query, o)
}

//----------------------------------------------------------------------------------------------------------------------------//

func mssqlTuner(query string, o tunerOptions) (string, error) {
	offset, count, err := o.getOffsetAndCount()
	if err != nil {
		return "", err
	}

	if count > 0 {
		query += fmt.Sprintf(` OFFSET %d ROWS FETCH NEXT %d ROWS ONLY`, offset, count)
	} else if offset > 0 {
		query += fmt.Sprintf(` OFFSET %d ROWS`, offset)
	}

	return query, nil
}

//----------------------------------------------------------------------------------------------------------------------------//

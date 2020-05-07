package dblib

import (
	"database/sql"
	"encoding/json"
)

//----------------------------------------------------------------------------------------------------------------------------//

type (
	// NullBool --
	NullBool sql.NullBool

	// NullString --
	NullString sql.NullString

	// NullInt64 --
	NullInt64 sql.NullInt64

	// NullInt32 --
	NullInt32 sql.NullInt32

	// NullFloat64 --
	NullFloat64 sql.NullFloat64

	// NullTime --
	NullTime sql.NullTime
)

//----------------------------------------------------------------------------------------------------------------------------//

// Scan --
func (v *NullBool) Scan(value interface{}) error {
	return (*sql.NullBool)(v).Scan(value)
}

// MarshalJSON for NullBool
func (v NullBool) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(v.Bool)
}

//----------------------------------------------------------------------------------------------------------------------------//

// Scan implements the Scanner interface.
func (v *NullString) Scan(value interface{}) error {
	return (*sql.NullString)(v).Scan(value)
}

// MarshalJSON for NullString
func (v NullString) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(v.String)
}

//----------------------------------------------------------------------------------------------------------------------------//

// Scan implements the Scanner interface.
func (v *NullInt64) Scan(value interface{}) error {
	return (*sql.NullInt64)(v).Scan(value)
}

// MarshalJSON for NullInt64
func (v NullInt64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(v.Int64)
}

//----------------------------------------------------------------------------------------------------------------------------//

// Scan implements the Scanner interface.
func (v *NullInt32) Scan(value interface{}) error {
	return (*sql.NullInt32)(v).Scan(value)
}

// MarshalJSON for NullInt32
func (v NullInt32) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(v.Int32)
}

//----------------------------------------------------------------------------------------------------------------------------//

// Scan implements the Scanner interface.
func (v *NullFloat64) Scan(value interface{}) error {
	return (*sql.NullFloat64)(v).Scan(value)
}

// MarshalJSON for NullFloat64
func (v NullFloat64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(v.Float64)
}

//----------------------------------------------------------------------------------------------------------------------------//

// Scan implements the Scanner interface.
func (v *NullTime) Scan(value interface{}) error {
	return (*sql.NullTime)(v).Scan(value)
}

// MarshalJSON for NullTime
func (v NullTime) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(v.Time)
}

//----------------------------------------------------------------------------------------------------------------------------//

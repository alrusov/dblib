package dblib

import (
	"testing"
)

//--------------------------------------------------------------------------------------------------------------------------------------------------------------//

func TestSecureString(t *testing.T) {
	srcs := []string{
		`aaa`,
		`\`,
		`'`,
		`'\'`,
		`'12\34\\56'78''9'`,
	}

	type data struct {
		db      *DB
		results []string
	}

	dbs := map[string]data{
		MYSQL: {
			&DB{driver: MYSQL},
			[]string{
				`aaa`,
				`\\`,
				`\'`,
				`\'\\\'`,
				`\'12\\34\\\\56\'78\'\'9\'`,
			},
		},
		PGSQL: {
			&DB{driver: PGSQL},
			[]string{
				`aaa`,
				`\\`,
				`\'`,
				`\'\\\'`,
				`\'12\\34\\\\56\'78\'\'9\'`,
			},
		},
		MSSQL: {
			&DB{driver: MSSQL},
			[]string{
				`aaa`,
				`\`,
				`''`,
				`''\''`,
				`''12\34\\56''78''''9''`,
			},
		},
	}

	for tp, db := range dbs {
		for i, src := range srcs {
			res := db.db.SecureString(src)
			goal := db.results[i]
			if res != goal {
				t.Errorf(`%s SecureString(%q): expect "%s", got "%s"`, tp, src, goal, res)
			}
		}
	}
}

//--------------------------------------------------------------------------------------------------------------------------------------------------------------//

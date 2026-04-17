package utils

import "testing"

func TestValidateExportSQL(t *testing.T) {
	cases := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{"empty", "", true},
		{"simple select", "SELECT id, name FROM users", false},
		{"select trailing semi", "SELECT id FROM users;", false},
		{"with cte", "WITH t AS (SELECT 1) SELECT * FROM t", false},
		{"multi statement", "SELECT 1; DROP TABLE users", true},
		{"block comment stack", "/* hi */ SELECT 1 /* */", false},
		{"line comment drop", "SELECT 1 -- DROP TABLE users", false},
		{"update root", "UPDATE users SET admin=1", true},
		{"delete root", "DELETE FROM users", true},
		{"union update via keyword", "SELECT 1 UNION update users", true},
		{"hidden drop via comment", "SELECT 1/*;*/;drop table users", true},
		{"insert into", "INSERT INTO u VALUES (1)", true},
		{"call proc", "CALL sp_admin()", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateExportSQL(c.sql)
			if (err != nil) != c.wantErr {
				t.Fatalf("%q: got err=%v want err=%v", c.sql, err, c.wantErr)
			}
		})
	}
}

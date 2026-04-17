package utils

import (
	"errors"
	"regexp"
	"strings"
)

// Exported templates ship raw SQL that is executed by Raw(). Any operator who
// can write templates therefore has DB-shell equivalent power; treat the SQL
// column as untrusted input and refuse anything that looks like a write or a
// multi-statement payload.

var (
	// Strip -- and /* */ comments before shape checks so an attacker cannot
	// hide a DROP/UPDATE behind a comment.
	sqlLineComment  = regexp.MustCompile(`--[^\n]*`)
	sqlBlockComment = regexp.MustCompile(`/\*[\s\S]*?\*/`)

	// Forbidden root statements. We only allow SELECT / WITH (CTE) templates.
	//nolint:lll
	forbiddenLeadingStmts = []string{
		"insert", "update", "delete", "drop", "alter", "truncate", "create",
		"grant", "revoke", "exec", "execute", "call", "replace", "merge",
		"rename", "load", "copy", "attach", "detach",
	}

	// Unquoted SQL identifiers used in safe concatenation (e.g. DROP TABLE name).
	plainSQLIdentifier = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
)

// ValidateExportSQL accepts a template's raw SQL and returns an error if the
// statement is not a read-only single statement. It is intentionally
// conservative: false negatives lose a feature, false positives lose data.
func ValidateExportSQL(raw string) error {
	sql := sqlBlockComment.ReplaceAllString(raw, " ")
	sql = sqlLineComment.ReplaceAllString(sql, " ")
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return errors.New("empty sql")
	}

	// Block multi-statement payloads (e.g. "select 1; drop table users").
	// We allow a single trailing semicolon so legitimate queries still work.
	trimmed := strings.TrimRight(sql, "; \t\r\n")
	if strings.Contains(trimmed, ";") {
		return errors.New("multi-statement sql is not allowed")
	}

	lower := strings.ToLower(trimmed)
	for _, kw := range forbiddenLeadingStmts {
		// Match either the statement root or a keyword appearing after a
		// UNION/(, which are the usual injection footholds inside a SELECT.
		if strings.HasPrefix(lower, kw+" ") || strings.HasPrefix(lower, kw+"(") {
			return errors.New("only read-only (SELECT/WITH) statements are allowed")
		}
		if strings.Contains(lower, " "+kw+" ") {
			return errors.New("forbidden keyword in sql: " + kw)
		}
	}

	if !(strings.HasPrefix(lower, "select") || strings.HasPrefix(lower, "with")) {
		return errors.New("only SELECT or WITH statements are allowed")
	}
	return nil
}

// ValidatePlainSQLIdentifier ensures s is safe to embed as an unquoted SQL
// identifier (table/column/schema names from codegen). This blocks injection
// via ";", spaces, quotes, and multi-byte tricks when concatenating into DDL.
//
// Allowed: [a-zA-Z_][a-zA-Z0-9_]* with a conservative length cap.
func ValidatePlainSQLIdentifier(s string) error {
	if s == "" {
		return errors.New("empty identifier")
	}
	if len(s) > 128 {
		return errors.New("identifier too long")
	}
	if !plainSQLIdentifier.MatchString(s) {
		return errors.New("identifier must be alphanumeric or underscore, and not start with a digit")
	}
	return nil
}

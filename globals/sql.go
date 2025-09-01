package globals

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var SqliteEngine = false
var PostgresEngine = false

type batch struct {
	Old   string
	New   string
	Regex bool
}

func batchReplace(sql string, batch []batch) string {
	for _, item := range batch {
		if item.Regex {
			sql = regexp.MustCompile(item.Old).ReplaceAllString(sql, item.New)
			continue
		}

		sql = strings.ReplaceAll(sql, item.Old, item.New)
	}
	return sql
}

func replaceDuplicateKey(sql string) string {
	// manual specification of conflict columns per table
	tableConflicts := map[string]string{
		"conversation": "user_id, conversation_id",
		"package":      "user_id, type",
		"quota":        "user_id",
		"subscription": "user_id",
		"sharing":      "hash",
		"redeem":       "code",
	}

	for table, conflictCols := range tableConflicts {
		// multiple-line flag needed
		pattern := fmt.Sprintf(`(?s)(INSERT INTO %s.*?)ON DUPLICATE KEY UPDATE`, table)
		replacement := fmt.Sprintf("${1}ON CONFLICT(%s) DO UPDATE SET", conflictCols)

		re := regexp.MustCompile(pattern)
		if re.MatchString(sql) {
			sql = re.ReplaceAllString(sql, replacement)
			break
		}
	}

	// convert VALUES(column) to EXCLUDED.column
	valuesRegex := regexp.MustCompile(`VALUES\(([^)]+)\)`)
	sql = valuesRegex.ReplaceAllString(sql, "EXCLUDED.$1")

	return sql
}

func replaceForSqlite(sql string) string {
	if strings.Contains(sql, "DUPLICATE KEY") {
		sql = replaceDuplicateKey(sql)
	}

	sql = batchReplace(sql, []batch{
		// KEYWORD REPLACEMENT
		{`INT `, `INTEGER `, false},
		{` AUTO_INCREMENT`, ` AUTOINCREMENT`, false},
		{`DATETIME`, `TEXT`, false},
		{`DECIMAL`, `REAL`, false},
		{`MEDIUMTEXT`, `TEXT`, false},
		{`VARCHAR`, `TEXT`, false},

		// TEXT(65535) -> TEXT, REAL(10,2) -> REAL
		{`TEXT\(\d+\)`, `TEXT`, true},
		{`REAL\(\d+,\d+\)`, `REAL`, true},

		// UNIQUE KEY -> UNIQUE
		{`UNIQUE KEY`, `UNIQUE`, false},
	})

	return sql
}

func replaceForPostgres(sql string) string {
	if strings.Contains(sql, "DUPLICATE KEY") {
		sql = replaceDuplicateKey(sql)
	}

	sql = batchReplace(sql, []batch{
		// KEYWORD REPLACEMENT
		{`DATETIME`, `TIMESTAMP`, false},
		{`MEDIUMTEXT`, `TEXT`, false},

		// TEXT(65535) -> TEXT, REAL(10,2) -> DECIMAL(10,2), DOUBLE(10,2) -> DECIMAL(10,2)
		{`TEXT\(\d+\)`, `TEXT`, true},
		{`REAL(`, `DECIMAL(`, false},
		{`DOUBLE(`, `DECIMAL(`, false},

		// UNIQUE KEY -> UNIQUE
		{`UNIQUE KEY`, `UNIQUE`, false},

		// INT PRIMARY KEY AUTO_INCREMENT -> SERIAL PRIMARY KEY
		{`INT PRIMARY KEY AUTO_INCREMENT`, `SERIAL PRIMARY KEY`, false},
	})

	return replaceQuestionMarks(sql)
}

func replaceQuestionMarks(sql string) string {
	var (
		result strings.Builder
		count  int
	)
	for i := 0; i < len(sql); i++ {
		if sql[i] == '?' {
			count++
			result.WriteString("$")
			result.WriteString(strconv.Itoa(count))
		} else {
			result.WriteByte(sql[i])
		}
	}
	return result.String()
}

func PreflightSql(sql string) string {
	// this is a simple way to adapt the sql to the sqlite and postgres engine
	// it's not a common way to use sqlite in production, just as polyfill

	if SqliteEngine {
		sql = replaceForSqlite(sql)
	} else if PostgresEngine {
		sql = replaceForPostgres(sql)
	}

	return sql
}

func ExecDb(db *sql.DB, sql string, args ...interface{}) (sql.Result, error) {
	sql = PreflightSql(sql)
	return db.Exec(sql, args...)
}

func PrepareDb(db *sql.DB, sql string) (*sql.Stmt, error) {
	sql = PreflightSql(sql)
	return db.Prepare(sql)
}

func QueryDb(db *sql.DB, sql string, args ...interface{}) (*sql.Rows, error) {
	sql = PreflightSql(sql)
	return db.Query(sql, args...)
}

func QueryRowDb(db *sql.DB, sql string, args ...interface{}) *sql.Row {
	sql = PreflightSql(sql)
	return db.QueryRow(sql, args...)
}

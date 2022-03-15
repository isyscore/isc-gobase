package database

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/isyscore/isc-gobase/isc"
)

type DatabaseType int

const (
	MySQL      DatabaseType = iota // import _ "github.com/go-sql-driver/mysql"
	Oracle                         // import _ "github.com/mattn/go-oci8"
	SqlServer                      // import _ "github.com/denisenkom/go-mssqldb"
	PostgreSql                     // import _ "github.com/lib/pq"
	Sqlite3                        // import _ "github.com/mattn/go-sqlite3"
)

const (
	//CONNECTION_STRING user:password@tcp(host:port)/databaseName
	CONNECTION_STRING = "%s:%s@tcp(%s:%d)/%s"
)

func Connect(dbType DatabaseType, connStr string) *sql.DB {
	return CustomConnect(dbTypeToString(dbType), connStr)
}

func CustomConnect(dbType string, connStr string) *sql.DB {
	// charset=utf8
	// parseTime=true
	innerParam := connStr
	if strings.Contains(connStr, "?") {
		// 已有参数
		if !strings.Contains(connStr, "charset=utf8") {
			innerParam += "&charset=utf8"
		}
		if !strings.Contains(connStr, "parseTime=true") {
			innerParam += "&parseTime=true"
		}
	} else {
		// 没有参数
		innerParam += "?charset=utf8&parseTime=true"
	}

	db, err := sql.Open(dbType, innerParam)
	if err != nil {
		log.Printf("初始化数据库失败(%v)\n", err)
		return nil
	}
	return db
}

func dbTypeToString(dbType DatabaseType) string {
	switch dbType {
	case MySQL:
		return "mysql"
	case Oracle:
		return "oci8"
	case SqlServer:
		return "mssql"
	case PostgreSql:
		return "postgres"
	case Sqlite3:
		return "sqlite3"
	default:
		log.Printf("不支持的数据库类型\n")
		return ""
	}
}

func Insert(db *sql.DB, sql string, args ...any) (int64, error) {
	var id int64
	var err error
	if strings.Contains(sql, " RETURNING ") {
		row := db.QueryRow(sql, args...)
		err = row.Scan(&id)
	} else {
		result, err1 := db.Exec(sql, args...)
		err = err1
		if err1 == nil {
			id, _ = result.LastInsertId()
		}
	}
	return id, err
}

func Update(db *sql.DB, sql string, args ...any) (int64, error) {
	var n int64
	var err error
	result, err := db.Exec(sql, args...)
	if err == nil {
		n, _ = result.RowsAffected()
	}
	return n, err
}

func Delete(db *sql.DB, sql string, args ...any) (int64, error) {
	return Update(db, sql, args...)
}

func Query(db *sql.DB, sql string, args ...any) ([]map[string]string, error) {
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	return fetchRows(rows, err)
}

func QueryRow(db *sql.DB, sql string, args ...any) (map[string]string, error) {
	rows, err := Query(db, sql, args...)
	if rows != nil && err == nil && len(rows) > 0 {
		return rows[0], err
	}
	return nil, err
}

func QueryScalar(db *sql.DB, sql string, key string, args ...any) (string, error) {
	rows, err := Query(db, sql, args...)
	if rows != nil && err == nil && len(rows) > 0 {
		row := rows[0]
		if value, ok := row[key]; ok {
			return value, err
		}
	}
	return "", err
}

// stmt 缓存
var stmtList = make(map[string]*sql.Stmt)

func PrepareSql(db *sql.DB, name, sql string) (*sql.Stmt, error) {
	stmt, bl := stmtList[name]
	if !bl {
		var err error
		stmt, err = db.Prepare(sql)
		if err != nil {
			return nil, err
		}
		stmtList[name] = stmt
	}
	return stmt, nil
}

func PrepareQuery(db *sql.DB, name, sql string, args ...any) ([]map[string]string, error) {
	stmt, err := PrepareSql(db, name, sql)
	if err != nil {
		return nil, err
	}
	rows, err1 := stmt.Query(args...)
	return fetchRows(rows, err1)
}

func PrepareQueryRow(db *sql.DB, name, sql string, args ...any) (map[string]string, error) {
	rows, err := PrepareQuery(db, name, sql, args...)
	if rows != nil && err == nil && len(rows) > 0 {
		return rows[0], err
	}
	return nil, err
}

func PrepareQueryScalar(db *sql.DB, name, sql string, args ...any) (string, error) {
	stmt, err := PrepareSql(db, name, sql)
	if err != nil {
		return "", err
	}
	var value string
	rows, err1 := stmt.Query(args...)
	if err1 != nil {
		return "", err1
	}
	if rows.Next() {
		_ = rows.Scan(&value)
	}
	_ = rows.Close()
	return value, err
}

func PrepareExec(db *sql.DB, name, sql string, args ...any) (int64, error) {
	var n int64
	stmt, err := PrepareSql(db, name, sql)
	if err != nil {
		return 0, err
	}
	if strings.Contains(sql, " RETURNING ") {
		row, err1 := stmt.Query(args...)
		if err1 != nil {
			return n, err1
		}
		row.Next()
		err = row.Scan(&n)
		_ = row.Close()
	} else {
		result, err1 := stmt.Exec(args...)
		if err1 != nil {
			return n, err1
		}
		if "INSERT" == strings.ToUpper(sql[0:6]) {
			// XXX: postgres不能用这个方法，处何处理待考虑
			n, err = result.LastInsertId()
		} else {
			n, err = result.RowsAffected()
		}
	}
	return n, err
}

type Rows struct {
	*sql.Rows
}

type DBValue struct {
	Value any
}

func (r *Rows) GetByName(fieldName string) *DBValue {
	cs, _ := r.Columns()
	index := isc.IndexOf(cs, fieldName)
	if index == -1 {
		return nil
	}
	count := len(cs)
	vals := make([]any, count)
	scans := make([]any, count)
	for i := range scans {
		scans[i] = &vals[i]
	}
	_ = r.Scan(scans...)
	if *(scans[index].(*any)) == nil {
		return nil
	} else {
		return &DBValue{Value: scans[index]}
	}
}

func (r *Rows) GetByNameDef(fieldName string, def any) *DBValue {
	v := r.GetByName(fieldName)
	if v == nil {
		i := def
		return &DBValue{
			Value: &i,
		}
	} else {
		return v
	}
}

func (r *Rows) GetByIndex(index int) *DBValue {
	cs, _ := r.Columns()
	count := len(cs)
	if index < 0 || index > count-1 {
		return nil
	}
	vals := make([]any, count)
	scans := make([]any, count)
	for i := range scans {
		scans[i] = &vals[i]
	}
	_ = r.Scan(scans...)
	if *(scans[index].(*any)) == nil {
		return nil
	} else {
		return &DBValue{Value: scans[index]}
	}
}

func (r *Rows) GetByIndexDef(index int, def any) *DBValue {
	v := r.GetByIndex(index)
	if v == nil {
		i := def
		return &DBValue{
			Value: &i,
		}
	} else {
		return v
	}
}

func (v *DBValue) ToString() string {
	return string((*(v.Value.(*any))).([]uint8))
}

func (v *DBValue) ToInt() int {
	return int((*(v.Value.(*any))).(int64))
}

func (v *DBValue) ToInt64() int64 {
	return (*(v.Value.(*any))).(int64)
}

func (v *DBValue) ToFloat() float32 {
	return (*(v.Value.(*any))).(float32)
}

func (v *DBValue) ToDouble() float64 {
	return float64((*(v.Value.(*any))).(float32))
}

func (v *DBValue) ToBoolean() bool {
	return (*(v.Value.(*any))).([]uint8)[0] == 1
}

func (v *DBValue) ToBytes() []byte {
	return (*(v.Value.(*any))).([]uint8)
}

func (v *DBValue) ToTime() time.Time {
	return (*(v.Value.(*any))).(time.Time)
}

func DBBoolean(b bool) []uint8 {
	if b {
		return []uint8{1}
	} else {
		return []uint8{0}
	}
}

func fetchRows(rows *sql.Rows, err error) ([]map[string]string, error) {
	if rows == nil || err != nil {
		return nil, err
	}

	fields, _ := rows.Columns()
	for k, v := range fields {
		fields[k] = camelCase(v)
	}
	columnsLength := len(fields)

	values := make([]string, columnsLength)
	args := make([]any, columnsLength)
	for i := 0; i < columnsLength; i++ {
		args[i] = &values[i]
	}

	index := 0
	listLength := 100
	lists := make([]map[string]string, listLength, listLength)
	for rows.Next() {
		if e := rows.Scan(args...); e == nil {
			row := make(map[string]string, columnsLength)
			for i, field := range fields {
				row[field] = values[i]
			}

			if index < listLength {
				lists[index] = row
			} else {
				lists = append(lists, row)
			}
			index++
		}
	}

	_ = rows.Close()

	return lists[0:index], nil
}

func camelCase(str string) string {
	if strings.Contains(str, "_") {
		items := strings.Split(str, "_")
		arr := make([]string, len(items))
		for k, v := range items {
			if 0 == k {
				arr[k] = v
			} else {
				arr[k] = strings.ToTitle(v)
			}
		}
		str = strings.Join(arr, "")
	}
	return str
}

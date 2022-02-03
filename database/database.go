package database

import (
	"database/sql"
	"isc-gobase/isc"
	"isc-gobase/logger"
	"strings"
	"time"
)

type DatabaseType int

const (
	MySQL      DatabaseType = iota // import _ "github.com/go-sql-driver/mysql"
	Oracle                         // import _ "github.com/mattn/go-oci8"
	SqlServer                      // import _ "github.com/denisenkom/go-mssqldb"
	PostgreSql                     // import _ "github.com/lib/pq"
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
		logger.Error("初始化数据库失败(%v)", err)
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
	default:
		logger.Error("不支持的数据库类型")
		return ""
	}
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

func (r *Rows) GetByNameDef(fieldName string, def string) *DBValue {
	v := r.GetByName(fieldName)
	if v == nil {
		i := any([]uint8(def))
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

func (r *Rows) GetByIndexDef(index int, def string) *DBValue {
	v := r.GetByIndex(index)
	if v == nil {
		i := any([]uint8(def))
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

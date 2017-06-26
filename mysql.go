package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var _DB *sql.DB

// Init database connect string, should Close when you exit
func Init(dbString string) error {
	// close init db before
	Close()
	db, err := sql.Open("mysql", dbString)
	if err != nil {
		return fmt.Errorf("[sqlOpen]->%s", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return fmt.Errorf("[db.Ping]->%s", err)
	}
	_DB = db
	return nil
}

// Close database
func Close() error {
	if _DB != nil {
		err := _DB.Close()
		_DB = nil
		return err
	}
	return nil
}

// Exec command, return the number of rows affected, and last insert id
func Exec(sqlString string, args ...interface{}) (int64, int64, error) {
	if _DB == nil {
		return 0, 0, fmt.Errorf("Don't init database")
	}
	stmt, err := _DB.Prepare(sqlString)
	if err != nil {
		return 0, 0, fmt.Errorf("[dbPrepare]->%s", err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, 0, fmt.Errorf("[stmtExec]->%s", err)
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return 0, 0, fmt.Errorf("[RowsAffected]->%s", err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, 0, fmt.Errorf("[LastInsertId]->%s", err)
	}
	return affect, lastId, nil
}

// Query command, return datas with [][]interface slice
func Query(sqlString string, args ...interface{}) ([][]sql.RawBytes, error) {
	if _DB == nil {
		return nil, fmt.Errorf("Don't init database")
	}
	stmt, err := _DB.Prepare(sqlString)
	if err != nil {
		return nil, fmt.Errorf("[dbPrepare]->%s", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, fmt.Errorf("[stmtQuery]->" + err.Error())
	}
	defer rows.Close()
	// check columns
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("[rows.Columns]->" + err.Error())
	}
	nCol := len(columns)
	if nCol < 1 {
		return nil, fmt.Errorf("columns len < 1")
	}
	// Make a slice for the values
	values := make([]sql.RawBytes, nCol)
	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, nCol)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var datas [][]sql.RawBytes
	for rows.Next() {
		if err = rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("[rows.Scan]->" + err.Error())
		}
		data := make([]sql.RawBytes, nCol)
		copy(data, values)
		datas = append(datas, data)
	}
	return datas, nil
}

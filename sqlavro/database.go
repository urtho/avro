package sqlavro

import (
	"database/sql"

	"github.com/urtho/avro"
)

// SQLDatabase2AVRO - fetch all tables of the given SQL database and translate them to avro schemas
func SQLDatabase2AVRO(db *sql.DB, dbName string) ([]avro.RecordSchema, error) {
	tables, err := GetTables(db, dbName)
	if err != nil {
		return nil, err
	}
	var (
		tableName string
		schema    *avro.RecordSchema
		schemas   = make([]avro.RecordSchema, 0, len(tables))
	)
	for _, tableName = range tables {
		schema, err = SQLTable2AVRO(db, "", dbName, tableName)
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, *schema)
	}
	return schemas, nil
}

// GetTables - returns table names of the given database
func GetTables(db *sql.DB, dbName string) ([]string, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if len(dbName) > 0 {
		rows, err = db.Query(
			`SELECT TABLE_NAME 
			 FROM INFORMATION_SCHEMA.TABLES 
			 WHERE TABLE_SCHEMA=?`,
			dbName,
		)
	} else {
		rows, err = db.Query(
			`SELECT TABLE_NAME 
			 FROM INFORMATION_SCHEMA.TABLES`,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		tableName string
		tables    = make([]string, 0, 50)
	)
	for rows.Next() {
		err = rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}

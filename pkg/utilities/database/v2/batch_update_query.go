package database

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// BatchUpdateQuery uses query without creating new table
// Usage: When parameter count stays within MaxParameter for Postgres (65535)
func BatchUpdateQuery(db *gorm.DB, tableName string, entries interface{}) (err error) { // nolint: gocognit
	slice, columnNames, _, err := validateBatchInput(db, entries)
	if err != nil {
		return err
	}

	pKeyColumnName := getPrimaryKeyColumnName(db, &slice)

	// nolint: gomnd
	amountOfTransactions := ((slice.Len() * len(columnNames)) / MaxParameters) + 2
	AmountOfEntriesInOneTransaction := slice.Len() / amountOfTransactions

	return Transact(db, func(tx *gorm.DB) error {
		for i := 0; i < amountOfTransactions; i++ {
			columnValueStrings := make(map[string][]string)
			columnValues := make(map[string][]interface{})
			for j := (i * AmountOfEntriesInOneTransaction); j < ((i+1)*AmountOfEntriesInOneTransaction) && j < slice.Len(); j++ {
				elem := slice.Index(j).Interface()
				fields := db.NewScope(elem).Fields()
				now := time.Now()
				for _, field := range fields {
					if shouldSkipField(field) {
						continue
					}
					// nolint: gocritic
					if field.IsBlank && field.HasDefaultValue {
						// If field is empty, use default value if possible.
						columnValueStrings[field.DBName] = append(columnValueStrings[field.DBName], "DEFAULT")
					} else if field.DBName == "created_at" || field.DBName == "updated_at" {
						columnValueStrings[field.DBName] = append(columnValueStrings[field.DBName], "?::timestamp")
						columnValues[field.DBName] = append(columnValues[field.DBName], now)
					} else {
						var v string
						scope := db.NewScope(field)
						sqlTag := scope.Dialect().DataTypeOf(field.StructField)
						v = fmt.Sprintf("?::%s", sqlTag)
						if t, ok := field.TagSettings["TYPE"]; ok {
							v = fmt.Sprintf("(?)::%s", t)
						}
						columnValueStrings[field.DBName] = append(columnValueStrings[field.DBName], v)
						columnValues[field.DBName] = append(columnValues[field.DBName], field.Field.Interface())
					}
				}
			}
			if err := updateInOneQuery(tx, tableName, columnNames, pKeyColumnName, columnValueStrings, columnValues); err != nil {
				return err
			}
		}
		return nil
	})
}

func getPrimaryKeyColumnName(db *gorm.DB, slice *reflect.Value) string {
	scope := db.NewScope(slice.Index(0).Interface())
	for _, field := range scope.Fields() {
		if field.IsPrimaryKey {
			return field.DBName
		}
	}
	return ""
}

func updateInOneQuery(tx *gorm.DB, tableName string, columns []string, pKeyColumnName string, valueStrings map[string][]string, columnValues map[string][]interface{}) error {
	var vals []interface{}
	var equalColumnNames []string
	for _, column := range columns {
		equalColumnNames = append(equalColumnNames, fmt.Sprintf("%s = tmptable.%s", column, column))
	}
	joinedColumns := strings.Join(equalColumnNames, ", ")

	var columnValueQuery []string
	for k, v := range valueStrings {
		vals = append(vals, columnValues[k]...)
		parameterString := strings.Join(v, ",")
		columnValueQuery = append(columnValueQuery, fmt.Sprintf(`unnest(array[%s]) as %s`, parameterString, k))
	}
	columnValueQueryJoined := strings.Join(columnValueQuery, ", ")

	// We can format here, since all columns and table names are
	// provided with quote identifiers
	stmt := fmt.Sprintf(`UPDATE %s 
		SET %s
		FROM
		(SELECT 
			%s
		) as tmptable
		WHERE %s.%s = tmptable.%s;`,
		pq.QuoteIdentifier(tableName),
		joinedColumns,
		columnValueQueryJoined,
		pq.QuoteIdentifier(tableName),
		pq.QuoteIdentifier(pKeyColumnName),
		pq.QuoteIdentifier(pKeyColumnName),
	)
	err := tx.Exec(stmt, vals...).Error
	return err
}

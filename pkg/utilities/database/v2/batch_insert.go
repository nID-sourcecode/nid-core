package database

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

// Error definitions
var (
	ErrEntriesShouldBeSlice error = fmt.Errorf("entries should be a slice")
)

const (
	// MaxParameters Maximum amount of parameters that postgres supports.
	MaxParameters int = 65535
	// ColumnNameEnd inidicates end column name
	ColumnNameEnd = "end"
)

// BatchInsert batch inserts efficiently inserts a list of entries in a database table
func BatchInsert(db *gorm.DB, tableName string, entries interface{}) (err error) {
	if entries == nil {
		return ErrEntriesShouldBeSlice
	}
	// Check if entries is a slice
	if reflect.TypeOf(entries).Kind() != reflect.Slice {
		return ErrEntriesShouldBeSlice
	}

	// Check if the slice is not empty
	slice := reflect.ValueOf(entries)
	if slice.Len() == 0 {
		return ErrEntriesShouldBeSlice
	}

	parameterCount := 0 // Counter to track how many parameters are in the insert query already.
	columnNames, defaultRowParameters := getColumnNames(db, &slice)
	for i, columnName := range columnNames {
		if columnName == ColumnNameEnd {
			columnNames[i] = fmt.Sprintf("\"%s\"", ColumnNameEnd)
		}
	}
	var valueStrings []string
	var values []interface{}
	rowParameters := make([]string, len(defaultRowParameters))

	// Run inserts in a transaction because could be split up into multiple inserts, if the
	// parameterCount exceeds the MaxParameters.
	return Transact(db, func(tx *gorm.DB) error {
		// Insert the elements.
		for i := 0; i < slice.Len(); i++ {
			// Postgres has a maximum amount of parameters. So if we're inserting so many records
			// that we can't do them in one insert, we'll split it up in multiple inserts.
			if parameterCount+len(columnNames) >= MaxParameters {
				err := insert(tx, tableName, columnNames, valueStrings, values)
				if err != nil {
					return errors.Wrap(err, "could not insert items")
				}
				valueStrings = []string{}
				values = []interface{}{}
				parameterCount = 0
			}

			// Restore default row parameters
			copy(rowParameters, defaultRowParameters)

			// Get item's values
			elem := slice.Index(i).Interface()
			itemValues, itemParameterString := getValuesForItem(db, elem, rowParameters)

			// Add the values of the current element to the statement.
			values = append(values, itemValues...)
			valueStrings = append(valueStrings, itemParameterString)

			parameterCount += len(itemValues)
		}

		err := insert(tx, tableName, columnNames, valueStrings, values)
		if err != nil {
			return errors.Wrap(err, "could not insert items")
		}
		return nil
	})
}

func getValuesForItem(db *gorm.DB, elem interface{}, rowParameters []string) (values []interface{}, parameterString string) {
	// Use gorm to get the model's fields with some
	// context, such as database field names and values.
	fields := db.NewScope(elem).Fields()
	now := time.Now()
	index := 0
	for _, field := range fields {
		if shouldSkipField(field) {
			continue
		}
		// nolint: gocritic
		if field.IsBlank && field.HasDefaultValue {
			// If field is empty, use default value if possible.
			rowParameters[index] = "DEFAULT"
		} else if field.DBName == "created_at" || field.DBName == "updated_at" {
			// Set created_at and updated_at to the current time.
			values = append(values, now)
		} else {
			// Normal field, just get the field's value.
			values = append(values, field.Field.Interface())
		}
		index++
	}
	// parameterString will be something like "(?, ?, DEFAULT, (?)::jsonb)"
	parameterString = fmt.Sprintf("(%s)", strings.Join(rowParameters, ", "))
	return values, parameterString
}

// Get the names of all columns of the table we're inserting into.
// And create a string slice to create the parameter string for each row. e.g: "(?, ?, (?)::json)"
func getColumnNames(db *gorm.DB, slice *reflect.Value) (columnNames, rowParameters []string) {
	scope := db.NewScope(slice.Index(0).Interface())
	for _, field := range scope.Fields() {
		if shouldSkipField(field) {
			continue
		}
		v := "?"
		// If the field has a type tag, try to cast it to that type.
		if t, ok := field.TagSettings["TYPE"]; ok {
			v = fmt.Sprintf("(?)::%s", t)
		}
		rowParameters = append(rowParameters, v)
		columnNames = append(columnNames, pq.QuoteIdentifier(field.DBName))
	}
	return columnNames, rowParameters
}

func insert(tx *gorm.DB, tableName string, columns, rowPlaceholders []string, values []interface{}) error {
	joinedColumns := strings.Join(columns, ", ")
	joinedRows := strings.Join(rowPlaceholders, ", ")
	// We can format here, since all columns and table names are
	// provided with quote identifiers
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		pq.QuoteIdentifier(tableName),
		joinedColumns,
		joinedRows,
	)
	err := tx.Exec(stmt, values...).Error
	return err
}

func shouldSkipField(field *gorm.Field) bool {
	return (field.IsPrimaryKey && field.IsBlank) || // Skip primary fields
		!field.IsNormal // Skip structs that are added because of a belongs_to relation for example.
}

package database

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

// BatchUpdateTmpTable create temp table with batch insert and use join table
// Usage: For updating a huge amount of rows in a table, multiple times of Maxparameters for postgres
func BatchUpdateTmpTable(db *gorm.DB, tableName string, entries interface{}) (err error) {
	_, columnNames, valueTypes, err := validateBatchInput(db, entries)
	if err != nil {
		return err
	}

	tNameTmp := getTmpTableName(tableName)
	return Transact(db, func(tx *gorm.DB) error {
		if err := createTmpTable(tx, tNameTmp, columnNames, valueTypes); err != nil {
			return errors.Wrap(err, "Creating temporary table failed")
		}
		if err := BatchInsert(tx, tNameTmp, entries); err != nil {
			return errors.Wrap(err, "Batch inserting entries in temporary table failed")
		}
		if err := update(tx, tableName, tNameTmp, columnNames); err != nil {
			return errors.Wrap(err, "Joined update with temporary table failed")
		}
		return deleteTmpTable(tx, tNameTmp)
	})
}

func getTmpTableName(tableName string) string {
	uid := uuid.Must(uuid.NewV4())
	stringified := uid.String()
	appended := strings.ReplaceAll(stringified, "-", "")
	return tableName + appended
}

func deleteTmpTable(tx *gorm.DB, tableName string) error {
	stmt := fmt.Sprintf(`DROP TABLE %s;`,
		pq.QuoteIdentifier(tableName),
	)
	err := tx.Exec(stmt).Error
	return err
}

func createTmpTable(tx *gorm.DB, tableName string, columns, valueTypes []string) error {
	var c []string
	for i, column := range columns {
		c = append(c, fmt.Sprintf("%s %s", column, valueTypes[i]))
	}
	joinedColumns := strings.Join(c, ", ")

	stmt := fmt.Sprintf(`CREATE TABLE %s (%s);`,
		pq.QuoteIdentifier(tableName),
		joinedColumns,
	)
	err := tx.Exec(stmt).Error
	return err
}

// getColumnNamesAndPostgresTypes gets postgres types for columns
func getColumnNamesAndPostgresTypes(db *gorm.DB, slice *reflect.Value) (columnNames, rowParameters []string) {
	scope := db.NewScope(slice.Index(0).Interface())
	for _, field := range scope.Fields() {
		if shouldSkipField(field) {
			continue
		}
		var v string
		// If the field has a type tag, try to cast it to that type.
		t, ok := field.TagSettings["TYPE"]
		if ok {
			v = t
		} else {
			v = getPostgresTypeForPrimitive(db, field)
		}
		if field.IsPrimaryKey {
			v = fmt.Sprintf("%s %s", v, "not null primary key")
		}
		rowParameters = append(rowParameters, v)
		columnNames = append(columnNames, pq.QuoteIdentifier(field.DBName))
	}
	return columnNames, rowParameters
}

// Get postgres types for primitives
func getPostgresTypeForPrimitive(db *gorm.DB, field *gorm.Field) string {
	scope := db.NewScope(field)
	sqlTag := scope.Dialect().DataTypeOf(field.StructField)
	return sqlTag
}

func update(tx *gorm.DB, tableName, tmpTableName string, columns []string) error {
	var columnWithValues []string
	for _, column := range columns {
		columnWithValues = append(columnWithValues, fmt.Sprintf("%s = %s.%s", column, tmpTableName, column))
	}
	joinedColumns := strings.Join(columnWithValues, ", ")
	// We can format here, since all columns and table names are
	// provided with quote identifiers
	stmt := fmt.Sprintf(`UPDATE %s SET %s
		FROM %s
		WHERE %s.id = %s.id`,
		pq.QuoteIdentifier(tableName),
		joinedColumns,
		pq.QuoteIdentifier(tmpTableName),
		pq.QuoteIdentifier(tableName),
		pq.QuoteIdentifier(tmpTableName))

	err := tx.Exec(stmt).Error
	return err
}

func validateBatchInput(db *gorm.DB, entries interface{}) (reflect.Value, []string, []string, error) {
	if entries == nil {
		return reflect.Value{}, nil, nil, ErrEntriesShouldBeSlice
	}
	// Check if entries is a slice
	if reflect.TypeOf(entries).Kind() != reflect.Slice {
		return reflect.Value{}, nil, nil, ErrEntriesShouldBeSlice
	}

	// Check if the slice is not empty
	slice := reflect.ValueOf(entries)
	if slice.Len() == 0 {
		return reflect.Value{}, nil, nil, ErrEntriesShouldBeSlice // Nothing to do.
	}

	columnNames, valueTypes := getColumnNamesAndPostgresTypes(db, &slice)
	for i, columnName := range columnNames {
		if columnName == "end" {
			columnNames[i] = "\"end\""
		}
	}
	return slice, columnNames, valueTypes, nil
}

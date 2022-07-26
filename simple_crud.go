package simple_crud

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
)

// Column name and column value struct type for searching specific rows.
type QueryHook struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Utilizes generic type for receiving different types of custom database
// struct type.
type Driver[T any] struct {
	db *sql.DB
}

func NewDriver[T any](db *sql.DB) *Driver[T] {
	return &Driver[T]{
		db: db,
	}
}

var (
	DuplicateRow    = errors.New("Duplicate row.")
	RowNotExist     = errors.New("Row doesn't exist.")
	RowUpdateFailed = errors.New("Row update failed.")
	RowDeleteFailed = errors.New("Row deletion failed.")
	TableNotExist   = errors.New("Table doesn't exist, created just now.")
)

// Create a table.
// Takes the table's name and creation query.
// Returns error if something wrong happened
func (d *Driver[T]) InitDB(tn string, q string) error {
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s(
	%s
	);
	`, tn, q)

	_, err := d.db.Exec(query)
	return err
}

// Drop a table.
// Takes the table's name and creation query.
// Returns error if something wrong happened
func (d *Driver[T]) DropTable(tn string) error {
	query := fmt.Sprintf("DROP TABLE %s;", tn)

	_, err := d.db.Exec(query)
	return err
}

// Create a row.
// Takes the table's name, field names, and values.
// Returns error if something wrong happened.
func (d *Driver[T]) CreateRow(tn string, fn string, vs string) error {
	query := fmt.Sprintf("INSERT INTO %s(%s) values(%s);", tn, fn, vs)
	_, err := d.db.Exec(query)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return DuplicateRow
			}
		}
		return err
	}
	return nil
}

// Get all rows from a table.
// Takes the table's name.
// Returns all rows in the struct type that was initialized with a nil as
// an error value or nil with an error value if something wrong happened.
func (d *Driver[T]) ReadAllRow(tn string) (*[]T, error) {
	rows, err := d.db.Query(fmt.Sprintf("SELECT * FROM %s;", tn))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []T

	// Get the columns dynamically.
	cols, _ := rows.Columns()

	res := make([][]byte, len(cols))

	for rows.Next() {
		if err := rows.Scan(DynamicScannerValues(res, cols)...); err != nil {
			return nil, err
		}

		// Transform results into structs of the specified custom database
		// struct type.
		var t T
		tmp, err := json.Marshal(ToStringifiedJSON(res, cols))
		if err != nil {
			log.Println(err)
			return nil, err
		}
		var s string
		err = json.Unmarshal(tmp, &s)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		log.Println(s)
		err = json.Unmarshal(s, &t)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		all = append(all, t)
	}

	return &all, nil
}

// Get certain row from a table.
// Takes the table's name, to be updated column's hook, and its value.
// Returns the row in the struct type that was initialized with a nil as
// an error value or nil with an error value if something wrong happened.
func (d *Driver[T]) ReadRow(tn string, qh *QueryHook) (*T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = \"%s\";", tn, qh.Name, qh.Value)
	row := d.db.QueryRow(query)

	// Get the columns dynamically.
	rows, err := d.db.Query(fmt.Sprintf("SELECT * FROM %s", tn))
	if err != nil {
		return nil, err
	}
	cols, _ := rows.Columns()
	rows.Close()

	res := make([][]byte, len(cols))

	if err := row.Scan(DynamicScannerValues(res, cols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, RowNotExist
		}
		return nil, err
	}

	// Transform the result into struct of the specified custom database
	// struct type.
	var single T
	err = json.Unmarshal([]byte(ToStringifiedJSON(res, cols)), &single)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &single, err
}

// Update certain row from a table.
// Takes the table's name, formatted values, to be updated column's hook,
// and its value.
// Returns error if something wrong happened.
func (d *Driver[T]) UpdateRow(tn string, fvs string, qh *QueryHook) error {
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = \"%s\";", tn, fvs, qh.Name, qh.Value)
	return UpdateDeleteHelper(d, query, RowUpdateFailed)
}

// Delete certain row from a table.
// Takes the table's name, to be deleted column's hook, and its value.
// Returns error if something wrong happened.
func (d *Driver[T]) DeleteRow(tn string, qh *QueryHook) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = \"%s\";", tn, qh.Name, qh.Value)
	return UpdateDeleteHelper(d, query, RowDeleteFailed)
}

// Delete multiple rows from a table.
// Takes the table's name, to be deleted column's hook, and its value with
// this format: `1, 2, 3, 4` or `"Article 1", "Article 2", "Article 3"`
// Returns error if something wrong happened.
func (d *Driver[T]) DeleteRows(tn string, qh *QueryHook) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s IN (\"%s\");", tn, qh.Name, qh.Value)
	return UpdateDeleteHelper(d, query, RowDeleteFailed)
}

// Reduce the repeating code for update and delete operations.
func UpdateDeleteHelper[T any](d *Driver[T], q string, customErr error) error {
	res, err := d.db.Exec(q)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return customErr
	}

	return nil
}

func DynamicScannerValues(row [][]byte, cols []string) []any {
	rowPtr := make([]any, len(cols))
	for i := range row {
		rowPtr[i] = &row[i]
	}

	return rowPtr
}

func ToStringifiedJSON(row [][]byte, cols []string) string {
	var s string
	for i, v := range row {
		if i == 0 {
			s += "{\n"
		}
		if i > 0 {
			s += ",\n"
		}
		s += "\t\"" + cols[i] + "\": \"" + string(v) + "\""
	}
	s += "\n}"
	return s
}

/*

Go Simple CRUD is a simple database CRUD operation API with dynamic row
scanning for Go.
Copyright (C) 2022  Aranggi J. Toar

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; only version 2 of the License.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License along
with this program; if not, write to the Free Software Foundation, Inc.,
51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

*/

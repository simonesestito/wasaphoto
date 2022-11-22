package database

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

// QueryStructRow executes the given query, expecting a single row to be returned,
// and it populates the struct referenced by destPointer.
// In case no rows are returned by the query, ErrNoResult is thrown.
// Other more low-level errors may be thrown as well.
func (db appSqlDatabase) QueryStructRow(destPointer any, query string, args ...any) error {
	rows, err := db.DB.Queryx(query, args...)
	if err != nil {
		return err
	}
	defer tryClosingRows(rows)

	if !rows.Next() {
		return ErrNoResult
	}

	return parseRow(rows, destPointer)
}

type StructRows struct {
	dest any
	rows *sqlx.Rows
}

// Next will parse and return the next row as a struct, not a pointer.
func (r *StructRows) Next() (any, error) {
	if !r.rows.Next() {
		r.Close()
		return nil, ErrNoResult
	}

	// Parse next row
	nextRow := reflect.New(reflect.TypeOf(r.dest)).Elem().Interface()
	err := parseRow(r.rows, &nextRow)
	if err != nil {
		return nil, err
	}

	return nextRow, nil
}

func (r *StructRows) Close() error {
	return r.rows.Close()
}

// QueryStructRows provides a mechanism to query multiple rows as a struct
//
// Example usage:
//
// rows, err := dao.Db.QueryStructRows(AStruct{}, query, args...)
// if err != nil { return err; }
//
//	for entity, err := rows.Next(); err == nil; entity, err = rows.Next() {
//		slice = append(slice, entity.(AStruct))
//	}
func (db appSqlDatabase) QueryStructRows(entityStruct any, query string, args ...any) (StructRows, error) {
	rows, err := db.DB.Queryx(query, args...)
	if err != nil {
		return StructRows{}, err
	}

	return StructRows{
		dest: entityStruct,
		rows: rows,
	}, nil
}

func parseRow(rows *sqlx.Rows, destPointer any) error {
	rowData := make(map[string]any)
	err := rows.MapScan(rowData)
	if err != nil {
		return err
	}

	// Decode to given struct
	decoderConfig := &mapstructure.DecoderConfig{
		ErrorUnused: true,
		ErrorUnset:  true,
		ZeroFields:  true,
		TagName:     "json",
		Result:      destPointer,
		Squash:      true,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return errors.New("error creating a decoder")
	}

	return decoder.Decode(rowData)
}

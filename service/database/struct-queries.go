package database

import (
	"database/sql"
	"errors"
	"github.com/mitchellh/mapstructure"
)

// QueryStructRow executes the given query, expecting a single row to be returned,
// and it populates the struct referenced by destPointer.
// In case no rows are returned by the query, sql.ErrNoRows is thrown.
// Other more low-level errors may be thrown as well.
func (db appSqlDatabase) QueryStructRow(destPointer any, query string, args ...any) error {
	rows, err := db.DB.Queryx(query, args...)
	if err != nil {
		return err
	}

	if !rows.Next() {
		return sql.ErrNoRows
	}

	rowData := make(map[string]any)
	err = rows.MapScan(rowData)
	if err != nil {
		return err
	}

	// Decode to given struct
	decoderConfig := &mapstructure.DecoderConfig{
		ErrorUnused: true,
		ZeroFields:  true,
		TagName:     "json",
		Result:      destPointer,
		Squash:      true,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return errors.New("error creating a decoder")
	}

	if err := decoder.Decode(rowData); err != nil {
		return err
	}

	return nil
}

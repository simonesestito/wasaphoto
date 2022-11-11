package database

func (db appSqlDatabase) QueryStructRow(destPointer any, query string, args ...any) error {
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return err
	}

	// TODO: Check no rows

	rows.Next()
	rows.Scan()
	return nil // TODO: QueryStructRow
}

package repository

import "database/sql"

func toNullString(jsonBytes []byte) sql.NullString {
	if string(jsonBytes) == "null" {
		return sql.NullString{} // Default is {String: "", Valid: false}
	}
	return sql.NullString{
		String: string(jsonBytes),
		Valid:  true,
	}
}

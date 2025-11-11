package data

import "database/sql"

type Report struct{}

type ReportModel struct {
	DB *sql.DB
}

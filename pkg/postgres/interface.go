package postgres

import (
	"fmt"
	"gorm.io/gorm/logger"
)

type ConfigureParams struct {
	User      string
	Password  string
	Host      string
	Database  string
	Params    string
	DebugMode logger.LogLevel
}

func GetPostgresUri(params ConfigureParams) string {
	return fmt.Sprintf("postgresql://%v:%v@%v/%v?%v", params.User, params.Password, params.Host, params.Database, params.Params)
}

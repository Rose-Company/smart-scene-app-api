package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	prefix          string
	db              *gorm.DB
	mode            logger.LogLevel
	params          ConfigureParams
	migrationTables []interface{}
}

func (s *Postgres) Get() interface{} {
	return s.db
}

func (s *Postgres) Run() error {
	uri := GetPostgresUri(s.params)
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger: logger.Default.LogMode(s.params.DebugMode),
	})
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Postgres) Configure(prefix string, params ConfigureParams) error {
	s.prefix = prefix
	s.params = params
	return nil
}

func (s *Postgres) GetPrefix() string {
	return s.prefix
}

func (s *Postgres) Stop() <-chan bool {
	stop := make(chan bool)
	go func() {
		stop <- true
	}()
	return stop
}

func (s *Postgres) SetMigrationTables(tables ...interface{}) {
	s.migrationTables = tables
}

func (s *Postgres) migrate(dst ...interface{}) {
	for _, table := range dst {
		err := s.db.AutoMigrate(table)
		if err != nil {
			panic(err)
		}
	}
}

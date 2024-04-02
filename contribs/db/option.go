package db

import (
	"database/sql"
	"fmt"
	"time"
)

// Option is database connection option
type Option struct {
	Driver             string        `mapstructure:"driver"`
	Address            string        `mapstructure:"address"`
	User               string        `mapstructure:"user"`
	Password           string        `mapstructure:"password"`
	Database           string        `mapstructure:"database"`
	Params             string        `mapstructure:"param"`
	MaxIdleConnections int           `mapstructure:"maxIdleConnections"`
	MaxOpenConnections int           `mapstructure:"maxOpenConnections"`
	MaxLifeTime        time.Duration `mapstructure:"maxLifeTime"`
	MaxIdleTime        time.Duration `mapstructure:"maxIdleTime"`
}

// Open returns a *sql.DB with the give Option
func Open(option Option) (*sql.DB, error) {

	// chose dsn
	var dsn string
	switch option.Driver {
	case Mysql:
		dsn = MysqlDsn(option)
	case Postgres:
		dsn = PostgresDsn(option)
	case Sqlite:
		dsn = SQLiteDsn(option)
	case Sqlserver:
		dsn = SQLiteDsn(option)
	default:
		return nil, fmt.Errorf("unsupported driver: %v", option.Driver)
	}

	// default configuration
	if option.MaxOpenConnections == 0 {
		option.MaxOpenConnections = 100
	}
	if option.MaxIdleConnections == 0 {
		option.MaxIdleConnections = 10
	}
	if option.MaxLifeTime == 0 {
		option.MaxLifeTime = time.Hour
	}
	if option.MaxIdleTime == 0 {
		option.MaxIdleTime = time.Minute * 10
	}

	// connect to the db
	db, err := sql.Open(option.Driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// set configuration
	db.SetMaxOpenConns(option.MaxOpenConnections)
	db.SetMaxIdleConns(option.MaxIdleConnections)
	db.SetConnMaxLifetime(option.MaxLifeTime)
	db.SetConnMaxIdleTime(option.MaxIdleTime)

	return db, nil
}

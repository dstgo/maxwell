package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/microsoft/go-mssqldb"
	_ "modernc.org/sqlite"

	"net/netip"
)

const (
	Mysql     = "mysql"
	Postgres  = "postgres"
	Sqlite    = "sqlite3"
	Sqlserver = "sqlserver"
)

func MysqlDsn(opt Option) string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?%s`, opt.User, opt.Password, opt.Address, opt.Database, opt.Params)
}

func PostgresDsn(opt Option) string {
	addrPort, _ := netip.ParseAddrPort(opt.Address)
	dsn := fmt.Sprintf("host=%s, port=%d user=%s password=%s %s", addrPort.Addr(), addrPort.Port(), opt.User, opt.Password, opt.Params)
	return dsn
}

func SQLiteDsn(opt Option) string {
	dsn := fmt.Sprintf("%s?%s", opt.Database, opt.Params)
	return dsn
}

func SQLServerDsn(opt Option) string {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s %s", opt.User, opt.Password, opt.Address, opt.Database, opt.Params)
	return dsn
}

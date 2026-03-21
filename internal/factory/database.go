package factory

import (
	"fmt"

	"github.com/canonflow/canonflow-go-ddd/internal/contract"
)

type mysqlDriver struct {
	Format string
}

type postgreSQLDriver struct {
	Format string
}

func newMysqlDriver(username string, password string, host string, port int, database string) *mysqlDriver {
	return &mysqlDriver{
		Format: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database),
	}
}

func (driver *mysqlDriver) GetDSN() string {
	return driver.Format
}

func newPostgreSQLDriver(username string, password string, host string, port int, database string) *postgreSQLDriver {
	return &postgreSQLDriver{
		Format: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", host, username, password, database, port),
	}
}

func (driver *postgreSQLDriver) GetDSN() string {
	return driver.Format
}

func NewDatabaseFactory(driver string, username string, password string, host string, port int, database string) contract.DatabaseContract {
	if driver == "mysql" {
		return newMysqlDriver(username, password, host, port, database)
	} else if driver == "postgres" {
		return newPostgreSQLDriver(username, password, host, port, database)
	}

	return nil
}

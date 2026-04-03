package configs

import "fmt"

type Database struct {
	Host     string `env:"HOST"     envDefault:"localhost"`
	Port     int    `env:"PORT"     envDefault:"5432"`
	User     string `env:"USER"     envDefault:"db-user"`
	Password string `env:"PASSWORD" envDefault:"db-password"`
	DBName   string `env:"DBNAME"   envDefault:"vna"`
}

func (db Database) BuildConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db.Host, db.Port, db.User, db.Password, db.DBName)
}

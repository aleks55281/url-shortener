package datab

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PostgrSql struct {
	Host     string
	Port     string
	User     string
	Dbname   string
	Password string
	Sslmode  string
}

func ConPostgrSql(p PostgrSql) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", p.Host, p.Port, p.User, p.Dbname, p.Password, p.Sslmode))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

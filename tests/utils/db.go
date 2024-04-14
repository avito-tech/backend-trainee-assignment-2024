package utils

import (
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"gta2024/pkg/repository"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	SSLMode  string
}

type Postgres struct {
	db *sqlx.DB
	c  *Config
}

func GetDSN(c Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.SSLMode,
	)
}

func NewPostgres(c Config) *Postgres {
	dsn := GetDSN(c)
	db, err := sqlx.Open("postgres", dsn)

	if err != nil {
		panic("can't connect to db")
	}

	err = db.Ping()
	if err != nil {
		logrus.Error(err.Error())
		panic(err)
	}

	return &Postgres{db: db, c: &c}
}

func GenerateDbName() string {
	return fmt.Sprintf("db_%s", strings.ReplaceAll(uuid.New().String(), "-", ""))
}

func (p *Postgres) SetUp() (*sqlx.DB, error) {
	dbName := GenerateDbName()

	_, err := p.db.Exec(fmt.Sprintf("create database %s", dbName))
	if err != nil {
		return nil, fmt.Errorf("can't create database: %w", err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     p.c.Host,
		Port:     p.c.Port,
		Username: p.c.Username,
		Password: p.c.Password,
		SSLMode:  p.c.SSLMode,
		DBName:   dbName,
	})
	if err != nil {
		return nil, fmt.Errorf("can't get connection: %w", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("error postgres.WithInstance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../schema",
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("error postgres.NewWithDatabaseInstance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("error run migrations: %w", err)
	}

	return db, err
}

func (p *Postgres) TearDown(db *sqlx.DB) {
	var dbName string

	row := db.QueryRow("select current_database()")
	if err := row.Scan(&dbName); err != nil {
		logrus.Errorf("can't get dbName: %s", err.Error())
		return
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("can't close connection with db %s: %s", dbName, err.Error())
		return
	}

	if _, err := p.db.Exec(fmt.Sprintf("drop database %s with (force)", dbName)); err != nil {
		logrus.Errorf("can't drop database with name %s: %s", dbName, err.Error())
		return
	}
}

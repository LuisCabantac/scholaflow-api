package main

import (
	postgres "github.com/LuisCabantac/scholaflow-api/internal/adapters/postgresql/sqlc"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wneessen/go-mail"
)

type dbConfig struct {
	dsn string
}

type config struct {
	addr   string
	env    string
	appUrl string
	db     dbConfig
}

type application struct {
	s3Client   *s3.Client
	mailClient *mail.Client
	db         *pgxpool.Pool
	queries    *postgres.Queries
	cfg        config
}

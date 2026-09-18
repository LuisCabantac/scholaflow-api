package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	postgres "github.com/LuisCabantac/scholaflow-api/internal/adapters/postgresql/sqlc"
	"github.com/LuisCabantac/scholaflow-api/internal/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	s3Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/wneessen/go-mail"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Debug(".env file not found, reading from system environment")
	}

	port := env.GetString("PORT", "8080")
	environment := env.GetString("APP_ENV", "development")
	appUrl := env.GetRequired("APP_URL")
	dbUrl := env.GetRequired("GOOSE_DBSTRING")
	smtpEmail := env.GetRequired("SMTP_EMAIL")
	smtpPassword := env.GetRequired("SMTP_PASSWORD")

	s3Endpoint := env.GetRequired("AWS_ENDPOINT_URL_S3")
	s3AccessKey := env.GetRequired("AWS_ACCESS_KEY_ID")
	s3SecretKey := env.GetRequired("AWS_SECRET_ACCESS_KEY")
	s3Region := env.GetRequired("AWS_REGION")

	ctx := context.Background()

	s3Cfg, err := s3Config.LoadDefaultConfig(ctx,
		s3Config.WithRegion(s3Region), s3Config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     s3AccessKey,
				SecretAccessKey: s3SecretKey,
			},
		}),
	)
	if err != nil {
		log.Fatalf("failed to load default s3 config: %v", err)
	}

	s3Client := s3.NewFromConfig(s3Cfg, func(o *s3.Options) {
		o.BaseEndpoint = &s3Endpoint
		o.UsePathStyle = true
	})

	mailClient, err := mail.NewClient(
		"smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(smtpEmail),
		mail.WithPassword(smtpPassword),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	if err != nil {
		log.Fatalf("failed to create mail client: %v", err)
	}

	poolCfg, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatalf("failed to parse db config: %v", err)
	}
	poolCfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	poolCfg.MaxConns = 5
	poolCfg.MinConns = 1
	poolCfg.MaxConnIdleTime = 5 * time.Minute
	poolCfg.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		log.Fatalf("failed to create to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	slog.Info("connected to database", "dsn", poolCfg.ConnConfig.Host)

	queries := postgres.New(pool)

	cfg := config{
		addr:   ":" + port,
		env:    environment,
		appUrl: appUrl,
		db: dbConfig{
			dsn: dbUrl,
		},
	}

	api := application{
		s3Client:   s3Client,
		mailClient: mailClient,
		db:         pool,
		queries:    queries,
		cfg:        cfg,
	}
}

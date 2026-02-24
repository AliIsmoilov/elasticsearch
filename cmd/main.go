package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"elasticsearch/api"
	"elasticsearch/config"
	"elasticsearch/storage"

	"github.com/casbin/casbin/v2"
	_ "github.com/golang-migrate/migrate/v4"                   // db automigration
	_ "github.com/golang-migrate/migrate/v4/database"          // db automigration
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // db automigration
	_ "github.com/golang-migrate/migrate/v4/source/file"       // db automigration
	_ "github.com/lib/pq"                                      // db driver
	_ "go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gorm.io/gorm/logger"
)

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()

	cfg := config.NewConfig(".")
	databaseUrl := buildDatabaseURL(&cfg)

	// m, err := migrate.New("file://migrations", databaseUrl)
	// m, err := migrate.New("file:///app/migrations", databaseUrl)
	// if err != nil {
	// 	log.Fatal("error in creating migrations: ", zap.Error(err))
	// }
	// fmt.Printf("")
	// if err = m.Up(); err != nil {
	// 	log.Println("error updating migrations: ", zap.Error(err))
	// }

	// Connect to PostgreSQL using GORM
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Log queries slower than 1s
				LogLevel:                  logger.Info, // Log all queries
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound errors
				Colorful:                  true,        // Colorize output
			},
		),
	})
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get underlying sql.DB:", err)
	}

	// Test DB connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("database ping failed:", err)
	}
	fmt.Println("Connection successfully established with GORM")

	strg := storage.New(db)

	enforcer, err := casbin.NewEnforcer(
		"config/model.conf",
		"config/policy.csv",
	)
	if err != nil {
		log.Fatal("failed to init casbin:", err)
	}

	// Start HTTP server
	engine := api.New(&api.Handler{
		Strg: strg,
		Cfg:  &cfg,
		Enf:  enforcer,
	})

	if err = engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func buildDatabaseURL(cfg *config.Config) string {
	// Railway / production
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	// Local fallback
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DB,
	)
}

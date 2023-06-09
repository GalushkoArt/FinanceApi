package main

import (
	"fmt"
	"github.com/GalushkoArt/simpleCache"
	"github.com/galushkoart/finance-api/internal/config"
	"github.com/galushkoart/finance-api/internal/handler"
	"github.com/galushkoart/finance-api/internal/logs"
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/galushkoart/finance-api/internal/repository"
	"github.com/galushkoart/finance-api/internal/service"
	"github.com/galushkoart/finance-api/pkg/conpool"
	pkg "github.com/galushkoart/finance-api/pkg/service"
	"github.com/galushkoart/finance-api/pkg/utils"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"time"
)

// @title						Finance API
// @version					1.0
// @description				Finance REST API for equities, fx and crypto rates.
// @securityDefinitions.apikey	ApiKeyAuth
// @name						Authorization
// @in							header
// @scope.client				Grants read access to resources
// @scope.admin				Grants read and write access to resources
func main() {
	config.Init()
	logs.Init(config.Conf.Logs.Level, config.Conf.Logs.Path)
	dbConf := config.Conf.Database
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.Name, dbConf.User, dbConf.Password))
	if err != nil {
		log.Fatal().Err(err)
	}
	if db.Ping() != nil {
		log.Fatal().Err(err)
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatal().Err(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:./db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal().Err(err)
	}
	if err = m.Up(); err != nil {
		log.Fatal().Err(err)
	}
	symbolRepository := repository.NewSymbolRepository(db)
	userRepository := repository.NewUserRepository(db)
	twelveDataConf := config.Conf.API.TwelveData
	twelveDataPool := conpool.NewTwelveDataPool(twelveDataConf.ApiKey, twelveDataConf.Host, twelveDataConf.Timeout, twelveDataConf.RateLimit, 1*time.Minute)
	auditConf := config.Conf.Audit
	auditClient, err := pkg.NewAuditClient(auditConf.GRPCEnabled, auditConf.GRPCAddress)
	utils.PanicOnError(err)
	auditPublisher := pkg.NewAuditPublisher(auditConf.QueueName)
	closeMq := auditPublisher.InitPublishChannel(auditConf.MQEnabled, auditConf.MQUri)
	auditService := service.NewAuditService(auditConf.GRPCEnabled, auditClient, auditConf.MQEnabled, auditPublisher)
	symbolService := service.NewSymbolService(symbolRepository, twelveDataPool, auditService)
	symbolCache := simpleCache.NewGenericConcurrentCache[model.Symbol](config.Conf.Cache.SymbolTTL)
	hasher := service.NewHasher(dbConf.Salt)
	jwtConf := config.Conf.JWT
	jwtProducer := service.NewJwtProducer(jwtConf.HMACSecret, jwtConf.ExpiryTimeout)
	jwtParser := service.NewJwtParser(jwtConf.HMACSecret)
	authService := service.NewAuthService(userRepository, hasher, jwtProducer, time.Duration(jwtConf.RefreshTimeoutDays)*24*time.Hour, auditService)

	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		WriteTimeout: config.Conf.Server.WriteTimeout,
		ReadTimeout:  config.Conf.Server.ReadTimeout,
		Prefork:      config.Conf.Server.Prefork,
		BodyLimit:    1 << 20,
		AppName:      "Finance App " + config.Conf.Server.Environment,
	})
	app.Use(requestid.New())
	httpHandler := handler.New(swagger.HandlerDefault, authService, symbolService, symbolCache, handler.RequestLogger(), handler.AuthMiddleware(jwtParser))
	httpHandler.InitRoutes(app)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)
	serverShutdown := make(chan struct{})
	go func() {
		<-exit
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
		serverShutdown <- struct{}{}
	}()

	log.Info().Msg("Starting server")
	if err := app.Listen(":" + config.Conf.Server.Port); err != nil {
		log.Fatal().Err(err).Msg("Error on running server!")
	}

	<-serverShutdown
	done := make(chan bool)

	go func() {
		utils.PanicOnError(auditClient.Close())
		utils.PanicOnError(closeMq())
		utils.PanicOnError(driver.Close())
		closeDb(db)
		done <- true
	}()
	select {
	case <-time.After(30 * time.Second):
		log.Error().Msg("Failed to shutdown in 30 seconds")
		os.Exit(1)
	case <-done:
		log.Info().Msg("Shutdown successfully")
	}
}

func closeDb(db *sqlx.DB) {
	log.Info().Msg("Closing DB connection")
	err := db.Close()
	if err != nil {
		log.Error().Err(err).Msg("Couldn't close db connection correctly!")
	}
}

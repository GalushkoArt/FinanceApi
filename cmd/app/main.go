package main

import (
	"FinanceApi/internal/config"
	"FinanceApi/internal/handler"
	"FinanceApi/internal/logs"
	"FinanceApi/internal/model"
	"FinanceApi/internal/repository"
	"FinanceApi/internal/service"
	"FinanceApi/pkg/connectionPool"
	"fmt"
	"github.com/GalushkoArt/simpleCache"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
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
	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err)
	}
	defer closeDb(db)
	symbolRepository := repository.NewSymbolRepository(db)
	userRepository := repository.NewUserRepository(db)
	twelveDataConf := config.Conf.API.TwelveData
	twelveDataPool := connectionPool.NewTwelveDataPool(twelveDataConf.ApiKey, twelveDataConf.Host, twelveDataConf.Timeout, twelveDataConf.RateLimit)
	symbolService := service.NewSymbolService(symbolRepository, twelveDataPool)
	symbolCache := simpleCache.NewGenericConcurrentCache[model.Symbol](config.Conf.Cache.SymbolTTL)
	hasher := service.NewHasher(dbConf.Salt)
	jwtConf := config.Conf.JWT
	jwtProducer := service.NewJwtProducer(jwtConf.HMACSecret, jwtConf.ExpiryTimeout)
	jwtParser := service.NewJwtParser(jwtConf.HMACSecret)
	authService := service.NewAuthService(userRepository, hasher, jwtProducer, time.Duration(jwtConf.RefreshTimeoutDays)*24*time.Hour)

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
	httpHandler := handler.New(authService, symbolService, symbolCache, handler.RequestLogger(), handler.AuthMiddleware(jwtParser))
	httpHandler.InitRoutes(app)
	log.Info().Msg("Starting server")
	log.Fatal().Err(app.Listen(":" + config.Conf.Server.Port)).Msg("Error on running server!")
}

func closeDb(db *sqlx.DB) {
	log.Info().Msg("Closing DB connection")
	err := db.Close()
	if err != nil {
		log.Error().Err(err).Msg("Couldn't close db connection correctly!")
	}
}

package main

import (
	"FinanceApi/pkg/config"
	"FinanceApi/pkg/connectionPool"
	"FinanceApi/pkg/handler"
	"FinanceApi/pkg/log"
	"FinanceApi/pkg/model"
	"FinanceApi/pkg/repository"
	"FinanceApi/pkg/server"
	"FinanceApi/pkg/service"
	"context"
	"fmt"
	"github.com/GalushkoArt/simpleCache"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"sync"
)

func main() {
	config.Init()
	log.Init(log.LevelFromString(config.Conf.Logs.Level), config.Conf.Logs.Path)
	dbConf := config.Conf.Database
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.Name, dbConf.User, dbConf.Password))
	if err != nil {
		log.Panic(err)
	}
	defer closeDb(db)
	stopWg := &sync.WaitGroup{}
	ctx := context.WithValue(context.Background(), "stopWg", stopWg)
	symbolRepository := repository.NewSymbolRepository(db)
	twelveDataPool := connectionPool.NewTwelveDataPool(ctx)
	symbolService := service.NewSymbolService(symbolRepository, twelveDataPool)
	symbolCache := simpleCache.NewGenericConcurrentCache[model.Symbol](config.Conf.Cache.SymbolTTL)
	httpHandler := handler.New(symbolService, symbolCache)
	srv := new(server.Server)
	if err := srv.Run(httpHandler.InitRoutes()); err != nil {
		log.Panic("Error on running server!", err)
	}
	stopWg.Wait()
}

func closeDb(db *sqlx.DB) {
	err := db.Close()
	if err != nil {
		log.Error("Couldn't close db connection correctly!", err)
	}
}

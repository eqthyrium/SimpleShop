package app

import (
	neo4j2 "SimpleShop/internal/repository/neo4j"
	"SimpleShop/internal/service/usecase"
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"

	"SimpleShop/internal/config"
	"SimpleShop/internal/transport/customHttp"
	"SimpleShop/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

func RunApplication() {
	conf := config.NewConfiguration()

	LoggerObjectHttp := logger.NewLogger().GetLoggerObject("../logging/info.log", "../logging/error.log", "../logging/debug.log", "HTTP")
	db, err := openDb(*conf.Dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close(context.Background())

	repositoryObject := neo4j2.NewRepository(db)
	serviceObject := usecase.NewUseCase(repositoryObject)
	httpTransport := customHttp.NewTransportHttpHandler(serviceObject, LoggerObjectHttp)

	router := httpTransport.Routering()
	message := fmt.Sprintf("The server is running at: http://localhost%s/\n", *conf.Addr)
	log.Print(message)
	httpTransport.InfoLog.Print(message)
	log.Fatalln(http.ListenAndServe(*conf.Addr, router))
}

func openDb(dsn string) (neo4j.DriverWithContext, error) {
	// Open Neo4j database connection using DSN (Data Source Name)
	driver, err := neo4j.NewDriverWithContext(dsn, neo4j.BasicAuth("neo4j", "Iphone12345", ""))
	if err != nil {
		return nil, err
	}

	// Test the connection
	ctx := context.Background()
	if err = driver.VerifyConnectivity(ctx); err != nil {
		return nil, err
	}

	return driver, nil
}

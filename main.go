package main

import (
	"context"
	"fmt"

	"github.com/estoyaburrido/rest-grpc-api-example/app"
	"github.com/estoyaburrido/rest-grpc-api-example/app/databases"
	"github.com/estoyaburrido/rest-grpc-api-example/app/tools"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("Host", "web-api")
	viper.SetDefault("HttpPort", "8080")
	viper.SetDefault("GrpcPort", "9111")

	viper.SetDefault("RedisDatabaseHost", "redis-db")
	viper.SetDefault("RedisDatabasePort", "6379")

	viper.SetConfigName("config/config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error with config file: %s", err))
	}
	viper.WatchConfig()
}

func main() {
	host := viper.GetString("Host")
	httpPort := viper.GetString("HttpPort")
	grpcPort := viper.GetString("GrpcPort")

	redisDBHost := viper.GetString("RedisDatabaseHost")
	redisDBPort := viper.GetString("RedisDatabasePort")

	redisAddr := fmt.Sprintf("%v:%v", redisDBHost, redisDBPort)
	redisDB := databases.NewRedis(redisAddr, context.Background()) // Not using password as it's a test task

	logger := logrus.New().WithField("Service", "Fibonacci")
	fibonacci := tools.NewFibonacci(redisDB, logger)

	app := app.NewApp(host, httpPort, grpcPort, fibonacci)

	app.Run()

	panic("Returned from app.Run()")
}

package main

import (
	"HomeWork_1/internal/app/web"
	"HomeWork_1/internal/config"
	"HomeWork_1/internal/metric"
	module "HomeWork_1/internal/module/order"
	"HomeWork_1/internal/pkg/database"
	"HomeWork_1/internal/pkg/kafka"
	"HomeWork_1/internal/service/events"
	"HomeWork_1/internal/service/hash_generator"
	"HomeWork_1/internal/service/input_validator"
	"HomeWork_1/internal/storage/postgres"
	"HomeWork_1/internal/tracing"
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	db, err := database.NewDatabase(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.GetPool(ctx).Close()
	ordersRepository := postgres.NewRepository(db)
	hashGenerator := hash_generator.NewHashGenerator()
	orderService := module.NewModule(ordersRepository, &hashGenerator)
	inputValidator := input_validator.NewValidator()

	configDB, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	producer, err := kafka.NewProducer(configDB.Brokers, configDB.KafkaTopic)
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := events.NewReader(configDB.KafkaEnable)
	consumer, err := kafka.NewConsumerGroup(configDB.Brokers, configDB.KafkaGroup, reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	go consumer.Run(configDB.KafkaTopic)

	sender := events.NewSender(configDB.KafkaEnable, producer)

	promMetrics := metric.RegisterMetrics()

	tracing.NewTracer()

	app := web.NewWeb(&orderService, inputValidator, sender, promMetrics)
	app.Run(ctx, configDB)
}

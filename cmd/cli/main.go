package main

import (
	"context"
	"fmt"

	"HomeWork_1/internal/app/cli"
	"HomeWork_1/internal/config"
	module "HomeWork_1/internal/module/order"
	"HomeWork_1/internal/pkg/database"
	"HomeWork_1/internal/pkg/kafka"
	"HomeWork_1/internal/service/events"
	"HomeWork_1/internal/service/hash_generator"
	"HomeWork_1/internal/service/input_validator"
	"HomeWork_1/internal/storage/postgres"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	commands := cli.NewCLI(&orderService, inputValidator, sender)
	if err := commands.Run(ctx); err != nil {
		fmt.Println(err)
	}

}

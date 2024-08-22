package config

import (
	"encoding/json"
	"os"
)

type ConfigInfo struct {
	Host        string   `json:"host"`
	Name        string   `json:"name"`
	Password    string   `json:"password"`
	Port        string   `json:"port"`
	User        string   `json:"user"`
	Brokers     []string `json:"brokers"`
	KafkaGroup  string   `json:"kafka_group"`
	KafkaTopic  string   `json:"kafka_topic"`
	KafkaEnable bool     `json:"kafka_enable"`
	GrpcPort    int      `json:"grpc_port"`
	HttpPort    int      `json:"http_port"`
	SwaggerPort int      `json:"swagger_port"`
}

func Read() (*ConfigInfo, error) {
	b, err := os.ReadFile("config/config.json")
	if err != nil {
		return nil, err
	}
	var databaseInfo ConfigInfo
	if err := json.Unmarshal(b, &databaseInfo); err != nil {
		return nil, err
	}

	return &databaseInfo, nil
}

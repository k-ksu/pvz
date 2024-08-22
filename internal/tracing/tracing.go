package tracing

import (
	"log"

	"github.com/uber/jaeger-client-go/config"
)

func NewTracer() {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	_, err := cfg.InitGlobalTracer("homework-1")
	if err != nil {
		log.Fatal("Cannot init tracing", err)
	}
}

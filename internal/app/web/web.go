package web

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/flowchartsman/swaggerui"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"HomeWork_1/internal/api"
	"HomeWork_1/internal/config"
	"HomeWork_1/internal/model"
	"HomeWork_1/pkg/api/pvz"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50051", "gRPC server endpoint")
)

type orderModule interface {
	GetOrderFromCourier(ctx context.Context, order *model.Order) error
	GiveOrder(ctx context.Context, orders []model.OrderID) (int, error)
	ReturnOrder(ctx context.Context, orderID model.OrderID) error
	ListOrders(ctx context.Context, clientID model.ClientID, action int) ([]model.Order, map[model.PackageType]model.Package, error)
	ReturnFromClient(ctx context.Context, orderID model.OrderID, clientID model.ClientID) error
	ListReturns(ctx context.Context) ([]model.Order, map[model.PackageType]model.Package, error)
	LoadPackagesToCheck(ctx context.Context) ([]model.Package, error)
	GiveOrderWithNewPackage(ctx context.Context, orders []model.OrderID, pack model.PackageType) (int, error)
}

type inputValidator interface {
	ValidateGetOrderFromCourier(orderID, clientID, date, pack string, price, weight int) (*model.Order, error)
	ValidateGiveOrder(orders, pack string, loadedPackages []model.Package) ([]model.OrderID, *model.PackageType, error)
	ValidateReturnOrder(orderID string) (model.OrderID, error)
	ValidateListOrders(clientID, action string) (model.ClientID, error)
	ValidateReturnFromClient(orderID, clientID string) (model.OrderID, model.ClientID, error)
	ValidateListReturns(pageSize, pageNumber int) (int, int, error)
	ValidatePackage(weight int, pack model.PackageType, loadedPackages []model.Package) error
}

type sender interface {
	SendMessage(event *model.EventMessage) error
}

type Web struct {
	pvz.UnimplementedPVZServer

	orderModule    orderModule
	inputValidator inputValidator
	sender         sender
	prom           *prometheus.Registry

	wg *sync.WaitGroup
}

func NewWeb(module orderModule, inputValidator inputValidator, sender sender, prom *prometheus.Registry) *Web {
	return &Web{
		orderModule:    module,
		inputValidator: inputValidator,
		sender:         sender,
		prom:           prom,

		wg: &sync.WaitGroup{},
	}
}

func (w *Web) Run(ctx context.Context, config *config.ConfigInfo) {
	go w.ServeSwaggerUI(config)
	go w.ServeHTTP(ctx, config)
	w.ServeGRPC(config)
}

func (w *Web) ServeGRPC(config *config.ConfigInfo) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcPort))
	if err != nil {
		fmt.Println(err)
		return
	}

	grpcServer := grpc.NewServer()

	pvz.RegisterPVZServer(
		grpcServer,
		api.NewHandler(w.orderModule, w.inputValidator, w.sender),
	)

	reflection.Register(grpcServer)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-sigCh
		log.Printf("got signal %v, graceful shutdown", s)
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Println(err)
		return
	}
}

func (w *Web) ServeHTTP(ctx context.Context, config *config.ConfigInfo) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := pvz.RegisterPVZHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		fmt.Println(err)
		return
	}

	gwServer := &http.Server{
		Addr: fmt.Sprintf(":%d", config.HttpPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			mux.ServeHTTP(w, r)
		}),
	}

	if err = gwServer.ListenAndServe(); err != nil {
		fmt.Println(err)
		return
	}
}

func (w *Web) ServeSwaggerUI(config *config.ConfigInfo) {
	spec, err := os.ReadFile("pkg/api/pvz/pvz.swagger.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	specWithHost, err := addHost(config, spec)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(specWithHost)))
	http.Handle("/metrics", promhttp.HandlerFor(w.prom, promhttp.HandlerOpts{EnableOpenMetrics: true}))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.SwaggerPort), nil))
}

func addHost(config *config.ConfigInfo, spec []byte) ([]byte, error) {
	var swagger map[string]any
	err := json.Unmarshal(spec, &swagger)
	if err != nil {
		return nil, err
	}

	swagger["host"] = fmt.Sprintf(":%d", config.HttpPort)
	return json.Marshal(swagger)
}

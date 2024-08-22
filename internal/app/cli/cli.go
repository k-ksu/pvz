package cli

import (
	"HomeWork_1/internal/model/errs"
	"context"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"

	"HomeWork_1/internal/model"
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

type CLI struct {
	module         orderModule
	inputValidator inputValidator
	sender         sender
	commandList    []command
	workerCount    int
}

func NewCLI(module orderModule, inputValidator inputValidator, sender sender) CLI {
	return CLI{
		module:         module,
		inputValidator: inputValidator,
		sender:         sender,
		commandList: []command{
			{
				name:        help,
				description: "\nсписок доступных функций, \nиспользование help\n",
			},
			{
				name:        getOrderFromCourier,
				description: "\nпринять заказ от курьера (вес товара в граммах, упаковки: plasticBag, box, film), \nиспользование get --orderID=78 --clientID=67 --date=28.08.2024-14:40 --price=45 --weight=7000 --package=film\n",
			},
			{
				name:        returnOrder,
				description: "\nвернуть заказ курьеру, \nиспользование return --orderID=78\n",
			},
			{
				name:        giveOrder,
				description: "\nвыдать заказы клиенту, \nесли не введете параметры для упаковки, останется предыдущая, \nможно выбрать только 1 новый тип упаковки, \nтипы упаковки: plasticBag, box, film \nиспользование give --orders=1,56,23 --package=film\n",
			}, {
				name: listOrders,
				description: "\nпосмотреть заказы \n1) последние N заказов \n2) заказы клиента, " +
					"находящиеся в нашем ПВЗ, \nиспользование listOrders --clientID=78 --action=1\n",
			},
			{
				name: returnFromClient,
				description: "\nвозврат заказа от клиента. Доступно если прошло не более 2х дней с момента получения," +
					" \nиспользование returnFromClient --clientID=45 --orderID=66\n",
			},
			{
				name: listReturns,
				description: "\nпосмотреть список возвратов (необходмо дополнительно ввести размер страницы и номер страницы)," +
					"\nиспользование listReturns --pageSize=10 --pageNumber=1\n",
			}, {
				name: exit,
				description: "\nвайти из программы,\n" +
					"использвание exit\n",
			},
			{
				name:        changeWorkerCount,
				description: "\nуправление кол-вом горутин, \nиспользование changeWorkerCount --number=1\n",
			},
		},
		workerCount: 1,
	}
}

func (c *CLI) Run(ctx context.Context) error {
	done := make(chan os.Signal, 1)
	wg := sync.WaitGroup{}
	args := os.Args[1:]
	logs := make(chan string)
	commands := make(chan []string)
	if len(args) == 0 {
		c.CliModeRun(ctx, done, commands, logs, &wg)
		return nil
	}
	var err error
	currentCommand := args[0]
	switch currentCommand {
	case help:
		wg.Add(1)
		c.help(ctx, &wg)
	case changeWorkerCount:
		wg.Add(1)
		c.changeWorkerCount(ctx, args[1:], done, commands, logs, &wg)
	case getOrderFromCourier:
		wg.Add(1)
		c.getOrderFromCourier(ctx, args[1:], logs, &wg)
	case returnOrder:
		wg.Add(1)
		c.returnOrder(ctx, args[1:], logs, &wg)
	case giveOrder:
		wg.Add(1)
		c.giveOrder(ctx, args[1:], logs, &wg)
	case listOrders:
		wg.Add(1)
		c.listOrders(ctx, args[1:], logs, &wg)
	case returnFromClient:
		wg.Add(1)
		c.returnFromClient(ctx, args[1:], logs, &wg)
	case listReturns:
		wg.Add(1)
		c.listReturns(ctx, args[1:], logs, &wg)
	default:
		err = fmt.Errorf("command isn't one of listed")
	}
	if err != nil {
		return err
	}

	return nil
}

func (c *CLI) help(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("list of available commands:")
	for _, com := range c.commandList {
		fmt.Println("", com.name, com.description)
	}
}

func (c *CLI) changeWorkerCount(ctx context.Context, args []string, done chan os.Signal, commands <-chan []string, logs chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	var number int

	fs := flag.NewFlagSet(getOrderFromCourier, flag.ContinueOnError)
	fs.IntVar(&number, "number", 1, "use --number=1")

	if err := fs.Parse(args); err != nil {
		logs <- err.Error()
		return
	}

	if number < 1 {
		logs <- "недопустимое количество горутин"
		return
	}

	for i := 0; i < number; i++ {
		wg.Add(1)
		go subscribe(ctx, c, commands, logs, done, wg)
	}

	c.workerCount += number
	logs <- fmt.Sprintf("changeWorkerCount successfully finished with number of workers: %v", c.workerCount)
}

func (c *CLI) getOrderFromCourier(ctx context.Context, args []string, logs chan<- string, wg *sync.WaitGroup) {
	var orderID, clientID, date, pack string
	var price, weight int

	fs := flag.NewFlagSet(getOrderFromCourier, flag.ContinueOnError)
	fs.StringVar(&orderID, "orderID", "", "use --orderID=6789")
	fs.StringVar(&clientID, "clientID", "", "use --clientID=456")
	fs.StringVar(&date, "date", "", "use --date=02.01.2006-15:04")
	fs.StringVar(&pack, "package", "", "use --package=film")
	fs.IntVar(&price, "price", 0, "use --price=456")
	fs.IntVar(&weight, "weight", 0, "use --weight=2067")

	if err := fs.Parse(args); err != nil {
		logs <- err.Error()
		return
	}

	order, err := c.inputValidator.ValidateGetOrderFromCourier(orderID, clientID, date, pack, price, weight)

	if err != nil && !errors.Is(err, errs.ErrPackageDoesNotSet) {
		logs <- err.Error()
		return
	}

	if order.Package != model.WithoutPackage {
		wg.Add(1)
		loadedPackages, err := c.module.LoadPackagesToCheck(ctx)
		wg.Done()
		if err != nil {
			logs <- err.Error()
			return
		}
		err = c.inputValidator.ValidatePackage(weight, order.Package, loadedPackages)
		if err != nil {
			logs <- "GetOrderFromCourier finished with error: " + err.Error()
			return
		}
	}

	go func() {
		defer wg.Done()

		order.MaxWeight = weight
		order.Price = price

		err := c.module.GetOrderFromCourier(ctx, order)
		if err != nil {
			logs <- "GetOrderFromCourier finished with error: " + err.Error()
		} else {
			logs <- "GetOrderFromCourier finished successfully"
		}
	}()
}

func (c *CLI) giveOrder(ctx context.Context, args []string, logs chan<- string, wg *sync.WaitGroup) {
	var orders, pack string

	fs := flag.NewFlagSet(giveOrder, flag.ContinueOnError)
	fs.StringVar(&orders, "orders", "", "use --orders=1,6,45")
	fs.StringVar(&pack, "package", "", "use --package=film")

	if err := fs.Parse(args); err != nil {
		logs <- err.Error()
		return
	}

	wg.Add(1)
	loadedPackages, err := c.module.LoadPackagesToCheck(ctx)
	wg.Done()
	if err != nil {
		logs <- err.Error()
		return
	}

	readyListOfOrders, readyPack, err := c.inputValidator.ValidateGiveOrder(orders, pack, loadedPackages)
	if err != nil && !errors.Is(err, errs.ErrPackageDoesNotSet) {
		logs <- err.Error()
		return
	}

	go func() {
		defer wg.Done()

		var summa int
		if errors.Is(err, errs.ErrPackageDoesNotSet) {
			summa, err = c.module.GiveOrder(ctx, readyListOfOrders)
		} else {
			summa, err = c.module.GiveOrderWithNewPackage(ctx, readyListOfOrders, *readyPack)
		}
		if err != nil {
			logs <- "giveOrder finished with error: " + err.Error()
		} else {
			logs <- fmt.Sprintf("giveOrder finished successfully. Amount to be paid %v", summa)
		}
	}()
}

func (c *CLI) returnOrder(ctx context.Context, args []string, logs chan<- string, wg *sync.WaitGroup) {
	var orderID string

	fs := flag.NewFlagSet(returnOrder, flag.ContinueOnError)
	fs.StringVar(&orderID, "orderID", "", "use --orderID=6789")
	if err := fs.Parse(args); err != nil {
		logs <- err.Error()
		return
	}

	readyOrderID, err := c.inputValidator.ValidateReturnOrder(orderID)

	if err != nil {
		logs <- err.Error()
		return
	}
	go func() {
		defer wg.Done()

		err := c.module.ReturnOrder(ctx, readyOrderID)
		if err != nil {
			logs <- "ReturnOrder finished with error: " + err.Error()
		} else {
			logs <- "ReturnOrder finished successfully"
		}
	}()
}

func (c *CLI) listOrders(ctx context.Context, args []string, logs chan<- string, wg *sync.WaitGroup) {
	var clientID string
	var action string

	fs := flag.NewFlagSet(listOrders, flag.ContinueOnError)
	fs.StringVar(&clientID, "clientID", "", "use --clientID=456")
	fs.StringVar(&action, "action", "", "use --action=1")
	if err := fs.Parse(args); err != nil {
		logs <- err.Error()
		return
	}

	readyClientID, err := c.inputValidator.ValidateListOrders(clientID, action)
	if err != nil {
		logs <- err.Error()
		return
	}

	go func() {
		defer wg.Done()

		actionInt, err := strconv.ParseInt(action, 10, 64)
		if err != nil {
			logs <- err.Error()
			return
		}
		listOfOrders, packages, err := c.module.ListOrders(ctx, readyClientID, int(actionInt))
		if err != nil {
			logs <- "ListOrders finished with error: " + err.Error()
			return
		}

		for _, order := range listOfOrders {
			var maxW string
			fmt.Printf("ID заказа: %s\nID клиента: %s\nСостояние заказа: %s\nСрок хранения до: %s\nДата выдачи: %s\nВес: %v\nЦена: %v\n",
				order.OrderID, order.ClientID, order.Condition, order.ArrivedAt, order.ReceivedAt, order.MaxWeight, order.Price)
			if packages[order.Package].PackageMaxWeight == -1 {
				fmt.Printf("Упаковка заказа: %s\nСтоимость упаковки: %v\nМаксимальный вес для упаковки: нет огрничения\n",
					order.Package, packages[order.Package].PackageSurcharge)
			} else {
				fmt.Printf("Упаковка заказа: %s\nСтоимость упаковки: %v\nМаксимальный вес для упаковки: %v\n",
					order.Package, packages[order.Package].PackageSurcharge, packages[order.Package].PackageMaxWeight)
			}
			fmt.Printf("Упаковка заказа: %s\nСтоимость упаковки: %v\nМаксимальный вес для упаковки: %v\n",
				order.Package, packages[order.Package].PackageSurcharge, maxW)
			fmt.Printf("Итоговая стоимость с упаковой: %v\n", packages[order.Package].PackageSurcharge+order.Price)
			fmt.Println()
		}
		logs <- "ListOrders finished successfully"
	}()
}

func (c *CLI) returnFromClient(ctx context.Context, args []string, logs chan<- string, wg *sync.WaitGroup) {
	var orderID, clientID string

	fs := flag.NewFlagSet(returnFromClient, flag.ContinueOnError)
	fs.StringVar(&orderID, "orderID", "", "use --orderID=6789")
	fs.StringVar(&clientID, "clientID", "", "use --clientID=456")
	if err := fs.Parse(args); err != nil {
		logs <- err.Error()
		return
	}

	readyOrderID, readyClientID, err := c.inputValidator.ValidateReturnFromClient(orderID, clientID)
	if err != nil {
		logs <- err.Error()
		return
	}
	go func() {
		defer wg.Done()

		err := c.module.ReturnFromClient(ctx, readyOrderID, readyClientID)
		if err != nil {
			logs <- "ReturnFromClient finished with error: " + err.Error()
		} else {
			logs <- "ReturnFromClient finished successfully"
		}
	}()
}

func (c *CLI) listReturns(ctx context.Context, args []string, logs chan<- string, wg *sync.WaitGroup) {
	var pageSize, pageNumber int

	fs := flag.NewFlagSet(listReturns, flag.ContinueOnError)
	fs.IntVar(&pageSize, "pageSize", 5, "use --pageSize=SomeSize")
	fs.IntVar(&pageNumber, "pageNumber", 1, "use --pageNumber=SomeNumber")
	if err := fs.Parse(args); err != nil {
		logs <- err.Error()
		return
	}

	pageSize, pageNumber, err := c.inputValidator.ValidateListReturns(pageSize, pageNumber)
	if err != nil {
		logs <- err.Error()
		return
	}

	go func() {
		defer wg.Done()

		returns, packages, err := c.module.ListReturns(ctx)
		if err != nil {
			logs <- "ListReturns finished with error: " + err.Error()
			return
		}
		maxPage := int(math.Ceil(float64(len(returns)) / float64(pageSize)))

		if pageNumber > maxPage {
			fmt.Printf("Всего возвратов: %v\nВсего страниц: %v\nРазмер страницы: %v\nНомер старницы: %v",
				len(returns), maxPage, pageSize, pageNumber)
			logs <- "ListReturns finished with error: " + errs.ErrPageNumberIsLarge.Error()
			return
		}

		fmt.Printf("Всего возвратов: %v\nВсего страниц: %v\nРазмер страницы: %v\nНомер старницы: %v\n",
			len(returns), maxPage, pageSize, pageNumber)
		for i := (pageNumber - 1) * pageSize; i < min(pageNumber*pageSize, len(returns)); i++ {
			fmt.Printf("ID заказа: %s\nID клиента: %s\nСостояние заказа: %s\nСрок хранения до: %s\nДата выдачи: %s\nВес: %v\nЦена: %v\n",
				returns[i].OrderID, returns[i].ClientID, returns[i].Condition, returns[i].ArrivedAt, returns[i].ReceivedAt, returns[i].MaxWeight, returns[i].Price)
			if packages[returns[i].Package].PackageMaxWeight == -1 {
				fmt.Printf("Упаковка заказа: %s\nСтоимость упаковки: %v\nМаксимальный вес для упаковки: нет отграничения\n",
					returns[i].Package, packages[returns[i].Package].PackageSurcharge)
			} else {
				fmt.Printf("Упаковка заказа: %s\nСтоимость упаковки: %v\nМаксимальный вес для упаковки: %v\n",
					returns[i].Package, packages[returns[i].Package].PackageSurcharge, packages[returns[i].Package].PackageMaxWeight)
			}
			fmt.Printf("Итоговая стоимость с упаковой: %v\n", packages[returns[i].Package].PackageSurcharge+returns[i].Price)
			fmt.Println()
		}

		logs <- "ListReturns finished successfully"
	}()
}

package cli

const (
	help                = "help"             //Список команд
	getOrderFromCourier = "get"              //Принять заказ от курьера
	returnOrder         = "return"           //Вернуть заказ курьеру
	giveOrder           = "give"             //Выдать заказ клиенту
	listOrders          = "listOrders"       //Получить список заказов
	returnFromClient    = "returnFromClient" //Принять возврат от клиента
	listReturns         = "listReturns"      //Получить список возвратов
	exit                = "exit"
	changeWorkerCount   = "changeWorkerCount"
)

type command struct {
	name        string
	description string
}

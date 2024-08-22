package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"HomeWork_1/internal/model"
)

func (c *CLI) CliModeRun(ctx context.Context, done chan os.Signal, commands chan []string, logs chan string, wg *sync.WaitGroup) {
	go logger(logs, wg)
	go publish(commands)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	subscribe(ctx, c, commands, logs, done, wg)
	wg.Wait() // дожидаемся выполнение всех горутин

	wg.Add(1)
	close(logs)
	wg.Wait() // дожидаемся вывода всех сообщения из логера

	close(commands)
	wg.Wait()
}

func logger(logs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for msg := range logs {
		fmt.Println()
		fmt.Println(msg)
	}
}

func publish(commands chan<- []string) {
	fmt.Print("Введите команду или exit чтобы выйти: ")
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		cmd := scanner.Text()
		cmds := strings.Split(cmd, " ")
		commands <- cmds
	}
}
func subscribe(ctx context.Context, c *CLI, commands <-chan []string, logs chan<- string, done chan os.Signal, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case cmd := <-commands:
			err := c.sender.SendMessage(&model.EventMessage{
				Method:    cmd[0],
				Args:      cmd[1:],
				TimeStamp: time.Now(),
			})
			if err != nil {
				return
			}

			switch cmd[0] {
			case changeWorkerCount:
				wg.Add(1)
				c.changeWorkerCount(ctx, cmd[1:], done, commands, logs, wg)
			case exit:
				done <- os.Kill
				return
			case help:
				wg.Add(1)
				c.help(ctx, wg)
			case getOrderFromCourier:
				wg.Add(1)
				c.getOrderFromCourier(ctx, cmd[1:], logs, wg)
			case returnOrder:
				wg.Add(1)
				c.returnOrder(ctx, cmd[1:], logs, wg)
			case giveOrder:
				wg.Add(1)
				c.giveOrder(ctx, cmd[1:], logs, wg)
			case listOrders:
				wg.Add(1)
				c.listOrders(ctx, cmd[1:], logs, wg)
			case returnFromClient:
				wg.Add(1)
				c.returnFromClient(ctx, cmd[1:], logs, wg)
			case listReturns:
				wg.Add(1)
				c.listReturns(ctx, cmd[1:], logs, wg)
			default:
				logs <- "Command is not one of the listed"
			}
		case d := <-done:
			done <- d // чтобы завершить всех воркеров
			logs <- "завершение работы всех горутин"
			return
		}
	}
}

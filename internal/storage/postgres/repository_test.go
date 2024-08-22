package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"HomeWork_1/internal/model"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	repo *Repository
}

func (s *RepositoryTestSuite) SetupSuite() {
	s.repo = testRepo
	time.Local = time.UTC
}

func (s *RepositoryTestSuite) TestAppendOrder() {
	var rowsCnt int
	var rowsCnt2 int
	ctx := context.Background()
	err := s.repo.db.Get(ctx, &rowsCnt, "SELECT COUNT(*) FROM orders")
	require.NoError(s.T(), err)

	err = s.repo.AppendOrder(ctx, model.Order{
		OrderID:    "test_id",
		ClientID:   "1",
		Condition:  "accepted",
		ArrivedAt:  time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC),
		ReceivedAt: time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC),
		Price:      0,
		Package:    "box",
		MaxWeight:  0,
	}, "someHash")
	defer s.clearTable(ctx, "orders", "id", "test_id")
	require.NoError(s.T(), err)

	err = s.repo.db.Get(ctx, &rowsCnt2, "SELECT COUNT(*) FROM orders")
	require.NoError(s.T(), err)

	require.Equal(s.T(), rowsCnt+1, rowsCnt2)
}

func (s *RepositoryTestSuite) TestOrderByID() {
	ctx := context.Background()
	wantOrder := model.Order{
		OrderID:    "test_id",
		ClientID:   "1",
		Condition:  "accepted",
		ArrivedAt:  time.Date(0, time.January, 1, 0, 0, 0, 0, time.Local),
		ReceivedAt: time.Date(0, time.January, 1, 0, 0, 0, 0, time.Local),
		Price:      0,
		Package:    "box",
		MaxWeight:  0,
	}

	err := s.repo.AppendOrder(ctx, wantOrder, "someHash")
	defer s.clearTable(ctx, "orders", "id", "test_id")
	require.NoError(s.T(), err)

	resOrder, err := s.repo.OrderByID(ctx, "test_id", "someHash")
	require.NoError(s.T(), err)
	// check everything except dates
	wantOrder.ArrivedAt = resOrder.ArrivedAt
	wantOrder.ReceivedAt = resOrder.ReceivedAt
	require.Equal(s.T(), wantOrder, *resOrder)
}

func (s *RepositoryTestSuite) TestUpdateOrders() {
	ctx := context.Background()
	wantOrder := model.Order{
		OrderID:    "test_id",
		ClientID:   "1",
		Condition:  "accepted",
		ArrivedAt:  time.Time{},
		ReceivedAt: time.Time{},
		Price:      0,
		Package:    "box",
		MaxWeight:  0,
	}

	newWantOrder := model.Order{
		OrderID:    "test_id",
		ClientID:   "1",
		Condition:  "accepted",
		ArrivedAt:  time.Time{},
		ReceivedAt: time.Time{},
		Price:      0,
		Package:    "film",
		MaxWeight:  0,
	}

	err := s.repo.AppendOrder(ctx, wantOrder, "someHash")
	defer s.clearTable(ctx, "orders", "id", "test_id")
	require.NoError(s.T(), err)

	err = s.repo.UpdateOrders(ctx, []model.Order{newWantOrder}, "someHash")
	require.NoError(s.T(), err)
	resOrder, err := s.repo.OrderByID(ctx, "test_id", "someHash")
	require.NoError(s.T(), err)
	// check everything except dates
	newWantOrder.ArrivedAt = resOrder.ArrivedAt
	newWantOrder.ReceivedAt = resOrder.ReceivedAt
	require.Equal(s.T(), newWantOrder, *resOrder)
}

func (s *RepositoryTestSuite) TestLoadOrders() {
	ctx := context.Background()
	wantOrder := model.Order{
		OrderID:    "test_id",
		ClientID:   "1",
		Condition:  "accepted",
		ArrivedAt:  time.Time{},
		ReceivedAt: time.Time{},
		Price:      0,
		Package:    "box",
		MaxWeight:  0,
	}

	err := s.repo.AppendOrder(ctx, wantOrder, "someHash")
	defer s.clearTable(ctx, "orders", "id", "test_id")
	require.NoError(s.T(), err)

	counter := 0
	orders, err := s.repo.LoadOrders(ctx, "someHash")
	require.NoError(s.T(), err)
	for _, order := range orders {
		if order.OrderID == wantOrder.OrderID {
			counter += 1
		}
	}
	require.Equal(s.T(), counter, 1)
}

func (s *RepositoryTestSuite) TestDeleteOrder() {
	ctx := context.Background()
	wantOrder := model.Order{
		OrderID:    "test_id",
		ClientID:   "1",
		Condition:  "accepted",
		ArrivedAt:  time.Time{},
		ReceivedAt: time.Time{},
		Price:      0,
		Package:    "box",
		MaxWeight:  0,
	}

	err := s.repo.AppendOrder(ctx, wantOrder, "someHash")
	require.NoError(s.T(), err)

	err = s.repo.DeleteOrder(ctx, "test_id", "someHash")
	require.NoError(s.T(), err)

	counter := 0
	orders, err := s.repo.LoadOrders(ctx, "someHash")
	require.NoError(s.T(), err)
	for _, order := range orders {
		if order.OrderID == wantOrder.OrderID {
			counter += 1
		}
	}
	require.Equal(s.T(), counter, 0)
}

func (s *RepositoryTestSuite) TestOrdersWithRefundCondition() {
	ctx := context.Background()
	wantOrder := model.Order{
		OrderID:    "test_id",
		ClientID:   "1",
		Condition:  "refund",
		ArrivedAt:  time.Time{},
		ReceivedAt: time.Time{},
		Price:      0,
		Package:    "box",
		MaxWeight:  0,
	}

	err := s.repo.AppendOrder(ctx, wantOrder, "someHash")
	defer s.clearTable(ctx, "orders", "id", "test_id")

	require.NoError(s.T(), err)

	counter := 0
	orders, err := s.repo.OrdersWithRefundCondition(ctx, "someHash")
	require.NoError(s.T(), err)
	for _, order := range orders {
		if order.OrderID == wantOrder.OrderID {
			counter += 1
		}
	}
	require.Equal(s.T(), counter, 1)
}

func (s *RepositoryTestSuite) TestOrderByClientIDOrderID() {
	ctx := context.Background()
	wantOrder := model.Order{
		OrderID:    "test_id",
		ClientID:   "1",
		Condition:  "refund",
		ArrivedAt:  time.Time{}.UTC(),
		ReceivedAt: time.Time{}.UTC(),
		Price:      0,
		Package:    "box",
		MaxWeight:  0,
	}

	err := s.repo.AppendOrder(ctx, wantOrder, "someHash")
	defer s.clearTable(ctx, "orders", "id", "test_id")
	require.NoError(s.T(), err)

	order, err := s.repo.OrderByClientIDOrderID("test_id", "1", ctx, "someHash")
	require.NoError(s.T(), err)
	// check everything except dates
	wantOrder.ArrivedAt = order.ArrivedAt
	wantOrder.ReceivedAt = order.ReceivedAt
	require.Equal(s.T(), wantOrder, *order)
}

func (s *RepositoryTestSuite) TestOrdersWithGivenClientID() {
	ctx := context.Background()
	wantOrder := model.Order{
		OrderID:    "test_id",
		ClientID:   "1",
		Condition:  "refund",
		ArrivedAt:  time.Time{},
		ReceivedAt: time.Time{},
		Price:      0,
		Package:    "box",
		MaxWeight:  0,
	}

	err := s.repo.AppendOrder(ctx, wantOrder, "someHash")
	defer s.clearTable(ctx, "orders", "id", "test_id")
	require.NoError(s.T(), err)

	counter := 0
	orders, err := s.repo.OrdersWithGivenClientID("1", ctx, "someHash")
	require.NoError(s.T(), err)
	for _, order := range orders {
		if order.ClientID == wantOrder.ClientID {
			counter += 1
		}
	}
	require.Equal(s.T(), counter, len(orders))
}

func (s *RepositoryTestSuite) TestOrdersWithGivenClientIDAndAcceptedCondition() {
	ctx := context.Background()
	wantOrder := model.Order{
		OrderID:    "test_id",
		ClientID:   "1",
		Condition:  "refund",
		ArrivedAt:  time.Time{},
		ReceivedAt: time.Time{},
		Price:      0,
		Package:    "box",
		MaxWeight:  0,
	}

	err := s.repo.AppendOrder(ctx, wantOrder, "someHash")
	defer s.clearTable(ctx, "orders", "id", "test_id")
	require.NoError(s.T(), err)

	counter := 0
	orders, err := s.repo.OrdersWithGivenClientIDAndAcceptedCondition("1", ctx, "someHash")
	require.NoError(s.T(), err)
	for _, order := range orders {
		if order.ClientID == wantOrder.ClientID && order.Condition == model.ConditionAccepted {
			counter += 1
		}
	}
	require.Equal(s.T(), counter, len(orders))
}

func (s *RepositoryTestSuite) TestOrdersByGivenOrderIDList() {
	ctx := context.Background()
	wantOrder := model.Order{
		OrderID:    "test_id",
		ClientID:   "1",
		Condition:  "refund",
		ArrivedAt:  time.Time{},
		ReceivedAt: time.Time{},
		Price:      0,
		Package:    "box",
		MaxWeight:  0,
	}

	err := s.repo.AppendOrder(ctx, wantOrder, "someHash")
	defer s.clearTable(ctx, "orders", "id", "test_id")
	require.NoError(s.T(), err)

	orders, err := s.repo.OrdersByGivenOrderIDList([]model.OrderID{"test_id"}, ctx, "someHash")
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, len(orders))
}

func (s *RepositoryTestSuite) TestLoadPackages() {
	ctx := context.Background()

	packs, err := s.repo.LoadPackages(ctx)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 4, len(packs))
}

func (s *RepositoryTestSuite) TestLoadPackByType() {
	ctx := context.Background()

	packs, err := s.repo.LoadPackByType(ctx, model.Box)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 30000, packs.PackageMaxWeight)
}

func (s *RepositoryTestSuite) clearTable(ctx context.Context, tableName, columnName, columnValue string) {
	_, err := s.repo.db.Cluster.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE %s = '%s'", tableName, columnName, columnValue))
	require.NoError(s.T(), err)
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

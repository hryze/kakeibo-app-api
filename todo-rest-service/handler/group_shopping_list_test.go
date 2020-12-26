package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/todo-rest-service/testutil"
)

type MockGroupShoppingListRepository struct{}

func (m MockGroupShoppingListRepository) GetGroupRegularShoppingList(groupID int) (model.GroupRegularShoppingList, error) {
	if dbCounter == 1 {
		atomic.AddInt64(&dbCounter, -1)

		return model.GroupRegularShoppingList{
			GroupRegularShoppingList: []model.GroupRegularShoppingItem{
				{
					ID:                   1,
					PostedDate:           time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
					UpdatedDate:          time.Date(2020, 12, 19, 20, 0, 0, 0, time.UTC),
					ExpectedPurchaseDate: model.Date{Time: time.Date(2020, 12, 25, 0, 0, 0, 0, time.UTC)},
					CycleType:            "monthly",
					Cycle:                model.NullInt{Int: 0, Valid: false},
					Purchase:             "米",
					Shop:                 model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
					Amount:               model.NullInt64{NullInt64: sql.NullInt64{Int64: 4000, Valid: true}},
					BigCategoryID:        2,
					BigCategoryName:      "",
					MediumCategoryID:     model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
					MediumCategoryName:   model.NullString{NullString: sql.NullString{String: "", Valid: false}},
					CustomCategoryID:     model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
					CustomCategoryName:   model.NullString{NullString: sql.NullString{String: "", Valid: false}},
					PaymentUserID:        model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
					TransactionAutoAdd:   true,
				},
				{
					ID:                   2,
					PostedDate:           time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
					UpdatedDate:          time.Date(2020, 12, 19, 20, 0, 0, 0, time.UTC),
					ExpectedPurchaseDate: model.Date{Time: time.Date(2020, 12, 25, 0, 0, 0, 0, time.UTC)},
					CycleType:            "weekly",
					Cycle:                model.NullInt{Int: 0, Valid: false},
					Purchase:             "トイレットペーパー",
					Shop:                 model.NullString{NullString: sql.NullString{String: "クリエイト", Valid: true}},
					Amount:               model.NullInt64{NullInt64: sql.NullInt64{Int64: 300, Valid: true}},
					BigCategoryID:        3,
					BigCategoryName:      "",
					MediumCategoryID:     model.NullInt64{NullInt64: sql.NullInt64{Int64: 13, Valid: true}},
					MediumCategoryName:   model.NullString{NullString: sql.NullString{String: "", Valid: false}},
					CustomCategoryID:     model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
					CustomCategoryName:   model.NullString{NullString: sql.NullString{String: "", Valid: false}},
					PaymentUserID:        model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
					TransactionAutoAdd:   true,
				},
			},
		}, nil
	}

	atomic.AddInt64(&dbCounter, 1)

	return model.GroupRegularShoppingList{
		GroupRegularShoppingList: []model.GroupRegularShoppingItem{
			{
				ID:                   1,
				PostedDate:           time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:          time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				ExpectedPurchaseDate: model.Date{Time: time.Date(2020, 12, 25, 0, 0, 0, 0, time.UTC)},
				CycleType:            "monthly",
				Cycle:                model.NullInt{Int: 0, Valid: false},
				Purchase:             "米",
				Shop:                 model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Amount:               model.NullInt64{NullInt64: sql.NullInt64{Int64: 4000, Valid: true}},
				BigCategoryID:        2,
				BigCategoryName:      "",
				MediumCategoryID:     model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				MediumCategoryName:   model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:     model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				CustomCategoryName:   model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:        model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				TransactionAutoAdd:   true,
			},
			{
				ID:                   2,
				PostedDate:           time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:          time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				ExpectedPurchaseDate: model.Date{Time: time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC)},
				CycleType:            "weekly",
				Cycle:                model.NullInt{Int: 0, Valid: false},
				Purchase:             "トイレットペーパー",
				Shop:                 model.NullString{NullString: sql.NullString{String: "クリエイト", Valid: true}},
				Amount:               model.NullInt64{NullInt64: sql.NullInt64{Int64: 300, Valid: true}},
				BigCategoryID:        3,
				BigCategoryName:      "",
				MediumCategoryID:     model.NullInt64{NullInt64: sql.NullInt64{Int64: 13, Valid: true}},
				MediumCategoryName:   model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:     model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName:   model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:        model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				TransactionAutoAdd:   true,
			},
		},
	}, nil
}

func (m MockGroupShoppingListRepository) GetGroupRegularShoppingItem(groupRegularShoppingItemID int) (model.GroupRegularShoppingItem, error) {
	return model.GroupRegularShoppingItem{
		ID:                   1,
		PostedDate:           time.Date(2020, 9, 6, 14, 4, 52, 0, time.UTC),
		UpdatedDate:          time.Date(2020, 9, 6, 14, 4, 52, 0, time.UTC),
		ExpectedPurchaseDate: model.Date{Time: time.Date(2020, 9, 13, 0, 0, 0, 0, time.UTC)},
		CycleType:            "weekly",
		Cycle:                model.NullInt{Int: 0, Valid: false},
		Purchase:             "トイレットペーパー",
		Shop:                 model.NullString{NullString: sql.NullString{String: "クリエイト", Valid: true}},
		Amount:               model.NullInt64{NullInt64: sql.NullInt64{Int64: 300, Valid: true}},
		BigCategoryID:        3,
		BigCategoryName:      "",
		MediumCategoryID:     model.NullInt64{NullInt64: sql.NullInt64{Int64: 13, Valid: true}},
		MediumCategoryName:   model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		CustomCategoryID:     model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
		CustomCategoryName:   model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		PaymentUserID:        model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
		TransactionAutoAdd:   true,
	}, nil
}

func (m MockGroupShoppingListRepository) GetGroupShoppingListRelatedToGroupRegularShoppingItem(todayGroupShoppingItemID int, laterThanTodayGroupShoppingItemID int) (model.GroupShoppingList, error) {
	return model.GroupShoppingList{
		GroupShoppingList: []model.GroupShoppingItem{
			{
				ID:                     1,
				PostedDate:             time.Date(2020, 9, 6, 14, 4, 52, 0, time.UTC),
				UpdatedDate:            time.Date(2020, 9, 6, 14, 4, 52, 0, time.UTC),
				ExpectedPurchaseDate:   model.Date{Time: time.Date(2020, 9, 6, 0, 0, 0, 0, time.UTC)},
				CompleteFlag:           false,
				Purchase:               "トイレットペーパー",
				Shop:                   model.NullString{NullString: sql.NullString{String: "クリエイト", Valid: true}},
				Amount:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 300, Valid: true}},
				BigCategoryID:          3,
				BigCategoryName:        "",
				MediumCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 13, Valid: true}},
				MediumCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				PaymentUserID:          model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
			{
				ID:                     2,
				PostedDate:             time.Date(2020, 9, 6, 14, 4, 52, 0, time.UTC),
				UpdatedDate:            time.Date(2020, 9, 6, 14, 4, 52, 0, time.UTC),
				ExpectedPurchaseDate:   model.Date{Time: time.Date(2020, 9, 13, 0, 0, 0, 0, time.UTC)},
				CompleteFlag:           false,
				Purchase:               "トイレットペーパー",
				Shop:                   model.NullString{NullString: sql.NullString{String: "クリエイト", Valid: true}},
				Amount:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 300, Valid: true}},
				BigCategoryID:          3,
				BigCategoryName:        "",
				MediumCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 13, Valid: true}},
				MediumCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				PaymentUserID:          model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
		},
	}, nil
}

func (m MockGroupShoppingListRepository) PostGroupRegularShoppingItem(groupRegularShoppingItem *model.GroupRegularShoppingItem, groupID int, today time.Time) (sql.Result, sql.Result, sql.Result, error) {
	return MockSqlResult{}, MockSqlResult{}, MockSqlResult{}, nil
}

func (m MockGroupShoppingListRepository) PutGroupRegularShoppingItem(groupRegularShoppingItem *model.GroupRegularShoppingItem, groupRegularShoppingItemID int, groupID int, today time.Time) (sql.Result, sql.Result, error) {
	return MockSqlResult{}, MockSqlResult{}, nil
}

func (m MockGroupShoppingListRepository) PutGroupRegularShoppingList(groupRegularShoppingList model.GroupRegularShoppingList, groupID int, today time.Time) error {
	return nil
}

func (m MockGroupShoppingListRepository) DeleteGroupRegularShoppingItem(groupRegularShoppingItemID int) error {
	return nil
}

func (m MockGroupShoppingListRepository) GetDailyGroupShoppingListByDay(date time.Time, groupID int) (model.GroupShoppingList, error) {
	return model.GroupShoppingList{
		GroupShoppingList: []model.GroupShoppingItem{
			{
				ID:                    1,
				PostedDate:            time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:           time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				ExpectedPurchaseDate:  model.Date{Time: time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC)},
				CompleteFlag:          true,
				Purchase:              "米",
				Shop:                  model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Amount:                model.NullInt64{NullInt64: sql.NullInt64{Int64: 4000, Valid: true}},
				BigCategoryID:         2,
				BigCategoryName:       "",
				MediumCategoryID:      model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				MediumCategoryName:    model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:      model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				CustomCategoryName:    model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				RegularShoppingListID: model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				PaymentUserID:         model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				TransactionAutoAdd:    true,
				RelatedTransactionData: &model.GroupTransactionData{
					ID: model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				},
			},
			{
				ID:                     2,
				PostedDate:             time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:            time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				ExpectedPurchaseDate:   model.Date{Time: time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC)},
				CompleteFlag:           false,
				Purchase:               "トイレットペーパー",
				Shop:                   model.NullString{NullString: sql.NullString{String: "クリエイト", Valid: true}},
				Amount:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 300, Valid: true}},
				BigCategoryID:          3,
				BigCategoryName:        "",
				MediumCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 13, Valid: true}},
				MediumCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 2, Valid: true}},
				PaymentUserID:          model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
		},
	}, nil
}

func (m MockGroupShoppingListRepository) GetDailyGroupShoppingListByCategory(date time.Time, groupID int) (model.GroupShoppingList, error) {
	return model.GroupShoppingList{
		GroupShoppingList: []model.GroupShoppingItem{
			{
				ID:                    1,
				PostedDate:            time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:           time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				ExpectedPurchaseDate:  model.Date{Time: time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC)},
				CompleteFlag:          true,
				Purchase:              "米",
				Shop:                  model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Amount:                model.NullInt64{NullInt64: sql.NullInt64{Int64: 4000, Valid: true}},
				BigCategoryID:         2,
				BigCategoryName:       "",
				MediumCategoryID:      model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				MediumCategoryName:    model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:      model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				CustomCategoryName:    model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				RegularShoppingListID: model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				PaymentUserID:         model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				TransactionAutoAdd:    true,
				RelatedTransactionData: &model.GroupTransactionData{
					ID: model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				},
			},
			{
				ID:                     2,
				PostedDate:             time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:            time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				ExpectedPurchaseDate:   model.Date{Time: time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC)},
				CompleteFlag:           false,
				Purchase:               "トイレットペーパー",
				Shop:                   model.NullString{NullString: sql.NullString{String: "クリエイト", Valid: true}},
				Amount:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 300, Valid: true}},
				BigCategoryID:          3,
				BigCategoryName:        "",
				MediumCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 13, Valid: true}},
				MediumCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 2, Valid: true}},
				PaymentUserID:          model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
		},
	}, nil
}

func (m MockGroupShoppingListRepository) GetMonthlyGroupShoppingListByDay(firstDay time.Time, lastDay time.Time, groupID int) (model.GroupShoppingList, error) {
	return model.GroupShoppingList{
		GroupShoppingList: []model.GroupShoppingItem{
			{
				ID:                    1,
				PostedDate:            time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:           time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				ExpectedPurchaseDate:  model.Date{Time: time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC)},
				CompleteFlag:          true,
				Purchase:              "米",
				Shop:                  model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Amount:                model.NullInt64{NullInt64: sql.NullInt64{Int64: 4000, Valid: true}},
				BigCategoryID:         2,
				BigCategoryName:       "",
				MediumCategoryID:      model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				MediumCategoryName:    model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:      model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				CustomCategoryName:    model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				RegularShoppingListID: model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				PaymentUserID:         model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				TransactionAutoAdd:    true,
				RelatedTransactionData: &model.GroupTransactionData{
					ID: model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				},
			},
			{
				ID:                     2,
				PostedDate:             time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:            time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				ExpectedPurchaseDate:   model.Date{Time: time.Date(2020, 12, 20, 0, 0, 0, 0, time.UTC)},
				CompleteFlag:           false,
				Purchase:               "トイレットペーパー",
				Shop:                   model.NullString{NullString: sql.NullString{String: "クリエイト", Valid: true}},
				Amount:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 300, Valid: true}},
				BigCategoryID:          3,
				BigCategoryName:        "",
				MediumCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 13, Valid: true}},
				MediumCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 2, Valid: true}},
				PaymentUserID:          model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
		},
	}, nil
}

func (m MockGroupShoppingListRepository) GetMonthlyGroupShoppingListByCategory(firstDay time.Time, lastDay time.Time, groupID int) (model.GroupShoppingList, error) {
	return model.GroupShoppingList{
		GroupShoppingList: []model.GroupShoppingItem{
			{
				ID:                    1,
				PostedDate:            time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:           time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				ExpectedPurchaseDate:  model.Date{Time: time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC)},
				CompleteFlag:          true,
				Purchase:              "米",
				Shop:                  model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Amount:                model.NullInt64{NullInt64: sql.NullInt64{Int64: 4000, Valid: true}},
				BigCategoryID:         2,
				BigCategoryName:       "",
				MediumCategoryID:      model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				MediumCategoryName:    model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:      model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				CustomCategoryName:    model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				RegularShoppingListID: model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				PaymentUserID:         model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				TransactionAutoAdd:    true,
				RelatedTransactionData: &model.GroupTransactionData{
					ID: model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				},
			},
			{
				ID:                     2,
				PostedDate:             time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:            time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				ExpectedPurchaseDate:   model.Date{Time: time.Date(2020, 12, 20, 0, 0, 0, 0, time.UTC)},
				CompleteFlag:           false,
				Purchase:               "トイレットペーパー",
				Shop:                   model.NullString{NullString: sql.NullString{String: "クリエイト", Valid: true}},
				Amount:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 300, Valid: true}},
				BigCategoryID:          3,
				BigCategoryName:        "",
				MediumCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 13, Valid: true}},
				MediumCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 2, Valid: true}},
				PaymentUserID:          model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
		},
	}, nil
}

func (m MockGroupShoppingListRepository) GetExpiredGroupShoppingList(dueDate time.Time, groupID int) (model.ExpiredGroupShoppingList, error) {
	return model.ExpiredGroupShoppingList{
		ExpiredGroupShoppingList: []model.GroupShoppingItem{
			{
				ID:                     1,
				PostedDate:             time.Date(2020, 10, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:            time.Date(2020, 11, 18, 14, 0, 0, 0, time.UTC),
				ExpectedPurchaseDate:   model.Date{Time: time.Date(2020, 11, 18, 0, 0, 0, 0, time.UTC)},
				CompleteFlag:           false,
				Purchase:               "米",
				Shop:                   model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Amount:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 4000, Valid: true}},
				BigCategoryID:          2,
				BigCategoryName:        "",
				MediumCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				MediumCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				CustomCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				PaymentUserID:          model.NullString{NullString: sql.NullString{String: "userID2", Valid: true}},
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
			{
				ID:                     2,
				PostedDate:             time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:            time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				ExpectedPurchaseDate:   model.Date{Time: time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC)},
				CompleteFlag:           false,
				Purchase:               "トイレットペーパー",
				Shop:                   model.NullString{NullString: sql.NullString{String: "クリエイト", Valid: true}},
				Amount:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 300, Valid: true}},
				BigCategoryID:          3,
				BigCategoryName:        "",
				MediumCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 13, Valid: true}},
				MediumCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 2, Valid: true}},
				PaymentUserID:          model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
			{
				ID:                     3,
				PostedDate:             time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:            time.Date(2020, 12, 19, 20, 0, 0, 0, time.UTC),
				ExpectedPurchaseDate:   model.Date{Time: time.Date(2020, 12, 25, 0, 0, 0, 0, time.UTC)},
				CompleteFlag:           false,
				Purchase:               "トイレットペーパー",
				Shop:                   model.NullString{NullString: sql.NullString{String: "クリエイト", Valid: true}},
				Amount:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 300, Valid: true}},
				BigCategoryID:          3,
				BigCategoryName:        "",
				MediumCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 13, Valid: true}},
				MediumCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				CustomCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 2, Valid: true}},
				PaymentUserID:          model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
		},
	}, nil
}

func (m MockGroupShoppingListRepository) GetGroupShoppingItem(groupShoppingItemID int) (model.GroupShoppingItem, error) {
	if groupShoppingItemID == 2 {
		return model.GroupShoppingItem{
			ID:                    1,
			PostedDate:            time.Date(2020, 12, 14, 16, 0, 0, 0, time.UTC),
			UpdatedDate:           time.Date(2020, 12, 15, 16, 0, 0, 0, time.UTC),
			ExpectedPurchaseDate:  model.Date{Time: time.Date(2020, 12, 15, 0, 0, 0, 0, time.UTC)},
			CompleteFlag:          true,
			Purchase:              "鶏肉3kg",
			Shop:                  model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
			Amount:                model.NullInt64{NullInt64: sql.NullInt64{Int64: 1000, Valid: true}},
			BigCategoryID:         2,
			BigCategoryName:       "",
			MediumCategoryID:      model.NullInt64{NullInt64: sql.NullInt64{Int64: 6, Valid: true}},
			MediumCategoryName:    model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			CustomCategoryID:      model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
			CustomCategoryName:    model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			RegularShoppingListID: model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
			PaymentUserID:         model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
			TransactionAutoAdd:    true,
			RelatedTransactionData: &model.GroupTransactionData{
				ID: model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
			},
		}, nil
	}

	return model.GroupShoppingItem{
		ID:                     1,
		PostedDate:             time.Date(2020, 12, 24, 16, 0, 0, 0, time.UTC),
		UpdatedDate:            time.Date(2020, 12, 24, 16, 0, 0, 0, time.UTC),
		ExpectedPurchaseDate:   model.Date{Time: time.Date(2020, 12, 25, 0, 0, 0, 0, time.UTC)},
		CompleteFlag:           false,
		Purchase:               "鶏肉3kg",
		Shop:                   model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
		Amount:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 1000, Valid: true}},
		BigCategoryID:          2,
		BigCategoryName:        "",
		MediumCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 6, Valid: true}},
		MediumCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		CustomCategoryID:       model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
		CustomCategoryName:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
		PaymentUserID:          model.NullString{NullString: sql.NullString{String: "userID1", Valid: true}},
		TransactionAutoAdd:     true,
		RelatedTransactionData: nil,
	}, nil
}

func (m MockGroupShoppingListRepository) PostGroupShoppingItem(groupShoppingItem *model.GroupShoppingItem, groupID int) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (m MockGroupShoppingListRepository) PutGroupShoppingItem(groupShoppingItem *model.GroupShoppingItem) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (m MockGroupShoppingListRepository) DeleteGroupShoppingItem(groupShoppingItemID int) error {
	return nil
}

func TestDBHandler_GetDailyGroupShoppingDataByDay(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetGroupShoppingItemCategoriesNameList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var mockCategoriesNameList []MockCategoriesName

		if serverCounter == 0 {
			atomic.AddInt64(&serverCounter, 1)

			mockCategoriesNameList = []MockCategoriesName{
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "食費", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
				},
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "日用品", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "消耗品", Valid: true}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				},
			}
		} else if serverCounter == 1 {
			atomic.AddInt64(&serverCounter, -1)

			mockCategoriesNameList = []MockCategoriesName{
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "食費", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
				},
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "日用品", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "消耗品", Valid: true}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				},
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&mockCategoriesNameList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	mockGetGroupShoppingItemRelatedTransactionDataList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupShoppingItemRelatedTransactionDataList := []*model.GroupTransactionData{
			{
				ID:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				TransactionDate:    "2020/12/18(金)",
				Shop:               model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Memo:               model.NullString{NullString: sql.NullString{String: "【買い物リスト】米", Valid: true}},
				Amount:             4000,
				PostedUserID:       "userID1",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID1",
				BigCategoryID:      2,
				BigCategoryName:    "食費",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
			},
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&groupShoppingItemRelatedTransactionDataList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/groups/{group_id:[0-9]+}/categories/names", mockGetGroupShoppingItemCategoriesNameList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/transactions/related-shopping-list", mockGetGroupShoppingItemRelatedTransactionDataList).Methods("GET")

	listener, err := net.Listen("tcp", accountHostURL)
	if err != nil {
		t.Fatalf("unexpected error by net.Listen() '%#v'", err)
	}

	ts := httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: router},
	}

	ts.Start()
	defer ts.Close()

	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupShoppingListRepo: MockGroupShoppingListRepository{},
		TimeManage:            MockTime{},
	}

	r := httptest.NewRequest("GET", "/groups/1/shopping-list/2020-12-18/daily", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"date":     "2020-12-18",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	dbMu.Lock()
	defer dbMu.Unlock()

	serverMu.Lock()
	defer serverMu.Unlock()

	h.GetDailyGroupShoppingDataByDay(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupShoppingDataByDay{}, &model.GroupShoppingDataByDay{})
}

func TestDBHandler_GetDailyGroupShoppingDataByCategory(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetGroupShoppingItemCategoriesNameList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var mockCategoriesNameList []MockCategoriesName

		if serverCounter == 0 {
			atomic.AddInt64(&serverCounter, 1)

			mockCategoriesNameList = []MockCategoriesName{
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "食費", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
				},
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "日用品", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "消耗品", Valid: true}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				},
			}
		} else if serverCounter == 1 {
			atomic.AddInt64(&serverCounter, -1)

			mockCategoriesNameList = []MockCategoriesName{
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "食費", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
				},
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "日用品", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "消耗品", Valid: true}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				},
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&mockCategoriesNameList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	mockGetGroupShoppingItemRelatedTransactionDataList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupShoppingItemRelatedTransactionDataList := []*model.GroupTransactionData{
			{
				ID:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				TransactionDate:    "2020/12/18(金)",
				Shop:               model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Memo:               model.NullString{NullString: sql.NullString{String: "【買い物リスト】米", Valid: true}},
				Amount:             4000,
				PostedUserID:       "userID1",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID1",
				BigCategoryID:      2,
				BigCategoryName:    "食費",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
			},
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&groupShoppingItemRelatedTransactionDataList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/groups/{group_id:[0-9]+}/categories/names", mockGetGroupShoppingItemCategoriesNameList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/transactions/related-shopping-list", mockGetGroupShoppingItemRelatedTransactionDataList).Methods("GET")

	listener, err := net.Listen("tcp", accountHostURL)
	if err != nil {
		t.Fatalf("unexpected error by net.Listen() '%#v'", err)
	}

	ts := httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: router},
	}

	ts.Start()
	defer ts.Close()

	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupShoppingListRepo: MockGroupShoppingListRepository{},
		TimeManage:            MockTime{},
	}

	r := httptest.NewRequest("GET", "/groups/1/shopping-list/2020-12-18/categories", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"date":     "2020-12-18",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	dbMu.Lock()
	defer dbMu.Unlock()

	serverMu.Lock()
	defer serverMu.Unlock()

	h.GetDailyGroupShoppingDataByCategory(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupShoppingDataByCategory{}, &model.GroupShoppingDataByCategory{})
}

func TestDBHandler_GetMonthlyGroupShoppingDataByDay(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetGroupShoppingItemCategoriesNameList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var mockCategoriesNameList []MockCategoriesName

		if serverCounter == 0 {
			atomic.AddInt64(&serverCounter, 1)

			mockCategoriesNameList = []MockCategoriesName{
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "食費", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
				},
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "日用品", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "消耗品", Valid: true}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				},
			}
		} else if serverCounter == 1 {
			atomic.AddInt64(&serverCounter, -1)

			mockCategoriesNameList = []MockCategoriesName{
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "食費", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
				},
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "日用品", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "消耗品", Valid: true}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				},
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&mockCategoriesNameList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	mockGetGroupShoppingItemRelatedTransactionDataList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupShoppingItemRelatedTransactionDataList := []*model.GroupTransactionData{
			{
				ID:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				TransactionDate:    "2020/12/18(金)",
				Shop:               model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Memo:               model.NullString{NullString: sql.NullString{String: "【買い物リスト】米", Valid: true}},
				Amount:             4000,
				PostedUserID:       "userID1",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID1",
				BigCategoryID:      2,
				BigCategoryName:    "食費",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
			},
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&groupShoppingItemRelatedTransactionDataList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/groups/{group_id:[0-9]+}/categories/names", mockGetGroupShoppingItemCategoriesNameList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/transactions/related-shopping-list", mockGetGroupShoppingItemRelatedTransactionDataList).Methods("GET")

	listener, err := net.Listen("tcp", accountHostURL)
	if err != nil {
		t.Fatalf("unexpected error by net.Listen() '%#v'", err)
	}

	ts := httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: router},
	}

	ts.Start()
	defer ts.Close()

	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupShoppingListRepo: MockGroupShoppingListRepository{},
		TimeManage:            MockTime{},
	}

	r := httptest.NewRequest("GET", "/groups/1/shopping-list/2020-12/daily", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id":   "1",
		"year_month": "2020-12",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	dbMu.Lock()
	defer dbMu.Unlock()

	serverMu.Lock()
	defer serverMu.Unlock()

	h.GetMonthlyGroupShoppingDataByDay(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupShoppingDataByDay{}, &model.GroupShoppingDataByDay{})
}

func TestDBHandler_GetMonthlyGroupShoppingDataByCategory(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetGroupShoppingItemCategoriesNameList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var mockCategoriesNameList []MockCategoriesName

		if serverCounter == 0 {
			atomic.AddInt64(&serverCounter, 1)

			mockCategoriesNameList = []MockCategoriesName{
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "食費", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
				},
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "日用品", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "消耗品", Valid: true}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				},
			}
		} else if serverCounter == 1 {
			atomic.AddInt64(&serverCounter, -1)

			mockCategoriesNameList = []MockCategoriesName{
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "食費", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
				},
				{
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "日用品", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "消耗品", Valid: true}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				},
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&mockCategoriesNameList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	mockGetGroupShoppingItemRelatedTransactionDataList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupShoppingItemRelatedTransactionDataList := []*model.GroupTransactionData{
			{
				ID:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				TransactionDate:    "2020/12/18(金)",
				Shop:               model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Memo:               model.NullString{NullString: sql.NullString{String: "【買い物リスト】米", Valid: true}},
				Amount:             4000,
				PostedUserID:       "userID1",
				UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				PaymentUserID:      "userID1",
				BigCategoryID:      2,
				BigCategoryName:    "食費",
				MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
			},
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&groupShoppingItemRelatedTransactionDataList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/groups/{group_id:[0-9]+}/categories/names", mockGetGroupShoppingItemCategoriesNameList).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/transactions/related-shopping-list", mockGetGroupShoppingItemRelatedTransactionDataList).Methods("GET")

	listener, err := net.Listen("tcp", accountHostURL)
	if err != nil {
		t.Fatalf("unexpected error by net.Listen() '%#v'", err)
	}

	ts := httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: router},
	}

	ts.Start()
	defer ts.Close()

	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupShoppingListRepo: MockGroupShoppingListRepository{},
		TimeManage:            MockTime{},
	}

	r := httptest.NewRequest("GET", "/groups/1/shopping-list/2020-12/categories", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id":   "1",
		"year_month": "2020-12",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	dbMu.Lock()
	defer dbMu.Unlock()

	serverMu.Lock()
	defer serverMu.Unlock()

	h.GetMonthlyGroupShoppingDataByCategory(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupShoppingDataByCategory{}, &model.GroupShoppingDataByCategory{})
}

func TestDBHandler_GetExpiredGroupShoppingList(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetGroupShoppingItemCategoriesNameList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockCategoriesNameList := []MockCategoriesName{
			{
				BigCategoryName:    model.NullString{NullString: sql.NullString{String: "食費", Valid: true}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "米", Valid: true}},
			},
			{
				BigCategoryName:    model.NullString{NullString: sql.NullString{String: "日用品", Valid: true}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "消耗品", Valid: true}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			},
			{
				BigCategoryName:    model.NullString{NullString: sql.NullString{String: "日用品", Valid: true}},
				MediumCategoryName: model.NullString{NullString: sql.NullString{String: "消耗品", Valid: true}},
				CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			},
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&mockCategoriesNameList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/groups/{group_id:[0-9]+}/categories/names", mockGetGroupShoppingItemCategoriesNameList).Methods("GET")

	listener, err := net.Listen("tcp", accountHostURL)
	if err != nil {
		t.Fatalf("unexpected error by net.Listen() '%#v'", err)
	}

	ts := httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: router},
	}

	ts.Start()
	defer ts.Close()

	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupShoppingListRepo: MockGroupShoppingListRepository{},
		TimeManage:            MockTime{},
	}

	r := httptest.NewRequest("GET", "/groups/1/shopping-list/expired", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	dbMu.Lock()
	defer dbMu.Unlock()

	serverMu.Lock()
	defer serverMu.Unlock()

	h.GetExpiredGroupShoppingList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.ExpiredGroupShoppingList{}, &model.ExpiredGroupShoppingList{})
}

func TestDBHandler_PostGroupRegularShoppingItem(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetGroupShoppingItemCategoriesName := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockCategoriesName := MockCategoriesName{
			BigCategoryName:    model.NullString{NullString: sql.NullString{String: "日用品", Valid: true}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "消耗品", Valid: true}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&mockCategoriesName); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/groups/{group_id:[0-9]+}/categories/name", mockGetGroupShoppingItemCategoriesName).Methods("GET")

	listener, err := net.Listen("tcp", accountHostURL)
	if err != nil {
		t.Fatalf("unexpected error by net.Listen() '%#v'", err)
	}

	ts := httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: router},
	}

	ts.Start()
	defer ts.Close()

	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupShoppingListRepo: MockGroupShoppingListRepository{},
		TimeManage:            MockTime{},
	}

	r := httptest.NewRequest("POST", "/groups/1/shopping-list/regular", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostGroupRegularShoppingItem(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res,
		&struct {
			GroupRegularShoppingItem model.GroupRegularShoppingItem `json:"regular_shopping_item"`
			model.GroupShoppingList
		}{},
		&struct {
			GroupRegularShoppingItem model.GroupRegularShoppingItem `json:"regular_shopping_item"`
			model.GroupShoppingList
		}{})
}

func TestDBHandler_PutGroupRegularShoppingItem(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetGroupShoppingItemCategoriesName := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockCategoriesName := MockCategoriesName{
			BigCategoryName:    model.NullString{NullString: sql.NullString{String: "日用品", Valid: true}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "消耗品", Valid: true}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&mockCategoriesName); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/groups/{group_id:[0-9]+}/categories/name", mockGetGroupShoppingItemCategoriesName).Methods("GET")

	listener, err := net.Listen("tcp", accountHostURL)
	if err != nil {
		t.Fatalf("unexpected error by net.Listen() '%#v'", err)
	}

	ts := httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: router},
	}

	ts.Start()
	defer ts.Close()

	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupShoppingListRepo: MockGroupShoppingListRepository{},
		TimeManage:            MockTime{},
	}

	r := httptest.NewRequest("PUT", "/groups/1/shopping-list/regular/1", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"id":       "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutGroupRegularShoppingItem(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res,
		&struct {
			GroupRegularShoppingItem model.GroupRegularShoppingItem `json:"regular_shopping_item"`
			model.GroupShoppingList
		}{},
		&struct {
			GroupRegularShoppingItem model.GroupRegularShoppingItem `json:"regular_shopping_item"`
			model.GroupShoppingList
		}{})
}

func TestDBHandler_DeleteGroupRegularShoppingItem(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupShoppingListRepo: MockGroupShoppingListRepository{},
	}

	r := httptest.NewRequest("DELETE", "/groups/1/shopping-list/regular/1", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"id":       "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.DeleteGroupRegularShoppingItem(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteContentMsg{}, &DeleteContentMsg{})
}

func TestDBHandler_PostGroupShoppingItem(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetGroupShoppingItemCategoriesName := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockCategoriesName := MockCategoriesName{
			BigCategoryName:    model.NullString{NullString: sql.NullString{String: "食費", Valid: true}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "食料品", Valid: true}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&mockCategoriesName); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/groups/{group_id:[0-9]+}/categories/name", mockGetGroupShoppingItemCategoriesName).Methods("GET")

	listener, err := net.Listen("tcp", accountHostURL)
	if err != nil {
		t.Fatalf("unexpected error by net.Listen() '%#v'", err)
	}

	ts := httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: router},
	}

	ts.Start()
	defer ts.Close()

	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupShoppingListRepo: MockGroupShoppingListRepository{},
	}

	r := httptest.NewRequest("POST", "/groups/1/shopping-list", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostGroupShoppingItem(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.GroupShoppingItem{}, &model.GroupShoppingItem{})
}

func TestDBHandler_PutGroupShoppingItem(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockPostGroupTransaction := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockGroupTransaction := model.GroupTransactionData{
			ID:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
			TransactionType:    "expense",
			PostedDate:         time.Date(2020, 12, 15, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 12, 15, 16, 0, 0, 0, time.UTC),
			TransactionDate:    "2020/12/15(火)",
			Shop:               model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
			Memo:               model.NullString{NullString: sql.NullString{String: "【買い物リスト】鶏肉3kg", Valid: true}},
			Amount:             1000,
			PostedUserID:       "userID1",
			UpdatedUserID:      model.NullString{NullString: sql.NullString{String: "", Valid: false}},
			PaymentUserID:      "userID1",
			BigCategoryID:      2,
			BigCategoryName:    "食費",
			MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 6, Valid: true}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "食料品", Valid: true}},
			CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(&mockGroupTransaction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	mockDeleteGroupTransaction := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mockGetGroupShoppingItemCategoriesName := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockCategoriesName := MockCategoriesName{
			BigCategoryName:    model.NullString{NullString: sql.NullString{String: "食費", Valid: true}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "食料品", Valid: true}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&mockCategoriesName); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/groups/{group_id:[0-9]+}/transactions", mockPostGroupTransaction).Methods("POST")
	router.HandleFunc("/groups/{group_id:[0-9]+}/transactions/{id:[0-9]+}", mockDeleteGroupTransaction).Methods("DELETE")
	router.HandleFunc("/groups/{group_id:[0-9]+}/categories/name", mockGetGroupShoppingItemCategoriesName).Methods("GET")

	listener, err := net.Listen("tcp", accountHostURL)
	if err != nil {
		t.Fatalf("unexpected error by net.Listen() '%#v'", err)
	}

	ts := httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: router},
	}

	ts.Start()
	defer ts.Close()

	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupShoppingListRepo: MockGroupShoppingListRepository{},
	}

	r := httptest.NewRequest("PUT", "/groups/1/shopping-list/2", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"id":       "2",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutGroupShoppingItem(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupShoppingItem{}, &model.GroupShoppingItem{})
}

func TestDBHandler_DeleteGroupShoppingItem(t *testing.T) {
	h := DBHandler{
		AuthRepo:              MockAuthRepository{},
		GroupShoppingListRepo: MockGroupShoppingListRepository{},
	}

	r := httptest.NewRequest("DELETE", "/groups/1/shopping-list/2", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
		"id":       "2",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.DeleteGroupShoppingItem(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteContentMsg{}, &DeleteContentMsg{})
}

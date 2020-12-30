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

	"github.com/google/go-cmp/cmp"

	"github.com/gorilla/mux"

	"github.com/google/uuid"
	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/todo-rest-service/testutil"
)

type MockCategoriesName struct {
	BigCategoryName    model.NullString `json:"big_category_name"`
	MediumCategoryName model.NullString `json:"medium_category_name"`
	CustomCategoryName model.NullString `json:"custom_category_name"`
}

type MockShoppingListRepository struct{}

func (m MockShoppingListRepository) GetRegularShoppingList(userID string) (model.RegularShoppingList, error) {
	if dbCounter == 1 {
		atomic.AddInt64(&dbCounter, -1)

		return model.RegularShoppingList{
			RegularShoppingList: []model.RegularShoppingItem{
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
					TransactionAutoAdd:   true,
				},
			},
		}, nil
	}

	atomic.AddInt64(&dbCounter, 1)

	return model.RegularShoppingList{
		RegularShoppingList: []model.RegularShoppingItem{
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
				TransactionAutoAdd:   true,
			},
		},
	}, nil
}

func (m MockShoppingListRepository) GetRegularShoppingItem(regularShoppingItemID int) (model.RegularShoppingItem, error) {
	return model.RegularShoppingItem{
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
		TransactionAutoAdd:   true,
	}, nil
}

func (m MockShoppingListRepository) GetShoppingListRelatedToPostedRegularShoppingItem(todayShoppingItemID int, laterThanTodayShoppingItemID int) (model.ShoppingList, error) {
	return model.ShoppingList{
		ShoppingList: []model.ShoppingItem{
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
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
		},
	}, nil
}

func (m MockShoppingListRepository) PostRegularShoppingItem(regularShoppingItem *model.RegularShoppingItem, userID string, today time.Time) (sql.Result, sql.Result, sql.Result, error) {
	return MockSqlResult{}, MockSqlResult{}, MockSqlResult{}, nil
}

func (m MockShoppingListRepository) GetShoppingListRelatedToUpdatedRegularShoppingItem(regularShoppingItemID int) (model.ShoppingList, error) {
	return model.ShoppingList{
		ShoppingList: []model.ShoppingItem{
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
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
		},
	}, nil
}

func (m MockShoppingListRepository) PutRegularShoppingItem(regularShoppingItem *model.RegularShoppingItem, regularShoppingItemID int, userID string, today time.Time) error {
	return nil
}

func (m MockShoppingListRepository) PutRegularShoppingList(regularShoppingList model.RegularShoppingList, userID string, today time.Time) error {
	return nil
}

func (m MockShoppingListRepository) DeleteRegularShoppingItem(regularShoppingItemID int) error {
	return nil
}

func (m MockShoppingListRepository) GetDailyShoppingListByDay(date time.Time, userID string) (model.ShoppingList, error) {
	return model.ShoppingList{
		ShoppingList: []model.ShoppingItem{
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
				TransactionAutoAdd:    true,
				RelatedTransactionData: &model.TransactionData{
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
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
		},
	}, nil
}

func (m MockShoppingListRepository) GetDailyShoppingListByCategory(date time.Time, userID string) (model.ShoppingList, error) {
	return model.ShoppingList{
		ShoppingList: []model.ShoppingItem{
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
				TransactionAutoAdd:    true,
				RelatedTransactionData: &model.TransactionData{
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
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
		},
	}, nil
}

func (m MockShoppingListRepository) GetMonthlyShoppingListByDay(firstDay time.Time, lastDay time.Time, userID string) (model.ShoppingList, error) {
	return model.ShoppingList{
		ShoppingList: []model.ShoppingItem{
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
				RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
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
				TransactionAutoAdd:    true,
				RelatedTransactionData: &model.TransactionData{
					ID: model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				},
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
				RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
		},
	}, nil
}

func (m MockShoppingListRepository) GetMonthlyShoppingListByCategory(firstDay time.Time, lastDay time.Time, userID string) (model.ShoppingList, error) {
	return model.ShoppingList{
		ShoppingList: []model.ShoppingItem{
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
				TransactionAutoAdd:    true,
				RelatedTransactionData: &model.TransactionData{
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
				RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
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
				RegularShoppingListID:  model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
		},
	}, nil
}

func (m MockShoppingListRepository) GetExpiredShoppingList(dueDate time.Time, userID string) (model.ExpiredShoppingList, error) {
	return model.ExpiredShoppingList{
		ExpiredShoppingList: []model.ShoppingItem{
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
				TransactionAutoAdd:     true,
				RelatedTransactionData: nil,
			},
		},
	}, nil
}

func (m MockShoppingListRepository) GetShoppingItem(shoppingItemID int) (model.ShoppingItem, error) {
	if shoppingItemID == 2 {
		return model.ShoppingItem{
			ID:                    2,
			PostedDate:            time.Date(2020, 12, 13, 16, 0, 0, 0, time.UTC),
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
			TransactionAutoAdd:    true,
			RelatedTransactionData: &model.TransactionData{
				ID: model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
			},
		}, nil
	}

	return model.ShoppingItem{
		ID:                     1,
		PostedDate:             time.Date(2020, 12, 13, 16, 0, 0, 0, time.UTC),
		UpdatedDate:            time.Date(2020, 12, 13, 16, 0, 0, 0, time.UTC),
		ExpectedPurchaseDate:   model.Date{Time: time.Date(2020, 12, 15, 0, 0, 0, 0, time.UTC)},
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
		TransactionAutoAdd:     true,
		RelatedTransactionData: nil,
	}, nil
}

func (m MockShoppingListRepository) PostShoppingItem(shoppingItem *model.ShoppingItem, userID string) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (m MockShoppingListRepository) PutShoppingItem(shoppingItem *model.ShoppingItem) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (m MockShoppingListRepository) DeleteShoppingItem(shoppingItemID int) error {
	return nil
}

func (m MockShoppingListRepository) PutShoppingListCustomCategoryIdToMediumCategoryId(mediumCategoryID int, customCategoryID int) error {
	return nil
}

func TestDBHandler_GetDailyShoppingDataByDay(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetCategoriesNameList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	mockGetShoppingItemRelatedTransactionDataList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shoppingItemRelatedTransactionDataList := []*model.TransactionData{
			{
				ID:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				TransactionDate:    "2020/12/18(金)",
				Shop:               model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Memo:               model.NullString{NullString: sql.NullString{String: "【買い物リスト】米", Valid: true}},
				Amount:             4000,
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
		if err := json.NewEncoder(w).Encode(&shoppingItemRelatedTransactionDataList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/categories/names", mockGetCategoriesNameList).Methods("GET")
	router.HandleFunc("/transactions/related-shopping-list", mockGetShoppingItemRelatedTransactionDataList).Methods("GET")

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
		AuthRepo:         MockAuthRepository{},
		ShoppingListRepo: MockShoppingListRepository{},
		TimeManage:       MockTime{},
	}

	r := httptest.NewRequest("GET", "/shopping-list/2020-12-18/daily", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"date": "2020-12-18",
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

	h.GetDailyShoppingDataByDay(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.ShoppingDataByDay{}, &model.ShoppingDataByDay{})
}

func TestDBHandler_GetDailyShoppingDataByCategory(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetCategoriesNameList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	mockGetShoppingItemRelatedTransactionDataList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shoppingItemRelatedTransactionDataList := []*model.TransactionData{
			{
				ID:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				TransactionDate:    "2020/12/18(金)",
				Shop:               model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Memo:               model.NullString{NullString: sql.NullString{String: "【買い物リスト】米", Valid: true}},
				Amount:             4000,
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
		if err := json.NewEncoder(w).Encode(&shoppingItemRelatedTransactionDataList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/categories/names", mockGetCategoriesNameList).Methods("GET")
	router.HandleFunc("/transactions/related-shopping-list", mockGetShoppingItemRelatedTransactionDataList).Methods("GET")

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
		AuthRepo:         MockAuthRepository{},
		ShoppingListRepo: MockShoppingListRepository{},
		TimeManage:       MockTime{},
	}

	r := httptest.NewRequest("GET", "/shopping-list/2020-12-18/categories", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"date": "2020-12-18",
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

	h.GetDailyShoppingDataByCategory(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.ShoppingDataByDay{}, &model.ShoppingDataByDay{})
}

func TestDBHandler_GetMonthlyShoppingDataByDay(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetCategoriesNameList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
					BigCategoryName:    model.NullString{NullString: sql.NullString{String: "日用品", Valid: true}},
					MediumCategoryName: model.NullString{NullString: sql.NullString{String: "消耗品", Valid: true}},
					CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
				},
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

	mockGetShoppingItemRelatedTransactionDataList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shoppingItemRelatedTransactionDataList := []*model.TransactionData{
			{
				ID:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				TransactionDate:    "2020/12/18(金)",
				Shop:               model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Memo:               model.NullString{NullString: sql.NullString{String: "【買い物リスト】米", Valid: true}},
				Amount:             4000,
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
		if err := json.NewEncoder(w).Encode(&shoppingItemRelatedTransactionDataList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/categories/names", mockGetCategoriesNameList).Methods("GET")
	router.HandleFunc("/transactions/related-shopping-list", mockGetShoppingItemRelatedTransactionDataList).Methods("GET")

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
		AuthRepo:         MockAuthRepository{},
		ShoppingListRepo: MockShoppingListRepository{},
		TimeManage:       MockTime{},
	}

	r := httptest.NewRequest("GET", "/shopping-list/2020-12/daily", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
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

	h.GetMonthlyShoppingDataByDay(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.ShoppingDataByDay{}, &model.ShoppingDataByDay{})
}

func TestDBHandler_GetMonthlyShoppingDataByCategory(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetCategoriesNameList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	mockGetShoppingItemRelatedTransactionDataList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shoppingItemRelatedTransactionDataList := []*model.TransactionData{
			{
				ID:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
				TransactionType:    "expense",
				PostedDate:         time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				UpdatedDate:        time.Date(2020, 12, 18, 14, 0, 0, 0, time.UTC),
				TransactionDate:    "2020/12/18(金)",
				Shop:               model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
				Memo:               model.NullString{NullString: sql.NullString{String: "【買い物リスト】米", Valid: true}},
				Amount:             4000,
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
		if err := json.NewEncoder(w).Encode(&shoppingItemRelatedTransactionDataList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/categories/names", mockGetCategoriesNameList).Methods("GET")
	router.HandleFunc("/transactions/related-shopping-list", mockGetShoppingItemRelatedTransactionDataList).Methods("GET")

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
		AuthRepo:         MockAuthRepository{},
		ShoppingListRepo: MockShoppingListRepository{},
		TimeManage:       MockTime{},
	}

	r := httptest.NewRequest("GET", "/shopping-list/2020-12/categories", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
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

	h.GetMonthlyShoppingDataByCategory(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.ShoppingDataByCategory{}, &model.ShoppingDataByCategory{})
}

func TestDBHandler_GetExpiredShoppingList(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetCategoriesNameList := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("/categories/names", mockGetCategoriesNameList).Methods("GET")

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
		AuthRepo:         MockAuthRepository{},
		ShoppingListRepo: MockShoppingListRepository{},
		TimeManage:       MockTime{},
	}

	r := httptest.NewRequest("GET", "/shopping-list/expired", nil)
	w := httptest.NewRecorder()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetExpiredShoppingList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.ExpiredShoppingList{}, &model.ExpiredShoppingList{})
}

func TestDBHandler_PostRegularShoppingItem(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetCategoriesName := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("/categories/name", mockGetCategoriesName).Methods("GET")

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
		AuthRepo:         MockAuthRepository{},
		ShoppingListRepo: MockShoppingListRepository{},
		TimeManage:       MockTime{},
	}

	r := httptest.NewRequest("POST", "/shopping-list/regular", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostRegularShoppingItem(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res,
		&struct {
			RegularShoppingItem model.RegularShoppingItem `json:"regular_shopping_item"`
			model.ShoppingList
		}{},
		&struct {
			RegularShoppingItem model.RegularShoppingItem `json:"regular_shopping_item"`
			model.ShoppingList
		}{})
}

func TestDBHandler_PutRegularShoppingItem(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetCategoriesName := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("/categories/name", mockGetCategoriesName).Methods("GET")

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
		AuthRepo:         MockAuthRepository{},
		ShoppingListRepo: MockShoppingListRepository{},
		TimeManage:       MockTime{},
	}

	r := httptest.NewRequest("PUT", "/shopping-list/regular/1", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutRegularShoppingItem(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res,
		&struct {
			RegularShoppingItem model.RegularShoppingItem `json:"regular_shopping_item"`
			model.ShoppingList
		}{},
		&struct {
			RegularShoppingItem model.RegularShoppingItem `json:"regular_shopping_item"`
			model.ShoppingList
		}{})
}

func TestDBHandler_DeleteRegularShoppingItem(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		ShoppingListRepo: MockShoppingListRepository{},
	}

	r := httptest.NewRequest("DELETE", "/shopping-list/regular/1", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.DeleteRegularShoppingItem(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteContentMsg{}, &DeleteContentMsg{})
}

func TestDBHandler_PostShoppingItem(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockGetCategoriesName := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("/categories/name", mockGetCategoriesName).Methods("GET")

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
		AuthRepo:         MockAuthRepository{},
		ShoppingListRepo: MockShoppingListRepository{},
		TimeManage:       MockTime{},
	}

	r := httptest.NewRequest("POST", "/shopping-list", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostShoppingItem(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.ShoppingItem{}, &model.ShoppingItem{})
}

func TestDBHandler_PutShoppingItem(t *testing.T) {
	if err := os.Setenv("ACCOUNT_HOST", "localhost"); err != nil {
		t.Fatalf("unexpected error by os.Setenv() '%#v'", err)
	}

	accountHost := os.Getenv("ACCOUNT_HOST")
	accountHostURL := fmt.Sprintf("%s:8081", accountHost)

	mockPostTransaction := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockTransaction := model.TransactionData{
			ID:                 model.NullInt64{NullInt64: sql.NullInt64{Int64: 1, Valid: true}},
			TransactionType:    "expense",
			PostedDate:         time.Date(2020, 12, 15, 16, 0, 0, 0, time.UTC),
			UpdatedDate:        time.Date(2020, 12, 15, 16, 0, 0, 0, time.UTC),
			TransactionDate:    "2020/12/15(火)",
			Shop:               model.NullString{NullString: sql.NullString{String: "コストコ", Valid: true}},
			Memo:               model.NullString{NullString: sql.NullString{String: "【買い物リスト】鶏肉3kg", Valid: true}},
			Amount:             1000,
			BigCategoryID:      2,
			BigCategoryName:    "食費",
			MediumCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 6, Valid: true}},
			MediumCategoryName: model.NullString{NullString: sql.NullString{String: "食料品", Valid: true}},
			CustomCategoryID:   model.NullInt64{NullInt64: sql.NullInt64{Int64: 0, Valid: false}},
			CustomCategoryName: model.NullString{NullString: sql.NullString{String: "", Valid: false}},
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(&mockTransaction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	mockDeleteTransaction := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mockGetCategoriesName := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("/transactions", mockPostTransaction).Methods("POST")
	router.HandleFunc("/transactions/{id:[0-9]+}", mockDeleteTransaction).Methods("DELETE")
	router.HandleFunc("/categories/name", mockGetCategoriesName).Methods("GET")

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
		AuthRepo:         MockAuthRepository{},
		ShoppingListRepo: MockShoppingListRepository{},
	}

	r := httptest.NewRequest("PUT", "/shopping-list/2", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"id": "2",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PutShoppingItem(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.ShoppingItem{}, &model.ShoppingItem{})
}

func TestDBHandler_DeleteShoppingItem(t *testing.T) {
	h := DBHandler{
		AuthRepo:         MockAuthRepository{},
		ShoppingListRepo: MockShoppingListRepository{},
	}

	r := httptest.NewRequest("DELETE", "/shopping-list/1", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.DeleteShoppingItem(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteContentMsg{}, &DeleteContentMsg{})
}

func TestDBHandler_PutShoppingListCustomCategoryIdToMediumCategoryId(t *testing.T) {
	h := DBHandler{
		ShoppingListRepo: MockShoppingListRepository{},
	}

	r := httptest.NewRequest("PUT", "/shopping-list/categories", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	h.PutShoppingListCustomCategoryIdToMediumCategoryId(w, r)

	res := w.Result()
	defer res.Body.Close()

	if diff := cmp.Diff(http.StatusOK, res.StatusCode); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

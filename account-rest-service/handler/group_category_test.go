package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
	"github.com/paypay3/kakeibo-app-api/account-rest-service/testutil"
)

type MockGroupCategoriesRepository struct{}

func (m MockGroupCategoriesRepository) GetGroupBigCategoriesList() ([]model.GroupBigCategory, error) {
	return []model.GroupBigCategory{
		{ID: 1, Name: "収入", TransactionType: "income"},
		{ID: 2, Name: "食費", TransactionType: "expense"},
		{ID: 3, Name: "日用品", TransactionType: "expense"},
		{ID: 4, Name: "趣味・娯楽", TransactionType: "expense"},
		{ID: 5, Name: "交際費", TransactionType: "expense"},
		{ID: 6, Name: "交通費", TransactionType: "expense"},
		{ID: 7, Name: "衣服・美容", TransactionType: "expense"},
		{ID: 8, Name: "健康・医療", TransactionType: "expense"},
		{ID: 9, Name: "通信費", TransactionType: "expense"},
		{ID: 10, Name: "教養・教育", TransactionType: "expense"},
		{ID: 11, Name: "住宅", TransactionType: "expense"},
		{ID: 12, Name: "水道・光熱費", TransactionType: "expense"},
		{ID: 13, Name: "自動車", TransactionType: "expense"},
		{ID: 14, Name: "保険", TransactionType: "expense"},
		{ID: 15, Name: "税金・社会保険", TransactionType: "expense"},
		{ID: 16, Name: "現金・カード", TransactionType: "expense"},
		{ID: 17, Name: "その他", TransactionType: "expense"},
	}, nil
}

func (m MockGroupCategoriesRepository) GetGroupMediumCategoriesList() ([]model.GroupAssociatedCategory, error) {
	return []model.GroupAssociatedCategory{
		{CategoryType: "MediumCategory", ID: 1, Name: "給与", BigCategoryID: 1},
		{CategoryType: "MediumCategory", ID: 2, Name: "賞与", BigCategoryID: 1},
		{CategoryType: "MediumCategory", ID: 3, Name: "一時所得", BigCategoryID: 1},
		{CategoryType: "MediumCategory", ID: 4, Name: "事業所得", BigCategoryID: 1},
		{CategoryType: "MediumCategory", ID: 5, Name: "その他収入", BigCategoryID: 1},
		{CategoryType: "MediumCategory", ID: 6, Name: "食料品", BigCategoryID: 2},
		{CategoryType: "MediumCategory", ID: 7, Name: "朝食", BigCategoryID: 2},
		{CategoryType: "MediumCategory", ID: 8, Name: "昼食", BigCategoryID: 2},
		{CategoryType: "MediumCategory", ID: 9, Name: "夕食", BigCategoryID: 2},
		{CategoryType: "MediumCategory", ID: 10, Name: "外食", BigCategoryID: 2},
		{CategoryType: "MediumCategory", ID: 11, Name: "カフェ", BigCategoryID: 2},
		{CategoryType: "MediumCategory", ID: 12, Name: "その他食費", BigCategoryID: 2},
		{CategoryType: "MediumCategory", ID: 13, Name: "消耗品", BigCategoryID: 3},
		{CategoryType: "MediumCategory", ID: 14, Name: "子育て用品", BigCategoryID: 3},
		{CategoryType: "MediumCategory", ID: 15, Name: "ペット用品", BigCategoryID: 3},
		{CategoryType: "MediumCategory", ID: 16, Name: "家具", BigCategoryID: 3},
		{CategoryType: "MediumCategory", ID: 17, Name: "家電", BigCategoryID: 3},
		{CategoryType: "MediumCategory", ID: 18, Name: "その他日用品", BigCategoryID: 3},
		{CategoryType: "MediumCategory", ID: 19, Name: "アウトドア", BigCategoryID: 4},
		{CategoryType: "MediumCategory", ID: 20, Name: "旅行", BigCategoryID: 4},
		{CategoryType: "MediumCategory", ID: 21, Name: "イベント", BigCategoryID: 4},
		{CategoryType: "MediumCategory", ID: 22, Name: "スポーツ", BigCategoryID: 4},
		{CategoryType: "MediumCategory", ID: 23, Name: "映画・動画", BigCategoryID: 4},
		{CategoryType: "MediumCategory", ID: 24, Name: "音楽", BigCategoryID: 4},
		{CategoryType: "MediumCategory", ID: 25, Name: "漫画", BigCategoryID: 4},
		{CategoryType: "MediumCategory", ID: 26, Name: "書籍", BigCategoryID: 4},
		{CategoryType: "MediumCategory", ID: 27, Name: "ゲーム", BigCategoryID: 4},
		{CategoryType: "MediumCategory", ID: 28, Name: "その他趣味・娯楽", BigCategoryID: 4},
		{CategoryType: "MediumCategory", ID: 29, Name: "飲み会", BigCategoryID: 5},
		{CategoryType: "MediumCategory", ID: 30, Name: "プレゼント", BigCategoryID: 5},
		{CategoryType: "MediumCategory", ID: 31, Name: "冠婚葬祭", BigCategoryID: 5},
		{CategoryType: "MediumCategory", ID: 32, Name: "その他交際費", BigCategoryID: 5},
		{CategoryType: "MediumCategory", ID: 33, Name: "電車", BigCategoryID: 6},
		{CategoryType: "MediumCategory", ID: 34, Name: "バス", BigCategoryID: 6},
		{CategoryType: "MediumCategory", ID: 35, Name: "タクシー", BigCategoryID: 6},
		{CategoryType: "MediumCategory", ID: 36, Name: "新幹線", BigCategoryID: 6},
		{CategoryType: "MediumCategory", ID: 37, Name: "飛行機", BigCategoryID: 6},
		{CategoryType: "MediumCategory", ID: 38, Name: "その他交通費", BigCategoryID: 6},
		{CategoryType: "MediumCategory", ID: 39, Name: "衣服", BigCategoryID: 7},
		{CategoryType: "MediumCategory", ID: 40, Name: "アクセサリー", BigCategoryID: 7},
		{CategoryType: "MediumCategory", ID: 41, Name: "クリーニング", BigCategoryID: 7},
		{CategoryType: "MediumCategory", ID: 42, Name: "美容院・理髪", BigCategoryID: 7},
		{CategoryType: "MediumCategory", ID: 43, Name: "化粧品", BigCategoryID: 7},
		{CategoryType: "MediumCategory", ID: 44, Name: "エステ・ネイル", BigCategoryID: 7},
		{CategoryType: "MediumCategory", ID: 45, Name: "その他衣服・美容", BigCategoryID: 7},
		{CategoryType: "MediumCategory", ID: 46, Name: "病院", BigCategoryID: 8},
		{CategoryType: "MediumCategory", ID: 47, Name: "薬", BigCategoryID: 8},
		{CategoryType: "MediumCategory", ID: 48, Name: "ボディケア", BigCategoryID: 8},
		{CategoryType: "MediumCategory", ID: 49, Name: "フィットネス", BigCategoryID: 8},
		{CategoryType: "MediumCategory", ID: 50, Name: "その他健康・医療", BigCategoryID: 8},
		{CategoryType: "MediumCategory", ID: 51, Name: "携帯電話", BigCategoryID: 9},
		{CategoryType: "MediumCategory", ID: 52, Name: "固定電話", BigCategoryID: 9},
		{CategoryType: "MediumCategory", ID: 53, Name: "インターネット", BigCategoryID: 9},
		{CategoryType: "MediumCategory", ID: 54, Name: "放送サービス", BigCategoryID: 9},
		{CategoryType: "MediumCategory", ID: 55, Name: "情報サービス", BigCategoryID: 9},
		{CategoryType: "MediumCategory", ID: 56, Name: "宅配・運送", BigCategoryID: 9},
		{CategoryType: "MediumCategory", ID: 57, Name: "切手・はがき", BigCategoryID: 9},
		{CategoryType: "MediumCategory", ID: 58, Name: "その他通信費", BigCategoryID: 9},
		{CategoryType: "MediumCategory", ID: 59, Name: "新聞", BigCategoryID: 10},
		{CategoryType: "MediumCategory", ID: 60, Name: "参考書", BigCategoryID: 10},
		{CategoryType: "MediumCategory", ID: 61, Name: "受験料", BigCategoryID: 10},
		{CategoryType: "MediumCategory", ID: 62, Name: "学費", BigCategoryID: 10},
		{CategoryType: "MediumCategory", ID: 63, Name: "習い事", BigCategoryID: 10},
		{CategoryType: "MediumCategory", ID: 64, Name: "塾", BigCategoryID: 10},
		{CategoryType: "MediumCategory", ID: 65, Name: "その他教養・教育", BigCategoryID: 10},
		{CategoryType: "MediumCategory", ID: 66, Name: "家賃", BigCategoryID: 11},
		{CategoryType: "MediumCategory", ID: 67, Name: "住宅ローン", BigCategoryID: 11},
		{CategoryType: "MediumCategory", ID: 68, Name: "リフォーム", BigCategoryID: 11},
		{CategoryType: "MediumCategory", ID: 69, Name: "その他住宅", BigCategoryID: 11},
		{CategoryType: "MediumCategory", ID: 70, Name: "水道", BigCategoryID: 12},
		{CategoryType: "MediumCategory", ID: 71, Name: "電気", BigCategoryID: 12},
		{CategoryType: "MediumCategory", ID: 72, Name: "ガス", BigCategoryID: 12},
		{CategoryType: "MediumCategory", ID: 73, Name: "その他水道・光熱費", BigCategoryID: 12},
		{CategoryType: "MediumCategory", ID: 74, Name: "自動車ローン", BigCategoryID: 13},
		{CategoryType: "MediumCategory", ID: 75, Name: "ガソリン", BigCategoryID: 13},
		{CategoryType: "MediumCategory", ID: 76, Name: "駐車場", BigCategoryID: 13},
		{CategoryType: "MediumCategory", ID: 77, Name: "高速料金", BigCategoryID: 13},
		{CategoryType: "MediumCategory", ID: 78, Name: "車検・整備", BigCategoryID: 13},
		{CategoryType: "MediumCategory", ID: 79, Name: "その他自動車", BigCategoryID: 13},
		{CategoryType: "MediumCategory", ID: 80, Name: "生命保険", BigCategoryID: 14},
		{CategoryType: "MediumCategory", ID: 81, Name: "医療保険", BigCategoryID: 14},
		{CategoryType: "MediumCategory", ID: 82, Name: "自動車保険", BigCategoryID: 14},
		{CategoryType: "MediumCategory", ID: 83, Name: "住宅保険", BigCategoryID: 14},
		{CategoryType: "MediumCategory", ID: 84, Name: "学資保険", BigCategoryID: 14},
		{CategoryType: "MediumCategory", ID: 85, Name: "その他保険", BigCategoryID: 14},
		{CategoryType: "MediumCategory", ID: 86, Name: "所得税", BigCategoryID: 15},
		{CategoryType: "MediumCategory", ID: 87, Name: "住民税", BigCategoryID: 15},
		{CategoryType: "MediumCategory", ID: 88, Name: "年金保険料", BigCategoryID: 15},
		{CategoryType: "MediumCategory", ID: 89, Name: "自動車税", BigCategoryID: 15},
		{CategoryType: "MediumCategory", ID: 90, Name: "その他税金・社会保険", BigCategoryID: 15},
		{CategoryType: "MediumCategory", ID: 91, Name: "現金引き出し", BigCategoryID: 16},
		{CategoryType: "MediumCategory", ID: 92, Name: "カード引き落とし", BigCategoryID: 16},
		{CategoryType: "MediumCategory", ID: 93, Name: "電子マネー", BigCategoryID: 16},
		{CategoryType: "MediumCategory", ID: 94, Name: "立替金", BigCategoryID: 16},
		{CategoryType: "MediumCategory", ID: 95, Name: "その他現金・カード", BigCategoryID: 16},
		{CategoryType: "MediumCategory", ID: 96, Name: "仕送り", BigCategoryID: 17},
		{CategoryType: "MediumCategory", ID: 97, Name: "お小遣い", BigCategoryID: 17},
		{CategoryType: "MediumCategory", ID: 98, Name: "使途不明金", BigCategoryID: 17},
		{CategoryType: "MediumCategory", ID: 99, Name: "雑費", BigCategoryID: 17},
	}, nil
}

func (m MockGroupCategoriesRepository) GetGroupCustomCategoriesList(groupID int) ([]model.GroupAssociatedCategory, error) {
	return []model.GroupAssociatedCategory{
		{CategoryType: "CustomCategory", ID: 14, Name: "株配当金", BigCategoryID: 1},
		{CategoryType: "CustomCategory", ID: 3, Name: "米", BigCategoryID: 2},
		{CategoryType: "CustomCategory", ID: 2, Name: "パン", BigCategoryID: 2},
		{CategoryType: "CustomCategory", ID: 1, Name: "調味料", BigCategoryID: 2},
		{CategoryType: "CustomCategory", ID: 6, Name: "歯磨き粉", BigCategoryID: 3},
		{CategoryType: "CustomCategory", ID: 5, Name: "トイレットペーパー", BigCategoryID: 3},
		{CategoryType: "CustomCategory", ID: 4, Name: "洗剤", BigCategoryID: 3},
	}, nil
}

func (m MockGroupCategoriesRepository) FindGroupCustomCategory(groupCustomCategory *model.GroupCustomCategory, groupID int) error {
	return sql.ErrNoRows
}

func (m MockGroupCategoriesRepository) PostGroupCustomCategory(groupCustomCategory *model.GroupCustomCategory, groupID int) (sql.Result, error) {
	return MockSqlResult{}, nil
}

func (m MockGroupCategoriesRepository) PutGroupCustomCategory(groupCustomCategory *model.GroupCustomCategory) error {
	return nil
}

func (m MockGroupCategoriesRepository) FindGroupCustomCategoryID(groupCustomCategoryID int) error {
	return nil
}

func (m MockGroupCategoriesRepository) GetBigCategoryID(groupCustomCategoryID int) (int, error) {
	return 2, nil
}

func (m MockGroupCategoriesRepository) DeleteGroupCustomCategory(previousGroupCustomCategoryID int, replaceMediumCategoryID int) error {
	return nil
}

func TestDBHandler_GetGroupCategoriesList(t *testing.T) {
	h := DBHandler{
		AuthRepo:            MockAuthRepository{},
		GroupCategoriesRepo: MockGroupCategoriesRepository{},
	}

	r := httptest.NewRequest("GET", "/groups/1/categories", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.GetGroupCategoriesList(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupCategoriesList{}, &model.GroupCategoriesList{})
}

func TestDBHandler_PostGroupCustomCategory(t *testing.T) {
	h := DBHandler{
		AuthRepo:            MockAuthRepository{},
		GroupCategoriesRepo: MockGroupCategoriesRepository{},
	}

	r := httptest.NewRequest("POST", "/groups/1/categories/custom-categories", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"group_id": "1",
	})

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: uuid.New().String(),
	}

	r.AddCookie(cookie)

	h.PostGroupCustomCategory(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusCreated)
	testutil.AssertResponseBody(t, res, &model.GroupCustomCategory{}, &model.GroupCustomCategory{})
}

func TestDBHandler_PutGroupCustomCategory(t *testing.T) {
	h := DBHandler{
		AuthRepo:            MockAuthRepository{},
		GroupCategoriesRepo: MockGroupCategoriesRepository{},
	}

	r := httptest.NewRequest("PUT", "/groups/1/categories/custom-categories/1", strings.NewReader(testutil.GetRequestJsonFromTestData(t)))
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

	h.PutGroupCustomCategory(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &model.GroupCustomCategory{}, &model.GroupCustomCategory{})
}

func TestDBHandler_DeleteGroupCustomCategory(t *testing.T) {
	h := DBHandler{
		AuthRepo:            MockAuthRepository{},
		GroupCategoriesRepo: MockGroupCategoriesRepository{},
	}

	r := httptest.NewRequest("DELETE", "/groups/1/categories/custom-categories/1", nil)
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

	h.DeleteGroupCustomCategory(w, r)

	res := w.Result()
	defer res.Body.Close()

	testutil.AssertResponseHeader(t, res, http.StatusOK)
	testutil.AssertResponseBody(t, res, &DeleteContentMsg{}, &DeleteContentMsg{})
}

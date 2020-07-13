package infrastructure

import (
	"database/sql"

	"github.com/garyburd/redigo/redis"
	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"
)

type DBRepository struct {
	*CategoriesRepository
}

type CategoriesRepository struct {
	*MySQLHandler
	*RedisHandler
}

func NewDBRepository(mysqlHandler *MySQLHandler, redisHandler *RedisHandler) *DBRepository {
	DBRepository := &DBRepository{
		&CategoriesRepository{mysqlHandler, redisHandler},
	}
	return DBRepository
}

func (r *CategoriesRepository) GetUserID(sessionID string) (string, error) {
	conn := r.RedisHandler.pool.Get()
	defer conn.Close()
	userID, err := redis.String(conn.Do("GET", sessionID))
	if err != nil {
		return userID, err
	}
	return userID, nil
}

func (r *CategoriesRepository) GetBigCategoriesList() ([]model.BigCategory, error) {
	query := "SELECT id, category_name FROM big_categories"
	rows, err := r.MySQLHandler.conn.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bigCategoriesList []model.BigCategory
	for rows.Next() {
		bigCategory := model.NewBigCategory()
		if err := rows.StructScan(&bigCategory); err != nil {
			return nil, err
		}
		bigCategoriesList = append(bigCategoriesList, bigCategory)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return bigCategoriesList, nil
}

func (r *CategoriesRepository) GetMediumCategoriesList() ([]model.MediumCategory, error) {
	query := "SELECT id, category_name, big_category_id FROM medium_categories"
	rows, err := r.MySQLHandler.conn.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mediumCategoriesList []model.MediumCategory
	for rows.Next() {
		mediumCategory := model.NewMediumCategory()
		if err := rows.StructScan(&mediumCategory); err != nil {
			return nil, err
		}
		mediumCategoriesList = append(mediumCategoriesList, mediumCategory)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return mediumCategoriesList, nil
}

func (r *CategoriesRepository) GetCustomCategoriesList(userID string) ([]model.CustomCategory, error) {
	query := "SELECT id, category_name, big_category_id FROM custom_categories WHERE user_id = ? ORDER BY id DESC"
	rows, err := r.MySQLHandler.conn.Queryx(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customCategoriesList []model.CustomCategory
	for rows.Next() {
		customCategory := model.NewCustomCategory()
		if err := rows.StructScan(&customCategory); err != nil {
			return nil, err
		}
		customCategoriesList = append(customCategoriesList, customCategory)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return customCategoriesList, nil
}

func (r *CategoriesRepository) FindCustomCategory(customCategory *model.CustomCategory, userID string) error {
	var dbCustomCategoryName string
	query := "SELECT category_name FROM custom_categories WHERE user_id = ? AND big_category_id = ? AND category_name = ?"
	if err := r.MySQLHandler.conn.QueryRowx(query, userID, customCategory.BigCategoryID, customCategory.Name).Scan(&dbCustomCategoryName); err != nil {
		return err
	}
	return nil
}

func (r *CategoriesRepository) PostCustomCategory(customCategory *model.CustomCategory, userID string) (sql.Result, error) {
	query := "INSERT INTO custom_categories(category_name, big_category_id, user_id) VALUES(?,?,?)"
	result, err := r.MySQLHandler.conn.Exec(query, customCategory.Name, customCategory.BigCategoryID, userID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *CategoriesRepository) PutCustomCategory(customCategory *model.CustomCategory, userID string) error {
	query := "UPDATE custom_categories SET category_name = ? WHERE user_id = ? AND id = ?"
	_, err := r.MySQLHandler.conn.Exec(query, customCategory.Name, userID, customCategory.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CategoriesRepository) DeleteCustomCategory(customCategory *model.CustomCategory, userID string) error {
	query := "DELETE FROM custom_categories WHERE user_id = ? AND id = ?"
	_, err := r.MySQLHandler.conn.Exec(query, userID, customCategory.ID)
	if err != nil {
		return err
	}
	return nil
}

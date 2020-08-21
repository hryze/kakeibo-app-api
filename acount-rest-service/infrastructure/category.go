package infrastructure

import (
	"database/sql"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"
)

type CategoriesRepository struct {
	*MySQLHandler
}

func NewCategoriesRepository(mysqlHandler *MySQLHandler) *CategoriesRepository {
	return &CategoriesRepository{mysqlHandler}
}

func (r *CategoriesRepository) GetBigCategoriesList() ([]model.BigCategory, error) {
	query := `
        SELECT
            id, category_name, transaction_type 
        FROM 
            big_categories
        ORDER BY
            id`

	rows, err := r.MySQLHandler.conn.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bigCategoriesList []model.BigCategory
	for rows.Next() {
		var bigCategory model.BigCategory
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
	query := `
        SELECT
            id, category_name, big_category_id 
        FROM
            medium_categories
        ORDER BY
            id`

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
	query := `
        SELECT
            id, category_name, big_category_id 
        FROM 
            custom_categories 
        WHERE 
            user_id = ? 
        ORDER BY 
            id 
        DESC`

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
	query := `
        SELECT 
            category_name 
        FROM 
            custom_categories 
        WHERE 
            user_id = ? 
        AND 
            big_category_id = ? 
        AND 
            category_name = ?`

	err := r.MySQLHandler.conn.QueryRowx(query, userID, customCategory.BigCategoryID, customCategory.Name).Scan(&dbCustomCategoryName)
	return err
}

func (r *CategoriesRepository) PostCustomCategory(customCategory *model.CustomCategory, userID string) (sql.Result, error) {
	query := `
        INSERT INTO custom_categories
            (category_name, big_category_id, user_id) 
        VALUES
            (?,?,?)`

	result, err := r.MySQLHandler.conn.Exec(query, customCategory.Name, customCategory.BigCategoryID, userID)
	return result, err
}

func (r *CategoriesRepository) PutCustomCategory(customCategory *model.CustomCategory) error {
	query := `
        UPDATE 
            custom_categories 
        SET 
            category_name = ? 
        WHERE 
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, customCategory.Name, customCategory.ID)
	return err
}

func (r *CategoriesRepository) DeleteCustomCategory(customCategoryID int) error {
	query := `
        DELETE 
        FROM 
            custom_categories 
        WHERE 
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, customCategoryID)
	return err
}

package infrastructure

import (
	"database/sql"
	"strings"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
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
            big_categories`

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

func (r *CategoriesRepository) GetMediumCategoriesList() ([]model.AssociatedCategory, error) {
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

	var mediumCategoriesList []model.AssociatedCategory
	for rows.Next() {
		mediumCategory := model.AssociatedCategory{CategoryType: "MediumCategory"}
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

func (r *CategoriesRepository) GetCustomCategoriesList(userID string) ([]model.AssociatedCategory, error) {
	query := `
        SELECT
            id, category_name, big_category_id 
        FROM 
            custom_categories 
        WHERE 
            user_id = ? 
        ORDER BY 
            id DESC`

	rows, err := r.MySQLHandler.conn.Queryx(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customCategoriesList []model.AssociatedCategory
	for rows.Next() {
		customCategory := model.AssociatedCategory{CategoryType: "CustomCategory"}
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

func (r *CategoriesRepository) GetBigCategoryID(customCategoryID int) (int, error) {
	query := `
        SELECT 
            big_category_id 
        FROM 
            custom_categories 
        WHERE 
            id = ?`

	var bigCategoryID int
	if err := r.MySQLHandler.conn.QueryRowx(query, customCategoryID).Scan(&bigCategoryID); err != nil {
		return bigCategoryID, err
	}

	return bigCategoryID, nil
}

func (r *CategoriesRepository) DeleteCustomCategory(previousCustomCategoryID int, replaceMediumCategoryID int) error {
	transactionQuery := `
        UPDATE
            transactions
        SET 
            medium_category_id = ?,
            custom_category_id = ?
        WHERE
            custom_category_id = ?`

	categoryQuery := `
        DELETE 
        FROM 
            custom_categories 
        WHERE 
            id = ?`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		if _, err := tx.Exec(transactionQuery, replaceMediumCategoryID, nil, previousCustomCategoryID); err != nil {
			return err
		}

		if _, err := tx.Exec(categoryQuery, previousCustomCategoryID); err != nil {
			return err
		}

		return nil
	}

	if err := transactions(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return nil
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CategoriesRepository) GetCategoriesName(categoriesID model.CategoriesID) (*model.CategoriesName, error) {
	var categoriesName model.CategoriesName
	var query string
	var categoryID int64

	if categoriesID.MediumCategoryID.Valid {
		query = `
            SELECT
                big_categories.category_name big_category_name,
                medium_categories.category_name medium_category_name
            FROM
                medium_categories
            LEFT JOIN
                big_categories
            ON
                medium_categories.big_category_id = big_categories.id
            WHERE
                medium_categories.id = ?`

		categoryID = categoriesID.MediumCategoryID.Int64
	} else if categoriesID.CustomCategoryID.Valid {
		query = `
            SELECT
                big_categories.category_name big_category_name,
                custom_categories.category_name custom_category_name
            FROM
                custom_categories
            LEFT JOIN
                big_categories
            ON
                custom_categories.big_category_id = big_categories.id
            WHERE
                custom_categories.id = ?`

		categoryID = categoriesID.CustomCategoryID.Int64
	}

	if err := r.MySQLHandler.conn.QueryRowx(query, categoryID).StructScan(&categoriesName); err != nil {
		return nil, err
	}

	return &categoriesName, nil
}

func (r *CategoriesRepository) GetCategoriesNameList(categoriesIDList []model.CategoriesID) ([]model.CategoriesName, error) {
	sliceQuery := make([]string, len(categoriesIDList))
	queryArgs := make([]interface{}, len(categoriesIDList))

	for i, categoriesID := range categoriesIDList {
		if categoriesID.MediumCategoryID.Valid {
			sliceQuery[i] = `
            SELECT
                big_categories.category_name big_category_name,
                medium_categories.category_name medium_category_name,
                NULL custom_category_name
            FROM
                medium_categories
            LEFT JOIN
                big_categories
            ON
                medium_categories.big_category_id = big_categories.id
            WHERE
                medium_categories.id = ?`

			queryArgs[i] = categoriesID.MediumCategoryID.Int64
		} else if categoriesID.CustomCategoryID.Valid {
			sliceQuery[i] = `
            SELECT
                big_categories.category_name big_category_name,
                NULL medium_category_name,
                custom_categories.category_name custom_category_name
            FROM
                custom_categories
            LEFT JOIN
                big_categories
            ON
                custom_categories.big_category_id = big_categories.id
            WHERE
                custom_categories.id = ?`

			queryArgs[i] = categoriesID.CustomCategoryID.Int64
		}
	}

	query := strings.Join(sliceQuery, " UNION ALL ")

	rows, err := r.MySQLHandler.conn.Queryx(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categoriesNameList := make([]model.CategoriesName, len(categoriesIDList))
	for i := 0; rows.Next(); i++ {
		var categoriesName model.CategoriesName
		if err := rows.StructScan(&categoriesName); err != nil {
			return nil, err
		}

		categoriesNameList[i] = categoriesName
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categoriesNameList, nil
}

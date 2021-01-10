package infrastructure

import (
	"database/sql"
	"strings"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
)

type GroupCategoriesRepository struct {
	*MySQLHandler
}

func NewGroupCategoriesRepository(mysqlHandler *MySQLHandler) *GroupCategoriesRepository {
	return &GroupCategoriesRepository{mysqlHandler}
}

func (r *GroupCategoriesRepository) GetGroupBigCategoriesList() ([]model.GroupBigCategory, error) {
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

	var groupBigCategoriesList []model.GroupBigCategory
	for rows.Next() {
		var groupBigCategory model.GroupBigCategory
		if err := rows.StructScan(&groupBigCategory); err != nil {
			return nil, err
		}

		groupBigCategoriesList = append(groupBigCategoriesList, groupBigCategory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupBigCategoriesList, nil
}

func (r *GroupCategoriesRepository) GetGroupMediumCategoriesList() ([]model.GroupAssociatedCategory, error) {
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

	var groupMediumCategoriesList []model.GroupAssociatedCategory
	for rows.Next() {
		groupMediumCategory := model.GroupAssociatedCategory{CategoryType: "MediumCategory"}
		if err := rows.StructScan(&groupMediumCategory); err != nil {
			return nil, err
		}

		groupMediumCategoriesList = append(groupMediumCategoriesList, groupMediumCategory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupMediumCategoriesList, nil
}

func (r *GroupCategoriesRepository) GetGroupCustomCategoriesList(groupID int) ([]model.GroupAssociatedCategory, error) {
	query := `
        SELECT
            id, category_name, big_category_id
        FROM
            group_custom_categories
        WHERE
            group_id = ?
        ORDER BY
            id DESC`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupCustomCategoriesList []model.GroupAssociatedCategory
	for rows.Next() {
		groupCustomCategory := model.GroupAssociatedCategory{CategoryType: "CustomCategory"}
		if err := rows.StructScan(&groupCustomCategory); err != nil {
			return nil, err
		}

		groupCustomCategoriesList = append(groupCustomCategoriesList, groupCustomCategory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupCustomCategoriesList, nil
}

func (r *GroupCategoriesRepository) FindGroupCustomCategory(groupCustomCategory *model.GroupCustomCategory, groupID int) error {
	query := `
        SELECT 
            category_name 
        FROM 
            group_custom_categories 
        WHERE 
            group_id = ? 
        AND 
            big_category_id = ? 
        AND 
            category_name = ?`

	var groupCustomCategoryName string
	err := r.MySQLHandler.conn.QueryRowx(query, groupID, groupCustomCategory.BigCategoryID, groupCustomCategory.Name).Scan(&groupCustomCategoryName)

	return err
}

func (r *GroupCategoriesRepository) PostGroupCustomCategory(groupCustomCategory *model.GroupCustomCategory, groupID int) (sql.Result, error) {
	query := `
        INSERT INTO group_custom_categories
            (category_name, big_category_id, group_id) 
        VALUES
            (?,?,?)`

	result, err := r.MySQLHandler.conn.Exec(query, groupCustomCategory.Name, groupCustomCategory.BigCategoryID, groupID)

	return result, err
}

func (r *GroupCategoriesRepository) PutGroupCustomCategory(groupCustomCategory *model.GroupCustomCategory) error {
	query := `
        UPDATE 
            group_custom_categories 
        SET 
            category_name = ? 
        WHERE 
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupCustomCategory.Name, groupCustomCategory.ID)

	return err
}

func (r *GroupCategoriesRepository) FindGroupCustomCategoryID(groupCustomCategoryID int) error {
	query := `
        SELECT 
            category_name 
        FROM 
            group_custom_categories 
        WHERE 
            id = ?`

	var groupCustomCategoryName string
	err := r.MySQLHandler.conn.QueryRowx(query, groupCustomCategoryID).Scan(&groupCustomCategoryName)

	return err
}

func (r *GroupCategoriesRepository) GetBigCategoryID(groupCustomCategoryID int) (int, error) {
	query := `
        SELECT 
            big_category_id 
        FROM 
            group_custom_categories 
        WHERE 
            id = ?`

	var bigCategoryID int
	if err := r.MySQLHandler.conn.QueryRowx(query, groupCustomCategoryID).Scan(&bigCategoryID); err != nil {
		return bigCategoryID, err
	}

	return bigCategoryID, nil
}

func (r *GroupCategoriesRepository) DeleteGroupCustomCategory(previousGroupCustomCategoryID int, replaceMediumCategoryID int) error {
	transactionQuery := `
        UPDATE
            group_transactions
        SET 
            medium_category_id = ?,
            custom_category_id = ?
        WHERE
            custom_category_id = ?`

	categoryQuery := `
        DELETE 
        FROM 
            group_custom_categories 
        WHERE 
            id = ?`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		if _, err := tx.Exec(transactionQuery, replaceMediumCategoryID, nil, previousGroupCustomCategoryID); err != nil {
			return err
		}

		if _, err := tx.Exec(categoryQuery, previousGroupCustomCategoryID); err != nil {
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

func (r *GroupCategoriesRepository) GetGroupCategoriesName(categoriesID model.CategoriesID) (*model.CategoriesName, error) {
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
                group_custom_categories.category_name custom_category_name
            FROM
                group_custom_categories
            LEFT JOIN
                big_categories
            ON
                group_custom_categories.big_category_id = big_categories.id
            WHERE
                group_custom_categories.id = ?`

		categoryID = categoriesID.CustomCategoryID.Int64
	}

	if err := r.MySQLHandler.conn.QueryRowx(query, categoryID).StructScan(&categoriesName); err != nil {
		return nil, err
	}

	return &categoriesName, nil
}

func (r *GroupCategoriesRepository) GetGroupCategoriesNameList(categoriesIDList []model.CategoriesID) ([]model.CategoriesName, error) {
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
                group_custom_categories.category_name custom_category_name
            FROM
                group_custom_categories
            LEFT JOIN
                big_categories
            ON
                group_custom_categories.big_category_id = big_categories.id
            WHERE
                group_custom_categories.id = ?`

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

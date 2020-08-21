package infrastructure

import (
	"database/sql"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"
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
            big_categories
        ORDER BY
            id`

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

func (r *GroupCategoriesRepository) GetGroupMediumCategoriesList() ([]model.GroupMediumCategory, error) {
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

	var groupMediumCategoriesList []model.GroupMediumCategory
	for rows.Next() {
		groupMediumCategory := model.NewGroupMediumCategory()
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

func (r *GroupCategoriesRepository) GetGroupCustomCategoriesList(groupID int) ([]model.GroupCustomCategory, error) {
	query := `
        SELECT
            id, category_name, big_category_id
        FROM
            group_custom_categories
        WHERE
            group_id = ?
        ORDER BY
            id
        DESC`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupCustomCategoriesList []model.GroupCustomCategory
	for rows.Next() {
		groupCustomCategory := model.NewGroupCustomCategory()
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

func (r *GroupCategoriesRepository) DeleteGroupCustomCategory(groupCustomCategoryID int) error {
	query := `
        DELETE 
        FROM 
            group_custom_categories 
        WHERE 
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupCustomCategoryID)
	return err
}

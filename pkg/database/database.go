package database

import (
	"fmt"

	"gorm.io/gorm"
)

type DBStructure struct {
	db *gorm.DB
}

type TabelDetail struct {
	SchemaName string
	TableName  string
	ColumnName string
	Position   string
	DataType   string
	MaxLeangth string
	IsNullable string
}

type FK struct {
	PKCOLUMN_NAME string
	FKTABLE_NAME  string
}

type Constraint struct {
	Name           string
	ConstraintType string
}

func New(db *gorm.DB) DBStructure {
	return DBStructure{db}
}

func (d DBStructure) GetTabelName() []string {
	var result []string
	d.db.Raw(`SELECT TABLE_NAME AS TableName
			FROM INFORMATION_SCHEMA.COLUMNS
			GROUP BY TABLE_NAME
			ORDER BY TableName`).Scan(&result)
	return result
}

func (d DBStructure) GetTabelDetail(name string) []TabelDetail {
	var result []TabelDetail
	d.db.Raw(fmt.Sprintf(`SELECT TABLE_SCHEMA AS SchemaName
			,TABLE_NAME AS TableName
			,COLUMN_NAME AS ColumnName
			,ORDINAL_POSITION AS Position
			,DATA_TYPE AS DataType
			,CHARACTER_MAXIMUM_LENGTH AS MaxLeangth
			,IS_NULLABLE AS IsNullable
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_NAME = '%s'
		ORDER BY
			SchemaName
		,TableName
		,ColumnName`, name)).Scan(&result)

	return result
}

func (d DBStructure) GetFK(name string) []FK {
	var result []FK
	d.db.Raw(fmt.Sprintf(`EXEC sp_fkeys @pktable_name = '%s', @pktable_owner = 'dbo'`, name)).Scan(&result)

	return result
}

func (d DBStructure) GetPrimaryKey(name string) []string {
	var result []string
	d.db.Raw(fmt.Sprintf(`SELECT column_name as PRIMARYKEYCOLUMN
				FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS TC 
				INNER JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS KU
				ON TC.CONSTRAINT_TYPE = 'PRIMARY KEY' 
				AND TC.CONSTRAINT_NAME = KU.CONSTRAINT_NAME 
				AND KU.table_name='%s'
				ORDER BY 
					KU.TABLE_NAME
					,KU.ORDINAL_POSITION`, name)).Scan(&result)
	return result
}

func (d DBStructure) GetConstraints(name string) []Constraint {
	var result []Constraint
	d.db.Raw(fmt.Sprintf(`SELECT CONSTRAINT_NAME AS ConstraintName
			,CONSTRAINT_TYPE AS Constrainttype 
			FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS where TABLE_NAME = '%s'`, name))

	return result
}

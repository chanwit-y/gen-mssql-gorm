package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/chanwit-y/gen-mssql-gorm.git/pkg/database"
	"github.com/chanwit-y/gen-mssql-gorm.git/pkg/env"
	"github.com/chanwit-y/gen-mssql-gorm.git/schemax"

	"github.com/samber/lo"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var dbStructure database.DBStructure

func init() {
	dsn := env.Env().CONNECTION_STRING
	db, _ := gorm.Open(sqlserver.Open(dsn))
	dbStructure = database.New(db)

	// var tripItems schemax.TripItems
	// db.Debug().Preload("Flight").Find(&tripItems, "TAI_ID = ?", 3874)
	// db.Debug().Preload(clause.Associations).Find(&tripItems, "TAI_ID = ?", 3876)
	var trip schemax.Trip
	db.Debug().Preload(clause.Associations).Find(&trip, "TA_ID = ?", 1950)
	fmt.Println(trip.TripItems[0].TaiId)

	taiId := trip.TripItems[0].TaiId
	var cars []schemax.Car
	db.Debug().Preload(clause.Associations).Find(&cars, "TAI_ID = ?", taiId)

	fmt.Println(cars[0].CarId)

}

func main() {
	// tabels := dbStructure.GetTabelName()
	// lo.ForEach(tabels, func(t string, i int) {
	// 	schema := genSchema(t)
	// 	createFile(fmt.Sprintf("./schemax/%s.go", strings.ToLower(t)), schema)
	// })

	// lo.ForEach([]string{"TRIP", "TRIP_ITEMS", "HOTEL", "FLIGHT", "CAR", "TRAIN"}, func(t string, i int) {
	// 	// lo.ForEach([]string{"TRIP_ITEMS", "FLIGHT"}, func(t string, i int) {
	// 	schema := genSchema(t)
	// 	createFile(fmt.Sprintf("./schemax/%s.go", strings.ToLower(t)), schema)
	// })

}

func genSchema(name string) []string {
	const packageName = "schemax"

	var schema []string

	schema = append(schema, fmt.Sprintf("package %s\n", packageName))
	schema = append(schema, fmt.Sprintf("type %s struct {\n", toCamelCase(name)))

	pks := dbStructure.GetPrimaryKey(name)

	details := dbStructure.GetTabelDetail(name)
	lo.ForEach(sortTabelDetail(details), func(t database.TabelDetail, i int) {
		pk := ternary(isPk(pks, t.ColumnName), "primaryKey", "")
		colName := toCamelCase(t.ColumnName)
		dataType := toGoType(t.DataType)
		schema = append(schema, fmt.Sprintf("	%s %s `gorm:\"column:%s;type:%s;%s\"`\n", colName, dataType, t.ColumnName, t.DataType, pk))
	})

	// contraints := dbStructure.GetConstraints(name)

	// lo.ForEach(contraints, func(t database.Constraint, i int) {

	// 	fmt.Println(t.ConstraintName)

	// 	constraintNames := dbStructure.GetUniqueConstraintName(t.ConstraintName)

	// 	fmt.Println(constraintNames)

	// 	lo.ForEach(constraintNames, func(t string, i int) {
	// 		colName := strings.ReplaceAll(t, "_PK", "")
	// 		schema = append(schema, fmt.Sprintf("	%s %s `gorm:\"foreignKey:TaId\"` \n", toCamelCase(colName), toCamelCase(colName)))
	// 	})
	// })

	fk := dbStructure.GetFK(name)
	lo.ForEach(fk, func(t database.FK, i int) {
		schema = append(schema, fmt.Sprintf("	%s []%s `gorm:\"foreignKey:%s;references:%s\"`\n",
			toCamelCase(t.FKTABLE_NAME),
			toCamelCase(t.FKTABLE_NAME),
			toCamelCase(t.PKCOLUMN_NAME),
			toCamelCase(t.PKCOLUMN_NAME)))
	})

	// for _, s := range schema {
	// 	fmt.Println(s)
	// }

	schema = append(schema, "}\n")

	schema = append(schema, fmt.Sprintf("func (%s) TableName() string {\n", toCamelCase(name)))
	schema = append(schema, fmt.Sprintf("	return \"%s\"\n", name))
	schema = append(schema, "}\n")

	return schema
}

func toCamelCase(text string) string {
	spName := strings.Split(text, "_")
	tabelName := lo.Reduce(spName, func(r string, t string, i int) string {
		return r + t[0:1] + strings.ToLower(t[1:])
	}, "")

	return tabelName
}

func isPk(pks []string, colName string) bool {
	_, f := lo.Find(pks, func(t string) bool {
		return t == colName
	})

	return f
}

func toGoType(dataType string) string {

	switch strings.ToLower(dataType) {
	case "nvarchar":
		return "string"
	case "bigint":
		return "int64"
	case "int":
		return "int64"
	case "bit":
		return "bool"
	case "datetime":
		return "time.Time"
	case "time":
		return "time.Time"
	case "decimal":
		return "float64"
	case "varchar":
		return "string"
	default:
		return ""
	}
}

func sortTabelDetail(tabelDetails []database.TabelDetail) []database.TabelDetail {
	sort.SliceStable(tabelDetails, func(i, j int) bool {
		return tabelDetails[i].Position < tabelDetails[j].Position
	})

	return tabelDetails
}

func createFile(name string, schema []string) {
	file, _ := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	datawrite := bufio.NewWriter(file)

	for _, data := range schema {
		_, _ = datawrite.WriteString(data)
	}

	datawrite.Flush()
	file.Close()
}

func ternary[T any](condition bool, v1 T, v2 T) T {
	if condition {
		return v1
	} else {
		return v2
	}

}

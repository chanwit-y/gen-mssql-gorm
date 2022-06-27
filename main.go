package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/chanwit-y/gen-mssql-gorm.git/pkg/database"
	"github.com/chanwit-y/gen-mssql-gorm.git/pkg/env"

	"github.com/samber/lo"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var dbStructure database.DBStructure

func init() {
	dsn := env.Env().CONNECTION_STRING
	db, _ := gorm.Open(sqlserver.Open(dsn))
	dbStructure = database.New(db)

	// var tripItems schemax.TripItems
	// db.First(&tripItems)
	// println(tripItems.CountryCode)
}

func main() {
	// tabels := dbStructure.GetTabelName()
	// lo.ForEach(tabels, func(t string, i int) {
	// 	detail := dbStructure.GetTabelDetail(t)
	// 	fmt.Println(detail)
	// })

	schema := genSchema("TRIP_ITEMS")
	createFile(fmt.Sprintf("./schemax/%s.go", strings.ToLower("TRIP_ITEMS")), schema)

}

func genSchema(name string) []string {
	var schema []string

	schema = append(schema, "package schemax\n")
	schema = append(schema, fmt.Sprintf("type %s struct {\n", toCamelCase(name)))

	pks := dbStructure.GetPrimaryKey(name)

	details := dbStructure.GetTabelDetail(name)
	lo.ForEach(sortTabelDetail(details), func(t database.TabelDetail, i int) {
		pk := ternary(isPk(pks, t.ColumnName), "primaryKey", "")
		colName := toCamelCase(t.ColumnName)
		dataType := toGoType(t.DataType)
		schema = append(schema, fmt.Sprintf("	%s %s `gorm:\"column:%s;type:%s;%s\"`\n", colName, dataType, t.ColumnName, t.DataType, pk))
	})

	for _, s := range schema {
		fmt.Println(s)
	}

	fk := dbStructure.GetFK(name)
	fmt.Println(fk)

	schema = append(schema, "}")

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

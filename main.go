package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Model represents the structure of your database table
type Model struct {
	gorm.Model
	Name    string
	Fields  []Field `gorm:"foreignKey:ModelID"`
	ModelID uint    // This field will serve as the foreign key

}

// Field represents the structure of a field in the model
type Field struct {
	Name    string
	Type    string
	ModelID uint
	gorm.Model
}

// Database connection
var db *gorm.DB

// Page represents the structure of a web page
type Page struct {
	Title string
}

func init() {
	timeZone := "Asia/Jakarta"
	var err error
	const (
		DbHost     = "localhost"
		DbPort     = 5432
		DbUsername = "postgres"
		DbPass     = "admin"
		DbName     = "crud_builder_simple"
	)

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password='%s' dbname=%s sslmode=disable TimeZone=%s",
		DbHost,
		DbPort,
		DbUsername,
		DbPass,
		DbName,
		timeZone,
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// AutoMigrate the model
	db.AutoMigrate(&Model{}, &Field{})
}

type CustomStructField struct {
	Name string
	Type reflect.Type
	Tag  reflect.StructTag
}

func camelToSnake(input string) string {
	regex := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := regex.ReplaceAllString(input, "${1}_${2}")
	return strings.ToLower(snake)
}

func createTable(model Model) {
	// Create a dynamic struct
	fields := []reflect.StructField{
		{Name: "Id", Type: reflect.TypeOf(0), Tag: reflect.StructTag(`gorm:"primaryKey"`)},
		{Name: "CreatedAt", Type: reflect.TypeOf(time.Time{}), Tag: reflect.StructTag(`gorm:"default:CURRENT_TIMESTAMP()"`)},
		{Name: "UpdatedAt", Type: reflect.TypeOf(time.Time{}), Tag: reflect.StructTag(`gorm:"default:CURRENT_TIMESTAMP()"`)},
		{Name: "DeletedAt", Type: reflect.TypeOf(gorm.DeletedAt{})},
	}
	for _, field := range model.Fields {
		var fieldType reflect.Type

		switch field.Type {
		case "number":
			fieldType = reflect.TypeOf(0)
		case "string":
			fieldType = reflect.TypeOf("")
		case "boolean":
			fieldType = reflect.TypeOf(true)
		default:
			// Handle unsupported types
			panic(fmt.Sprintf("Unsupported type for field %s", camelToSnake(field.Name)))
		}

		fields = append(fields, reflect.StructField{
			Name: field.Name,
			Type: fieldType,
			Tag:  reflect.StructTag(fmt.Sprintf(`gorm:"column:%s"`, camelToSnake(field.Name))),
		})
	}

	resultType := reflect.StructOf(fields)

	// Create an instance of the dynamic struct
	instance := reflect.New(resultType).Interface()

	tableName := model.Name
	// AutoMigrate the dynamic struct
	db.Table(tableName).AutoMigrate(instance)

}

func listRecords(c *gin.Context) {
	tableName := c.Param("table")
	var records []map[string]interface{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted_at IS NULL", tableName)

	// Executing the raw SQL query
	rows, err := db.Raw(query).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a slice to store values for each row
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// Scan each row and store the values in a map
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowData := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			rowData[col] = val
		}

		records = append(records, rowData)
	}

	c.JSON(http.StatusOK, records)
}

func createRecord(c *gin.Context) {
	tableName := c.Param("table")
	var requestData map[string]interface{}
	var result interface{}
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&result)
	if err != nil {
		panic(err)
	}

	// Assuming requestData is meant to be used instead of result
	requestData = result.(map[string]interface{})

	// Build the INSERT query dynamically
	var columns, values []string
	for key, value := range requestData {
		columns = append(columns, key)
		values = append(values, fmt.Sprintf("'%v'", value))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(values, ", "))

	// Execute the raw SQL query
	if err := db.Exec(query).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": requestData})
}

func main() {

	r := gin.Default()

	// Define API routes
	r.POST("/models", createModel)
	r.GET("/models", listModels)
	r.GET("/models/:id", getModel)

	// Define Table
	r.GET("/tables/:table", listRecords)
	r.POST("/tables/:table", createRecord)

	// Run the server
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}

// Handlers

func createModel(c *gin.Context) {
	var newModel Model
	if err := c.BindJSON(&newModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the model
	if err := db.Create(&newModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createTable(newModel)
	c.JSON(http.StatusOK, newModel)
}

func listModels(c *gin.Context) {
	var models []Model
	if err := db.Preload(clause.Associations).Find(&models).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models)
}

func getModel(c *gin.Context) {
	id := c.Param("id")
	var model Model
	if err := db.Preload(clause.Associations).First(&model, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	c.JSON(http.StatusOK, model)
}

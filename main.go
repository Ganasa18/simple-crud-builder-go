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
		{Name: "CreatedAt", Type: reflect.TypeOf(time.Time{}), Tag: reflect.StructTag(`gorm:"default:current_timestamp"`)},
		{Name: "UpdatedAt", Type: reflect.TypeOf(time.Time{}), Tag: reflect.StructTag(`gorm:"default:current_timestamp"`)},
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

// Handlers Model

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

// Handler Table

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

func getRecord(c *gin.Context) {
	tableName := c.Param("table")
	idRecord := c.Param("id")

	// Define a dynamic struct to hold the record data
	record := make(map[string]interface{})

	// Query the database to retrieve the record
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tableName)
	db.Raw(query, idRecord).Scan(&record)

	// Check if the record exists
	if len(record) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	// Return the record data
	c.JSON(http.StatusOK, gin.H{"data": record})
}

func deleteRecord(c *gin.Context) {
	tableName := c.Param("table")
	idRecord := c.Param("id")

	// Check if the record exists before attempting to delete
	var count int64
	if err := db.Table(tableName).Where("id = ?", idRecord).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	// Perform the delete operation
	if err := db.Table(tableName).Where("id = ?", idRecord).Delete(nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"data": "success delete"})
}

func softDeleteRecord(c *gin.Context) {
	tableName := c.Param("table")
	idRecord := c.Param("id")

	// Check if the record exists before attempting to soft delete
	var count int64
	if err := db.Table(tableName).Where("id = ?", idRecord).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	// Perform the soft delete operation
	if err := db.Table(tableName).Where("id = ?", idRecord).Updates(map[string]interface{}{"deleted_at": time.Now()}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"data": "success soft delete"})
}

func updateRecord(c *gin.Context) {
	tableName := c.Param("table")
	idRecord := c.Param("id")

	// Check if the record exists before attempting to update
	var count int64
	if err := db.Table(tableName).Where("id = ?", idRecord).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	// Decode the JSON request body into a map
	var requestData map[string]interface{}
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the 'updated_at' field to the current timestamp
	requestData["updated_at"] = time.Now()

	// Perform the update operation
	if err := db.Table(tableName).Where("id = ?", idRecord).Updates(requestData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch the updated record from the database as a map
	var updatedRecord map[string]interface{}
	if err := db.Table(tableName).Where("id = ?", idRecord).Scan(&updatedRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated record as a map
	c.JSON(http.StatusOK, gin.H{"data": updatedRecord})
}

func main() {

	r := gin.Default()

	// Define API routes Model
	r.POST("/models", createModel)
	r.GET("/models", listModels)
	r.GET("/models/:id", getModel)

	// Define API routes Table
	r.GET("/tables/:table", listRecords)
	r.POST("/tables/:table", createRecord)
	r.GET("/tables/:table/:id", getRecord)
	r.PATCH("/tables/:table/:id", updateRecord)
	r.DELETE("/tables/:table/:id", deleteRecord)
	r.DELETE("/tables-soft/:table/:id", softDeleteRecord)

	// Run the server
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}

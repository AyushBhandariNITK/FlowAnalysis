package main

import (
	"flowanalysis/pkg/inmemory"
	"flowanalysis/pkg/service"
	"fmt"

	"github.com/labstack/echo/v4"
)

func init() {
	go inmemory.StartMap()

	// Define the query to get the list of tables in the public schema
	//	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';`

	// Execute the query
	// tables, err := instance.Query(query)
	// fmt.Printf("Output of tables %s\n", tables)
	// if err != nil {
	// 	log.Fatalf("Failed to execute query: %v", err)
	// }

	// // Type assertion to get the result as []string
	// if tableNames, ok := tables.([]string); ok {
	// 	// Print the list of tables
	// 	fmt.Println("List of tables in the database:")
	// 	for _, tableName := range tableNames {
	// 		fmt.Println(tableName)
	// 	}
	// } else {
	// 	log.Fatal("Unexpected result type")
	// }
	fmt.Println("Database connection established successfully.")
}

func main() {

	e := echo.New()
	e.GET("/api/verve/accept", service.GetAcceptHandler)
	e.Start(":5010")
}

package main

import (
	"fmt"

	t_ "design_patterns/template"
)

func main() {
	fmt.Println("=== CSV Migration ===")
	csv := t_.NewCSVMigration("users.csv", []string{
		"ID, Name, City",
		"1, Alice, New York",
		"2, Bob, Los Angeles",
	})
	if err := t_.RunMigration(csv); err != nil {
		fmt.Println("CSV migration failed:", err)
	}

	fmt.Println()
	fmt.Println("=== JSON Migration ===")
	json := t_.NewJSONMigration("https://api.example.com/users", []t_.Record{
		{"id": "1", "name": "  Charlie  ", "city": "  Chicago  "},
		{"id": "2", "name": "  Diana  ", "city": "  Seattle  "},
	})
	if err := t_.RunMigration(json); err != nil {
		fmt.Println("JSON migration failed:", err)
	}
}

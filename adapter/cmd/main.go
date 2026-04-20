package main

import (
	"fmt"

	a "design_patterns/adapter"
)

// fileService depends only on CloudStorage — it has no idea an FTP client is underneath.
func fileService(store a.CloudStorage) {
	store.Upload("data/report.csv", "id,name\n1,alice")

	content, err := store.Download("data/report.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Downloaded:", content)

	store.Remove("data/report.csv")

	_, err = store.Download("data/report.csv")
	fmt.Println("After remove:", err)
}

func main() {
	// Wire the adapter once at the entry point.
	ftpClient := a.NewLegacyFTPClient()
	store := a.NewFTPAdapter(ftpClient)

	fileService(store)
}

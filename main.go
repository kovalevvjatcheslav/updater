package main

import (
	"fmt"
	"log"
	"os"
	_ "updater/types"
	"updater/utils"
)

func main() {
	file_path := "fixtures/Коды ответа API - Документация Подсказок - Confluence.html"
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatalf("Cannot open file %s", file_path)
	}
	defer file.Close()
	test_table := utils.ParseHtml(file)

	if len(os.Args[1:]) < 2 {
		log.Fatal("Not enough arguments passed")
	}
	path_to_secret := os.Args[1]
	doc_id := os.Args[2]
	fmt.Println(path_to_secret, doc_id)

	utils.UpdateDoc(path_to_secret, doc_id, &test_table)

}

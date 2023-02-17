package main

import (
	"log"
	"net/http"
	"os"
	"updater/utils"
)

func main() {
	if len(os.Args[1:]) < 3 {
		log.Fatal("Not enough arguments passed")
	}
	url_conf_doc := os.Args[1]
	path_to_secret := os.Args[2]
	doc_id := os.Args[3]

	res, err := http.Get(url_conf_doc)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	table := utils.ParseHtml(res.Body)

	utils.UpdateDoc(path_to_secret, doc_id, &table)
}

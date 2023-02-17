package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"updater/types"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

func getClient(config *oauth2.Config) *http.Client {
	tok_file := "token.json"
	tok, err := tokenFromFile(tok_file)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tok_file, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	auth_url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", auth_url)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache OAuth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func UpdateDoc(path_to_secret string, doc_id string, table *types.Table) {
	ctx := context.Background()
	b, err := os.ReadFile(path_to_secret)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/documents")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := docs.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Docs client: %v", err)
	}

	is_table_exists := clearData(srv, doc_id, table)
	if !is_table_exists {
		insertTable(srv, doc_id, table)
	}
	updateTable(srv, doc_id, !is_table_exists, table)
}

func clearData(srv *docs.Service, doc_id string, table *types.Table) bool {
	doc, err := srv.Documents.Get(doc_id).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from document: %v", err)
	}
	req := &docs.BatchUpdateDocumentRequest{}
	req.Requests = []*docs.Request{}
	for _, content := range doc.Body.Content {
		if content.Table != nil {
			if content.Table.Rows == 1 {
				return true
			}
			for i := content.Table.Rows - 1; i > 0; i-- {
				cell_location := &docs.TableCellLocation{RowIndex: i, TableStartLocation: &docs.Location{Index: content.StartIndex}}
				delete_req := &docs.DeleteTableRowRequest{TableCellLocation: cell_location}
				req.Requests = append(req.Requests, &docs.Request{DeleteTableRow: delete_req})
			}
			for i := 1; i < len(table.Rows); i++ {
				cell_location := &docs.TableCellLocation{RowIndex: 0, TableStartLocation: &docs.Location{Index: content.StartIndex}}
				insert_req := &docs.InsertTableRowRequest{TableCellLocation: cell_location, InsertBelow: true}
				req.Requests = append(req.Requests, &docs.Request{InsertTableRow: insert_req})
			}
		}
	}
	if len(req.Requests) > 0 {
		_, err := srv.Documents.BatchUpdate(doc_id, req).Do()
		if err != nil {
			log.Fatalf("Unable to update document: %v", err)
		}
		return true
	}
	return false
}

func insertTable(srv *docs.Service, doc_id string, table *types.Table) {
	req := &docs.BatchUpdateDocumentRequest{}
	table_location := &docs.Location{Index: 1}
	columns := int64(len(table.Rows[0].Cols))
	rows := int64(len(table.Rows))
	req.Requests = []*docs.Request{
		{InsertTable: &docs.InsertTableRequest{Columns: columns, Rows: rows, Location: table_location}},
	}
	if len(req.Requests) > 0 {
		_, err := srv.Documents.BatchUpdate(doc_id, req).Do()
		if err != nil {
			log.Fatalf("Unable to update document: %v", err)
		}
	}
}

func getDocsTable(srv *docs.Service, doc_id string) *docs.Table {
	doc, err := srv.Documents.Get(doc_id).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from document: %v", err)
	}
	for _, content := range doc.Body.Content {
		if content.Table != nil {
			return content.Table
		}
	}
	return nil
}

func updateTable(srv *docs.Service, doc_id string, with_header bool, table *types.Table) {
	docs_table := getDocsTable(srv, doc_id)
	req := &docs.BatchUpdateDocumentRequest{}
	req.Requests = []*docs.Request{}
	var min_row_id int
	if with_header {
		min_row_id = 0
	} else {
		min_row_id = 1
	}
	for row_id := len(table.Rows) - 1; row_id >= min_row_id; row_id-- {
		for col_id := len(table.Rows[0].Cols) - 1; col_id >= 0; col_id-- {
			location := &docs.Location{Index: docs_table.TableRows[row_id].TableCells[col_id].StartIndex + 1}
			text := table.Rows[row_id].Cols[col_id]
			req.Requests = append(req.Requests, &docs.Request{InsertText: &docs.InsertTextRequest{Location: location, Text: text}})
		}
	}

	if len(req.Requests) > 0 {
		_, err := srv.Documents.BatchUpdate(doc_id, req).Do()
		if err != nil {
			log.Fatalf("Unable to update document: %v", err)
		}
	}
}

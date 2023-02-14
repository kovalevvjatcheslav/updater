package utils

import (
	"io"
	"log"
	"updater/types"

	"github.com/PuerkitoBio/goquery"
)

func ParseHtml(r io.Reader) types.Table {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	table_class := ".confluenceTable"
	html_table := doc.Find(table_class).First()
	_, exists := html_table.Attr("class")
	if !exists {
		log.Fatalf("Element with class %s not found", table_class)
	}
	table := types.Table{}
	html_table.Find("thead").EachWithBreak(func(i int, s *goquery.Selection) bool { return addHeader(&table, i, s) })
	html_table.Find("tbody").Each(func(i int, s *goquery.Selection) { addBody(&table, i, s) })
	return table
}

func addHeader(table *types.Table, i int, thead *goquery.Selection) bool {
	style, _ := thead.Attr("style")
	if style != "display: none;" {
		thead.Find("th").Each(func(i int, th *goquery.Selection) {
			table.Header.Cols = append(table.Header.Cols, th.Text())
		})
		return false
	}
	return true
}

func addBody(table *types.Table, i int, tbody *goquery.Selection) {
	style, _ := tbody.Attr("style")
	if style != "display: none;" {
		tbody.Find("tr").Each(func(i int, tr *goquery.Selection) {
			cols := []string{}
			tr.Find("td").Each(func(i int, td *goquery.Selection) {
				cols = append(cols, td.Text())
			})
			table.Rows = append(table.Rows, types.Row{Cols: cols})
		})
	}
}

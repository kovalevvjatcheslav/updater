package utils

import (
	"io"
	"log"
	"updater/types"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
)

func ParseHtml(data io.Reader) types.Table {
	doc, err := goquery.NewDocumentFromReader(data)
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
		cols := []string{}
		thead.Find("th").Each(func(i int, th *goquery.Selection) {
			cols = append(cols, th.Text())
		})
		table.Rows = append(table.Rows, types.Row{Cols: cols})
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
				cols = append(cols, buildText(td.Nodes[0].FirstChild))
			})
			table.Rows = append(table.Rows, types.Row{Cols: cols})
		})
	}
}

func buildText(elem *html.Node) string {
	if elem == nil {
		return ""
	} else if elem.Type == html.TextNode {
		return elem.Data + buildText(elem.NextSibling)
	} else if elem.Data == "p" {
		return buildText(elem.FirstChild) + "\n" + buildText(elem.NextSibling)
	} else if elem.Data == "ul" || elem.Data == "span" || elem.Data == "code" {
		return buildText(elem.FirstChild) + buildText(elem.NextSibling)
	} else if elem.Data == "li" {
		return "\t-" + buildText(elem.FirstChild) + "\n" + buildText(elem.NextSibling)
	}

	return ""
}

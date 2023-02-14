package main

import (
	"fmt"
	_ "log"
	"reflect"
	"updater/types"
	_ "updater/utils"
)

func main() {
	var t1 = types.Table{}
	t1.Header.Cols = []string{"code", "text"}
	t1.Rows = []types.Row{
		{
			Cols: []string{"100", "test"},
		},
	}
	var t2 = types.Table{}
	t2.Header.Cols = []string{"code", "text"}
	t2.Rows = []types.Row{
		{
			Cols: []string{"100", "test"},
		},
	}
	fmt.Println(reflect.DeepEqual(t1, t2))

}

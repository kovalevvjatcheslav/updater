package main

import (
	"os"
	"reflect"
	"testing"
	"updater/types"
	"updater/utils"
)

func TestParseHtml(t *testing.T) {
	file_path := "fixtures/Коды ответа API - Документация Подсказок - Confluence.html"
	file, err := os.Open(file_path)
	if err != nil {
		t.Errorf("Cannot open file %s", file_path)
	}
	defer file.Close()
	expected_table := types.Table{}
	expected_table.Header = types.Row{Cols: []string{"HTTP-код ответа", "Описание"}}
	expected_table.Rows = []types.Row{
		{Cols: []string{"200", "Запрос успешно обработан"}},
		{Cols: []string{"400", "Некорректный запрос (невалидный JSON или XML)"}},
		{Cols: []string{"405", "Запрос сделан с методом, отличным от GET или POST"}},
		{Cols: []string{"413", "Нарушены ограничения:длина параметра query больше 300 символовили количество ограничений в параметре locations больше 100"}},
		{Cols: []string{"500", "Произошла внутренняя ошибка сервиса"}},
		{Cols: []string{"503", "Нет лицензии на запрошенный сервис"}},
	}
	test_table := utils.ParseHtml(file)
	if !reflect.DeepEqual(expected_table, test_table) {
		t.Fail()
	}
}

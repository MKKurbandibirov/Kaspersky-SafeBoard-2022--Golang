// Package csv_task
// Типы и функции для чтения и создания HTML таблицы на основе .csv файла
//
//		CsvFile - Тип хранящий все ячейки .csv файла
//
//		String - Метод для имплементация интерфейса Stringer для CsvFile
//
//		CsvFileRead - Функция чтения из .csv файла
//
//		FileToPersons - Функция сериализации данных из файла в сущности Person
//
//		GenerateDBHeader - Функция создания заголовка таблицы
//
//		GenerateDBTable - Функция создания самой таблицы
//
//		GenerateHTMLFile - Функция сохрнения таблицы в .html файл
//
package csv_task

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"test_task/internal/domain"
)

type CsvFile [][]string

func (c *CsvFile) String() string {
	builder := new(strings.Builder)
	for _, line := range *c {
		builder.WriteString(fmt.Sprint("|"))
		for _, val := range line {
			builder.WriteString(fmt.Sprintf("%25s|", val))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func CsvFileRead(file *os.File) (CsvFile, error) {
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 6

	result, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return result[1:], nil
}

func FileToPersons(file CsvFile) ([]*domain.Person, error) {
	persons := make([]*domain.Person, 0)
	for _, line := range file {
		person, err := domain.NewPerson(line, "02/01/2006")
		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}
	return persons, nil
}

func GenerateDBHeader() string {
	builder := new(strings.Builder)
	builder.WriteString("  <tr>\n")
	builder.WriteString(fmt.Sprintf("   <th>Name</th>\n"))
	builder.WriteString(fmt.Sprintf("   <th>Address</th>\n"))
	builder.WriteString(fmt.Sprintf("   <th>Postcode</th>\n"))
	builder.WriteString(fmt.Sprintf("   <th>Mobile</th>\n"))
	builder.WriteString(fmt.Sprintf("   <th>Limit</th>\n"))
	builder.WriteString(fmt.Sprintf("   <th>Birthday</th>\n"))
	builder.WriteString("  </tr>\n")
	return builder.String()
}

func GenerateDBTable(file CsvFile) (string, error) {
	result := "  <table>\n"
	result += GenerateDBHeader()
	persons, err := FileToPersons(file)
	if err != nil {
		return "", err
	}
	for _, val := range persons {
		result += val.ToHTML("02/01/2006")
	}
	result += "  </table>\n"
	return result, nil
}

func GenerateHTMLFile(file CsvFile) error {
	table, err := GenerateDBTable(file)
	if err != nil {
		return err
	}

	csvHTML, err := os.Create("csvDB.html")
	if err != nil {
		return err
	}
	csvHTML.WriteString(table)
	csvHTML.Close()
	return nil
}

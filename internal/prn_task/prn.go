// Package prn_task
// Типы и функции для чтения и создания HTML таблицы на основе .prn файла
//
//		PrnFile - Тип хранящий все ячейки .prn файла
//
//		String - Метод для имплементация интерфейса Stringer для PrnFile
//
//		DelimitersHelper - Функция-помощник для оптимального разделения столбцов в .prn файле
//
//		GetPrnDelimiters - Функция для создания меток вдоль которых стоит делить данные из файла
//
//		PrnFileRead - Функция для чтения данных из .prn файла
//
//		FileToPersons - Функция сериализации данных из файла в сущности Person
//
//		GenerateDBHeader - Функция создания заголовка таблицы
//
//		GenerateDBTable - Функция создания самой таблицы
//
//		GenerateHTMLFile - Функция сохрнения таблицы в .html файл
//
package prn_task

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"test_task/internal/domain"
)

type PrnFile [][]string

func (p *PrnFile) String() string {
	builder := new(strings.Builder)
	for _, line := range *p {
		builder.WriteString(fmt.Sprint("|"))
		for _, val := range line {
			builder.WriteString(fmt.Sprintf("%25s|", val))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func DelimitersHelper(file *os.File) (string, []string, error) {
	defer file.Close()

	reader := bufio.NewReader(file)
	header, err := reader.ReadString('\n')
	if err != nil {
		return "", nil, err
	}
	helper := make([]string, 0)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", nil, err
		}
		helper = append(helper, line)
	}
	return header, helper, nil
}

func GetPrnDelimiters(header string, helper []string) []int {
	delimPositions := make([]int, 0)
	for _, val := range [...]string{"First name", "Address", "Postcode", "Mobile", "Limit", "Birthday"} {
		delimPositions = append(delimPositions, strings.Index(header, val))
	}
	delimPositions = append(delimPositions, len(header)-1)
	for i := delimPositions[3]; i < delimPositions[4]; i++ {
		j := 0
		for ; j < len(helper); j++ {
			if helper[j][i] != ' ' {
				break
			}
		}
		if j == len(helper) {
			delimPositions[4] = i
			break
		}
	}
	return delimPositions
}

func PrnFileRead(file *os.File) (PrnFile, error) {
	header, helper, err := DelimitersHelper(file)
	if err != nil {
		return nil, err
	}

	delims := GetPrnDelimiters(header, helper)

	var result = make([][]string, 0)
	for _, val := range helper {
		var tmp = make([]string, 0)
		for i := 1; i < len(delims); i++ {
			elem := string([]byte(val)[delims[i-1]:delims[i]])
			elem = strings.Trim(elem, " ")
			tmp = append(tmp, elem)
		}
		result = append(result, tmp)
	}
	return result, nil
}

func FileToPersons(file PrnFile) ([]*domain.Person, error) {
	persons := make([]*domain.Person, 0)
	for _, line := range file {
		person, err := domain.NewPerson(line, "20060102")
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
	builder.WriteString(fmt.Sprintf("   <th>First name</th>\n"))
	builder.WriteString(fmt.Sprintf("   <th>Address</th>\n"))
	builder.WriteString(fmt.Sprintf("   <th>Postcode</th>\n"))
	builder.WriteString(fmt.Sprintf("   <th>Mobile</th>\n"))
	builder.WriteString(fmt.Sprintf("   <th>Limit</th>\n"))
	builder.WriteString(fmt.Sprintf("   <th>Birthday</th>\n"))
	builder.WriteString("  </tr>\n")
	return builder.String()
}

func GenerateDBTable(file PrnFile) (string, error) {
	result := "  <table>\n"
	result += GenerateDBHeader()
	persons, err := FileToPersons(file)
	if err != nil {
		return "", err
	}
	for _, val := range persons {
		result += val.ToHTML("20060102")
	}
	result += "  </table>\n"
	return result, nil
}

func GenerateHTMLFile(file PrnFile) error {
	table, err := GenerateDBTable(file)
	if err != nil {
		return err
	}

	csvHTML, err := os.Create("prnDB.html")
	if err != nil {
		return err
	}
	csvHTML.WriteString(table)
	csvHTML.Close()
	return nil
}

// Package domain
// Типы для реализации сущностей из файлов
//
//		Person - Сущность лежащая в каждом файле
//
//		NewPerson - Конструктор для типа Person
//
//		String - Метод для имплементация интерфейса Stringer для Person
//
//		ToHTML - Метод превращающий сущность в строку HTML-таблицы
//
package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Person struct {
	Name     string
	Address  string
	Postcode string
	Mobile   string
	Limit    float64
	Birthday time.Time
}

func NewPerson(csvLine []string, layout string) (*Person, error) {
	limit, err := strconv.ParseFloat(csvLine[4], 64)
	if err != nil {
		return nil, err
	}
	birthday, err := time.Parse(layout, csvLine[5])
	if err != nil {
		return nil, err
	}
	return &Person{
		Name:     csvLine[0],
		Address:  csvLine[1],
		Postcode: csvLine[2],
		Mobile:   csvLine[3],
		Limit:    limit,
		Birthday: birthday,
	}, nil
}

func (p *Person) String() string {
	return fmt.Sprintf("|%30s|%30s|%30s|%30s|%30f|%30v|\n",
		p.Name, p.Address, p.Postcode, p.Mobile, p.Limit, p.Birthday)
}

func (p *Person) ToHTML(layout string) string {
	builder := new(strings.Builder)
	builder.WriteString("  <tr>\n")
	builder.WriteString(fmt.Sprintf("   <td>%s</td>\n", p.Name))
	builder.WriteString(fmt.Sprintf("   <td>%s</td>\n", p.Address))
	builder.WriteString(fmt.Sprintf("   <td>%s</td>\n", p.Postcode))
	builder.WriteString(fmt.Sprintf("   <td>%s</td>\n", p.Mobile))
	builder.WriteString(fmt.Sprintf("   <td>%.2f</td>\n", p.Limit))
	builder.WriteString(fmt.Sprintf("   <td>%s</td>\n", p.Birthday.Format(layout)))
	builder.WriteString("  </tr>\n")
	return builder.String()
}

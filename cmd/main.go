// Package main
// Главный пакет
//
//		OpenFile - Функция для открытия файла
//
//		Start - Функция для асинхронного преобразования файлов
//
package main

import (
	"flag"
	"log"
	"os"
	"sync"
	csvpkg "test_task/internal/csv_task"
	prnpkg "test_task/internal/prn_task"
)

func OpenFile(name string) (*os.File, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func Start(csvFile, prnFile *os.File) {
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()
		csv, err := csvpkg.CsvFileRead(csvFile)
		if err != nil {
			log.Fatal(err)
		}
		err = csvpkg.GenerateHTMLFile(csv)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		prn, err := prnpkg.PrnFileRead(prnFile)
		if err != nil {
			log.Fatal(err)
		}
		err = prnpkg.GenerateHTMLFile(prn)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}

func main() {
	// Получение файлов через соответствующие флаги
	prnFileName := flag.String("prn", "", "Enter a .prn file path")
	csvFileName := flag.String("csv", "", "Enter a .csv file path")
	flag.Parse()

	csvFile, err := OpenFile(*csvFileName)
	if err != nil {
		log.Fatal(err)
	}

	prnFile, err := OpenFile(*prnFileName)
	if err != nil {
		log.Fatal(err)
	}
	Start(csvFile, prnFile)
}

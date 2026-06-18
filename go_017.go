package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	/*
		///////////////////////////////////////
		// Стандартный ввод с клавиатуры
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			fmt.Printf("input = %v\n", input.Text())
		}
		//input = dfghjjkkkkkkkkk
	*/

	/*
		///////////////////////////////////////
		// Аргументы командной строки
		args := os.Args[0:len(os.Args)]
		fmt.Println("len = ", len(os.Args))  // len = 2
		fmt.Println("elem[0] = ", args[0])   // go_017
		fmt.Println("elem[1] = ", args[1])   // help
	*/

	///////////////////////////////////////
	// Read Files
	files := os.Args[1:len(os.Args)]
	if len(files) == 0 {
		fmt.Println("No file")
	} else {
		fmt.Println("Number of files: ", len(files))
		for _, arg01 := range files {
			file01, err01 := os.Open(arg01)
			if err01 != nil {
				fmt.Println("Error reading file: ", err01)
				continue
			} else {
				fmt.Println("New file: ", file01.Name())
			}
			input01 := bufio.NewScanner(file01)
			// Построчное чтение файла
			for input01.Scan() {
				fmt.Println(input01.Text()) // Построчный вывод
			}
			file01.Close()
		}
	}
	/*
		//Example
		F:\URA\PROG\GO\go_017>go_017 t c:\temp2\text02.txt
		Number of files:  2
		Error reading file:  open t: The system cannot find the file specified.
		New file:  c:\temp2\text02.txt
		asd
		fgh
	*/
}

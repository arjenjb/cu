package main

import (
	"fmt"
	"strings"
)

func main() {
	b := strings.Builder{}
	for i := 0; i < 10; i++ {
		b.WriteString(fmt.Sprintf("%d---|----|", i))
	}
	fmt.Println(b.String())
	//cursorLineUp(1)
	//cursorColumn(0)
	cursorPosition(3, 3)
	fmt.Println("###")
}

func cursorPosition(line int, col int) {
	fmt.Printf("\033[%d;%dH", line, col)
}

func cursorColumn(i int) {
	fmt.Printf("\033[%dG", i)
}

func cursorLineUp(i int) {
	fmt.Printf("\033[%dA", i)
}

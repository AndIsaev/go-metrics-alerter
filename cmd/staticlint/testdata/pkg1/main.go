package main

import "os"

func main() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	os.Exit(1) // want "direct call to os.Exit found in ..."
}

// Package main содержит программу, которая использует multichecker для статического анализа кода.
// Она включает несколько анализаторов, в том числе и пользовательский анализатор ExitCheckAnalyzer.

package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	sa := make([]*analysis.Analyzer, 0, 100)

	// Добавляем все анализаторы класса SA из пакета staticcheck
	for _, analyzer := range staticcheck.Analyzers {
		sa = append(sa, analyzer.Analyzer)
	}

	// Запускаем multichecker с набором анализаторов, включая пользовательский ExitCheckAnalyzer
	multichecker.Main(
		append(sa,
			ExitCheckAnalyzer,  // Запрещающий использовать прямой вызов os.Exit в функции main пакета main
			printf.Analyzer,    // Печатает предупреждения о неправильных функциях форматирования
			shadow.Analyzer,    // Ищет элементы переменной, которые скрывают друг друга
			structtag.Analyzer, // Проверяет наличие ошибок в структурных тегах
			shift.Analyzer,     // Проверяет использование операций сдвига
		)...,
	)
}

package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// ExitCheckAnalyzer проверяет наличие прямых вызовов os.Exit в функции main пакета main.
var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for direct calls to os.Exit in main function of package main",
	Run:  run,
}

// run выполняет анализатор ExitCheckAnalyzer. Он проходит по файлам,
// принадлежащим пакету main и ищет вызовы os.Exit в функции main.
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Убедимся, что анализируем только пакет main
		if pass.Pkg.Name() != "main" {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			// Проверяем, что вызов относится к os.Exit
			pkgIdent, ok := selExpr.X.(*ast.Ident)
			if !ok || pkgIdent.Name != "os" || selExpr.Sel.Name != "Exit" {
				return true
			}

			pos := pass.Fset.Position(callExpr.Pos())
			pass.Reportf(callExpr.Pos(), "direct call to os.Exit found in %s:%d", pos.Filename, pos.Line)
			return false
		})
	}
	return nil, nil
}

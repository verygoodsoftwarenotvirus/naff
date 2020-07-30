package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func getExpectedOutputForTest(testName string) string {
	cmd := exec.Command("go", "test", "gitlab.com/verygoodsoftwarenotvirus/naff/templates/database/v1/client", "-run", testName)

	outputBytes, err := cmd.CombinedOutput()
	if err != nil && fmt.Sprintf("%v", err) != "exit status 1" { // lol that second condition is groooooooooossssssss
		log.Printf("error occurred getting combined output: %v", err)
	}
	output := string(outputBytes)

	var (
		actualStartLine int
		actualEndLine   int
	)
	actualStartRegexp := regexp.MustCompile(`^\s+actual\s+:\s+`)
	actualEndRegexp := regexp.MustCompile(`^\s+Diff:`)

	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if actualStartRegexp.MatchString(line) {
			actualStartLine = i
		}
		if actualStartLine != 0 && actualEndRegexp.MatchString(line) {
			actualEndLine = i - 1
			break
		}
	}

	if actualStartLine > 0 && actualEndLine > actualStartLine {
		expectedOutput := actualStartRegexp.ReplaceAllString(strings.Join(lines[actualStartLine:actualEndLine], ""), "")
		if strings.HasPrefix(expectedOutput, `"`) {
			expectedOutput = expectedOutput[1:]
		}
		if strings.HasSuffix(expectedOutput, `"`) {
			expectedOutput = expectedOutput[:len(expectedOutput)-1]
		}
		expectedOutput = strings.ReplaceAll(expectedOutput, "\\n", "\n")
		expectedOutput = strings.ReplaceAll(expectedOutput, "\\t", "\t")
		expectedOutput = strings.ReplaceAll(expectedOutput, `\"`, `"`)
		expectedOutput = strings.ReplaceAll(expectedOutput, "`", "`+\"`\"+`")

		return expectedOutput
	}

	return ""
}

type foundTest struct {
	name       string
	lineNumber int
}

func listTestsForFile(filePath string) []foundTest {
	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	fset := token.NewFileSet()
	code, err := parser.ParseFile(fset, filePath, src, parser.AllErrors)
	if err != nil {
		log.Fatalf("error parsing file: %v", err)
	}

	functions := []foundTest{}

	for _, f := range code.Decls {
		if fn, ok := f.(*ast.FuncDecl); ok {
			fun := foundTest{
				name: fn.Name.Name,
			}
			functions = append(functions)
			if len(fn.Body.List) == 2 {
				if tDotRunCall, tdrExprOk := fn.Body.List[1].(*ast.ExprStmt); tdrExprOk {
					if subtestCallExpr, subtestCallExprOk := tDotRunCall.X.(*ast.CallExpr); subtestCallExprOk {
						if len(subtestCallExpr.Args) == 2 {
							if funcLit, funcLitOk := subtestCallExpr.Args[1].(*ast.FuncLit); funcLitOk {
								if len(funcLit.Body.List) == 6 {
									if assStmt, assStmtOk := funcLit.Body.List[3].(*ast.AssignStmt); assStmtOk {
										expectedAssign := assStmt.Lhs[0].(*ast.Ident)
										expectedAssignVal := assStmt.Rhs[0].(*ast.BasicLit)
										if expectedAssign.Name == "expected" && expectedAssignVal.Value == "``" {
											fun.lineNumber = fset.Position(expectedAssign.Pos()).Line
											functions = append(functions, fun)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return functions
}

func correctTestsForFile(filePath string, lineNumber int, results string) {
	fileAsBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("error reading test file to transform: %v", err)
	}
	lines := strings.Split(string(fileAsBytes), "\n")

	newFile := strings.Join(append(
		append(lines[:lineNumber-1], fmt.Sprintf("\t\texpected := `%s`", results)),
		lines[lineNumber:]...,
	), "\n")

	if err := ioutil.WriteFile(filePath, []byte(newFile), 0644); err != nil {
		log.Fatalf("error writing new file %q: %v", filePath, err)
	}
}

func listTestFilesInDirectory(dirPath string) []string {
	out := []string{}

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("error reading directory: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "_test.go") {
			out = append(out, filepath.Join(dirPath, file.Name()))
		}
	}

	return out
}

func main() {
	for _, filePath := range listTestFilesInDirectory("templates/database/v1/client") {
		testsForFile := listTestsForFile(filePath)
		for i := 0; i < len(testsForFile); i++ {
			tests := listTestsForFile(filePath)
			if len(tests) > 0 {
				test := tests[0]
				if test.lineNumber > 0 {
					result := getExpectedOutputForTest(fmt.Sprintf("%s/obligatory", test.name))
					if result != "" {
						correctTestsForFile(filePath, test.lineNumber, result)
					}
				}
			}
		}
	}
}

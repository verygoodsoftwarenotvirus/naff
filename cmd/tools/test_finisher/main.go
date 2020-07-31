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

func getExpectedOutputForTest(packagePath, testName string) string {
	cmd := exec.Command("go", "test", filepath.Join("gitlab.com/verygoodsoftwarenotvirus/naff", packagePath), "-run", testName)

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
		expectedOutput = strings.ReplaceAll(expectedOutput, `\\`, `\`)

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
			for _, decl := range fn.Body.List {
				if tDotRunCall, tdrExprOk := decl.(*ast.ExprStmt); tdrExprOk {
					if subtestCallExpr, subtestCallExprOk := tDotRunCall.X.(*ast.CallExpr); subtestCallExprOk {
						if len(subtestCallExpr.Args) == 2 {
							if subtestName, subtestNameOk := subtestCallExpr.Args[0].(*ast.BasicLit); subtestNameOk {
								fun.name = fmt.Sprintf("%s/%s", fn.Name.Name, strings.ReplaceAll(subtestName.Value, `"`, ""))
							}
							if funcLit, funcLitOk := subtestCallExpr.Args[1].(*ast.FuncLit); funcLitOk {
								for _, expr := range funcLit.Body.List {
									if assStmt, assStmtOk := expr.(*ast.AssignStmt); assStmtOk {
										if expectedAssign, expectedAssignOk := assStmt.Lhs[0].(*ast.Ident); expectedAssignOk {
											if expectedAssignVal, expectedAssignValOk := assStmt.Rhs[0].(*ast.BasicLit); expectedAssignValOk {
												if expectedAssign.Name == "expected" && expectedAssignVal.Value == "``" {
													fun.lineNumber = fset.Position(expectedAssign.Pos()).Line
													functions = append(functions, fun)
													break
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
	const packagePath = "templates/database/v1/queriers"

	for _, filePath := range listTestFilesInDirectory(packagePath) {
		testsForFile := listTestsForFile(filePath)
		for i := 0; i < len(testsForFile); i++ {
			tests := listTestsForFile(filePath)
			if len(tests) > 0 {
				test := tests[0]
				if test.lineNumber > 0 {
					result := getExpectedOutputForTest(packagePath, test.name)
					if result != "" {
						log.Printf("correcting test: %q", test.name)
						correctTestsForFile(filePath, test.lineNumber, result)
					}
				}
			}
		}
	}
}

package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Project структура
type Project struct {
	Name        string
	Content     string
	Comments    string
	Identifiers Identifiers
	ControlFlow ControlFlow
	Functions   FunctionAnalysis
	Imports     ImportAnalysis
	Formatting  FormatAnalysis
	Tokens      TokenAnalysis
	FilePath    string
	Language    string // Добавляем определение языка программирования
}

// ComparisonResult хранит результат сравнения двух проектов
type ComparisonResult struct {
	Project1              string
	Project2              string
	Similarity            float64
	CommentSimilarity     float64
	IdentifierSimilarity  float64
	ControlFlowSimilarity float64
	FunctionSimilarity    float64
	ImportSimilarity      float64
	FormatSimilarity      float64
	TokenSimilarity       float64
	Language              string
}

// Добавляем структуру для хранения идентификаторов
type Identifiers struct {
	Variables  []string
	Functions  []string
	Classes    []string
	Interfaces []string
	Constants  []string
}

// Добавляем структуру для анализа потока управления
type ControlFlow struct {
	IfCount        int
	ForCount       int
	WhileCount     int
	SwitchCount    int
	MaxNesting     int
	ControlPattern string // паттерн последовательности управляющих конструкций
}

// Добавляем структуру для анализа функций
type FunctionAnalysis struct {
	ParamCount    map[string]int    // количество параметров для каждой функции
	ReturnTypes   map[string]string // типы возвращаемых значений
	FunctionSizes map[string]int    // размеры функций
	DeclareOrder  []string          // порядок объявления функций
}

// Добавляем структуру для анализа импортов
type ImportAnalysis struct {
	ImportList    []string          // список импортов
	ImportOrder   string            // порядок импортов
	UsagePatterns map[string]string // паттерны использования импортированных функций
}

// Добавляем структуру для анализа форматирования
type FormatAnalysis struct {
	IndentStyle    string
	SpacingPattern string
	LineBreaks     string
}

// Добавляем структуру для анализа токенов
type TokenAnalysis struct {
	OperatorSequence string
	TokenPatterns    []string
}

// Добавим список поддерживаемых расширений файлов
var supportedExtensions = map[string]string{
	".go":   "golang",
	".py":   "python",
	".java": "java",
	".js":   "javascript",
	".cpp":  "c++",
	".c":    "c",
	".cs":   "c#",
	".php":  "php",
	".rb":   "ruby",
	".rs":   "rust",
}

// Добавьте новую структуру для HTML отчета
type HtmlReport struct {
	GeneratedTime         string
	GeneratedDate         string // Добавляем дату
	TotalProjects         int
	TotalComparisons      int
	Results               []ComparisonResult
	AverageSimilarity     float64
	HighSimilarityCount   int
	MediumSimilarityCount int
	LowSimilarityCount    int
}

func main() {
	projectsDir := "./projects" // директория с проектами студентов
	projects, err := loadProjects(projectsDir)
	if err != nil {
		fmt.Printf("О��ибка при загрузке проектов: %v\n", err)
		return
	}

	results := compareAllProjects(projects)
	printResults(results)
}

// loadProjects загружает все проекты из указанной директории
func loadProjects(dir string) ([]Project, error) {
	var projects []Project
	fmt.Println("Начинаю загрузку проектов...")
	fmt.Printf("Поиск проектов в директории: %s\n", dir)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if lang, ok := supportedExtensions[ext]; ok {
				fmt.Printf("Обнаружен файл: %s (язык: %s)\n", path, lang)
				content, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				rawContent := string(content)
				projectName := filepath.Base(filepath.Dir(path))
				fmt.Printf("Анализирую проект: %s\n", projectName)

				project := Project{
					Name:        projectName,
					Content:     normalizeCode(rawContent, lang),
					Comments:    extractComments(rawContent, lang),
					Identifiers: extractIdentifiers(rawContent, lang),
					FilePath:    path,
					Language:    lang,
				}
				projects = append(projects, project)
				fmt.Printf("Проект %s успешно загружен и проанализирован\n", projectName)
				fmt.Println("----------------------------------------")
			}
		}
		return nil
	})

	fmt.Printf("Загружено проектов: %d\n\n", len(projects))
	return projects, err
}

// normalizeCode подготавливает код для сравнения
func normalizeCode(code, language string) string {
	// Удаляем комментарии в зависимости от языка
	switch language {
	case "python":
		// Удаляем Python-комментарии
		lines := strings.Split(code, "\n")
		var filtered []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "#") && line != "" {
				filtered = append(filtered, line)
			}
		}
		code = strings.Join(filtered, " ")
	default:
		// Для остальных языков удаляем стандартные комментарии
		lines := strings.Split(code, "\n")
		var filtered []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "/*") &&
				!strings.HasSuffix(line, "*/") && line != "" {
				filtered = append(filtered, line)
			}
		}
		code = strings.Join(filtered, " ")
	}

	// Общая нормализация
	code = strings.ToLower(code)

	// Удаляем строковые литералы
	code = removeStringLiterals(code)

	// Удаляем лишние пробелы
	code = strings.Join(strings.Fields(code), " ")

	return code
}

// Добавляем функцию для удаления строковых литералов
func removeStringLiterals(code string) string {
	// Удаляем строки в двойных кавычках
	doubleQuoteRegex := regexp.MustCompile(`"[^"]*"`)
	code = doubleQuoteRegex.ReplaceAllString(code, "\"\"")

	// Удаляем строки в одинарных кавычках
	singleQuoteRegex := regexp.MustCompile(`'[^']*'`)
	code = singleQuoteRegex.ReplaceAllString(code, "''")

	return code
}

// Добавляем функцию для извлечения комментариев
func extractComments(code, language string) string {
	var comments []string
	lines := strings.Split(code, "\n")

	switch language {
	case "python":
		// Извлекаем Python-комментарии
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "#") {
				comment := strings.TrimPrefix(line, "#")
				comments = append(comments, strings.TrimSpace(comment))
			}
		}
	default:
		// Для остальных языков извлекаем стандартные комментарии
		inMultilineComment := false
		for _, line := range lines {
			line = strings.TrimSpace(line)

			// Обработка многострочных комментариев
			if strings.Contains(line, "/*") && strings.Contains(line, "*/") {
				// Проверяем корректность индексов
				startIdx := strings.Index(line, "/*")
				endIdx := strings.Index(line, "*/")
				if startIdx < endIdx { // Добавляем проверку
					comment := line[startIdx+2 : endIdx]
					comments = append(comments, strings.TrimSpace(comment))
				}
			} else if strings.Contains(line, "/*") {
				inMultilineComment = true
				startIdx := strings.Index(line, "/*")
				if startIdx < len(line)-2 { // Добавляем проверку
					comment := line[startIdx+2:]
					if comment != "" {
						comments = append(comments, strings.TrimSpace(comment))
					}
				}
			} else if strings.Contains(line, "*/") {
				inMultilineComment = false
				endIdx := strings.Index(line, "*/")
				if endIdx > 0 { // Добавляем проверку
					comment := line[:endIdx]
					if comment != "" {
						comments = append(comments, strings.TrimSpace(comment))
					}
				}
			} else if inMultilineComment {
				if line != "" {
					comments = append(comments, strings.TrimSpace(line))
				}
			} else if strings.HasPrefix(line, "//") {
				comment := strings.TrimPrefix(line, "//")
				comments = append(comments, strings.TrimSpace(comment))
			}
		}
	}

	// Нормализуем комментарии
	normalizedComments := strings.Join(comments, " ")
	normalizedComments = strings.ToLower(normalizedComments)
	normalizedComments = strings.Join(strings.Fields(normalizedComments), " ")

	return normalizedComments
}

// Добавляем функцию для извлечения идентификаторов
func extractIdentifiers(code, language string) Identifiers {
	var ids Identifiers

	// Удаляем комментарии и строковые литералы для чистого анализа
	code = removeComments(code, language)
	code = removeStringLiterals(code)

	switch language {
	case "golang":
		// Извлекаем переменные (var name type, name :=)
		varRegex := regexp.MustCompile(`(?m)(var\s+(\w+))|(\w+\s*:=)`)
		for _, match := range varRegex.FindAllStringSubmatch(code, -1) {
			if match[2] != "" {
				ids.Variables = append(ids.Variables, match[2])
			} else if match[3] != "" {
				name := strings.TrimSpace(strings.TrimSuffix(match[3], ":="))
				ids.Variables = append(ids.Variables, name)
			}
		}

		// Извлекаем функции (func name)
		funcRegex := regexp.MustCompile(`func\s+(\w+)\s*\(`)
		for _, match := range funcRegex.FindAllStringSubmatch(code, -1) {
			ids.Functions = append(ids.Functions, match[1])
		}

		// Извлекаем типы/структуры (type name)
		typeRegex := regexp.MustCompile(`type\s+(\w+)\s+struct`)
		for _, match := range typeRegex.FindAllStringSubmatch(code, -1) {
			ids.Classes = append(ids.Classes, match[1])
		}

		// Извлекаем константы (const name)
		constRegex := regexp.MustCompile(`const\s+(\w+)`)
		for _, match := range constRegex.FindAllStringSubmatch(code, -1) {
			ids.Constants = append(ids.Constants, match[1])
		}

	case "python":
		// Извлекаем переменные (name =)
		varRegex := regexp.MustCompile(`(\w+)\s*=\s*[^\=]`)
		for _, match := range varRegex.FindAllStringSubmatch(code, -1) {
			ids.Variables = append(ids.Variables, match[1])
		}

		// Извлекаем функции (def name)
		funcRegex := regexp.MustCompile(`def\s+(\w+)\s*\(`)
		for _, match := range funcRegex.FindAllStringSubmatch(code, -1) {
			ids.Functions = append(ids.Functions, match[1])
		}

		// Извлекаем классы (class name)
		classRegex := regexp.MustCompile(`class\s+(\w+)`)
		for _, match := range classRegex.FindAllStringSubmatch(code, -1) {
			ids.Classes = append(ids.Classes, match[1])
		}

	case "java":
		// Извлекаем переменные (type name)
		varRegex := regexp.MustCompile(`(?m)(private|public|protected)?\s+\w+\s+(\w+)\s*[;=]`)
		for _, match := range varRegex.FindAllStringSubmatch(code, -1) {
			ids.Variables = append(ids.Variables, match[2])
		}

		// Извлекаем методы
		methodRegex := regexp.MustCompile(`(?m)(private|public|protected)?\s+\w+\s+(\w+)\s*\(`)
		for _, match := range methodRegex.FindAllStringSubmatch(code, -1) {
			ids.Functions = append(ids.Functions, match[2])
		}

		// Извлекаем классы
		classRegex := regexp.MustCompile(`class\s+(\w+)`)
		for _, match := range classRegex.FindAllStringSubmatch(code, -1) {
			ids.Classes = append(ids.Classes, match[1])
		}

		// Извлекаем интерфейсы
		interfaceRegex := regexp.MustCompile(`interface\s+(\w+)`)
		for _, match := range interfaceRegex.FindAllStringSubmatch(code, -1) {
			ids.Interfaces = append(ids.Interfaces, match[1])
		}
	}

	return ids
}

// Добавляем функцию для сравнения идентификаторов
func compareIdentifiers(ids1, ids2 Identifiers) float64 {
	var matches, total float64

	// Сравниваем переменные
	matches += float64(len(findCommonElements(ids1.Variables, ids2.Variables)))
	total += float64(len(ids1.Variables))

	// Сравниваем функции
	matches += float64(len(findCommonElements(ids1.Functions, ids2.Functions)))
	total += float64(len(ids1.Functions))

	// Сравниваем классы
	matches += float64(len(findCommonElements(ids1.Classes, ids2.Classes)))
	total += float64(len(ids1.Classes))

	// Сравниваем интерфейсы
	matches += float64(len(findCommonElements(ids1.Interfaces, ids2.Interfaces)))
	total += float64(len(ids1.Interfaces))

	// Сравниваем константы
	matches += float64(len(findCommonElements(ids1.Constants, ids2.Constants)))
	total += float64(len(ids1.Constants))

	if total == 0 {
		return 0
	}

	return (matches / total) * 100
}

// Вспомогательная функция для поиска общих элементов
func findCommonElements(slice1, slice2 []string) []string {
	var common []string
	seen := make(map[string]bool)

	for _, item := range slice1 {
		seen[item] = true
	}

	for _, item := range slice2 {
		if seen[item] {
			common = append(common, item)
		}
	}

	return common
}

// Добавляем функцию для анализа потока управления
func analyzeControlFlow(code, language string) ControlFlow {
	var cf ControlFlow

	switch language {
	case "golang":
		// Подсчет if
		ifRegex := regexp.MustCompile(`\bif\b`)
		cf.IfCount = len(ifRegex.FindAllString(code, -1))

		// Подсчет for
		forRegex := regexp.MustCompile(`\bfor\b`)
		cf.ForCount = len(forRegex.FindAllString(code, -1))

		// Подсчет switch
		switchRegex := regexp.MustCompile(`\bswitch\b`)
		cf.SwitchCount = len(switchRegex.FindAllString(code, -1))

		// Анализ вложенности
		cf.MaxNesting = analyzeNesting(code)

		// Создание паттерна управляющих конструкций
		cf.ControlPattern = createControlPattern(code)

	case "python":
		// Аналогичный анализ для Python
		ifRegex := regexp.MustCompile(`\bif\b`)
		cf.IfCount = len(ifRegex.FindAllString(code, -1))

		forRegex := regexp.MustCompile(`\bfor\b`)
		cf.ForCount = len(forRegex.FindAllString(code, -1))

		whileRegex := regexp.MustCompile(`\bwhile\b`)
		cf.WhileCount = len(whileRegex.FindAllString(code, -1))
	}

	return cf
}

// Функция для анализа вложенности
func analyzeNesting(code string) int {
	lines := strings.Split(code, "\n")
	maxNesting := 0
	currentNesting := 0

	for _, line := range lines {
		indent := countIndentation(line)
		currentNesting = indent / 4 // предполагаем отступ в 4 пробела
		if currentNesting > maxNesting {
			maxNesting = currentNesting
		}
	}

	return maxNesting
}

// Функция для анализа функций
func analyzeFunctions(code, language string) FunctionAnalysis {
	fa := FunctionAnalysis{
		ParamCount:    make(map[string]int),
		ReturnTypes:   make(map[string]string),
		FunctionSizes: make(map[string]int),
	}

	switch language {
	case "golang":
		// Анализ функций Go
		funcRegex := regexp.MustCompile(`func\s+(\w+)\s*\((.*?)\)\s*(.*?)\s*{`)
		matches := funcRegex.FindAllStringSubmatch(code, -1)

		for _, match := range matches {
			funcName := match[1]
			params := match[2]
			returnType := match[3]

			fa.DeclareOrder = append(fa.DeclareOrder, funcName)
			fa.ParamCount[funcName] = len(strings.Split(params, ","))
			fa.ReturnTypes[funcName] = returnType
			fa.FunctionSizes[funcName] = calculateFunctionSize(code, funcName)
		}
	}

	return fa
}

// Функция для анализа импортов
func analyzeImports(code, language string) ImportAnalysis {
	ia := ImportAnalysis{
		UsagePatterns: make(map[string]string),
	}

	switch language {
	case "golang":
		// Анализ импортов Go
		importRegex := regexp.MustCompile(`import\s*\((.*?)\)`)
		matches := importRegex.FindStringSubmatch(code)
		if len(matches) > 1 {
			imports := strings.Split(matches[1], "\n")
			for _, imp := range imports {
				imp = strings.TrimSpace(imp)
				if imp != "" {
					ia.ImportList = append(ia.ImportList, imp)
				}
			}
		}
	}

	ia.ImportOrder = strings.Join(ia.ImportList, ",")
	return ia
}

// Функция для анализа форматирования
func analyzeFormatting(code string) FormatAnalysis {
	var fa FormatAnalysis

	// Анализ стиля отступов
	if strings.Contains(code, "\t") {
		fa.IndentStyle = "tabs"
	} else {
		fa.IndentStyle = "spaces"
	}

	// Анализ пробелов вокруг операторов
	spacingRegex := regexp.MustCompile(`\w+\s*[+\-*/=]\s*\w+`)
	matches := spacingRegex.FindAllString(code, -1)
	fa.SpacingPattern = strings.Join(matches, ";")

	// Анализ переносов строк
	fa.LineBreaks = analyzeLineBreaks(code)

	return fa
}

// Функция для сравнения потока управления
func compareControlFlow(cf1, cf2 ControlFlow) float64 {
	total := 0.0
	matches := 0.0

	// Сравнение ��оличества управляющих конструкций
	if cf1.IfCount == cf2.IfCount {
		matches++
	}
	total++

	if cf1.ForCount == cf2.ForCount {
		matches++
	}
	total++

	if cf1.WhileCount == cf2.WhileCount {
		matches++
	}
	total++

	if cf1.SwitchCount == cf2.SwitchCount {
		matches++
	}
	total++

	// Сравнение вложенности
	if cf1.MaxNesting == cf2.MaxNesting {
		matches++
	}
	total++

	// Сравнение аттерна управляющих конструкций
	if cf1.ControlPattern == cf2.ControlPattern {
		matches += 2 // придаем больший вес этому критерию
	}
	total += 2

	return (matches / total) * 100
}

// Обновляем функцию compareProjects
func compareProjects(p1, p2 Project) ComparisonResult {
	result := ComparisonResult{
		Project1: p1.Name,
		Project2: p2.Name,
		Language: p1.Language,
	}

	// Базовое сравнение кода
	result.Similarity = compareTexts(p1.Content, p2.Content)

	// Сравнение комментариев
	if p1.Comments != "" && p2.Comments != "" {
		result.CommentSimilarity = compareTexts(p1.Comments, p2.Comments)
	}

	// Сравнение идентификаторов
	result.IdentifierSimilarity = compareIdentifiers(p1.Identifiers, p2.Identifiers)

	// Сравнение потока управления
	result.ControlFlowSimilarity = compareControlFlow(p1.ControlFlow, p2.ControlFlow)

	// Сравнение функций
	result.FunctionSimilarity = compareFunctions(p1.Functions, p2.Functions)

	// Сравнение импортов
	result.ImportSimilarity = compareImports(p1.Imports, p2.Imports)

	// Сравнение форматирования
	result.FormatSimilarity = compareFormatting(p1.Formatting, p2.Formatting)

	return result
}

// Добавляем вспомогательную функцию для сравнения текстов
func compareTexts(text1, text2 string) float64 {
	s1 := strings.Split(text1, " ")
	s2 := strings.Split(text2, " ")

	matches := 0
	totalTokens := len(s1)

	for _, token := range s1 {
		if token == "" {
			continue
		}
		for _, token2 := range s2 {
			if token == token2 {
				matches++
				break
			}
		}
	}

	if totalTokens == 0 {
		return 0
	}

	return float64(matches) / float64(totalTokens) * 100
}

// Обновим функцию compareAllProjects для вывода прогресса
func compareAllProjects(projects []Project) []ComparisonResult {
	var results []ComparisonResult
	totalComparisons := (len(projects) * (len(projects) - 1)) / 2
	currentComparison := 0

	fmt.Printf("Начинаю сравнение %d проектов (%d сравнений)...\n",
		len(projects), totalComparisons)

	for i := 0; i < len(projects); i++ {
		for j := i + 1; j < len(projects); j++ {
			currentComparison++
			fmt.Printf("\rПрогресс: %d/%d сравнений", currentComparison, totalComparisons)

			if projects[i].Language == projects[j].Language {
				fmt.Printf("\nСравниваю проекты:\n")
				fmt.Printf("  1. %s (%s)\n", projects[i].Name, projects[i].FilePath)
				fmt.Printf("  2. %s (%s)\n", projects[j].Name, projects[j].FilePath)

				result := compareProjects(projects[i], projects[j])
				if result.Similarity > 50 {
					results = append(results, result)
					fmt.Printf("  Обнаружена схожесть: %.2f%%\n", result.Similarity)
				} else {
					fmt.Printf("  Схожесть ниже порога (%.2f%%)\n", result.Similarity)
				}
				fmt.Println("----------------------------------------")
			}
		}
	}
	fmt.Printf("\nЗавершено сравнение проектов\n\n")
	return results
}

// Обновляем функцию printResults для более подробного вывода
func printResults(results []ComparisonResult) {
	if len(results) == 0 {
		fmt.Println("Подозрительных совпадений не обнаружено")
		return
	}

	// Вычисляем общую статистику
	var totalSimilarity float64
	highSim := 0
	mediumSim := 0
	lowSim := 0

	for _, result := range results {
		// Вычисляем среднюю схожесть для каждого сравнения
		avgSim := (result.Similarity +
			result.CommentSimilarity +
			result.IdentifierSimilarity +
			result.ControlFlowSimilarity +
			result.FunctionSimilarity +
			result.ImportSimilarity +
			result.FormatSimilarity) / 7.0

		totalSimilarity += avgSim

		// Подсчитываем количество разных уровней схожести
		if avgSim >= 80 {
			highSim++
		} else if avgSim >= 60 {
			mediumSim++
		} else {
			lowSim++
		}
	}

	averageSimilarity := totalSimilarity / float64(len(results))

	// Выводим краткую статистику в консоль
	fmt.Printf("\nОБЩАЯ СТАТИСТИКА:\n")
	fmt.Printf("Всего сравнений: %d\n", len(results))
	fmt.Printf("Средняя схожесть: %.2f%%\n", averageSimilarity)
	fmt.Printf("Высокая схожесть (>80%%): %d проектов\n", highSim)
	fmt.Printf("Средняя схожесть (60-80%%): %d проектов\n", mediumSim)
	fmt.Printf("Низкая схожесть (<60%%): %d проектов\n", lowSim)

	// Создаем директорию для отчетов, если она не существует
	reportDir := "./reports"
	if err := os.MkdirAll(reportDir, os.ModePerm); err != nil {
		fmt.Printf("Ошибка при создании директории отчетов: %v\n", err)
		return
	}

	// Создаем файл отчета с новым именем в директории отчетов
	reportFileName := fmt.Sprintf("%s/plagiarism_report_%s.html", reportDir, time.Now().Format("2006-01-02_15-04-05"))

	// Обновляем HTML шаблон, добавляя отображение даты и времени
	const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Отчет о проверке на плагиат</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; border-radius: 5px; }
        .summary { background-color: #e9ecef; padding: 15px; margin: 20px 0; border-radius: 5px; }
        .results { margin: 20px 0; }
        table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        th, td { padding: 10px; border: 1px solid #dee2e6; text-align: left; }
        th { background-color: #f8f9fa; }
        .similarity-bar {
            background-color: #e9ecef;
            height: 20px;
            border-radius: 10px;
            overflow: hidden;
        }
        .similarity-fill {
            height: 100%;
            background-color: #28a745;
            transition: width 0.3s ease;
        }
        .high-similarity { background-color: #dc3545; }
        .medium-similarity { background-color: #ffc107; }
        .low-similarity { background-color: #28a745; }
        .datetime {
            font-size: 1.1em;
            color: #666;
            margin-bottom: 15px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Отчет о проверке на плагиат</h1>
        <div class="datetime">
            <p>Дата: {{.GeneratedDate}}</p>
            <p>Время: {{.GeneratedTime}}</p>
        </div>
    </div>

    <div class="summary">
        <h2>Общая статистика</h2>
        <p>Всего проверено сравнений: {{.TotalComparisons}}</p>
        <p>Средняя схожесть: {{printf "%.2f" .AverageSimilarity}}%</p>
        <p>Высокая схожесть (>80%): {{.HighSimilarityCount}} проектов</p>
        <p>Средняя схожесть (60-80%): {{.MediumSimilarityCount}} проектов</p>
        <p>Низкая схожесть (<60%): {{.LowSimilarityCount}} проектов</p>
    </div>

    <div class="results">
        <h2>Подробные результаты</h2>
        <table>
            <tr>
                <th>№</th>
                <th>Проекты</th>
                <th>Язык</th>
                <th>Общая схожесть</th>
                <th>Схожесть кода</th>
                <th>Схожесть комментариев</th>
                <th>Схожесть идентификаторов</th>
            </tr>
            {{range $i, $r := .Results}}
            <tr>
                <td>{{inc $i}}</td>
                <td>{{$r.Project1}} и {{$r.Project2}}</td>
                <td>{{$r.Language}}</td>
                <td>
                    <div class="similarity-bar">
                        <div class="similarity-fill {{similarityClass $r.Similarity}}"
                             style="width: {{$r.Similarity}}%"></div>
                    </div>
                    {{printf "%.2f" $r.Similarity}}%
                </td>
                <td>{{printf "%.2f" $r.Similarity}}%</td>
                <td>{{printf "%.2f" $r.CommentSimilarity}}%</td>
                <td>{{printf "%.2f" $r.IdentifierSimilarity}}%</td>
            </tr>
            {{end}}
        </table>
    </div>

    <div class="summary">
        <h2>Выводы</h2>
        {{if gt .HighSimilarityCount 0}}
        <p style="color: #dc3545">⚠️ Обнаружено {{.HighSimilarityCount}} случаев высокой схожести (более 80%)</p>
        {{end}}
        {{if gt .MediumSimilarityCount 0}}
        <p style="color: #ffc107">⚠️ Обнаружено {{.MediumSimilarityCount}} случаев средней схожести (60-80%)</p>
        {{end}}
        {{if gt .LowSimilarityCount 0}}
        <p style="color: #28a745">ℹ️ Обнаружено {{.LowSimilarityCount}} случаев низкой схожести (менее 60%)</p>
        {{end}}
    </div>
</body>
</html>
`

	// Создаем функции для шаблона
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"similarityClass": func(similarity float64) string {
			if similarity >= 80 {
				return "high-similarity"
			} else if similarity >= 60 {
				return "medium-similarity"
			}
			return "low-similarity"
		},
	}

	// Создаем и выполняем шаблон
	tmpl := template.Must(template.New("report").Funcs(funcMap).Parse(htmlTemplate))

	// Создаем переменную report типа HtmlReport
	report := HtmlReport{
		GeneratedTime:         time.Now().Format("15:04:05"),
		GeneratedDate:         time.Now().Format("2006-01-02"),
		TotalProjects:         len(results),
		TotalComparisons:      len(results),
		Results:               results,
		AverageSimilarity:     averageSimilarity,
		HighSimilarityCount:   highSim,
		MediumSimilarityCount: mediumSim,
		LowSimilarityCount:    lowSim,
	}

	// Создаем файл отчета с новым именем
	file, err := os.Create(reportFileName)
	if err != nil {
		fmt.Printf("Ошибка при создании файла отчета: %v\n", err)
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, report) // Исправлено: теперь передаем report
	if err != nil {
		fmt.Printf("Ошибка при генерации отчета: %v\n", err)
		return
	}

	fmt.Printf("\nПодробный отчет сохранен в файл: %s\n", reportFileName)
}

// Функция для создания паттерна управляющих конструкций
func createControlPattern(code string) string {
	var pattern []string
	lines := strings.Split(code, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "if") {
			pattern = append(pattern, "if")
		} else if strings.HasPrefix(line, "for") {
			pattern = append(pattern, "for")
		} else if strings.HasPrefix(line, "switch") {
			pattern = append(pattern, "switch")
		} else if strings.HasPrefix(line, "while") {
			pattern = append(pattern, "while")
		}
	}

	return strings.Join(pattern, "->")
}

// Функция для подсчета отступов
func countIndentation(line string) int {
	indent := 0
	for _, char := range line {
		if char == ' ' {
			indent++
		} else if char == '\t' {
			indent += 4
		} else {
			break
		}
	}
	return indent
}

// Функция для вычисления размера функции
func calculateFunctionSize(code, funcName string) int {
	lines := strings.Split(code, "\n")
	size := 0
	inFunction := false
	bracketCount := 0

	for _, line := range lines {
		if strings.Contains(line, "func "+funcName) {
			inFunction = true
		}

		if inFunction {
			size++
			bracketCount += strings.Count(line, "{")
			bracketCount -= strings.Count(line, "}")

			if bracketCount == 0 && size > 1 {
				break
			}
		}
	}

	return size
}

// Функция для анализа переносов строк
func analyzeLineBreaks(code string) string {
	lines := strings.Split(code, "\n")
	var pattern []string

	for i := 0; i < len(lines)-1; i++ {
		currentLine := strings.TrimSpace(lines[i])
		nextLine := strings.TrimSpace(lines[i+1])

		if currentLine == "" && nextLine == "" {
			pattern = append(pattern, "double")
		} else if currentLine == "" {
			pattern = append(pattern, "single")
		}
	}

	return strings.Join(pattern, ",")
}

// Функция для удаления комментариев
func removeComments(code, language string) string {
	switch language {
	case "python":
		lines := strings.Split(code, "\n")
		var filtered []string
		for _, line := range lines {
			if !strings.HasPrefix(strings.TrimSpace(line), "#") {
				filtered = append(filtered, line)
			}
		}
		return strings.Join(filtered, "\n")

	default:
		// Удаляем многострочные комментарии
		multilineRegex := regexp.MustCompile(`/\*[\s\S]*?\*/`)
		code = multilineRegex.ReplaceAllString(code, "")

		// Удаляем однострочные комментарии
		lines := strings.Split(code, "\n")
		var filtered []string
		for _, line := range lines {
			if !strings.Contains(line, "//") {
				filtered = append(filtered, line)
			}
		}
		return strings.Join(filtered, "\n")
	}
}

// Функция для сравнения функций
func compareFunctions(f1, f2 FunctionAnalysis) float64 {
	total := 0.0
	matches := 0.0

	// Сравниваем количество параметров
	for funcName, count1 := range f1.ParamCount {
		if count2, exists := f2.ParamCount[funcName]; exists && count1 == count2 {
			matches++
		}
		total++
	}

	// Сравниваем типы возвращаемых значений
	for funcName, type1 := range f1.ReturnTypes {
		if type2, exists := f2.ReturnTypes[funcName]; exists && type1 == type2 {
			matches++
		}
		total++
	}

	// Сравниваем размеры функций
	for funcName, size1 := range f1.FunctionSizes {
		if size2, exists := f2.FunctionSizes[funcName]; exists && size1 == size2 {
			matches++
		}
		total++
	}

	// Сравниваем порядок объявления
	if len(f1.DeclareOrder) == len(f2.DeclareOrder) {
		matching := true
		for i := range f1.DeclareOrder {
			if f1.DeclareOrder[i] != f2.DeclareOrder[i] {
				matching = false
				break
			}
		}
		if matching {
			matches++
		}
	}
	total++

	if total == 0 {
		return 0
	}

	return (matches / total) * 100
}

// Функция для сравнения импортов
func compareImports(i1, i2 ImportAnalysis) float64 {
	total := 3.0 // три критерия сравнения
	matches := 0.0

	// Сравниваем списки импортов
	commonImports := findCommonElements(i1.ImportList, i2.ImportList)
	if len(commonImports) > 0 {
		matches += float64(len(commonImports)) / float64(len(i1.ImportList))
	}

	// Сравниваем пордок импортов
	if i1.ImportOrder == i2.ImportOrder {
		matches++
	}

	// Сравниваем паттерны использования
	commonPatterns := 0
	for pattern, usage1 := range i1.UsagePatterns {
		if usage2, exists := i2.UsagePatterns[pattern]; exists && usage1 == usage2 {
			commonPatterns++
		}
	}
	if len(i1.UsagePatterns) > 0 {
		matches += float64(commonPatterns) / float64(len(i1.UsagePatterns))
	}

	return (matches / total) * 100
}

// Функция для сравнения форматирования кода
func compareFormatting(f1, f2 FormatAnalysis) float64 {
	total := 3.0 // три критерия сравнения
	matches := 0.0

	// Сравниваем стиль отступов
	if f1.IndentStyle == f2.IndentStyle {
		matches++
	}

	// Сравниваем паттерны пробелов
	if f1.SpacingPattern == f2.SpacingPattern {
		matches++
	}

	// Сравниваем паттерны переносов строк
	if f1.LineBreaks == f2.LineBreaks {
		matches++
	}

	return (matches / total) * 100
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	colorReset  = "\x1b[0m"
	colorRed    = "\x1b[31m"
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
	colorCyan   = "\x1b[36m"
)

type GameResult struct {
	Date     string `json:"date"`
	Result   string `json:"result"`
	Attempts int    `json:"attempts"`
}

func printMenu() {
	fmt.Println(`
Добро пожаловать в игру! Твоя задача угадать число.
Доступные операции:
1) Запустить игру
2) Показать историю всех раундов из файла
3) Выход`)
}

func input(scanner *bufio.Scanner, prompt string) (int, bool) {
	for {
		fmt.Print(prompt)
		if !scanner.Scan() {
			return 0, false
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		num, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println(colorRed + "Ошибка: введите корректное число!" + colorReset)
			continue
		}
		return num, true
	}
}

func generateSecretNumber(max int) int {
	return rand.Intn(max) + 1
}

func giveHint(usernumber int, secret int) {
	difference := usernumber - secret
	if difference < 0 {
		difference = -difference
	}

	if difference <= 5 {
		fmt.Println("🔥 горячо")
	} else if difference <= 15 {
		fmt.Println("🙂 тепло")
	} else {
		fmt.Println("❄️ холодно")
	}
}

func compareNumbers(usernumber int, secret int) bool {
	if usernumber == secret {
		fmt.Println(colorGreen + "Молодец! Ты выиграл.🙌" + colorReset)
		return true
	} else if usernumber > secret {
		fmt.Println("Секретное число меньше👇")
	} else {
		fmt.Println("Секретное число больше👆")
	}
	return false
}

func loadResults(filename string) []GameResult {
	file, err := os.ReadFile(filename)
	if err != nil {
		return []GameResult{}
	}

	var results []GameResult
	err = json.Unmarshal(file, &results)
	if err != nil {
		return []GameResult{}
	}

	return results
}

func saveResults(filename string, results []GameResult) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("Ошибка при подготовке данных для записи:", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
	}
}

func main() {
	const filename = "results.json"
	scanner := bufio.NewScanner(os.Stdin)
	historymain := loadResults(filename)

MainMenu:
	for {
		printMenu()
		choice, ok := input(scanner, "Выберите пункт меню: ")
		if !ok {
			fmt.Println("\nВвод завершен или прерван. Выход.")
			break
		}

		if choice == 3 {
			fmt.Println("Выход из игры.")
			break
		}

		switch choice {
		case 1:
			fmt.Println(`
Выберите уровень сложности:
1) Easy (числа от 1 до 50, 15 попыток)
2) Medium (числа от 1 до 100, 10 попыток)
3) Hard (числа от 1 до 200, 5 попыток)`)

			level, ok := input(scanner, "Ваш выбор: ")
			if !ok {
				fmt.Println("\nВвод прерван. Выход.")
				break MainMenu
			}

			var maxNumber int
			var maxAttempts int

			switch level {
			case 1:
				maxNumber = 50
				maxAttempts = 15
			case 2:
				maxNumber = 100
				maxAttempts = 10
			case 3:
				maxNumber = 200
				maxAttempts = 5
			default:
				fmt.Println(colorRed + "Ошибка: неверный уровень сложности! Возврат в главное меню." + colorReset)
				continue MainMenu
			}

			history := []int{}
			n := generateSecretNumber(maxNumber)
			fmt.Printf("\nИгра 'Угадай число' - от 1 до %d началась!\n", maxNumber)
			fmt.Printf("Угадайте число за %d попыток!\n", maxAttempts)

			gameOutcome := "проигрыш"
			usedAttempts := maxAttempts

			for i := 1; i <= maxAttempts; {
				prompt := fmt.Sprintf("%sПопытка #%d%s - Введи число: ", colorYellow, i, colorReset)
				usernumber, ok := input(scanner, prompt)
				if !ok {
					fmt.Println("\nВвод прерван. Выход из игры.")
					break MainMenu
				}

				history = append(history, usernumber)
				fmt.Printf("%sПрошлые попытки раунда: %v%s\n", colorCyan, history, colorReset)

				if compareNumbers(usernumber, n) {
					gameOutcome = "победа"
					usedAttempts = i
					break
				}

				giveHint(usernumber, n)

				if usernumber > n {
					fmt.Println("Секретное число меньше, чем твое")
				} else {
					fmt.Println("Секретное число больше, чем твое")
				}

				if i == maxAttempts {
					fmt.Printf("%sВы проиграли!😢 Секретное число было: %d%s\n", colorRed, n, colorReset)
				}
				i++
			}

			newResult := GameResult{
				Date:     time.Now().Format("2006-01-02 15:04:05"),
				Result:   gameOutcome,
				Attempts: usedAttempts,
			}
			historymain = append(historymain, newResult)
			saveResults(filename, historymain)

			fmt.Println("\nИгра закончена!")

		case 2:
			if len(historymain) == 0 {
				fmt.Println("История пока пуста. Сыграйте хотя бы один раунд!")
			} else {
				fmt.Println("\nИстория всех раундов из файла:")
				for idx, res := range historymain {
					resColor := colorRed
					if res.Result == "победа" {
						resColor = colorGreen
					}
					fmt.Printf("Раунд %d | Дата: %s | Исход: %s%s%s | Попыток: %d\n",
						idx+1, res.Date, resColor, res.Result, colorReset, res.Attempts)
				}
			}
		default:
			fmt.Println(colorRed + "Неверный пункт меню!" + colorReset)
		}
	}
}

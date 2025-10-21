package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var countID int
var expenseMap map[int]ExpenseInfo

type ExpenseInfo struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
}

func writeToAJSONFile() {
	newData, _ := json.MarshalIndent(expenseMap, "", "	")
	os.WriteFile("expense.json", newData, 0644)
}
func readFromAJSONFile() {
	data, err := os.ReadFile("expense.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(data, &expenseMap)
}

func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func executeAddCommand(args []string) {
	/*
		args[0] = add
		args[1] = --description
		args[2] = "Dinner"
		args[3] = --amount
		args[4] = 10
	*/

	/*
		add --description 								"Lunch" --amount 20
		add --description "Buy gifts for Kim Anh" --amount 	6
	*/

	if len(args) != 5 {
		fmt.Println("Command isn't valid (The number of parameters is incorrect)")
		return
	}

	description := strings.ReplaceAll(args[2], "\"", "")
	amount, err := strconv.ParseFloat(args[4], 64)
	if args[1] != "--description" || args[3] != "--amount" || err != nil {
		fmt.Println("Command isn't valid (The attributes are incorrect)")
		return
	}

	countID++
	expenseMap[countID] = ExpenseInfo{
		ID:          countID,
		Date:        time.Now(),
		Description: description,
		Amount:      amount,
	}

	result := "Expense added successfully (ID: " + strconv.Itoa(countID) + ")"
	fmt.Printf("%v\n", result)
}

func executeListCommand() {
	fmt.Printf("%-3s %-12s %-26s %s\n", "ID", "Date", "Description", "Amount")
	for _, expense := range expenseMap {
		fmt.Printf("%-3v %-12v %-26v 💲%v\n", expense.ID, expense.Date.Format("2006-01-02"), expense.Description, expense.Amount)
	}
}

func executeSummaryCommand(args []string) {
	/*
		summary
			+ args[0] = summary

		summary --month 8
			+ args[0] = summary
			+ args[1] = --month
			+ args[2] = 8
	*/

	summary := 0.0

	if len(args) == 1 {
		for _, expense := range expenseMap {
			summary += expense.Amount
		}
		fmt.Printf("Total expenses: 💲%v\n", summary)
		return
	}

	if len(args) == 3 {
		month, err := strconv.Atoi(args[2])
		if args[1] != "--month" || err != nil || month < 1 || 12 < month {
			fmt.Println("The command isn't valid (The attribute is incorrect)")
			return
		}

		for _, expense := range expenseMap {
			if expense.Date.Month() == time.Month(month) {
				summary += expense.Amount
			}
		}

		fmt.Printf("Total expense for %v: 💲%v\n", time.Month(month), summary)
		return
	}

	fmt.Println("The command isn't valid")
}

func executeDeleteCommand(args []string) {
	/*
		delete --id 2
			+ args[0] = delete
			+ args[1] = --id
			+ args[2] = 2
	*/

	if len(args) != 3 {
		fmt.Println("Command isn't valid (The number of parameters is incorrect)")
		return
	}
	if args[1] != "--id" {
		fmt.Println("Command isn't valid (Must be --id)")
	}

	deletedID, err := strconv.Atoi(args[2])
	if err != nil || deletedID < 0 {
		fmt.Println("Command isn't valid (The deleted id isn't valid)")
		return
	}

	_, exists := expenseMap[deletedID]
	if !exists {
		fmt.Println("The deleted id doesn't exist")
		return
	}

	delete(expenseMap, deletedID)
	fmt.Println("✖️ Expense deleted successfully")
}

func separateField(input string) []string {
	input += " "

	field := ""
	args := []string{}
	runeInput := []rune(input)
	lenOfInput := len(runeInput)

	for i := 0; i < lenOfInput; i++ {
		character_i := string(runeInput[i])

		if character_i == "\"" {
			for j := i + 1; j < lenOfInput; j++ {
				character_j := string(runeInput[j])
				if character_j == "\"" {
					i = j + 1
					break
				}
				field += character_j
			}

			args = append(args, field)
			field = ""
			continue
		}

		if character_i == " " {
			args = append(args, field)
			field = ""
			continue
		}

		field += character_i
	}

	return args
}

func headerOfCLI() {
	fmt.Println("Expense Tracker CLI")
	fmt.Println("--------------------------")
}

func setUp() {
	expenseMap = make(map[int]ExpenseInfo)
	readFromAJSONFile()

	for _, expense := range expenseMap {
		countID = max(countID, expense.ID)
	}
}

func main() {
	setUp()

	reader := bufio.NewReader(os.Stdin)
	headerOfCLI()

	running := true
	for running {
		fmt.Print("expense-tracker> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		args := separateField(input)

		switch {
		case input == "":
			continue

		case input == "exit" || input == "quit":
			fmt.Println("Bye bye!")
			writeToAJSONFile()
			running = false

		case input == "clear":
			clearScreen()
			headerOfCLI()

		case input != "" && args[0] == "add":
			executeAddCommand(args)

		case input != "" && args[0] == "list":
			executeListCommand()

		case input != "" && args[0] == "summary":
			executeSummaryCommand(args)

		case input != "" && args[0] == "delete":
			executeDeleteCommand(args)

		default:
			fmt.Println("Command isn't valid")
		}
	}
}

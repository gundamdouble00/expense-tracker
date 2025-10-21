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

var mapFlag bool
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

	if len(args) != 5 {
		fmt.Println("Command isn't valid")
		return
	}

	description := strings.ReplaceAll(args[2], "\"", "")
	amount, err := strconv.ParseFloat(args[4], 64)
	if args[1] != "--description" || args[3] != "--amount" || err != nil {
		fmt.Println("Command isn't valid")
		return
	}

	mapFlag = true
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
	fmt.Printf("%-3s %-12s %-12s %s\n", "ID", "Date", "Description", "Amount")
	for _, expense := range expenseMap {
		fmt.Printf("%-3v %-12v %-12v %v\n", expense.ID, expense.Date.Format("2006-01-02"), expense.Description, expense.Amount)
	}

	if mapFlag {
		writeToAJSONFile()
		mapFlag = false
	}
}

func headerOfCLI() {
	fmt.Println("Expense Tracker CLI")
	fmt.Println("--------------------------")
}

func setUp() {
	expenseMap = make(map[int]ExpenseInfo)
	readFromAJSONFile()
	mapFlag = false

	for _, expense := range expenseMap {
		countID = max(countID, expense.ID)
	}
}

func main() {
	// fmt.Printf("%v\n\n", exec.Command("clear"))
	setUp()

	reader := bufio.NewReader(os.Stdin)
	headerOfCLI()

	running := true
	for running {
		fmt.Print("expense-tracker> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		args := strings.Fields(input)

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

		default:
			for i, parameter := range args {
				fmt.Printf("%v. %v\n", i, parameter)
			}
		}
	}
}

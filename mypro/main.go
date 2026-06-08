package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Calculator")

	// Entry widget to display calculations
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter expression here")
	entry.Editable = false

	// Function to handle button click events
	handleClick := func(button string) {
		switch button {
		case "C":
			entry.SetText("")
		default:
			entry.AppendText(button)
		}
	}

	// Button labels
	buttonLabels := []string{"7", "8", "9", "/", "4", "5", "6", "*", "1", "2", "3", "-", "0", ".", "=", "+"}

	// Create buttons and add them to the grid
	grid := container.NewGridWithColumns(4)
	for _, label := range buttonLabels {
		label := label
		button := widget.NewButton(label, func() {
			if label == "=" {
				result, err := evaluateExpression(entry.Text())
				if err != nil {
					entry.SetText(err.Error())
				} else {
					entry.SetText(fmt.Sprintf("%.2f", result))
				}
			} else {
				handleClick(label)
			}
		})
		grid.Add(button)
	}

	// Layout
	w.SetContent(container.NewVBox(
		entry,
		grid,
	))

	w.Resize(fyne.NewSize(300, 400))
	w.ShowAndRun()
}

func evaluateExpression(expression string) (float64, error) {
	type token struct {
		isOperator bool
		value      string
	}

	var tokens []token
	currentNumber := ""

	for _, r := range expression {
		switch {
		case r == ' ':
			continue
		case (r >= '0' && r <= '9') || r == '.':
			currentNumber += string(r)
		case strings.ContainsRune("+-*/", r):
			if currentNumber == "" {
				return 0, fmt.Errorf("invalid expression")
			}
			tokens = append(tokens, token{false, currentNumber})
			currentNumber = ""
			tokens = append(tokens, token{true, string(r)})
		default:
			return 0, fmt.Errorf("invalid character: %c", r)
		}
	}

	if currentNumber != "" {
		tokens = append(tokens, token{false, currentNumber})
	}

	if len(tokens)%2 == 0 || len(tokens) < 3 {
		return 0, fmt.Errorf("invalid expression")
	}

	values := []float64{}
	operators := []string{}
	prec := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}

	applyOperator := func() error {
		if len(values) < 2 || len(operators) == 0 {
			return fmt.Errorf("invalid expression")
		}

		right := values[len(values)-1]
		left := values[len(values)-2]
		values = values[:len(values)-2]
		op := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		switch op {
		case "+":
			values = append(values, left+right)
		case "-":
			values = append(values, left-right)
		case "*":
			values = append(values, left*right)
		case "/":
			if right == 0 {
				return fmt.Errorf("division by zero")
			}
			values = append(values, left/right)
		default:
			return fmt.Errorf("unsupported operator: %s", op)
		}
		return nil
	}

	for _, tok := range tokens {
		if !tok.isOperator {
			num, err := strconv.ParseFloat(tok.value, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %s", tok.value)
			}
			values = append(values, num)
			continue
		}

		for len(operators) > 0 && prec[operators[len(operators)-1]] >= prec[tok.value] {
			if err := applyOperator(); err != nil {
				return 0, err
			}
		}
		operators = append(operators, tok.value)
	}

	for len(operators) > 0 {
		if err := applyOperator(); err != nil {
			return 0, err
		}
	}

	if len(values) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	return values[0], nil
}

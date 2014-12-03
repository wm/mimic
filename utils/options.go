package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wm/mimic/application"
)

// PromptUserToSelectFrom prompts the user to select one or more items,
// returning a slice of those selected. If 0 is selected it returns all items.
//
// Each item it displayed on the screen with a number associated with it.
func PromptUserToSelectFrom(apps []application.App) (selectedApps []application.App, err error) {
	displayOptions(apps)
	options, err := readOptionsSelected(len(apps) + 1)
	if err != nil {
		return nil, err
	}

	found := make(map[int64]bool)
	for _, option := range options {
		if option == 0 { // All
			return apps, nil
		}
		if !found[option] { // We only want each app once
			selectedApps = append(selectedApps, apps[option-1])
			found[option] = true
		}
	}

	return selectedApps, nil
}

func displayOptions(apps []application.App) {
	fmt.Println("Please select which applications you'd like to replace")
	fmt.Println()
	fmt.Println("    0: all")

	for i, app := range apps {
		fmt.Printf("    %v: %s\n", i+1, app)
	}

	fmt.Println("")
	fmt.Println("e.g. 1,2")
	fmt.Println("")
}

func readOptionsSelected(length int) (options []int64, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	numbers := strings.Split(scanner.Text(), ",")
	options = make([]int64, len(numbers))

	for i, number := range numbers {
		num, err := strconv.ParseInt(strings.Trim(number, " "), 0, 64)
		if err != nil {
			return nil, err
		}

		if num > int64(length-1) {
			return nil, fmt.Errorf("Invalid option")
		}

		options[i] = num
	}

	return options, nil
}

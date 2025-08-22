package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadStringInto(prompt string, dest *string) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	*dest = strings.TrimSpace(input)
}

func ReadStringSliceInto(prompt string, dest *[]string) {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	for _, value := range strings.Split(strings.TrimSpace(input), ",") {
		value = strings.TrimSpace(value)
		// Skip empty values
		if value == "" {
			continue
		}
		// Trim spaces and add to the destination slice
		*dest = append(*dest, value)
	}
}

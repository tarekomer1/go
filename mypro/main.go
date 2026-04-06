package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var result []string
	content, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	result = strings.Fields(string(content))
	fmt.Println(result)
	for i := 0; i < len(result); i++ {

		if result[i] == "(hex)" && i > 0 {
			result[i-1] = hexToDecimal(result[i-1])
			result = append(result[:i], result[i+1:]...)
			i--
		} else if result[i] == "(bin)" && i > 0 {
			result[i-1] = binToDecimal(result[i-1])
			result = append(result[:i], result[i+1:]...)
			i--
		} else if strings.HasPrefix(result[i], "(up") {
			result, i = applyTransform(result, i, toUpper)
		} else if strings.HasPrefix(result[i], "(low") {
			result, i = applyTransform(result, i, toLower)
		} else if strings.HasPrefix(result[i], "(cap") {
			result, i = applyTransform(result, i, capitalize)
		}
	}
	//result = string(content)
	fmt.Println(result)
}

func hexToDecimal(s string) string {
	if val, err := strconv.ParseInt(s, 16, 64); err == nil {
		return strconv.Itoa(int(val))
	}
	return s
}

func binToDecimal(s string) string {
	if val, err := strconv.ParseInt(s, 2, 64); err == nil {
		return strconv.Itoa(int(val))
	}
	return s
}

func toUpper(s string) string {
	return strings.ToUpper(s)
}

func toLower(s string) string {
	return strings.ToLower(s)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func applyTransform(result []string, i int, transform func(string) string) ([]string, int) {
	marker := result[i]
	parts := strings.Split(strings.Trim(marker, "()"), ",")
	num := 1
	if len(parts) > 1 {
		if n, err := strconv.Atoi(parts[1]); err == nil {
			num = n
		}
	}
	for j := 1; j <= num && i-j >= 0; j++ {
		result[i-j] = transform(result[i-j])
	}
	result = append(result[:i], result[i+1:]...)
	return result, i - 1
}

package main

import (
	//"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		return
	}

	content, _ := os.ReadFile(os.Args[1])
	words := splitWords(string(content))

	// 1. Process Modifiers (hex, bin, up, low, cap)
	for i := 0; i < len(words); i++ {
		w := words[i]
		if w == "(hex)" && i > 0 {
			words[i-1] = hexToDec(words[i-1])
			words = remove(words, i)
			i--
		} else if w == "(bin)" && i > 0 {
			words[i-1] = binToDec(words[i-1])
			words = remove(words, i)
			i--
		} else if isTag(w, "(up") {
			words, i = applyTransform(words, i, toUpper)
		} else if isTag(w, "(low") {
			words, i = applyTransform(words, i, toLower)
		} else if isTag(w, "(cap") {
			words, i = applyTransform(words, i, capitalize)
		}
	}

	// 2. Grammar: A to An
	for i := 0; i < len(words)-1; i++ {
		if (words[i] == "a" || words[i] == "A") && isVowelOrH(words[i+1][0]) {
			if words[i] == "a" {
				words[i] = "an"
			} else {
				words[i] = "An"
			}
		}
	}

	// 3. Reconstruct and Format Punctuation
	raw := join(words)
	result := formatPunctuation(raw)

	os.WriteFile(os.Args[2], []byte(result), 0644)
}

// --- Helper Functions (Replacing restricted packages) ---

func hexToDec(s string) string {
	var res int64
	for _, c := range s {
		res *= 16
		if c >= '0' && c <= '9' {
			res += int64(c - '0')
		} else if c >= 'a' && c <= 'f' {
			res += int64(c - 'a' + 10)
		} else if c >= 'A' && c <= 'F' {
			res += int64(c - 'A' + 10)
		}
	}
	return intToStr(res)
}

func binToDec(s string) string {
	var res int64
	for _, c := range s {
		res *= 2
		if c == '1' {
			res += 1
		}
	}
	return intToStr(res)
}

func intToStr(n int64) string {
	if n == 0 {
		return "0"
	}
	var res []byte
	for n > 0 {
		res = append([]byte{byte(n%10 + '0')}, res...)
		n /= 10
	}
	return string(res)
}

func applyTransform(words []string, i int, fn func(string) string) ([]string, int) {
	n := 1
	tag := words[i]
	if tag[len(tag)-1] != ')' { // case: (up, 2)
		numStr := ""
		for _, c := range words[i+1] {
			if c >= '0' && c <= '9' {
				numStr += string(c)
			}
		}
		n = int(strToInt(numStr))
		words = remove(words, i+1)
	}
	for j := 1; j <= n && i-j >= 0; j++ {
		words[i-j] = fn(words[i-j])
	}
	return remove(words, i), i - 1
}

func formatPunctuation(s string) string {
	var res []byte
	quoteOpen := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		// Remove space before punctuation
		if c == ' ' && i+1 < len(s) && isPunc(s[i+1]) {
			continue
		}
		// Handle single quotes
		if c == '\'' {
			if !quoteOpen { // Opening
				if len(res) > 0 && res[len(res)-1] != ' ' {
					res = append(res, ' ')
				}
				res = append(res, '\'')
				for i+1 < len(s) && s[i+1] == ' ' {
					i++
				}
				quoteOpen = true
				continue
			} else { // Closing
				for len(res) > 0 && res[len(res)-1] == ' ' {
					res = res[:len(res)-1]
				}
				res = append(res, '\'')
				quoteOpen = false
				continue
			}
		}
		res = append(res, c)
		// Ensure space after punctuation
		if isPunc(c) && i+1 < len(s) && s[i+1] != ' ' && !isPunc(s[i+1]) {
			res = append(res, ' ')
		}
	}
	return string(res)
}

// --- Basic Utilities ---

func splitWords(s string) []string {
	var words []string
	var cur []byte
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\n' || s[i] == '\t' {
			if len(cur) > 0 {
				words = append(words, string(cur))
				cur = nil
			}
		} else {
			cur = append(cur, s[i])
		}
	}
	if len(cur) > 0 {
		words = append(words, string(cur))
	}
	return words
}

func toUpper(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'a' && b[i] <= 'z' {
			b[i] -= 32
		}
	}
	return string(b)
}

func toLower(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 32
		}
	}
	return string(b)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	b := []byte(toLower(s))
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 32
	}
	return string(b)
}

func isVowelOrH(c byte) bool {
	l := c
	if l >= 'A' && l <= 'Z' {
		l += 32
	}
	return l == 'a' || l == 'e' || l == 'i' || l == 'o' || l == 'u' || l == 'h'
}

func isPunc(c byte) bool {
	return c == '.' || c == ',' || c == '!' || c == '?' || c == ':' || c == ';'
}

func isTag(w, prefix string) bool {
	if len(w) < len(prefix) {
		return false
	}
	return w[:len(prefix)] == prefix
}

func remove(slice []string, i int) []string {
	return append(slice[:i], slice[i+1:]...)
}

func join(slice []string) string {
	res := ""
	for i, s := range slice {
		res += s
		if i < len(slice)-1 {
			res += " "
		}
	}
	return res
}

func strToInt(s string) int64 {
	var res int64
	for _, c := range s {
		res = res*10 + int64(c-'0')
	}
	return res
}

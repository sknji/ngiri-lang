package lexer

import "io/ioutil"

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && 'Z' >= ch || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func getFileContent(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}

	return string(b)
}

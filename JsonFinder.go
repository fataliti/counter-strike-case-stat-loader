package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetJsonString(json_key string, document *goquery.Document) string {
	var scriptContent string
	document.Find("script").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, json_key) {
			scriptContent = text
		}
	})

	jsonString := ""
	if scriptContent != "" {
		balance := 0
		in_string := false
		escape := false

		start := strings.Index(scriptContent, json_key+" = ")
		start_index := start + len(json_key+" = ")
		end_index := start_index

		for end_index < len(scriptContent) {
			char := scriptContent[end_index]
			switch char {
			case '"':
				if !escape {
					in_string = !in_string
				}
				escape = false
			case '\\':
				escape = !escape
			case '}', ']':
				if in_string {
					balance -= 1
				}
			case '{', '[':
				if in_string {
					balance += 1
				}
			}

			if balance == 0 && !in_string {
				next_char := scriptContent[end_index+1]
				if next_char == ';' {
					jsonString = scriptContent[start_index-1 : end_index+1]
					break
				}
			}

			if escape && char != '\\' {
				escape = false
			}

			end_index += 1
		}
	}

	return jsonString
}

func FinsString(key string, doc *goquery.Document) string {
	var finded_string = GetJsonString(key, doc)
	if finded_string != "" {
		return finded_string[2 : len(finded_string)-1]
	}
	return ""
}

func GetFixedJsonString(content []byte) string {
	var fixed_json strings.Builder
	for _, r := range content {
		if r != '\n' {
			fixed_json.WriteByte(r)
		}
	}

	return fixed_json.String()
}

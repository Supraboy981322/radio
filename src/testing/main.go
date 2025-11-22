package main

import (
	"fmt"
)

var (
	library = map[string]string{
		"foo": "bar",
		"baz": "quz",
	}
)

func main() {
	for _, line := range buildJSONlibrary() {
		fmt.Println(line)
	}
}

func buildJSONlibrary() []string {
	final := []string{"["}
	for key, value := range library {
		final = append(final, "  [")
		final = append(final, fmt.Sprintf("    \"%s\",", key))
		final = append(final, fmt.Sprintf("    \"%s\"", value))
		final = append(final, "  ],")
	}
	//remove last char from last line (a comma)
	final[len(final)-1] = final[len(final)-1][:len(final[len(final)-1])-1]
	final = append(final, "]")
	return final
}

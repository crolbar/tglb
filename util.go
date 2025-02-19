package main

import (
	"strings"

	"github.com/crolbar/lipbalm"
)

func getAsAscii(n string) string {
	ascii := ""
	for _, c := range n {
		// fmt.Println(c - '0')
		num := c - '0'

		if num == -2 {
			ascii = lipbalm.JoinHorizontal(lipbalm.Left,
				ascii,
				digits[10],
			)
			continue
		}
		ascii = lipbalm.JoinHorizontal(lipbalm.Left,
			ascii,
			digits[num],
		)
	}
	return ascii
}

func getMaxWidth(s string) int {
	lines := strings.Split(s, "\n")
	maxWidth := 0
	for _, l := range lines {
		if len(l) > maxWidth {
			maxWidth = len(l)
		}
	}

	return maxWidth
}

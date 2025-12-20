package main

import (
	"fmt"
	"strings"
)

type ParsingSupplierBalance struct {
	MemberId string
	Pin      string
	Password string
}

func main() {
	parsing := fmt.Sprintf("%s", "memberId=peci28727|user=woyi64780|pin=683907|password=66d4a904d01336.39135")

	parts := strings.Split(parsing, "|")

	var result ParsingSupplierBalance
	for _, v := range parts {
		parts := strings.Split(v, "=")
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			switch key {
			case "memberId":
				result.MemberId = value
			case "pin":
				result.Pin = value
			case "password":
				result.Password = value
			}
		}
	}

	fmt.Println(result)
}

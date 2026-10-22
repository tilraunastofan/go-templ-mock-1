package pages

import (
	_ "embed"
	"encoding/json"
	"strings"
)

// Person is one record in the mock search dataset (search_data.json).
type Person struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
	City       string `json:"city"`
	Status     string `json:"status"`
}

//go:embed search_data.json
var searchDataJSON []byte

var people = mustLoadPeople()

func mustLoadPeople() []Person {
	var ps []Person
	if err := json.Unmarshal(searchDataJSON, &ps); err != nil {
		panic("search_data.json: " + err.Error())
	}
	return ps
}

const searchLimit = 10

// SearchPeople returns up to searchLimit people matching q, plus the total
// number of matches. Every whitespace-separated term must appear
// (case-insensitively) in some field; an empty query matches everyone.
func SearchPeople(q string) (results []Person, total int) {
	terms := strings.Fields(strings.ToLower(q))
	for _, p := range people {
		hay := strings.ToLower(strings.Join([]string{p.Name, p.Email, p.Role, p.Department, p.City, p.Status}, " "))
		match := true
		for _, t := range terms {
			if !strings.Contains(hay, t) {
				match = false
				break
			}
		}
		if !match {
			continue
		}
		total++
		if len(results) < searchLimit {
			results = append(results, p)
		}
	}
	return results, total
}

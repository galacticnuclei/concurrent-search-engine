package robots

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Rules struct {
	Disallow []string
}

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func Fetch(host string) (*Rules, error) {
	resp, err := client.Get("https://" + host + "/robots.txt")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rules := &Rules{}

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "Disallow:") {
			path := strings.TrimSpace(
				strings.TrimPrefix(line, "Disallow:"),
			)

			if path != "" {
				rules.Disallow = append(
					rules.Disallow,
					path,
				)
			}
		}
	}

	return rules, nil
}

func (r *Rules) Allowed(path string) bool {
	for _, disallowed := range r.Disallow {
		if strings.HasPrefix(path, disallowed) {
			return false
		}
	}

	return true
}

func (r *Rules) Print() {
	fmt.Println("Disallowed Paths:")

	for _, path := range r.Disallow {
		fmt.Println(path)
	}
}

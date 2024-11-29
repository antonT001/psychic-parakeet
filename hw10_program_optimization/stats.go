package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)
	for scanner.Scan() {
		email := fastjson.GetString(scanner.Bytes(), "Email")
		if ok := strings.HasSuffix(strings.ToLower(email), domain); ok {
			name := strings.ToLower(strings.SplitN(email, "@", 2)[1])
			result[name]++
		}
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}

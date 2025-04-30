package hw10programoptimization

import (
	"bufio"
	"fmt"
	"github.com/mailru/easyjson"
	"io"
	"strings"
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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	result = make(users, 0, 100000)
	var user User
	for scanner.Scan() {
		if err = easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return
		}
		result = append(result, user)
	}

	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched := strings.HasSuffix(user.Email, "."+domain)

		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}

package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ParseSipAddress1(str string) (isSip bool, data map[string]interface{}) {
	str = strings.ReplaceAll(str, " ", "")
	data = make(map[string]interface{})

	prefix := ""
	if strings.HasPrefix(str, "sip:") {
		prefix = "sip:"
		str = str[len(prefix):]
	} else if strings.HasPrefix(str, "sips:") {
		prefix = "sips:"
		str = str[len(prefix):]
	} else {
		isSip = false
		return
	}

	parts := strings.FieldsFunc(str, func(r rune) bool {
		return r == '?' || r == ';' || r == '&'
	})

	leftPart := strings.Split(parts[0], "@")

	userAndPass := strings.Split(leftPart[0], ":")
	data["displayName"] = userAndPass[0]
	data["userName"] = userAndPass[1]

	if len(leftPart) > 1 {
		domainAndPort := strings.Split(leftPart[1], ":")
		if net.ParseIP(domainAndPort[0]) != nil {
			data["ip"] = domainAndPort[0]
		} else {
			data["domain"] = domainAndPort[0]
		}
		if len(domainAndPort) > 1 {
			port, err := strconv.Atoi(domainAndPort[1])
			if err == nil {
				data["port"] = port
			}
		}
	}

	if len(parts) > 1 {
		params := make(map[string]string, len(parts)-1)
		data["params"] = params
		sipParameters := parts[1:]
		for _, sipParameter := range sipParameters {
			keyValue := strings.SplitN(sipParameter, "=", 2)
			if len(keyValue) > 1 {
				params[keyValue[0]] = keyValue[1]
			} else {
				params[keyValue[0]] = ""
			}
		}
	}

	isSip = true
	return
}

func main() {
	str := "LittleGuy<sips:admin@10.0.0.3:5061>;tag=123"
	isSip, data := ParseSipAddress1(str)
	fmt.Println("isSip:", isSip)
	fmt.Println("data:", data)
}

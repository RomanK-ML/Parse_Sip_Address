package main

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	UnitTest()
	duration := time.Since(start)
	fmt.Printf("Время выполнения: %d мс\n", duration.Microseconds())
}

// ParseSipAddress разбирает SIP-адрес и возвращает флаг успешного разбора и данные адреса
func ParseSipAddress(str string) (isSip bool, data map[string]interface{}) {
	// Создаем пустой map для хранения разобранных значений
	data = make(map[string]interface{})

	// Удаление всех пробелов из строки
	str = strings.ReplaceAll(str, " ", "")

	// Проверяем первый символ строки на наличие <
	if strings.HasPrefix(str, "<") {
		str = strings.ReplaceAll(str[1:], ">", "")
	} else {
		displayNameParts := strings.Split(str, "<")
		if len(displayNameParts) > 1 {
			data["displayName"] = displayNameParts[0]
			str = strings.ReplaceAll(displayNameParts[1], ">", "")
		}
	}

	prefix := ""
	// Проверяем префикс "sip:"
	if strings.HasPrefix(str, "sip:") {
		// Удаление префикса "sip:"
		prefix = "sip:"
		str = str[len(prefix):]
		// Проверяем префикс "sips:"
	} else if strings.HasPrefix(str, "sips:") {
		// Удаление префикса "sips:"
		prefix = "sips:"
		str = str[len(prefix):]
	} else {
		isSip = false
		return
	}

	// Разделение строки на подстроки с использованием символов "?", ";", и "&"
	parts := strings.FieldsFunc(str, func(r rune) bool {
		return r == '?' || r == ';' || r == '&'
	})

	// Разбор левой части
	leftPart := strings.Split(parts[0], "@")

	// Разбираем имя пользователя и пароль
	userAndPass := strings.Split(leftPart[0], ":")
	data["userName"] = userAndPass[0]
	if len(userAndPass) > 1 {
		data["userPassword"] = userAndPass[1]
	}

	// Разбираем домен и порт
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

	// Разбираем параметры
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

	// Установка флага isSip в true, чтобы указать успешный разбор адреса
	isSip = true
	return
}

// UnitTest Функция для тестирования ParseSipAddress
func UnitTest() {

	testData := []struct {
		str         string
		expectedSip bool
		expected    map[string]interface{}
	}{
		{
			str:         "sip:alice@example.com:5060;transport=udp",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "alice",
				"domain":   "example.com",
				"port":     5060,
				"params": map[string]string{
					"transport": "udp",
				},
			},
		},
		{
			str:         "sip:bob@example.com",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "bob",
				"domain":   "example.com",
			},
		},
		{
			str:         "sip:user123:pass456@domain.com:5060?param1=value1",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName":     "user123",
				"userPassword": "pass456",
				"domain":       "domain.com",
				"port":         5060,
				"params": map[string]string{
					"param1": "value1",
				},
			},
		},
		{
			str:         "sips:admin@10.0.0.1",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "admin",
				"ip":       "10.0.0.1",
			},
		},
		{
			str:         "sip:carol@example.com:5080",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "carol",
				"domain":   "example.com",
				"port":     5080,
			},
		},
		{
			str:         "invalid_sip_address",
			expectedSip: false,
			expected:    map[string]interface{}{},
		},
		{
			str:         "sip:user1@domain.com",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "user1",
				"domain":   "domain.com",
			},
		},
		{
			str:         "sip:user2:pass@192.168.0.1",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName":     "user2",
				"userPassword": "pass",
				"ip":           "192.168.0.1",
			},
		},
		{
			str:         "sips:user3@domain.org:1234;param1=value1;param2=value2",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "user3",
				"domain":   "domain.org",
				"port":     1234,
				"params": map[string]string{
					"param1": "value1",
					"param2": "value2",
				},
			},
		},
		{
			str:         "sip:john@example.com",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "john",
				"domain":   "example.com",
			},
		},
		{
			str:         "sip:jane@192.168.1.100:5062;param1=value1",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "jane",
				"ip":       "192.168.1.100",
				"port":     5062,
				"params": map[string]string{
					"param1": "value1",
				},
			},
		},
		{
			str:         "sips:guest@example.org",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "guest",
				"domain":   "example.org",
			},
		},
		{
			str:         "sip:user4:password@10.0.0.2:5080?param1=value1&param2=value2",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName":     "user4",
				"userPassword": "password",
				"ip":           "10.0.0.2",
				"port":         5080,
				"params": map[string]string{
					"param1": "value1",
					"param2": "value2",
				},
			},
		},
		{
			str:         "sip:test@example.com:5060;param1=value1;param2=value2",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "test",
				"domain":   "example.com",
				"port":     5060,
				"params": map[string]string{
					"param1": "value1",
					"param2": "value2",
				},
			},
		},
		{
			str:         "sips:user5@192.0.2.1",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "user5",
				"ip":       "192.0.2.1",
			},
		},
		{
			str:         "sips:user5@192.0.2.1?param1=value",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "user5",
				"ip":       "192.0.2.1",
				"params": map[string]string{
					"param1": "value",
				},
			},
		},
		{
			str:         "sip:user6:pass@domain.net:6000",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName":     "user6",
				"userPassword": "pass",
				"domain":       "domain.net",
				"port":         6000,
			},
		},
		{
			str:         "sips:admin@10.0.0.3:5061",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "admin",
				"ip":       "10.0.0.3",
				"port":     5061,
			},
		},
		{
			str:         "LittleGuy<sips:admin@10.0.0.3:5061>;tag=123",
			expectedSip: true,
			expected: map[string]interface{}{
				"displayName": "LittleGuy",
				"userName":    "admin",
				"ip":          "10.0.0.3",
				"port":        5061,
				"params": map[string]string{
					"tag": "123",
				},
			},
		},
		{
			str:         "LittleGuy<sips:admin@10.0.0.3:5061;tag=123>",
			expectedSip: true,
			expected: map[string]interface{}{
				"displayName": "LittleGuy",
				"userName":    "admin",
				"ip":          "10.0.0.3",
				"port":        5061,
				"params": map[string]string{
					"tag": "123",
				},
			},
		},
		{
			str:         "<sips:admin@10.0.0.3:5061>;tag=123",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "admin",
				"ip":       "10.0.0.3",
				"port":     5061,
				"params": map[string]string{
					"tag": "123",
				},
			},
		},
		{
			str:         "<sips:admin@10.0.0.3:5061?tag=123>",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "admin",
				"ip":       "10.0.0.3",
				"port":     5061,
				"params": map[string]string{
					"tag": "123",
				},
			},
		},
		{
			str:         "sip:user7@example.com;param1=value1",
			expectedSip: true,
			expected: map[string]interface{}{
				"userName": "user7",
				"domain":   "example.com",
				"params": map[string]string{
					"param1": "value1",
				},
			},
		},
	}

	for _, test := range testData {
		isSip, data := ParseSipAddress(test.str)

		if isSip != test.expectedSip {
			fmt.Printf("Expected isSip to be %v, but got %v for str: %s \n", test.expectedSip, isSip, test.str)
		}

		if !reflect.DeepEqual(data, test.expected) {
			fmt.Printf("Expected data to be %v, but got %v for str: %s \n", test.expected, data, test.str)
		}
	}
}

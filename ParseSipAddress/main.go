package main

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type SipAddress struct {
	DisplayName  string
	UserName     string
	UserPassword string
	Domain       string
	Ip           string
	Port         int
	Params       map[string]string
}

func main() {
	start := time.Now()

	//data := ParseSipAddress("text1 text2 'Little Guy' <user7:mypass> text3 text4 ")
	//PrintData(data)

	//for i := 0; i < 10000; i++ {
	//	UnitTest()
	//}

	UnitTest()

	duration := time.Since(start)
	fmt.Printf("Время выполнения: %d мкс\n", duration.Microseconds())
}

func ParseSipAddress(str string) (data SipAddress) {

	// Разбиваем строку на компоненты, используя пробел в качестве разделителя
	components := strings.Split(str, " ")

	// Функция для разбора SIP-адреса
	parseSipAddress := func(sipAddress string) {
		// Разделяем адрес на части по символу "@"
		parts := strings.Split(sipAddress, "@")

		// Проверяем условие для выполнения разбора
		if !(len(parts) > 1 || strings.HasPrefix(parts[0], "<") || strings.HasPrefix(parts[0], "sip:") || strings.HasPrefix(parts[0], "sips:")) {
			return
		}

		// Обработка имени пользователя
		userNameStr := strings.TrimPrefix(parts[0], "<")
		if strings.HasPrefix(userNameStr, "sip:") {
			userNameStr = userNameStr[4:]
		} else if strings.HasPrefix(userNameStr, "sips:") {
			userNameStr = userNameStr[5:]
		}

		// Разделяем имя пользователя и пароль
		userPasswordStr := strings.SplitN(userNameStr, ":", 2)
		if len(userPasswordStr) > 1 {
			data.UserPassword = userPasswordStr[1]
		}
		if strings.HasSuffix(userPasswordStr[0], ">") {
			// Удаляем символ ">" в конце имени пользователя
			data.UserName = userPasswordStr[0][:len(userPasswordStr[0])-1]
		} else {
			data.UserName = userPasswordStr[0]
		}

		if len(parts) > 1 {
			// Обработка домена и параметров
			domainParamsStr := strings.TrimSuffix(parts[1], ">")
			// Разделяем строку домена и параметров на компоненты
			componentsStr := strings.FieldsFunc(domainParamsStr, func(r rune) bool {
				return r == '?' || r == ';' || r == '>' || r == '&'
			})

			// Обработка домена и порта
			domainPortComponents := strings.SplitN(componentsStr[0], ":", 2)
			if len(domainPortComponents) > 1 {
				// Преобразуем порт в число
				port, err := strconv.Atoi(domainPortComponents[1])
				if err == nil {
					data.Port = port
				}
			}

			if net.ParseIP(domainPortComponents[0]) != nil {
				data.Ip = domainPortComponents[0]
			} else {
				data.Domain = domainPortComponents[0]
			}

			if len(componentsStr) > 1 {
				// Обработка параметров
				data.Params = make(map[string]string)
				for i := 1; i < len(componentsStr); i++ {
					if componentsStr[i] != "" {
						// Разделяем ключ и значение параметра
						keyValue := strings.SplitN(componentsStr[i], "=", 2)
						if len(keyValue) > 1 {
							data.Params[keyValue[0]] = keyValue[1]
						} else {
							data.Params[keyValue[0]] = ""
						}
					}
				}
			}
		}
	}

	// Проходим по каждому компоненту
	for i := 0; i < len(components); i++ {
		component := components[i]

		if strings.HasPrefix(component, "\"") || strings.HasPrefix(component, "'") {
			// Обработка отображаемого имени
			if strings.HasSuffix(component, "\"") || strings.HasSuffix(component, "'") {
				data.DisplayName = component[1 : len(component)-1]
			} else {
				displayNameStr := component[1:]
				for z := i + 1; z < len(components); z++ {
					displayNameStr += " " + components[z]
					if strings.HasSuffix(components[z], "\"") || strings.HasSuffix(components[z], "'") {
						data.DisplayName = displayNameStr[:len(displayNameStr)-1]
						i = z
						break
					}
				}
			}
		} else if net.ParseIP(component) != nil {
			data.Ip = component
		} else {
			parseSipAddress(component)
		}
	}
	return
}

// UnitTest Функция для тестирования ParseSipAddress
func UnitTest() {

	testData := []struct {
		str      string
		expected SipAddress
	}{
		{
			str: "sip:alice@example.com:5060;transport=udp",
			expected: SipAddress{
				UserName: "alice",
				Domain:   "example.com",
				Port:     5060,
				Params: map[string]string{
					"transport": "udp",
				},
			},
		},
		{
			str: "sip:bob@example.com",
			expected: SipAddress{
				UserName: "bob",
				Domain:   "example.com",
			},
		},
		{
			str: "sip:user123:pass456@domain.com:5060?param1=value1",
			expected: SipAddress{
				UserName:     "user123",
				UserPassword: "pass456",
				Domain:       "domain.com",
				Port:         5060,
				Params: map[string]string{
					"param1": "value1",
				},
			},
		},
		{
			str: "sips:admin@10.0.0.1",
			expected: SipAddress{
				UserName: "admin",
				Ip:       "10.0.0.1",
			},
		},
		{
			str: "sip:carol@example.com:5080",
			expected: SipAddress{
				UserName: "carol",
				Domain:   "example.com",
				Port:     5080,
			},
		},
		{
			str: "sip:user1@domain.com",
			expected: SipAddress{
				UserName: "user1",
				Domain:   "domain.com",
			},
		},
		{
			str: "sip:user2:pass@192.168.0.1",
			expected: SipAddress{
				UserName:     "user2",
				UserPassword: "pass",
				Ip:           "192.168.0.1",
			},
		},
		{
			str: "sips:user3@domain.org:1234;param1=value1;param2=value2",
			expected: SipAddress{
				UserName: "user3",
				Domain:   "domain.org",
				Port:     1234,
				Params: map[string]string{
					"param1": "value1",
					"param2": "value2",
				},
			},
		},
		{
			str: "sip:john@example.com",
			expected: SipAddress{
				UserName: "john",
				Domain:   "example.com",
			},
		},
		{
			str: "sip:jane@192.168.1.100:5062;param1=value1",
			expected: SipAddress{
				UserName: "jane",
				Ip:       "192.168.1.100",
				Port:     5062,
				Params: map[string]string{
					"param1": "value1",
				},
			},
		},
		{
			str: "sips:guest@example.org",
			expected: SipAddress{
				UserName: "guest",
				Domain:   "example.org",
			},
		},
		{
			str: "sip:user4:password@10.0.0.2:5080?param1=value1&param2=value2",
			expected: SipAddress{
				UserName:     "user4",
				UserPassword: "password",
				Ip:           "10.0.0.2",
				Port:         5080,
				Params: map[string]string{
					"param1": "value1",
					"param2": "value2",
				},
			},
		},
		{
			str: "sip:test@example.com:5060;param1=value1;param2=value2",
			expected: SipAddress{
				UserName: "test",
				Domain:   "example.com",
				Port:     5060,
				Params: map[string]string{
					"param1": "value1",
					"param2": "value2",
				},
			},
		},
		{
			str: "sips:user5@192.0.2.1",
			expected: SipAddress{
				UserName: "user5",
				Ip:       "192.0.2.1",
			},
		},
		{
			str: "sips:user5@192.0.2.1?param1=value",
			expected: SipAddress{
				UserName: "user5",
				Ip:       "192.0.2.1",
				Params: map[string]string{
					"param1": "value",
				},
			},
		},
		{
			str: "sip:user6:pass@domain.net:6000",
			expected: SipAddress{
				UserName:     "user6",
				UserPassword: "pass",
				Domain:       "domain.net",
				Port:         6000,
			},
		},
		{
			str: "sips:admin@10.0.0.3:5061",
			expected: SipAddress{
				UserName: "admin",
				Ip:       "10.0.0.3",
				Port:     5061,
			},
		},
		{
			str: "'LittleGuy' <sips:admin@10.0.0.3:5061>;tag=123",
			expected: SipAddress{
				DisplayName: "LittleGuy",
				UserName:    "admin",
				Ip:          "10.0.0.3",
				Port:        5061,
				Params: map[string]string{
					"tag": "123",
				},
			},
		},
		{
			str: "text1 'LittleGuy' <sips:admin@10.0.0.3:5061>;tag=123",
			expected: SipAddress{
				DisplayName: "LittleGuy",
				UserName:    "admin",
				Ip:          "10.0.0.3",
				Port:        5061,
				Params: map[string]string{
					"tag": "123",
				},
			},
		},
		{
			str: "text1 text2 \"LittleGuy\" <sips:admin@10.0.0.3:5061>;tag=123 text3",
			expected: SipAddress{
				DisplayName: "LittleGuy",
				UserName:    "admin",
				Ip:          "10.0.0.3",
				Port:        5061,
				Params: map[string]string{
					"tag": "123",
				},
			},
		},
		{
			str: "<sips:admin@10.0.0.3:5061>;tag=123",
			expected: SipAddress{
				UserName: "admin",
				Ip:       "10.0.0.3",
				Port:     5061,
				Params: map[string]string{
					"tag": "123",
				},
			},
		},
		{
			str: "textx 'LittleGuy' <sip:lg@domain.net>",
			expected: SipAddress{
				DisplayName: "LittleGuy",
				UserName:    "lg",
				Domain:      "domain.net",
			},
		},
		{
			str: "textx \"Little Guy\" <sip:lg@domain.net>",
			expected: SipAddress{
				DisplayName: "Little Guy",
				UserName:    "lg",
				Domain:      "domain.net",
			},
		},
		{
			str: "<sips:admin@10.0.0.3:5061?tag=123>",
			expected: SipAddress{
				UserName: "admin",
				Ip:       "10.0.0.3",
				Port:     5061,
				Params: map[string]string{
					"tag": "123",
				},
			},
		},
		{
			str: "sip:2011@192.168.1.150:5060;alias=192.168.1.151~5060~1",
			expected: SipAddress{
				UserName: "2011",
				Ip:       "192.168.1.150",
				Port:     5060,
				Params: map[string]string{
					"alias": "192.168.1.151~5060~1",
				},
			},
		},
		{
			str: "sip:user7@example.com;param1=value1",
			expected: SipAddress{
				UserName: "user7",
				Domain:   "example.com",
				Params: map[string]string{
					"param1": "value1",
				},
			},
		},
		{
			str: "'Little Guy' <sip:user7>",
			expected: SipAddress{
				DisplayName: "Little Guy",
				UserName:    "user7",
			},
		},
		{
			str: "text1 192.168.34.25 text2",
			expected: SipAddress{
				Ip: "192.168.34.25",
			},
		},
		{
			str: "text1 text2   192.168.34.25  text2",
			expected: SipAddress{
				Ip: "192.168.34.25",
			},
		},
	}

	for _, test := range testData {
		data := ParseSipAddress(test.str)

		if !reflect.DeepEqual(data, test.expected) {
			fmt.Printf("Expected data to be %v, but got %v for str: %s \n", test.expected, data, test.str)
		}
	}
}

func PrintData(data SipAddress) {
	fmt.Printf("\nDisplayName: ", data.DisplayName)
	fmt.Printf("\nUserName: ", data.UserName)
	fmt.Printf("\nUserPassword: ", data.UserPassword)
	fmt.Printf("\nDomain: ", data.Domain)
	fmt.Printf("\nIp: ", data.Ip)
	fmt.Printf("\nPort: ", data.Port)
	fmt.Printf("\nParams: ", data.Params)
}

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
	//data := ParseSipAddress("'Little Guy' <sip:user7> ")
	//PrintData(data)
	UnitTest()
	duration := time.Since(start)
	fmt.Printf("Время выполнения: %d мс\n", duration.Microseconds())
}

func ParseSipAddress(str string) (data SipAddress) {

	// Разбиваем строку на компоненты, используя пробел в качестве разделителя
	components := strings.Split(str, " ")

	parseSipAddress := func(sipAddress string) {
		parts := strings.Split(sipAddress, "@")
		if !(len(parts) > 1 || strings.HasPrefix(parts[0], "<") || strings.HasPrefix(parts[0], "sip:") || strings.HasPrefix(parts[0], "sips:")) {
			return
		}

		userNameStr := strings.TrimPrefix(parts[0], "<")
		if strings.HasPrefix(userNameStr, "sip:") {
			userNameStr = userNameStr[4:]
		} else if strings.HasPrefix(userNameStr, "sips:") {
			userNameStr = userNameStr[5:]
		}

		userPasswordStr := strings.Split(userNameStr, ":")
		if len(userPasswordStr) > 1 {
			data.UserPassword = userPasswordStr[1]
		}
		if strings.HasSuffix(userPasswordStr[0], ">") {
			data.UserName = userPasswordStr[0][:len(userPasswordStr[0])-1]
		} else {
			data.UserName = userPasswordStr[0]
		}

		if len(parts) > 1 {
			domainParamsStr := strings.TrimSuffix(parts[1], ">")
			componentsStr := strings.FieldsFunc(domainParamsStr, func(r rune) bool {
				return r == '?' || r == ';' || r == '>' || r == '&'
			})

			domainPortComponents := strings.Split(componentsStr[0], ":")
			if len(domainPortComponents) > 1 {
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
				data.Params = make(map[string]string)
				for i := 1; i < len(componentsStr); i++ {
					if componentsStr[i] != "" {
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

//func NewParseSipAddress1(str string) (data SipAddress) {
//
//	// Разбиваем строку на компоненты, используя пробел в качестве разделителя
//	components := strings.Split(str, " ")
//
//	// Проходим по каждому компоненту
//	for i := 0; i < len(components); i++ {
//
//		// Проверяем, начинается ли компонент с кавычек
//		if strings.HasPrefix(components[i], "\"") || strings.HasPrefix(components[i], "'") {
//			// Если компонент заканчивается кавычками, то он содержит полное отображаемое имя
//			if strings.HasSuffix(components[i], "\"") || strings.HasSuffix(components[i], "'") {
//				// Удаляем кавычки и сохраняем отображаемое имя
//				data.DisplayName = components[i][1 : len(components[i])-1]
//			} else if i < len(components) && (strings.HasSuffix(components[i+1], "\"") || strings.HasSuffix(components[i+1], "'")) {
//				// Если следующий компонент также заканчивается кавычками, объединяем текущий и следующий компоненты
//				data.DisplayName = components[i][1:] + " " + components[i+1][:len(components[i+1])]
//			}
//
//		} else if net.ParseIP(components[i]) != nil {
//			// Если компонент является IP-адресом, сохраняем его в поле "Ip"
//			data.Ip = components[i]
//		} else if strings.HasPrefix(components[i], "sip:") || strings.HasPrefix(components[i], "sips:") {
//			// Если компонент начинается с префикса "sip:" или "sips:", разбираем его как SIP-адрес
//			parts := strings.Split(components[i], "@")
//			if len(parts) > 1 {
//				userNameStr := parts[0]
//				if strings.HasPrefix(parts[0], "sip:") {
//					// Удаляем префикс "sip:" из имени пользователя
//					userNameStr = parts[0][4:]
//				} else if strings.HasPrefix(parts[0], "sips:") {
//					// Удаляем префикс "sips:" из имени пользователя
//					userNameStr = parts[0][5:]
//				}
//
//				// Разбиваем имя пользователя и пароль
//				userPasswordStr := strings.Split(userNameStr, ":")
//				if len(userPasswordStr) == 1 {
//					// Если пароль отсутствует, сохраняем только имя пользователя
//					data.UserName = userPasswordStr[0]
//				} else {
//					// Если есть пароль, сохраняем имя пользователя и пароль
//					data.UserName = userPasswordStr[0]
//					data.UserPassword = userPasswordStr[1]
//				}
//
//				// Разбиваем оставшуюся часть на компоненты, используя символы "?", ";", ">" и "&" в качестве разделителей
//				componentsStr := strings.FieldsFunc(parts[1], func(r rune) bool {
//					return r == '?' || r == ';' || r == '>' || r == '&'
//				})
//				domainParamsStr := make([]string, 0, len(components))
//				for _, component := range componentsStr {
//					if component != "" {
//						// Игнорируем пустые компоненты
//						domainParamsStr = append(domainParamsStr, component)
//					}
//				}
//
//				// Разбираем домен и порт
//				domainPortStr := strings.Split(domainParamsStr[0], ":")
//				if len(domainPortStr) == 1 {
//					if net.ParseIP(domainPortStr[0]) != nil {
//						// Если домен является IP-адресом, сохраняем его в поле "Ip"
//						data.Ip = domainPortStr[0]
//					} else {
//						// Если домен является именем хоста, сохраняем его в поле "Domain"
//						data.Domain = domainPortStr[0]
//					}
//				} else {
//					if net.ParseIP(domainPortStr[0]) != nil {
//						// Если домен является IP-адресом, сохраняем его в поле "Ip"
//						data.Ip = domainPortStr[0]
//					} else {
//						// Если домен является именем хоста, сохраняем его в поле "Domain"
//						data.Domain = domainPortStr[0]
//					}
//
//					// Преобразуем порт в целое число и сохраняем его в поле "Port"
//					port, err := strconv.Atoi(domainPortStr[1])
//					if err == nil {
//						data.Port = port
//					}
//				}
//				if len(domainParamsStr) > 1 {
//					// Разбор параметров
//					data.Params = make(map[string]string)
//					for i := 1; i < len(domainParamsStr); i++ {
//						if domainParamsStr[i] == "" {
//							continue
//						}
//
//						// Разбиваем параметр на ключ и значение
//						keyValue := strings.SplitN(domainParamsStr[i], "=", 2)
//						if len(keyValue) > 1 {
//							// Если параметр имеет значение, сохраняем его в поле "Params"
//							data.Params[keyValue[0]] = keyValue[1]
//						} else {
//							// Если параметр не имеет значения, сохраняем пустую строку в поле "Params"
//							data.Params[keyValue[0]] = ""
//						}
//					}
//				}
//			}
//		}
//
//		// Обработка компонента, начинающегося с "<"
//		if strings.HasPrefix(components[i], "<") {
//			parts := strings.Split(components[i], "@")
//			if len(parts) > 1 {
//				userNameStr := parts[0]
//				if strings.HasPrefix(parts[0], "<sip:") {
//					// Удаляем префикс "<sip:" из имени пользователя
//					userNameStr = parts[0][5:]
//				} else if strings.HasPrefix(parts[0], "<sips:") {
//					// Удаляем префикс "<sips:" из имени пользователя
//					userNameStr = parts[0][6:]
//				} else {
//					// Удаляем префикс "<" из имени пользователя
//					userNameStr = parts[0][1:]
//				}
//
//				// Разбиваем имя пользователя и пароль
//				userPasswordStr := strings.Split(userNameStr, ":")
//				if len(userPasswordStr) == 1 {
//					// Если пароль отсутствует, сохраняем только имя пользователя
//					data.UserName = userPasswordStr[0]
//				} else {
//					// Если есть пароль, сохраняем имя пользователя и пароль
//					data.UserName = userPasswordStr[0]
//					data.UserPassword = userPasswordStr[1]
//				}
//
//				//domainParamsStr := strings.Split(parts[1], ">")
//
//				// Разбиваем оставшуюся часть на компоненты, используя символы "?", ";", ">" в качестве разделителей
//				componentsStr := strings.FieldsFunc(parts[1], func(r rune) bool {
//					return r == '?' || r == ';' || r == '>'
//				})
//				domainParamsStr := make([]string, 0, len(components))
//				for _, component := range componentsStr {
//					if component != "" {
//						domainParamsStr = append(domainParamsStr, component)
//					}
//				}
//
//				// Разбираем домен и порт
//				domainPortStr := strings.Split(domainParamsStr[0], ":")
//				if len(domainPortStr) == 1 {
//					if net.ParseIP(domainPortStr[0]) != nil {
//						// Если домен является IP-адресом, сохраняем его в поле "Ip"
//						data.Ip = domainPortStr[0]
//					} else {
//						// Если домен является именем хоста, сохраняем его в поле "Domain"
//						data.Domain = domainPortStr[0]
//					}
//				} else {
//					if net.ParseIP(domainPortStr[0]) != nil {
//						// Если домен является IP-адресом, сохраняем его в поле "Ip"
//						data.Ip = domainPortStr[0]
//					} else {
//						// Если домен является именем хоста, сохраняем его в поле "Domain"
//						data.Domain = domainPortStr[0]
//					}
//					// Преобразуем порт в целое число и сохраняем его в поле "Port"
//					port, err := strconv.Atoi(domainPortStr[1])
//					if err == nil {
//						data.Port = port
//					}
//				}
//
//				// Разбор параметров
//				if len(domainParamsStr) > 1 {
//					data.Params = make(map[string]string)
//					for i := 1; i < len(domainParamsStr); i++ {
//						if domainParamsStr[i] == "" {
//							continue
//						}
//						// Разбиваем параметр на ключ и значение
//						keyValue := strings.SplitN(domainParamsStr[i], "=", 2)
//						if len(keyValue) > 1 {
//							// Если параметр имеет значение, сохраняем его в params
//							data.Params[keyValue[0]] = keyValue[1]
//						} else {
//							// Если параметр не имеет значения, сохраняем пустую строку в params
//							data.Params[keyValue[0]] = ""
//						}
//					}
//				}
//			}
//
//		}
//
//	}
//	return
//}
//
//// ParseSipAddress разбирает SIP-адрес и возвращает флаг успешного разбора и данные адреса
//func ParseSipAddress(str string) (isSip bool, data SipAddress) {
//	// Создаем пустой map для хранения разобранных значений
//	data = SipAddress{
//		Params: make(map[string]string),
//	}
//
//	// Удаление всех пробелов в начале и конце строки
//	str = strings.TrimSpace(str)
//
//	// Проверяем, начинается ли строка с символа "<"
//	if strings.HasPrefix(str, "<") {
//		// Удаляем символы < и >
//		str = strings.ReplaceAll(str[1:], ">", "")
//	} else {
//		// Если строка не начинается с "<", разбиваем ее на две части по символу "<"
//		displayNameParts := strings.Split(str, "<")
//		if len(displayNameParts) > 1 {
//			allowedChars := "+-!$"
//			var cleanedDisplayName strings.Builder
//			for _, char := range displayNameParts[0] {
//				if unicode.IsLetter(char) || unicode.IsDigit(char) || strings.ContainsRune(allowedChars, char) {
//					cleanedDisplayName.WriteRune(char)
//				}
//				if strings.ContainsRune(allowedChars, char) {
//					cleanedDisplayName.WriteRune(char)
//				}
//			}
//
//			// Если есть разделитель "<", сохраняем часть до "<" в поле "displayName"
//			//data.DisplayName = strings.TrimSpace(displayNameParts[0])
//			data.DisplayName = strings.TrimSpace(cleanedDisplayName.String())
//			// Сохраняем часть после "<" в переменную str и удаляем символ ">"
//			str = strings.ReplaceAll(displayNameParts[1], ">", "")
//		}
//	}
//
//	// Удаление всех пробелов из строки
//	str = strings.ReplaceAll(str, " ", "")
//
//	// Проверяем префикс "sip:" или "sips:" в адресе.
//	if strings.HasPrefix(str, "sip:") {
//		// Если адрес начинается с "sip:", удаляем этот префикс
//		str = str[4:]
//	} else if strings.HasPrefix(str, "sips:") {
//		// Если адрес начинается с "sips:", удаляем этот префикс
//		str = str[5:]
//	} else {
//		// Если префикс не найден, устанавливаем флаг isSip в false и возвращаем данные
//		isSip = false
//		return
//	}
//
//	// Разделение строки на подстроки с использованием символов "?", ";", и "&"
//	parts := strings.FieldsFunc(str, func(r rune) bool {
//		return r == '?' || r == ';' || r == '&'
//	})
//
//	// Разбор левой части
//	leftPart := strings.Split(parts[0], "@")
//
//	// Разбираем имя пользователя и пароль
//	userAndPass := strings.Split(leftPart[0], ":")
//	data.UserName = userAndPass[0]
//	if len(userAndPass) > 1 {
//		data.UserPassword = userAndPass[1]
//	}
//
//	// Разбираем домен и порт
//	if len(leftPart) > 1 {
//		domainAndPort := strings.Split(leftPart[1], ":")
//		if net.ParseIP(domainAndPort[0]) != nil {
//			// Если домен является IP-адресом, сохраняем его в data["ip"]
//			data.Ip = domainAndPort[0]
//		} else {
//			// Если домен является именем хоста, сохраняем его в data["domain"]
//			data.Domain = domainAndPort[0]
//		}
//		if len(domainAndPort) > 1 {
//			// Преобразуем порт в целое число и сохраняем его в data["port"]
//			port, err := strconv.Atoi(domainAndPort[1])
//			if err == nil {
//				data.Port = port
//			}
//		}
//	}
//
//	// Разбираем параметры
//	if len(parts) > 1 {
//		sipParameters := parts[1:]
//		for _, sipParameter := range sipParameters {
//			keyValue := strings.SplitN(sipParameter, "=", 2)
//			if len(keyValue) > 1 {
//				// Если параметр имеет значение, сохраняем его в params
//				data.Params[keyValue[0]] = keyValue[1]
//			} else {
//				// Если параметр не имеет значения, сохраняем пустую строку в params
//				data.Params[keyValue[0]] = ""
//			}
//		}
//
//	}
//
//	// Установка флага isSip в true, чтобы указать успешный разбор адреса
//	isSip = true
//	return
//}

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
			str: "textx 'Little Guy' <sip:lg@domain.net>",
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
	}

	for _, test := range testData {
		data := ParseSipAddress(test.str)

		if !reflect.DeepEqual(data, test.expected) {
			fmt.Printf("Expected data to be %v, but got %v for str: %s \n", test.expected, data, test.str)
		}
	}
}

func PrintData(data SipAddress) {
	fmt.Printf("\nUserName: ", data.UserName)
	fmt.Printf("\nUserPassword: ", data.UserPassword)
	fmt.Printf("\nDisplayName: ", data.DisplayName)
	fmt.Printf("\nDomain: ", data.Domain)
	fmt.Printf("\nIp: ", data.Ip)
	fmt.Printf("\nPort: ", data.Port)
	fmt.Printf("\nParams: ", data.Params)
}

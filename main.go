package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var (
	// URL stores the host
	URL = "https://play.google.com"
	// Attributes store values of fields
	nameAttr     = "name"
	authorAttr   = "document-subtitle primary"
	categoryAttr = "document-subtitle category"
	descriptAttr = "description"
	ratingAttr   = "tiny-star star-rating-non-editable-container"
	numberAttr   = "reviews-stats"
	updateAttr   = "datePublished"
	appAttr      = "card-click-target"
)

// Application contains all the required fields
type Application struct {
	name     string
	url      string
	author   string
	category string
	descript string
	rating   string
	number   string
	update   string
}

// GetAppData get data from app page
func GetAppData(link, word string) (Application, bool, error) {
	flag := false                         // check if the app contains our word
	doc, err := goquery.NewDocument(link) // make goquery document
	if err != nil {
		var app Application
		return app, flag, fmt.Errorf("сan't open the app page, please check your Internet connection")
	}

	var data Application
	data.url = link
	wg := &sync.WaitGroup{} // create WaitGroup
	wg.Add(1)
	// the goroutine find selectors with tag <a>
	// and read the information about author and category
	go func(d *goquery.Document) {
		defer wg.Done()
		d.Find("a").Each(func(i int, s *goquery.Selection) {
			if class, ok := s.Attr("class"); ok {
				if class == authorAttr {
					data.author = s.Text()
				} else if class == categoryAttr {
					data.category = s.Text()
				}
			}
		})
	}(doc)

	wg.Add(1)
	// the goroutine find selectors with tag <div>
	// and read the remained information about app
	go func(d *goquery.Document) {
		defer wg.Done()
		d.Find("div").Each(func(i int, s *goquery.Selection) {
			class, _ := s.Attr("class")
			if class == "main-content" { // in this area we will find our information
				s.Children().Find("div").Each(func(j int, kid *goquery.Selection) {
					if info, ok := kid.Attr("itemprop"); ok {
						if info == descriptAttr {
							data.descript = kid.Text()
							if strings.Contains(data.descript, word) {
								flag = true
							}
						} else if info == nameAttr {
							data.name = kid.Text()
							if strings.Contains(data.name, word) {
								flag = true
							}
						} else if info == updateAttr {
							data.update = kid.Text()
						}
					} else if info, ok = kid.Attr("class"); ok {
						if info == ratingAttr {
							data.rating, _ = kid.Attr("aria-label")
						} else if info == numberAttr {
							data.number = kid.Text()
						}
					}
				})
			}
		})
	}(doc)
	wg.Wait() // Wait ending of work each goroutine

	return data, flag, nil
}

// ParseURL - main function of parsing
func ParseURL(path, word string) ([]Application, error) {
	doc, err := goquery.NewDocument(path) // make goquery document
	if err != nil {
		return nil, fmt.Errorf("сan't open the website, please check your Internet connection")
	}

	allinfo := []Application{}
	wg := &sync.WaitGroup{}
	// find selectors with tag <a>
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		// find the Application link and go through it
		if class == appAttr {
			link, _ := s.Attr("href")
			link = URL + link
			wg.Add(1)
			go func() {
				defer wg.Done()
				// parse app page and check if it contains our word
				if data, ok, err := GetAppData(link, word); ok {
					if err == nil {
						allinfo = append(allinfo, data)
					}
				}
			}()
		}
	})
	wg.Wait()
	return allinfo, nil
}

// MakeRequestURL create search-url based on entered word
func MakeRequestURL(word string) string {
	binary := []byte(word)
	encodedStr := strings.ToUpper(hex.EncodeToString(binary)) // string like D0A1D0B1...
	re := regexp.MustCompile("..")                            // make regex expression
	slice := re.FindAllString(encodedStr, -1)                 // string split like [D0 A1 D0 B1 ...]
	encodedStr = strings.Join(slice, "%")                     // now the string is D0%A1%D0%B1%...
	requestURL := URL + "/store/search?q=%" + encodedStr + "&c=apps"
	return requestURL
}

// ConvertData -> to json
func ConvertData(data []Application) []byte {
	dict := make(map[int][]string, len(data))
	for idx, app := range data {
		dict[idx] = append(dict[idx], app.name)
		dict[idx] = append(dict[idx], app.url)
		dict[idx] = append(dict[idx], app.author)
		dict[idx] = append(dict[idx], app.category)
		dict[idx] = append(dict[idx], app.descript)
		dict[idx] = append(dict[idx], app.rating)
		dict[idx] = append(dict[idx], app.number)
		dict[idx] = append(dict[idx], app.update)
	}
	result, _ := json.Marshal(&dict)
	return result
}

func main() {
	stop := false // flag, which stops main loop
	// main loop of the program
	for !stop {
		fmt.Println("Введите название приложения:")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		fmt.Println("Идет сканирование Google Play ...")
		path := MakeRequestURL(scanner.Text())
		data, err := ParseURL(path, scanner.Text())
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Сканирование завершено, число найденных приложений =", len(data))
		fmt.Println("Сохранить данные в формате json в файл? [yes/no]")
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if strings.ToLower(scanner.Text()) == "yes" {
			flag := false
			for !flag {
				fmt.Println("Введите имя файла")
				scanner = bufio.NewScanner(os.Stdin)
				scanner.Scan()
				if scanner.Text() != "" {
					file, _ := os.Create(scanner.Text())
					// convert data to json
					jsonData := ConvertData(data)
					file.Write(jsonData)
					file.Close()
					flag = true
				}
			}
		}
		fmt.Println("Вывести данные на экран? [yes/no]")
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if strings.ToLower(scanner.Text()) == "yes" {
			for _, info := range data {
				fmt.Println("Название:             ", info.name)
				fmt.Println("URL:                   ", info.url)
				fmt.Println("Автор:                ", info.author)
				fmt.Println("Категория:            ", info.category)
				fmt.Println("Описание:\n", info.descript)
				fmt.Println("     ", info.rating)
				fmt.Println("            ", info.number)
				fmt.Println("Последнее обновление: ", info.update)
				fmt.Println()
			}
			fmt.Println("Число результатов =", len(data))
		}
		fmt.Println("Получить данные по другому запросу? [yes/no]")
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if strings.ToLower(scanner.Text()) == "no" {
			stop = true
		}
	}
}

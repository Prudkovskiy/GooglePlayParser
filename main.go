package main

import (
	"bufio"
	"crypto/md5"
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
	word1 := strings.ToUpper(word)
	word2 := strings.ToLower(word)
	word3 := strings.ToUpper(word[0:2]) + word[2:] // сбербанк -> Сбербанк (ru)
	word4 := strings.ToUpper(word[0:1]) + word[1:] // sberbank -> Sberbank (en)
	flag := false                                  // check if the app contains our word
	doc, err := goquery.NewDocument(link)          // make goquery document
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
							ok1 := strings.Contains(data.descript, word1)
							ok2 := strings.Contains(data.descript, word2)
							ok3 := strings.Contains(data.descript, word3)
							ok4 := strings.Contains(data.descript, word4)
							ok5 := strings.Contains(data.descript, word)
							if ok1 || ok2 || ok3 || ok4 || ok5 {
								flag = true
							}
						} else if info == nameAttr {
							data.name = kid.Text()
							ok1 := strings.Contains(data.name, word1)
							ok2 := strings.Contains(data.name, word2)
							ok3 := strings.Contains(data.name, word3)
							ok4 := strings.Contains(data.name, word4)
							ok5 := strings.Contains(data.name, word)
							if ok1 || ok2 || ok3 || ok4 || ok5 {
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
	dict := make(map[string][]string, len(data))

	for _, app := range data {
		key := fmt.Sprintf("%x", md5.Sum([]byte(app.url+app.name)))
		if _, ok := dict[key]; ok {
			continue
		}
		dict[key] = append(dict[key], app.name)
		dict[key] = append(dict[key], app.url)
		dict[key] = append(dict[key], app.author)
		dict[key] = append(dict[key], app.category)
		dict[key] = append(dict[key], app.descript)
		dict[key] = append(dict[key], app.rating)
		dict[key] = append(dict[key], app.number)
		dict[key] = append(dict[key], app.update)
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
		fmt.Println("Сохранить данные в формате json в файл? [Y/N]")
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if strings.ToLower(scanner.Text()) == "y" {
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
		fmt.Println("Вывести данные на экран? [Y/N]")
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if strings.ToLower(scanner.Text()) == "y" {
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
		fmt.Println("Получить данные по другому запросу? [Y/N]")
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if strings.ToLower(scanner.Text()) == "n" {
			stop = true
		}
	}
}

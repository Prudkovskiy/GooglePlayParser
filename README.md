Парсер магазина приложений Google Play
---

В целях защиты бренда часто возникает задача найти мобильные приложения, содержащие определенное ключевое слово. Как правило, магазины не предоставляют API для этих целей, вследствие чего приходится парсить сайт.  
 
Данное консольное приложение предназначено для того, чтобы по заданному ключевому слову возвращать информацию о найденных приложениях магазина, которые в названии или описании содержат данное слово.  
Search-URL для Google Play выглядит следующим образом: `https://play.google.com/store/search?q=<%ключевое слово%>&c=apps`  
По этому url происходит первоначальный переход для дальнейшей обработки всех приложений на данной странице.  

Сама работа с консольным приложением выглядит следующим образом:
```
Введите название приложения:  
_Сбербанк_  
Идет сканирование Google Play ...  
Сканирование завершено, число найденных приложений = 120  
Сохранить данные в формате json в файл? [yes/no]  
_yes_ (в случае _no_ продолжаем без записи в файл)  
Введите имя файла  
_sberbank.txt_  
Вывести данные на экран? [yes/no]  
_yes_  

...  
...  
Название:               Сбербанк Онлайн Казахстан!  
URL:                    https://play.google.com/store/apps/details?id=com.rssl.sboldb  
Автор:                  Sberbank Kazakhstan  
Категория:              Финансы  
Описание:  
texttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttext  
texttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttext  
texttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttexttext  
                Средняя оценка: 4,7 из 5  
                Оценок: 4 543  
  Последнее обновление:  29 января 2018 г.  
...  
...  
Число результатов = 120  
Получить данные по другому запросу? [yes/no]  
_no_ (в случае _yes_ идем в начало)  

```


Чтобы собрать исполняемый файл  
`go build main.go`  
Для запуска без сборки проекта  
`go run main.go`  
Для запуска тестов  
`go test -v`  
Для запуска бенчмарков  
`go test -bench . main_test.go main.go`  

Из дополнительных пакетов необходимо установить GoQuery  
`go get github.com/PuerkitoBio/goquery`  


Бенчмарки на трех запросах (Сбербанк, Viber, Yandex)
```
$ go test -bench . main_test.go main.go
goos: windows
goarch: amd64
BenchmarkSberbank-8            1        4491305000 ns/op
BenchmarkViber-8               1        4666631600 ns/op
BenchmarkYandex-8              1        5258527300 ns/op
PASS
ok      command-line-arguments  29.659s
```


При запуске тестов:  
```
$ go test -v
=== RUN   TestMakeRequestURL
--- PASS: TestMakeRequestURL (0.00s)
=== RUN   TestConvertData
--- PASS: TestConvertData (0.61s)
=== RUN   TestGetAppData
--- PASS: TestGetAppData (1.19s)
=== RUN   TestParseURL
--- PASS: TestParseURL (13.04s)
PASS
ok      GooglePlayParser        14.998s
```


`cover.html` <- здесь можно посмотреть тестовое покрытие кода (96.2 %)
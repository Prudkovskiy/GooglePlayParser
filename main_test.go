package main

import (
	"reflect"
	"testing"
)

type TestCase struct {
	word string
	URL  string
}

type TestCaseApp struct {
	word    string
	link    string
	contain bool
}

func BenchmarkSberbank(b *testing.B) {
	word := "Сбербанк"
	url := MakeRequestURL(word)
	ParseURL(url, word)
}

func BenchmarkViber(b *testing.B) {
	word := "Viber"
	url := MakeRequestURL(word)
	ParseURL(url, word)
}

func BenchmarkYandex(b *testing.B) {
	word := "Yandex"
	url := MakeRequestURL(word)
	ParseURL(url, word)
}

func TestMakeRequestURL(t *testing.T) {
	cases := []TestCase{
		TestCase{"Сбербанк", "https://play.google.com/store/search?q=%D0%A1%D0%B1%D0%B5%D1%80%D0%B1%D0%B0%D0%BD%D0%BA&c=apps"},
		TestCase{"Тинькофф", "https://play.google.com/store/search?q=%D0%A2%D0%B8%D0%BD%D1%8C%D0%BA%D0%BE%D1%84%D1%84&c=apps"},
		TestCase{"Альфа-Банк", "https://play.google.com/store/search?q=%D0%90%D0%BB%D1%8C%D1%84%D0%B0%2D%D0%91%D0%B0%D0%BD%D0%BA&c=apps"},
	}
	for caseNum, item := range cases {
		result := MakeRequestURL(item.word)
		if !reflect.DeepEqual(result, item.URL) {
			t.Errorf("[%d] wrong results: got \n%+v \nexpected \n%+v",
				caseNum, result, item.URL)
		}
	}
}

func TestConvertData(t *testing.T) {
	cases := []TestCase{
		TestCase{"Instagram", "https://play.google.com/store/apps/details?id=com.instagram.android&hl=ru"},
		TestCase{"Вконтакте", "https://play.google.com/store/apps/details?id=com.vkontakte.android&hl=ru"},
		TestCase{"Facebook", "https://play.google.com/store/apps/details?id=com.facebook.katana&hl=ru"},
	}
	var info []Application
	for caseNum, item := range cases {
		data, _, err := GetAppData(item.URL, item.word)
		if err != nil {
			t.Errorf("[%d] unexpected error %v", caseNum, err)
		}
		info = append(info, data)
	}
	byteInfo := ConvertData(info)
	if len(byteInfo) == 0 {
		t.Errorf("incorrect marshalling in json")
	}
}

func TestGetAppData(t *testing.T) {
	cases := []TestCaseApp{
		TestCaseApp{"яндекс", "https://play.google.com/store/apps/details?id=com.cleanmaster.mguard", false},
		TestCaseApp{"сбербанк", "https://play.google.com/store/apps/details?id=com.aksioma.aksiomapitanie", false},
		TestCaseApp{"яндекс", "https://play.google.com/store/apps/details?id=ru.indipartner.connectdriver&hl=ru", true},
		TestCaseApp{"сбербанк", "https://play.google.com/store/apps/details?id=mobi.sevenwinds.sberbank", true},
		TestCaseApp{"viber", "https://play.google.com/store/apps/details?id=com.whatsapp&hl=ru", false},
		TestCaseApp{"WhatsApp", "https://play.google.com/store/apps/details?id=com.viber.voip&hl=ru", false},
		TestCaseApp{"viber", "https://play.google.com/store/apps/details?id=com.viber.voip&hl=ru", true},
		TestCaseApp{"WhatsApp", "https://play.google.com/store/apps/details?id=com.whatsapp&hl=ru", true},
	}

	for caseNum, item := range cases {
		_, contain, err := GetAppData(item.link, item.word)
		if !item.contain && contain {
			t.Errorf("[%d] the app doesn't contains this word ", caseNum)
		}
		if item.contain == false && contain == true {
			t.Errorf("[%d] the app contains this word ", caseNum)
		}
		if err != nil {
			t.Errorf("[%d] unexpected error %v", caseNum, err)
		}
	}
}

func TestParseURL(t *testing.T) {
	cases := []TestCase{
		TestCase{"Starbucks", "https://play.google.com/store/search?q=Starbucks&c=apps"},
		TestCase{"WhatsApp", "https://play.google.com/store/search?q=WhatsApp&c=apps"},
		TestCase{"вконтакте", "https://play.google.com/store/search?q=%D0%B2%D0%BA%D0%BE%D0%BD%D1%82%D0%B0%D0%BA%D1%82%D0%B5&c=apps"},
	}

	for caseNum, item := range cases {
		app, err := ParseURL(item.URL, item.word)
		if err != nil {
			t.Errorf("[%d] unexpected error %v", caseNum, err)
		}
		if len(app) == 0 {
			t.Errorf("[%d] have no information for this request", caseNum)
		}
	}
}

/*
    to create test coverage:
	go test -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

	to start all benchmarks:
	go test -bench . main_test.go main.go
	go test -bench . -benchmem main_test.go main.go

*/

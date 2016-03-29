package main

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
	"strings"
)

var searchingKeyword string = "สถานที่ท่องเที่ยวในไทย"

func main(){
	url := fmt.Sprintf("https://www.google.co.th/search?q=%s&oq=%s", searchingKeyword)
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	list := getListOfUrl(body)

	for i:=0;i<len(list);i++{
		pageResponse := getResponseFromUrl(list[i][1])
		groupOfWord := cutOffHtmlTag(pageResponse)
		fmt.Println(strings.Split(groupOfWord, " "))
	}
}

func getListOfUrl(body []byte) [][]string{
	re := regexp.MustCompile("href=\"/url[?]q=([https://][^w][^e][^b][^c][^a][^c][^h][^e].*?)\"")
	result := re.FindAllStringSubmatch(string(body), -1)
	return result
}

func getResponseFromUrl(url string) string{
	resp, err := http.Get(url)
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		panic(err)
	}

	return string(body)
}


func cutOffHtmlTag(body string) string{
	re := regexp.MustCompile("<(.*?)>")
	cutOff := re.ReplaceAllLiteralString(body, " ")
	return cutOff
}
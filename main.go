package main

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

var searchingKeyword string = "สถานที่ท่องเที่ยวในไทย"

func main(){
	for start := 0 ;start < 1000 ; start+=10 {
		url := fmt.Sprintf("https://www.google.co.th/search?q=%s&oq=%s&start=%d", searchingKeyword, searchingKeyword, start)
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

		for i := 0; i < len(list); i++ {
			pageResponse := getResponseFromUrl(list[i][1])
			cutOffJavaScriptString := cutOffJavaScript(pageResponse)
			cutOffHtmlTagString := cutOffHtmlTag(cutOffJavaScriptString)
			completeString := strings.Replace(cutOffHtmlTagString, " ", "", -1)
			completeString = strings.Replace(completeString, "\t", "", -1)
			completeString = strings.Replace(completeString, "\r\n", "\n", -1)
			arrayString := strings.Split(completeString, "\n")
			arrayString = delete_empty(arrayString)
			for j := 0; j < len(arrayString); j++ {
				fmt.Printf("%d, %s\n", j, arrayString[j])
			}
		}
		time.Sleep(30*time.Second)
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
	re := regexp.MustCompile("<[^>]*>")
	cutOff := re.ReplaceAllLiteralString(body, `</\/>`)
	resultCutOff := strings.Replace(cutOff, `</\/>`, "", -1)
	return resultCutOff
}

func cutOffJavaScript(body string) string{
	re := regexp.MustCompile(`<script\b[^>]*>([\s\S]*?)<\/script>`)
	cutOff := re.ReplaceAllLiteralString(body, `</\/>`)
	resultCutOff := strings.Replace(cutOff, `</\/>`, "", -1)
	return resultCutOff
}

func delete_empty (s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
package main

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
	"strings"
//"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"log"
)

var searchingKeyword string = "สถานที่ท่องเที่ยวในภูเก็ต"

type Word struct{
	Spell string	`bson:spell`
	Frequency int	`bson:frequency`
}

func main(){
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("tc").C("phuket")

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
			if(pageResponse == ""){
				continue
			}
			cutOffJavaScriptString := cutOffJavaScript(pageResponse)
			cutOffHtmlTagString := cutOffHtmlTag(cutOffJavaScriptString)
			deletedEnglish := deleteEnglish(cutOffHtmlTagString)
			completeString := strings.Replace(deletedEnglish, "\t", "", -1)
			completeString = strings.Replace(completeString, "\r\n", "\n", -1)
			arrayString := strings.Split(completeString, "\n")
			arrayString = splitInnerMember(arrayString)
			arrayString = delete_empty(arrayString)
			for j := 0; j < len(arrayString); j++ {
				//fmt.Printf("%d, %s\n", j, arrayString[j])
				var word Word
				err = c.Find(bson.M{"spell":arrayString[j]}).One(&word)

				if err != nil{
					var w Word
					w.Frequency = 1
					w.Spell = arrayString[j]
					c.Insert(w)
				}else{
					fmt.Println("found same " + word.Spell)
					changeWord := word
					changeWord.Frequency = changeWord.Frequency+1
					err = c.Update(word, changeWord)
				}
			}
		}
		//time.Sleep(30*time.Second)
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
		//log.Fatal(err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		//log.Fatal(err)
		return ""
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
		rstr := strings.Replace(str, " ", "", -1)
		if rstr == "" {
			str = ""
		}

		if str != "" && str != " " {
			r = append(r, str)
		}
	}
	return r
}

func splitInnerMember(s []string) []string{
	var r []string
	for i:=0;i<len(s);i++ {
		str := s[i]
		isContainSpace := strings.Contains(str, " ")
		if isContainSpace {
			spliced := strings.Split(str, " ")
			r = append(r, spliced...)
			s = append(s[:i], s[i+1:]...)
			i = i-1
		}
	}

	s = append(s, r...)
	return s
}


func deleteEnglish(s string) string{
	re := regexp.MustCompile(`[\w;:/)(.,!@&><|}{#'\"\]\[+*′%?=-]+`)
	result := re.ReplaceAllLiteralString(s, `</\/>`)
	result = strings.Replace(result, `</\/>`, "", -1)
	return result
}
package main

import (
	"fmt"
	"net/http"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

type Word struct{
	Spell string	`bson:spell`
	Frequency int	`bson:frequency`
}

func mostSayAboutIslandInPhuketDescCsv(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("tc").C("phuket")
	var words []Word
	err = c.Find(
		bson.M{
			"spell":bson.M{
				"$regex":".*เกาะ.*",
			}}).Sort("-frequency").All(&words)
	if err != nil{
		panic(err)
	}
	//b, err := json.Marshal(words)
	//if err != nil{
	//	panic(err)
	//}
	for i:=0;i<len(words);i++{
		fmt.Fprintf(w, "%s, %d\n", words[i].Spell, words[i].Frequency)
	}
}

func mostSayAboutIslandInPhuketDescJson(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("tc").C("phuket")
	var words []Word
	err = c.Find(
		bson.M{
			"spell":bson.M{
				"$regex":".*เกาะ.*",
			}}).Sort("-frequency").All(&words)
	if err != nil{
		panic(err)
	}
	b, err := json.Marshal(words)
	if err != nil{
		panic(err)
	}

	fmt.Fprintf(w, "%s", b)

}

func main() {
	http.HandleFunc("/phuket/csv", mostSayAboutIslandInPhuketDescCsv) // set router
	http.HandleFunc("/phuket/json", mostSayAboutIslandInPhuketDescJson) // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
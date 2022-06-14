package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Recipes struct {
	XMLName xml.Name `xml:"recipes" json:"-"`
	Cake    []Cake   `xml:"cake" json:"cake"`
}

type Cake struct {
	Name        string `xml:"name" json:"name"`
	StoveTime   string `xml:"stovetime" json:"time"`
	Ingredients []Item `xml:"ingredients>item" json:"ingredients"`
}

type Item struct {
	ItemName  string `xml:"itemname" json:"ingredient_name"`
	ItemCount string `xml:"itemcount" json:"ingredient_count"`
	ItemUnit  string `xml:"itemunit,omitempty" json:"ingredient_unit,omitempty"`
}

type DBReader interface {
	ReadFile() string
}

func Process(r DBReader) {
	fmt.Println(r.ReadFile())
}

type XMLFile string
type JSONFile string

func (x XMLFile) ReadFile() string {
	path := getPath()
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file!")
		os.Exit(1)
	}
	defer f.Close()
	decoder := xml.NewDecoder(f)
	var data Recipes
	for {
		tok, err := decoder.Token()
		if tok == nil {
			break
		}
		if err != nil {
			fmt.Println("Error decoder token!")
			os.Exit(1)
		}
		switch tp := tok.(type) {
		case xml.StartElement:
			if tp.Name.Local == "cake" {
				var c Cake
				decoder.DecodeElement(&c, &tp)
				data.Cake = append(data.Cake, c)
			}
		}
	}
	u, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("Error converting into JSON!")
		os.Exit(1)
	}
	return string(u)
}

func (j JSONFile) ReadFile() string {
	path := getPath()
	var data Recipes
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file!")
		os.Exit(1)
	}
	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading file!")
		os.Exit(1)
	}
	err = json.Unmarshal([]byte(byteValue), &data)
	if err != nil {
		fmt.Println("Error!")
		os.Exit(1)
	}
	u, err := xml.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("Error converting into JSON!")
		os.Exit(1)
	}
	return string(u)
}

func getPath() string {
	var path string
	if len(os.Args) == 3 {
		path = "../../datasets/" + os.Args[2]
	} else {
		path = os.Args[1]
	}
	return path
}

func main() {
	path := getPath()
	if strings.Contains(path, "xml") {
		var x XMLFile
		Process(x)
	} else if strings.Contains(path, "json") {
		var j JSONFile
		Process(j)
	} else {
		fmt.Println("Incorrect filename!")
	}
}

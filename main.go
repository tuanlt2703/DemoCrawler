package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gocolly/colly"
)

type NetworkSpec struct {
	Url        string `json:"Url"`
	Technology string `json:"Technology"`
	TwoG       string `json:"2G"`
	ThreeG     string `json:"3G"`
	FourG      string `json:"4G"`
	Speed      string `json:"Speed"`
}

func writeToJSON(data []NetworkSpec) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = ioutil.WriteFile("NetworkSpec.json", file, 0644)
}

func crawlGSMDetailPage(url string) {
	var netSpec NetworkSpec
	netSpec.Url = url

	/* Refer: http://go-colly.org/docs/examples/basic/ */
	c := colly.NewCollector(
		colly.AllowedDomains("gsmarena.com", "www.gsmarena.com"),
	)

	c.OnHTML("#specs-list", func(e *colly.HTMLElement) {
		e.ForEach("table tbody", func(_ int, e_table_tbody *colly.HTMLElement) {
			//switch for test crawling only network spec
			switch e_table_tbody.ChildText("tr th") {
			case "Network":
				e_table_tbody.ForEach("tr", func(_ int, e_table_tbody_tr *colly.HTMLElement) {
					//ingore the ttl empty fornow
					switch e_table_tbody_tr.ChildText(".ttl") {
					case "Technology":
						netSpec.Technology = e_table_tbody_tr.ChildText(".nfo a")
					case "2G bands":
						netSpec.TwoG = e_table_tbody_tr.ChildText(".nfo")
					case "3G bands":
						netSpec.ThreeG = e_table_tbody_tr.ChildText(".nfo")
					case "4G bands":
						netSpec.FourG = e_table_tbody_tr.ChildText(".nfo")
					case "Speed":
						netSpec.Speed = e_table_tbody_tr.ChildText(".nfo")
					}
				})

				netSpecs = append(netSpecs, netSpec)
			}
		})
	})

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Crawling", request.URL.String())
	})

	c.Visit(url)
}

var netSpecs []NetworkSpec

func main() {
	crawlGSMDetailPage("https://www.gsmarena.com/samsung_galaxy_note10+-9732.php")
	crawlGSMDetailPage("https://www.gsmarena.com/apple_iphone_xs-9318.php")

	writeToJSON(netSpecs)
}

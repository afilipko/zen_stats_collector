package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/net/html"

	log "github.com/inconshreveable/log15"
	conf "github.com/olebedev/config"
)

// DefaultConfigPath default conf file path for server and cli client
const DefaultConfigPath string = "conf.yml"
const ZenBasePublisherURL string = "https://zen.yandex.ru/media/id/"

var confPath string
var publisherInfoJson []byte
var publisherData map[string]interface{}

// Config contains env name and *config.Config
type Config struct {
	PublisherID   string
	CookieContent string
}

// ParseConfig read conf from CONF_PATH
func parseConfig(confPath string) *Config {
	cfg, err := conf.ParseYamlFile(confPath)
	if err != nil {
		fmt.Println(err)

	}
	publisherID, _ := cfg.String("publisher_id")
	cookie, _ := cfg.String("cookie_content")
	collectorCfg := Config{PublisherID: publisherID, CookieContent: cookie}
	return &collectorCfg
}

func buildPublisherPageRequest(cfg *Config) (*http.Request, error) {
	publisherURL := ZenBasePublisherURL + cfg.PublisherID
	req, error := http.NewRequest("GET", publisherURL, nil)
	if error != nil {
		return nil, error
	}
	req.Header.Add("Cookie", cfg.CookieContent)
	req.Header.Add("Cache-Control", "no-cache")

	return req, nil
}

func getPublisherInfoFromResponse(resp *http.Response) []byte {
	defer resp.Body.Close()
	tokenizer := html.NewTokenizer(resp.Body)
	inTextArea := false
	for {
		token := tokenizer.Next()
		switch token {
		case html.ErrorToken:
			fmt.Println(tokenizer.Err())
			log.Error("No publisher info")
			return make([]byte, 0)
		case html.TextToken:
			// emitBytes should copy the []byte it receives,
			if inTextArea {
				log.Info("found")
				inTextArea = false
				return tokenizer.Text()[:]
			}
		case html.StartTagToken, html.EndTagToken:
			tn := tokenizer.Token()
			if tn.Data == "textarea" {
				for _, attr := range tn.Attr {
					if attr.Key == "id" && attr.Val == "init_data" {
						inTextArea = true
					}
				}
			}
		}
	}
	return make([]byte, 0)
}

func main() {
	if confPath = os.Getenv("CONF_PATH"); confPath == "" {
		confPath = DefaultConfigPath
	}
	log.Info("Config path " + confPath)
	cfg := parseConfig(confPath)

	client := &http.Client{}
	req, _ := buildPublisherPageRequest(cfg)

	resp, _ := client.Do(req)

	if resp.StatusCode == http.StatusOK {
		publisherInfoJson = getPublisherInfoFromResponse(resp)
		if err := json.Unmarshal(publisherInfoJson, &publisherData); err != nil {
			log.Error("Unmarshal error")
			panic(err)
		}
		fmt.Println(publisherData["publications"])
	} else {
		log.Error(strconv.Itoa(resp.StatusCode) + " code during main page request")
	}

}

// {
// 	"id":"59f71abc7ddde8f2f478b679",
// 	"publisherId":"59301217e3cda86683de2230",
// 	"content":{
// 	   "type":"article",
// 	   "embedVideoContents":[

// 	   ],
// 	   "preview":{
// 		  "title":"Перекресток Миллера",
// 		  "snippet":"G",
// 		  "tags":[

// 		  ]
// 	   },
// 	   "modTime":1509366471380
// 	},
// 	"privateData":{
// 	   "hasChanges":true,
// 	   "hasPublished":false,
// 	   "snippetFrozen":false,
// 	   "statistics":{
// 		  "feedShows":0,
// 		  "shows":0,
// 		  "likes":0,
// 		  "dislikes":0,
// 		  "shares":0
// 	   }
// 	},
// 	"itemId":"-11325337430083931",
// 	"publisherItemId":"4696001067498724411"
//  }

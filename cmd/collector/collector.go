package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"zen_stats_collector/data"

	"golang.org/x/net/html"

	log "github.com/inconshreveable/log15"
	conf "github.com/olebedev/config"
)

type PublicationStatistics struct {
	Shares         int
	Views          int
	ViewsTillEnd   int
	SumViewTimeSec int
	FeedShows      int
	Shows          int
	Likes          int
	Dislikes       int
}

type Publication struct {
	Id          string
	PublisherId int
	Privatedata *PublicationStatistics
}

// DefaultConfigPath default conf file path for server and cli client
const DefaultConfigPath string = "conf.yml"
const ZenBasePublisherURL string = "https://zen.yandex.ru/media/id/"
const ZenPublicationsURL string = "https://zen.yandex.ru/media-api/publisher-publications-next-page?publisherId=:publisherId&pageSize=:pageSize"
const PageSize = 10

// &lastPublicationId=:last

var confPath string
var publisherInfoJson []byte
var publisherData []*PublicationStatistics
var publications *data.PublicationsInfo

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

func buildPublicationsRequest(cfg *Config, lastPublicationId string) (*http.Request, error) {
	fmt.Println(ZenPublicationsURL)
	publisherScopedUrl := strings.Replace(ZenPublicationsURL, ":publisherId", cfg.PublisherID, 1)
	fmt.Println(publisherScopedUrl)
	pageSizeScoped := strings.Replace(publisherScopedUrl, "pageSize", "Ab", 1)
	fmt.Println(publisherScopedUrl)
	finalUrl := pageSizeScoped
	// strconv.Itoa(PageSize)
	// if len(lastPublicationId) > 0 {
	// 	finalUrl = finalUrl + "&lastPublicationId=" + lastPublicationId
	// }

	req, error := http.NewRequest("GET", finalUrl, nil)
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
	// req, _ := buildPublisherPageRequest(cfg)
	lastPublicationId := ""
	for {
		req, _ := buildPublicationsRequest(cfg, lastPublicationId)
		resp, _ := client.Do(req)
		fmt.Println(resp.Body)
		if resp.StatusCode == http.StatusOK {
			if err := json.Unmarshal(publisherInfoJson, &publications); err != nil {
				fmt.Print(publications.Publications[0].Content.Preview.Title)
			}
		}
		break
	}

	// 	publisherInfoJson = getPublisherInfoFromResponse(resp)
	// 	// n := bytes.IndexByte(publisherInfoJson, 0)
	// 	// m, err := objx.FromJSON(string(publisherInfoJson))
	// 	// if err != nil {
	// 	// fmt.Print("Err->")
	// 	// }
	// 	// publications := objx.MustFromJSON(string(publisherInfoJson))
	// 	// for key, value := range publications {
	// 	// 	if key == "publications" {
	// 	// 		fmt.Print(value[0].title)
	// 	// 	}
	// 	// }
	// 	// if err := json.Unmarshal(publisherInfoJson, &publisherData); err != nil {
	// 	// 	log.Error("Unmarshal error")
	// 	// 	panic(err)
	// 	// }
	// 	// fmt.Println(publisherData)
	// 	// t := publisherData["publications"]
	// 	// for _, publicationInfo := range t {

	// 	// }
	// } else {
	// 	log.Error(strconv.Itoa(resp.StatusCode) + " code during main page request")
	// }

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

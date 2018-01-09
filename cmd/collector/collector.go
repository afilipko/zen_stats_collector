package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"zen_stats_collector/data"

	"golang.org/x/net/html"

	log "github.com/inconshreveable/log15"
	client "github.com/influxdata/influxdb/client/v2"
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
	PublisherID    string
	CookieContent  string
	InfluxdbURL    string
	InfluxdbUser   string
	InfluxdbPass   string
	InfluxdbSeries string
	InfluxdbName   string
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

// ParseConfig read conf from CONF_PATH
func parseConfig(confPath string) *Config {
	cfg, err := conf.ParseYamlFile(confPath)
	if err != nil {
		fmt.Println(err)

	}
	publisherID, _ := cfg.String("publisher_id")
	cookie, _ := cfg.String("cookie_content")
	influxdbURL, _ := cfg.String("influxdb_url")
	influxdbUser, _ := cfg.String("influxdb_user")
	influxdbPass, _ := cfg.String("influxdb_pass")
	influxdbSeries, _ := cfg.String("influxdb_series")
	influxdbName, _ := cfg.String("influxdb_name")
	collectorCfg := Config{PublisherID: publisherID, CookieContent: cookie,
		InfluxdbURL: influxdbURL, InfluxdbUser: influxdbUser,
		InfluxdbPass: influxdbPass, InfluxdbSeries: influxdbSeries,
		InfluxdbName: influxdbName}
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
	publisherScopedUrl := strings.Replace(ZenPublicationsURL, ":publisherId", cfg.PublisherID, 1)
	pageSizeScoped := strings.Replace(publisherScopedUrl, ":pageSize", strconv.Itoa(PageSize), 1)
	finalUrl := pageSizeScoped
	// strconv.Itoa(PageSize)
	if len(lastPublicationId) > 0 {
		finalUrl = finalUrl + "&lastPublicationId=" + lastPublicationId
	}

	req, error := http.NewRequest("GET", finalUrl, nil)
	if error != nil {
		return nil, error
	}
	req.Header.Add("Cookie", cfg.CookieContent)
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
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

func sendStatToDB(publiction data.PulicationObj, cfg *Config) {

	fmt.Println(publiction.Content.Preview.Title)
	fmt.Println("==========>")
	fmt.Println("+++++++++=")
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     cfg.InfluxdbURL,
		Username: cfg.InfluxdbUser,
		Password: cfg.InfluxdbPass,
	})
	if err != nil {
		fmt.Println(err)
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  cfg.InfluxdbName,
		Precision: "h",
	})
	if err != nil {
		fmt.Println(err)
	}

	// Create a point and add to batch
	tags := map[string]string{"publicationId": publiction.ID, "publishedId": publiction.PublisherID}
	fields := map[string]interface{}{
		"feedShows":      publiction.PrivateData.Statistics.FeedShows,
		"shows":          publiction.PrivateData.Statistics.Shows,
		"likes":          publiction.PrivateData.Statistics.Likes,
		"dislikes":       publiction.PrivateData.Statistics.Dislikes,
		"shares":         publiction.PrivateData.Statistics.Shares,
		"views":          publiction.PrivateData.Statistics.Views,
		"viewsTillEnd":   publiction.PrivateData.Statistics.ViewsTillEnd,
		"sumViewTimeSec": publiction.PrivateData.Statistics.SumViewTimeSec,
	}

	pt, err := client.NewPoint(cfg.InfluxdbSeries, tags, fields, time.Now())
	if err != nil {
		fmt.Println(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		fmt.Println(err)
	}
	wg.Done()
}

var wg sync.WaitGroup

func processPageResponse(resp *http.Response, cfg *Config) string {
	defer resp.Body.Close()
	var localLastID string
	fmt.Println("----------> --> ")
	if resp.StatusCode == http.StatusOK {
		fmt.Println("----------> --> ---->")
		err := json.NewDecoder(resp.Body).Decode(&publications)
		// err := json.Unmarshal(publisherInfoJson, &publications);
		if err != nil {
			fmt.Println("ERR")
			return ""
		}
		if len(publications.Publications) == 0 {
			return ""
		}
		for _, publicationObj := range publications.Publications {
			fmt.Println(publicationObj.PrivateData)
			if publicationObj.PrivateData.HasPublished == true {
				wg.Add(1)
				go sendStatToDB(publicationObj, cfg)
			}
			localLastID = publicationObj.ID
		}
		fmt.Println("Page ended ----------------------->")
	}
	return localLastID
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
	// done := make(chan bool, 1)
	for {
		req, _ := buildPublicationsRequest(cfg, lastPublicationId)
		// fmt.Println(formatRequest(req))
		fmt.Println("---------->")
		resp, errresp := client.Do(req)
		if errresp != nil {
			fmt.Println(errresp)
		}
		lastPublicationId = processPageResponse(resp, cfg)
		if len(lastPublicationId) == 0 {
			break
		}
	}
	wg.Wait()
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

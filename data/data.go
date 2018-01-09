package data

type PulicationObj struct {
	ID          string `json:"id"`
	PublisherID string `json:"publisherId"`
	AddTime     int64  `json:"addTime"`
	Content     struct {
		Type               string        `json:"type"`
		EmbedVideoContents []interface{} `json:"embedVideoContents"`
		Preview            struct {
			Title   string        `json:"title"`
			Snippet string        `json:"snippet"`
			Tags    []interface{} `json:"tags"`
		} `json:"preview"`
		ModTime int64 `json:"modTime"`
	} `json:"content"`
	PrivateData struct {
		HasChanges    bool `json:"hasChanges"`
		HasPublished  bool `json:"hasPublished"`
		SnippetFrozen bool `json:"snippetFrozen"`
		Statistics    struct {
			FeedShows      int `json:"feedShows"`
			Shows          int `json:"shows"`
			Likes          int `json:"likes"`
			Dislikes       int `json:"dislikes"`
			Shares         int `json:"shares"`
			Views          int `json:"views"`
			ViewsTillEnd   int `json:"viewsTillEnd"`
			SumViewTimeSec int `json:"sumViewTimeSec"`
		} `json:"statistics"`
	} `json:"privateData"`
	ItemID          string `json:"itemId"`
	PublisherItemID string `json:"publisherItemId"`
}

type PublicationsInfo struct {
	Publications []PulicationObj `json:"publications"`
	Images       struct {
		Five9352A2Fe3Cda85Cf4156F7E struct {
			ID        string `json:"id"`
			Namespace string `json:"namespace"`
			GroupID   int    `json:"groupId"`
			ImageName string `json:"imageName"`
			MetaRaw   string `json:"metaRaw"`
			Sizes     struct {
				MinMWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_m_webp"`
				MinXh struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_xh"`
				MinXxhWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_xxh_webp"`
				H struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"h"`
				Xxh struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"xxh"`
				MinH struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_h"`
				M struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"m"`
				Xh struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"xh"`
				XxhWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"xxh_webp"`
				MinM struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_m"`
				MinHWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_h_webp"`
				MinXhWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_xh_webp"`
				MinXxh struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_xxh"`
				Orig struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"orig"`
				HWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"h_webp"`
				MWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"m_webp"`
				XhWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"xh_webp"`
			} `json:"sizes"`
		} `json:"59352a2fe3cda85cf4156f7e"`
		Five9352A49E3Cda85Cf4156F7F struct {
			ID        string `json:"id"`
			Namespace string `json:"namespace"`
			GroupID   int    `json:"groupId"`
			ImageName string `json:"imageName"`
			MetaRaw   string `json:"metaRaw"`
			Sizes     struct {
				MinMWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_m_webp"`
				MinXh struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_xh"`
				MinXxhWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_xxh_webp"`
				H struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"h"`
				Xxh struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"xxh"`
				MinH struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_h"`
				M struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"m"`
				Xh struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"xh"`
				XxhWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"xxh_webp"`
				MinM struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_m"`
				MinHWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_h_webp"`
				MinXhWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_xh_webp"`
				MinXxh struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"min_xxh"`
				Orig struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"orig"`
				HWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"h_webp"`
				MWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"m_webp"`
				XhWebp struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"xh_webp"`
			} `json:"sizes"`
		} `json:"59352a49e3cda85cf4156f7f"`
	} `json:"images"`
}

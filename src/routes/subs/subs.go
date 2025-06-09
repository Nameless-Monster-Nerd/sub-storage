package sub

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/gin-gonic/gin"
	utils "github.com/nameless-Monster-Nerd/subtitle/src/modules"
	"github.com/nameless-Monster-Nerd/subtitle/src/modules/fetchRabbit"
	"github.com/nameless-Monster-Nerd/subtitle/src/modules/psql"
	uploadsubtitles "github.com/nameless-Monster-Nerd/subtitle/src/modules/uploadSubtitles"
)

func Sub(c *gin.Context) {
	scheme := "https"
	if utils.Env != "PRO"{
		scheme = "http"
	}
	host := c.Request.Host
	origin := fmt.Sprintf("%s://%s", scheme, host)

	id := c.Param("id")

	var ssPtr, epPtr *string
	if ss := c.Query("ss"); ss != "" {
		ssPtr = &ss
	}
	if ep := c.Query("ep"); ep != "" {
		epPtr = &ep
	}
	var common []psql.Sub
	subList, err := psql.BatchSearch(id, ssPtr, epPtr, true)
	fmt.Println("samir gupta ")
	fmt.Println(subList)
	if err != nil {
		fmt.Println(err)
	}
	if  len(subList) != 0  {
		common = subList

	}

	if len(subList) == 0 {
		
		result := fetchRabbit.FetchRabbit(id, ssPtr, epPtr)
		fmt.Println(result)
		var wg sync.WaitGroup
		var mu sync.Mutex
		subs := []psql.Sub{}

		for _, o := range result.Tracks {
			wg.Add(1)
			go func(lang, urlStr string) {
				defer wg.Done()
				uploadInfo := uploadsubtitles.UploadSubtitle(lang, urlStr, id, ssPtr, epPtr, true)
				fmt.Println(uploadInfo.Key)

				mu.Lock()
				subs = append(subs, psql.Sub{
					ID:     id,
					SS:     ssPtr,
					EP:     epPtr,
					Key:    uploadInfo.Key,
					Bucket: uploadInfo.Bucket,
					Lang:   lang,
					Flix:   true,
				})
				mu.Unlock()
			}(o.Lang, o.URL)
		}

		wg.Wait()

		psql.BatchUpload(subs)

		subList, err = psql.BatchSearch(id, ssPtr, epPtr, true)
		
		if err != nil {
			fmt.Println(err)
		}
		common = subList
	}

	outResult := []OutPut{}

	for _, o := range common {
		outResult = append(outResult, OutPut{
			Lang: o.Lang,
			URL:  fmt.Sprintf("%s/proxy.vtt?key=%s", origin, url.QueryEscape(o.Key)),
		})
	}
	c.JSON(200, outResult)
}

type OutPut struct {
	Lang string `json:"lang"`
	URL  string `json:"url"`
}

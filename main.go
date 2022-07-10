package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/eduncan911/podcast"
)

type Source struct {
	EpisodeFrequency       string    `json:"episode_frequency"`
	EstimatedNextEpisodeAt time.Time `json:"estimated_next_episode_at"`
	HasSeasons             bool      `json:"has_seasons"`
	SeasonCount            int       `json:"season_count"`
	EpisodeCount           int       `json:"episode_count"`
	HasMoreEpisodes        bool      `json:"has_more_episodes"`
	Podcast                struct {
		Url         string      `json:"url"`
		Title       string      `json:"title"`
		Author      interface{} `json:"author"`
		Description string      `json:"description"`
		Category    string      `json:"category"`
		Audio       bool        `json:"audio"`
		ShowType    interface{} `json:"show_type"`
		Uuid        string      `json:"uuid"`
		Episodes    []struct {
			Uuid      string      `json:"uuid"`
			Title     string      `json:"title"`
			Url       string      `json:"url"`
			FileType  string      `json:"file_type"`
			FileSize  int         `json:"file_size"`
			Duration  interface{} `json:"duration"`
			Published time.Time   `json:"published"`
			Type      string      `json:"type"`
		} `json:"episodes"`
		Paid      int `json:"paid"`
		Licensing int `json:"licensing"`
	} `json:"podcast"`
}

func main() {
	resp, err := http.Get("https://cache.pocketcasts.com/mobile/podcast/full/5bcb4950-ab15-013a-d8be-0acc26574db2")
	if err != nil {
		panic(err)
	}
	source := new(Source)
	err = json.NewDecoder(resp.Body).Decode(source)
	if err != nil {
		panic(err)
	}
	// build
	now := time.Now()
	p := podcast.New(
		"梁文道·八分 - 第四季",
		"https://shop.vistopia.com.cn/detail?id=HLHQu",
		"《八分》第四季2022年1月14日起已开始更新，欢迎下载「看理想」App或喜马拉雅独家收听。\n「看理想」App，一个陪你成长的知识剧场。\n「八分」是由梁文道和看理想团队共同打造的一档全新文化类音频节目。梁文道将和你一起从理想走到现实，重新审视文化现象，社会趋势和热点话题。\n欢迎下载「看理想」App和梁文道互动，每集节目中，梁文道会回答一位听众提问。",
		nil, &now,
	)
	p.Language = "zh-cn"
	p.AddAuthor("梁文道", "klx@vistopia.com.cn")
	p.AddImage("https://cdn.vistopia.com.cn/img/podcast-bafen.jpg")
	for _, episode := range source.Podcast.Episodes {
		item := podcast.Item{
			Title:       episode.Title,
			Description: "no description",
			Link:        episode.Url,
		}
		resp, err := http.Head(episode.Url)
		if err != nil {
			panic(err)
		}
		item.AddEnclosure(episode.Url, podcast.MP3, resp.ContentLength)
		t, err := http.ParseTime(resp.Header.Get("last-modified"))
		if err == nil {
			item.AddPubDate(&t)
		}
		_, err = p.AddItem(item)
		if err != nil {
			panic(err)
		}
	}
	f, err := os.Create("dist/s4.xml")
	if err != nil {
		panic(err)
	}
	err = p.Encode(f)
	if err != nil {
		panic(err)
	}
}

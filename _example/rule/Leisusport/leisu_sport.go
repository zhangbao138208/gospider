package leisusport

import (
	"github.com/nange/gospider/spider"
	log "github.com/sirupsen/logrus"
	"time"
	"fmt"
)
func init() {
	spider.Register(rule)
}
var (
	outputFields=[]string{"province", "area", "aqi", "quality_grade", "pm10", "pm25", "no2", "so2", "o3", "co", "tip", "publish_time"}
	constraints= spider.NewConstraints(outputFields,
		"VARCHAR(16) NOT NULL DEFAULT ''",
		"VARCHAR(16)  NOT NULL DEFAULT ''",
		"VARCHAR(128) NOT NULL DEFAULT ''",
		"VARCHAR(8) NOT NULL DEFAULT ''",
		"VARCHAR(16) NOT NULL DEFAULT ''",
		"VARCHAR(16) NOT NULL DEFAULT ''",
		"VARCHAR(16)  NOT NULL DEFAULT ''",
		"VARCHAR(16)  NOT NULL DEFAULT ''",
		"VARCHAR(16)  NOT NULL DEFAULT ''",
		"VARCHAR(16)  NOT NULL DEFAULT ''",
		"VARCHAR(256) NOT NULL DEFAULT ''",
		"VARCHAR(256) NOT NULL DEFAULT ''",
	)
)
var rule = &spider.TaskRule{
	Name:           	"雷速体育足球比赛",
	Description:    	"抓雷速体育足球比赛信息",
	Namespace:      	"leisu_sport_competition",
	DisableCookies: 	true,
	OutputFields:   	outputFields,
	OutputConstraints:	constraints,
	Rule: &spider.Rule{
		Head: func(ctx *spider.Context) error { // 定义入口
			return ctx.VisitForNext("https://live.leisu.com/wanchang")
		},
		Nodes: map[int]*spider.Node{
			0: step1, // 第一步: 找到全国各省城市区县的链接

		},
	},
}

var step1 = &spider.Node{
	OnRequest: func(ctx *spider.Context, req *spider.Request) {
		log.Infof("Visiting %s", req.URL.String())
	},
	OnError: func(ctx *spider.Context, res *spider.Response, err error) error {
		log.Errorf("Visiting failed! url:%s, err:%s", res.Request.URL.String(), err.Error())
		// 出错时重试三次
		return Retry(ctx, 3)
	},
	OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
		`.city_list a`: func(ctx *spider.Context, el *spider.HTMLElement) error {
			link := el.Attr("href")
			return ctx.Visit(link)
		},
		`#list div ul`: func(ctx *spider.Context, ul *spider.HTMLElement) error {
			ul.ForEach(`li`, func(i int, li *spider.HTMLElement) {
				homeTeam := li.ChildText(`.lab-events a span`)
				println("homeTeam:",homeTeam)
				round := li.ChildText(`.lab-round`)
				println("homeTeam:",round)
				time := li.ChildText(`.lab-time`)
				println("homeTeam:",time)
				status := li.ChildText(`.lab-status span`)
				println("homeTeam:",status)
			})
			return nil
		},
	},
}
func Retry(ctx *spider.Context, count int) error {
	req := ctx.GetRequest()
	key := fmt.Sprintf("err_req_%s", req.URL.String())

	var et int
	if errCount := ctx.GetAnyReqContextValue(key); errCount != nil {
		et = errCount.(int)
		if et >= count {
			return fmt.Errorf("exceed %d counts", count)
		}
	}
	log.Infof("errCount:%d, we wil retry url:%s, after 1 second", et+1, req.URL.String())
	time.Sleep(time.Second)
	ctx.PutReqContextValue(key, et+1)
	ctx.Retry()

	return nil
}


package main

import (
	"fmt"
	model "homeplay/Model"
	repository "homeplay/Repository"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/queue"
	"github.com/robertkrimen/otto"
)

func main() {
	model.Startup()
	URL := "http://halihali2.com/"
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"),
		colly.MaxDepth(3),
		colly.Debugger(&debug.LogDebugger{}))

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 1 * time.Second,
		Delay:       1 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	q, _ := queue.New(
		5, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)
	q.AddURL(URL)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("异常信息为:", err)
		}
	}()

	vm := otto.New()

	//首页大类爬取
	c.OnHTML("div[class='h1 clearfix']", func(h *colly.HTMLElement) {
		defer func() {
			if err := recover(); err != nil {
				//	fmt.Println("异常信息为:", err)
				log.Panic(err)
			}
		}()
		href := h.ChildAttr("a[class='more']", "href")
		text := h.ChildText("a[href]>span")
		if href != "" {
			l := strings.Split(href, "/")[1]
			url := fmt.Sprintf("http://121.4.190.96:9991/getsortdata_all_z.php?action=%v&page=1&id=y2021&class=0&year=2021&area=all", l)
			q.AddURL(url)
			fmt.Printf("类型:%v,地址:%v\n", text, url)
		}
	})

	//爬取集数和链接
	c.OnHTML("a[class='li-hv']", func(h *colly.HTMLElement) {
		defer func() {
			if err := recover(); err != nil {
				//	fmt.Printf("异常信息为:%v\n", err)
				log.Panic(err)
			}
		}()
		video := new(model.VideoInfo)
		video.Name = h.Attr("title")
		video.ImgSrc = h.ChildAttr("div[class='img']>img[class='lazy']", "data-original")
		href := h.Attr("href")
		video.Url = fmt.Sprint(URL, href)
		video.VideoNo = strings.Split(href, "/")[2]
		video.Last = h.ChildText("div[class='img']>p[class='bz']")
		go repository.InsertInfos(video)
		fmt.Printf("%+v\n", video)
	})

	//取出集数
	c.OnHTML("script[src]", func(h *colly.HTMLElement) {
		defer func() {
			if err := recover(); err != nil {
				//	fmt.Println("异常信息为:", err)
				log.Panic(err)
			}
		}()
		src := h.Attr("src")
		if find := strings.Contains(src, "ne2"); find {
			fmt.Printf("地址:%v\n", src)
			res, err := http.Get(src)
			if err != nil {
				//	fmt.Printf("err:%v\n", err)
				log.Panic(err)
			}
			defer res.Body.Close()
			bytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			vm.Run(bytes)
			arr := "playarr"
			for i := 0; i < 11; i++ {
				playarr := arr
				if i != 0 {
					playarr = arr + "_" + fmt.Sprint(i)
				}
				v, _ := vm.Get(playarr)
				ids, _ := vm.Get("pl_id")
				no, _ := ids.Export()
				keys, err := v.Export()
				if err != nil {
					panic(err)
				}
				if keys == nil {
					continue
				}
				fmt.Printf("now load: %v\n", playarr)
				keyArray := keys.([]string)
				details := make([]model.VideoDetail, len(keyArray))
				for _, str := range keyArray {
					if !strings.Contains(str, "http") {
						str = "https://jsap.ahfuqi.net/?url=" + str
					}
					details = append(details, model.VideoDetail{
						VideoNo:  no.([]string)[0],
						SourceNo: i,
						Url:      strings.Split(str, ",")[0],
						Episode:  strings.Split(str, ",")[2],
					})
					//fmt.Printf("keys: %v\n", str)
				}
				go repository.InsertDetails(details)
			}
		}
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("err:%v", e.Error())
		log.Panic(e)
	})
	c.OnResponse(func(r *colly.Response) {
		defer func() {
			if err := recover(); err != nil {
				//	fmt.Println("异常信息为:", err)
				log.Panic(err)
			}
		}()
		pas := r.Request.URL.Query()
		if len(pas) > 0 {

			year, _ := strconv.Atoi(pas["year"][0])
			page, _ := strconv.Atoi(pas["page"][0])
			action := pas["action"][0]
			if len(r.Body) == 0 && year > time.Now().Year()-1 {
				year--
				url := fmt.Sprintf("http://121.4.190.96:9991/getsortdata_all_z.php?action=%v&page=0&id=y%v&class=0&year=%v&area=all", action, year, year)
				q.AddURL(url)
			} else if len(r.Body) > 0 {
				page++
				url := fmt.Sprintf("http://121.4.190.96:9991/getsortdata_all_z.php?action=%v&page=%v&id=y%v&class=0&year=%v&area=all", action, page, year, year)
				q.AddURL(url)
			}
		}
	})
	err := q.Run(c)
	if err != nil {
		log.Panic(err)
	}
}

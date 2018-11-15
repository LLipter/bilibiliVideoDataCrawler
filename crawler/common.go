package crawler

import (
	"github.com/LLipter/bilibiliVideoDataCrawler/conf"
	"github.com/LLipter/bilibiliVideoDataCrawler/proxy"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	wg           sync.WaitGroup
	curCrawlerNo chan bool
)

func init() {
	curCrawlerNo = make(chan bool, conf.VideoCrawlerConfig.MaxCrawlerNum)
}

func getResp(addr string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", addr, nil)
	req.Header.Add("User-Agent", conf.NetworkConfig.UserAgent)
	if conf.NetworkConfig.UseProxy {
		length := len(proxy.ProxyPool)
		proxyAddr := proxy.ProxyPool[rand.Intn(length)]
		urlProxy, err := url.Parse("http://" + proxyAddr)
		if err != nil {
			return nil, err
		}
		client = http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			},
			Timeout: time.Second * 30,
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

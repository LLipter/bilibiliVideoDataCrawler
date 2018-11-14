package proxy

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	ProxyPool []string
)

type data struct {
	Count      int
	Proxy_list []string
}

type proxyJson struct {
	Msg  string
	Code int
	Data data
}

func GetProxy() error {
	apiAddr := "https://dev.kdlapi.com/api/getproxy/?orderid=904212196080767&num=1000&b_pcchrome=1&b_pcie=1&b_pcff=1&protocol=2&method=2&an_tr=1&an_an=1&an_ha=1&sp1=1&sp2=1&quality=1&format=json&sep=1"
	resp, err := http.Get(apiAddr)
	if err != nil {
		return errors.New("get proxy failed, " + err.Error())
	}

	var proxyObj proxyJson
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	err = json.Unmarshal(buf, &proxyObj)
	if err != nil {
		return err
	}

	if proxyObj.Code != 0 {
		return errors.New("illegal get proxy parameters")
	}

	ProxyPool = proxyObj.Data.Proxy_list
	return nil
}

// change to your own codes to get proxy
func GetProxyRoutine() []string {
	for {
		err := GetProxy()
		if err != nil {
			log.Println(err)
			// retry after 10 seconds
			time.Sleep(time.Second * 10)
			continue
		}

		// refresh proxy pool every minute
		time.Sleep(time.Minute)
	}
}

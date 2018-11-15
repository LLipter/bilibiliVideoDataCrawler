package conf

type Conf struct {
	DB           DBConf `json:"database"`
	Network      NetworkConf
	VideoCrawler VideoCrawlerConf
}

type DBConf struct {
	User        string
	Passwd      string
	Host        string
	DBname      string
	MaxOpenConn int
}

type NetworkConf struct {
	UseProxy   bool
	UserAgent  string
	RetryTimes int
}

type VideoCrawlerConf struct {
	StartAid      int
	EndAid        int
	IsDaemon      bool
	MaxCrawlerNum int
}

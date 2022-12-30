package initialization

import (
	"github.com/olivere/elastic/v7"
	"go.backend/conf"
	"log"
)

type Es struct {
	*elastic.Client
}

func InitEs(c conf.EsConfig) *Es {
	//默认不记录日志

	// 本地启动时，需要设置为false，具体原因参考 https://juejin.cn/post/6895371414179446798
	sniffOpt := elastic.SetSniff(false)
	es, err := elastic.NewClient(elastic.SetURL(c.Host), sniffOpt)
	//ES, err = elastic.NewClient(elastic.SetURL(host))
	// 打印查询ES的语句
	//ES, err = elastic.NewClient(elastic.SetURL(host), elastic.SetTraceLog(new(TraceLog)))

	if err != nil {
		log.Println("ElasticSearch连接失败。。。", err)
	} else {
		log.Println("ElasticSearch已连接 >>> ")
	}
	return &Es{
		es,
	}
}

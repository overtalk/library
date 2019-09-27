package metric

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/robfig/cron"
)

type Config struct {
	AccessKey string `env:"ALIYUN_AK"`
	SecretKey string `env:"ALIYUN_SK"`
	Region    string `env:"ALIYUN_REGION"`
	GroupID   int64  `env:"ALIYUN_GROUPID"`
	Interval  string `env:"ALIYUN_CRON"`
}

// metricClient defines metric client interface
type MetricClient interface {
	GetCustomMetricMetricList() cms.PutCustomMetricMetricList
}

// CmsCron defines cms cron
type cmsGroup struct {
	config        Config
	cron          *cron.Cron
	cmsClient     *cms.Client
	metricClients []MetricClient
}

// cmsGroup is the constructor of cms group
func CmsGroup(config Config, metricClients ...MetricClient) (*cmsGroup, error) {
	cmsClient, err := cms.NewClientWithAccessKey(config.Region, config.AccessKey, config.SecretKey)
	if err != nil {
		return nil, err
	}
	return &cmsGroup{
		config:        config,
		cron:          cron.New(),
		cmsClient:     cmsClient,
		metricClients: metricClients,
	}, nil
}

// Start is to push data to aliyun cms
func (c *cmsGroup) Start() {
	fmt.Println("开始推送")
	c.cron.AddFunc(c.config.Interval, c.putCustomMetricList)
	c.cron.Start()
}

func (c *cmsGroup) putCustomMetricList() {
	req := cms.CreatePutCustomMetricRequest()
	req.MetricList = c.getCustomMetricMetricList()
	resp, err := c.cmsClient.PutCustomMetric(req)
	if err != nil || resp.Code != "200" {
		// put error
		fmt.Println("推送失败, err = ", err)
		fmt.Println("resp = ", resp)
		return
	}
	fmt.Println("推送成功")
}

// getCustomMetricMetricList get all custom metric list for the group
func (c *cmsGroup) getCustomMetricMetricList() *[]cms.PutCustomMetricMetricList {
	lists := make([]cms.PutCustomMetricMetricList, 0)
	for _, metricClient := range c.metricClients {
		list := metricClient.GetCustomMetricMetricList()
		if list.MetricName != "" && list.Values != "" && list.Dimensions != "" {
			list.Time = strconv.FormatInt(time.Now().Unix()*1000, 10)
			list.Type = "0" // free version
			list.GroupId = fmt.Sprintf("%d", c.config.GroupID)
			lists = append(lists, list)
		}
	}
	return &lists
}

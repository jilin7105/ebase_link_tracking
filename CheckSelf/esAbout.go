package CheckSelf

import (
	"context"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/jilin7105/ebase"
	"github.com/jilin7105/ebase/logger"
	"net/http"
)

func CheckIndex() error {
	es := ebase.GetEsV7("linkTracking")
	exist, err := CheckIndexExist(es)
	if err != nil {

		return err
	}
	if exist {
		logger.Info("创建索引linkTracking存在")
	} else {
		logger.Info("创建索引linkTracking 不存在")
		createdIndex(es)
	}
	return nil
}

func CheckIndexExist(es *elasticsearch.Client) (bool, error) {

	// 构造一个IndexExists请求
	req := esapi.IndicesExistsRequest{
		Index: []string{"auto_link_tracking"},
	}

	// 执行IndexExists请求
	res, err := req.Do(context.Background(), es)
	if err != nil {

		return false, err
	}
	defer res.Body.Close()

	es.Indices.Delete([]string{"auto_link_tracking"})
	return false, nil
	// 检查响应状态码
	if res.IsError() {
		if res.StatusCode == http.StatusNotFound {
			return false, nil

		} else {
			return true, nil
		}

	}
	return true, nil
}

func createdIndex(es *elasticsearch.Client) {

	create, err := es.Indices.Create("auto_link_tracking")
	if err != nil {
		logger.Info("创建索引失败 %s", err.Error())
		return

	}
	logger.Info("Error creating index: %v", create)
	return

}

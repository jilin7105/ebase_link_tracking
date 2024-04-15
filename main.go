package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/jilin7105/ebase"
	"github.com/jilin7105/ebase/kafka/ConsumerAbout"
	"github.com/jilin7105/ebase/logger"
	"github.com/jilin7105/ebase_link_tracking/CheckSelf"
	"os"
)

func Action(msg *sarama.ConsumerMessage, ctx context.Context) error {
	es_log := struct {
		LinkTrackID         string `json:"link_track_id"`          //追踪 id
		LinkTrackParentID   string `json:"link_track_parent_id"`   //父级 id
		LinkTrackSpan       string `json:"link_track_span"`        //事件类型
		LinkTrackDesc       string `json:"link_track_desc"`        //事件描述
		LinkTrackTime       string `json:"link_track_time"`        //触发时间
		LinkTrackActionTime string `json:"link_track_action_time"` //执行时间
		ServerName          string `json:"server_name"`            //服务名称
		ServerType          string `json:"server_type"`            //服务类型
	}{}

	err := json.Unmarshal(msg.Value, &es_log)
	if err != nil {
		return err
	}

	esdb := ebase.GetEsV7("linkTracking")
	if esdb == nil {
		return fmt.Errorf("")
	}

	res, err := esdb.Index("auto_link_tracking", bytes.NewReader(msg.Value), func(request *esapi.IndexRequest) {
		request.DocumentID = es_log.LinkTrackID

	})
	if err != nil {
		logger.Info("res %v, err %v", res, err)
		return err
	}
	logger.Info("res %v, err %v", res, err)

	return nil
}

func setup() error {
	return nil
}

// 使用go run main.go  启动测试服务
func main() {
	path, _ := os.Getwd()

	ebase.SetProjectPath(path)
	ebase.Init()
	eb := ebase.GetEbInstance()
	err := CheckSelf.CheckIndex()
	if err != nil {
		logger.Info("检查索引失败 %s", err.Error())
		return
	}

	eb.EasyRegisterKafkaHandle("linkTrackingKafka", ConsumerAbout.SetActionMessages(Action), ConsumerAbout.SetSetup(setup))

	eb.Run()
	select {}
}

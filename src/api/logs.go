package api

import (
	"encoding/json"
	"sort"
	"time"
)

type HistoricalLog struct {
	Host      string
	Message   string
	Timestamp time.Time
}

type HistoricalLogs []HistoricalLog

// Len is part of sort.Interface.
func (hl HistoricalLogs) Len() int {
	return len(hl)
}

// Swap is part of sort.Interface.
func (hl HistoricalLogs) Swap(i, j int) {
	hl[i], hl[j] = hl[j], hl[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (hl HistoricalLogs) Less(i, j int) bool {
	return hl[i].Timestamp.Before(hl[j].Timestamp)
}

func GetHistoricalLogs(deploymentId string, oauthToken string) (historicalLogs HistoricalLogs, err error) {
	body, err := rest_get(api_url("/deployments/"+deploymentId+"/historical_logs?size=200&sort=desc&grep_o=connection&grep=query"), oauthToken)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	for host, logs := range result {
		for _, log := range logs.(map[string]interface{})["logs"].([]interface{}) {
			ts := log.(map[string]interface{})["ts"].(string)
			timestamp, _ := time.Parse("2006-01-02T15:04:05Z", ts)
			historicalLogs = append(historicalLogs, HistoricalLog{Host: host, Message: log.(map[string]interface{})["message"].(string), Timestamp: timestamp})
		}
	}
	sort.Sort(historicalLogs)
	return historicalLogs, err
}

package works

import (
	"encoding/json"
	"fmt"
	"github.com/dalconoid/kiddy-lp/pkg/storage"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	baseballLineURL = "/api/v1/lines/baseball"
	footballLineURL = "/api/v1/lines/football"
	soccerLineURL = "/api/v1/lines/soccer"

	baseballKey = "baseball"
	footballKey = "football"
	soccerKey = "soccer"
)

// WorkManager work manager
type WorkManager struct {
	storage storage.Storage
	linesProviderURL string
}

// NewWorkManager creates work manager
func NewWorkManager(st storage.Storage, addr string) *WorkManager {
	return &WorkManager{
		storage: st,
		linesProviderURL: addr,
	}
}

// StartWorks starts works in background
func (wm *WorkManager) StartWorks(bt, ft, st float64) {
	go wm.watchBaseballLine(bt)
	go wm.watchFootballLine(ft)
	go wm.watchSoccerLine(st)
}


func (wm *WorkManager) watchBaseballLine(t float64) {
	type sport struct {
		Baseball string `json:"BASEBALL"`
	}
	type line struct {
		Lines *sport
	}
	for {
		data, err := getResponseData(wm.linesProviderURL + baseballLineURL)
		if err != nil {
			log.Error(err)
			continue
		}
		l := &line{}
		if err = json.Unmarshal(data, l); err != nil {
			log.Error(err)
			continue
		}
		k, err := strconv.ParseFloat(l.Lines.Baseball, 32)
		if err != nil {
			log.Error(err)
			continue
		}
		if err = wm.storage.WriteLineRate(k, baseballKey); err != nil {
			log.Error(err)
			continue
		}
		log.Debugf("BASEBALL line updated: k=[%.3f]", k)
		ms := int(1000 * t)
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

func (wm *WorkManager) watchFootballLine(t float64) {
	type sport struct {
		Football string `json:"FOOTBALL"`
	}
	type line struct {
		Lines *sport
	}
	for {
		data, err := getResponseData(wm.linesProviderURL + footballLineURL)
		if err != nil {
			log.Error(err)
			continue
		}
		l := &line{}
		if err = json.Unmarshal(data, l); err != nil {
			log.Error(err)
			continue
		}
		k, err := strconv.ParseFloat(l.Lines.Football, 32)
		if err != nil {
			log.Error(err)
			continue
		}
		if err = wm.storage.WriteLineRate(k, footballKey); err != nil {
			log.Error(err)
			continue
		}
		log.Debugf("FOOTBALL line updated: k=[%.3f]", k)
		ms := int(1000 * t)
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

func (wm *WorkManager) watchSoccerLine(t float64) {
	type sport struct {
		Soccer string `json:"SOCCER"`
	}
	type line struct {
		Lines *sport
	}
	for {
		data, err := getResponseData(wm.linesProviderURL + soccerLineURL)
		if err != nil {
			log.Error(err)
			continue
		}
		l := &line{}
		if err = json.Unmarshal(data, l); err != nil {
			log.Error(err)
			continue
		}
		k, err := strconv.ParseFloat(l.Lines.Soccer, 32)
		if err != nil {
			log.Error(err)
			continue
		}
		if err = wm.storage.WriteLineRate(k, soccerKey); err != nil {
			log.Error(err)
			continue
		}
		log.Debugf("SOCCER line updated: k=[%.3f]", k)
		ms := int(1000 * t)
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

func getResponseData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response: status code %v", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}

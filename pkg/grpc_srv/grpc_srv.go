package grpc_srv

import (
	"fmt"
	"github.com/dalconoid/kiddy-lp/api"
	"github.com/dalconoid/kiddy-lp/pkg/storage"
	"github.com/dalconoid/kiddy-lp/pkg/utils"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

const (
	baseball = "baseball"
	football = "football"
	soccer = "soccer"
)

// GRPCServer is a grpc server
type GRPCServer struct {
	api.UnimplementedKiddyServer
	storage storage.Storage
}

// New creates a new GRPCServer
func New(storage storage.Storage) *GRPCServer {
	return &GRPCServer{
		storage: storage,
	}
}

// SubscribeOnSportsLines is a rpc method
func (s *GRPCServer) SubscribeOnSportsLines(stream api.Kiddy_SubscribeOnSportsLinesServer) error {
	if err := s.storage.CheckConnection(); err != nil {
		return err
	}

	requestChanged := make(chan *api.SubscribeRequest)
	killSender := make(chan int)
	firstMsgReceived := false
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			killSender <- 1
			return nil
		}
		if err != nil {
			killSender <- 1
			return err
		}

		if !firstMsgReceived {
			firstMsgReceived = true
			//start goroutine
			go func() {
				lineRates := map[string]float64{}
				lineDeltas := map[string]float64{}
				sportLines := make([]string, 0, 3)
				var nextMsgTime time.Time
				var waitTime time.Duration

				waitTime = time.Duration(in.Time*1000) * time.Millisecond
				nextMsgTime = time.Now().Add(waitTime)
				sportLines = append(sportLines, in.Lines...)

				fmt.Println("FIRST SEND")
				if err := getLineDeltas(sportLines, lineRates, lineDeltas, s.storage); err != nil {
					log.Errorf("grpc: subscribe on sports lines: sender: %s", err.Error())
					return
				}
				if err := stream.Send(&api.LinesDeltas{LinesDeltas: lineDeltas}); err != nil {
					log.Errorf("grpc: subscribe on sports lines: sender: %s", err.Error())
					return
				}

				for {
					select {
					case  <-killSender:
						fmt.Println("KILL SENDER")
						return
					case r := <-requestChanged:
						fmt.Println("REQUEST CHANGED")
						if !utils.SameStringSlices(sportLines, r.Lines) {
							lineRates = map[string]float64{}
							lineDeltas = map[string]float64{}
							sportLines = make([]string, 0, 3)
							sportLines = append(sportLines, r.Lines...)
						}
						waitTime = time.Duration(r.Time*1000) * time.Millisecond
						nextMsgTime = time.Now().Add(waitTime)
						if err := getLineDeltas(sportLines, lineRates, lineDeltas, s.storage); err != nil {
							log.Errorf("grpc: subscribe on sports lines: sender: %s", err.Error())
							return
						}
						if err := stream.Send(&api.LinesDeltas{LinesDeltas: lineDeltas}); err != nil {
							log.Errorf("grpc: subscribe on sports lines: sender: %s", err.Error())
							return
						}
					default:
						if time.Now().After(nextMsgTime) {
							nextMsgTime = time.Now().Add(waitTime)
							fmt.Printf("Time: %v | Wait: %v | Next message time: %v", time.Now(), waitTime, nextMsgTime)
							if err := getLineDeltas(sportLines, lineRates, lineDeltas, s.storage); err != nil {
								log.Errorf("grpc: subscribe on sports lines: sender: %s", err.Error())
								return
							}
							if err := stream.Send(&api.LinesDeltas{LinesDeltas: lineDeltas}); err != nil {
								log.Errorf("grpc: subscribe on sports lines: sender: %s", err.Error())
								return
							}
						}
					}
				}
			}()
		} else {
			fmt.Println("SEND REQUEST CHANGED SIGNAL")
			requestChanged <- in
		}
	}
}



func getLineDeltas(lines []string, rates, deltas map[string]float64, st storage.Storage) error {
	for _, line := range lines {
		switch line {
		case baseball:
			k, err := st.GetLineRate(baseball)
			if err != nil {
				return err
			}
			pk := rates[baseball]
			deltas[baseball] = utils.RoundF64ToPrecision(k - pk, 3)
			rates[baseball] = k
		case football:
			k, err := st.GetLineRate(football)
			if err != nil {
				return err
			}
			pk := rates[football]
			deltas[football] = utils.RoundF64ToPrecision(k - pk, 3)
			rates[football] = k
		case soccer:
			k, err := st.GetLineRate(soccer)
			if err != nil {
				return err
			}
			pk := rates[soccer]
			deltas[soccer] =  utils.RoundF64ToPrecision(k - pk, 3)
			rates[soccer] = k
		}
	}
	return nil
}

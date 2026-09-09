package realtime

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/infrastructure/eventbus"
	"github.com/nougght/monitoring-system/server/internal/model"
	"github.com/nougght/monitoring-system/server/internal/model/event"
	realtime_model "github.com/nougght/monitoring-system/server/internal/model/realtime"
	dto "github.com/nougght/monitoring-system/server/internal/transport/dto/types"
	"github.com/nougght/monitoring-system/shared/go/util"
)

type RealtimeService struct {
	cfg        *config.Config
	transactor model.Transactor
	hub        *Hub
	upgrader   *websocket.Upgrader
	bus        *eventbus.EventBus
}

func NewRealtimeService(cfg *config.Config,
	transactor model.Transactor,
	bus *eventbus.EventBus,
) (*RealtimeService, error) {
	s := &RealtimeService{
		cfg:        cfg,
		transactor: transactor,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		bus: bus,
	}
	s.hub = NewHub(s.HandleSubscription, s.HandleUnsubscription)
	go s.hub.Run()
	return s, nil
}

func (s *RealtimeService) HandleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade to websocket: ", err)
		return
	}
	s.hub.RegisterClient(conn)
}

func (s *RealtimeService) HandleSubscription(subject string) {
	log.Printf("handle sub %s", subject)
	err := s.bus.Subscribe(subject, func(ctx context.Context, e event.Event) {
		switch e.Name() {
		case "AgentMetricsEvent":
			metricsEvent, ok := e.(*event.AgentMetricsEvent)
			if !ok {
				log.Printf("failed to convert event: %s", e)
			}
			log.Printf("received metrics event for agent %s", metricsEvent.AgentID)

			raw, err := json.Marshal(&dto.SeriesDTO{
				Key:    metricsEvent.Metric.Kind,
				Label:  metricsEvent.Metric.Label,
				Unit:   "",
				Ts:     []int64{},
				Values: []*float64{util.Ptr(metricsEvent.Metric.Value)},
			})
			if err != nil {
				log.Printf("failed to marshal json metric sample")
			}
			s.hub.SendMessage(subject, &realtime_model.Message{
				Type:    realtime_model.MessageTypeSeries,
				AgentID: metricsEvent.AgentID,
				Payload: raw,
			})
		}
	}, 1000, nil)
	if err != nil {
		log.Printf("failed to sub bus: %s", err.Error())
	}
}

func (s *RealtimeService) HandleUnsubscription(subject string) {

}

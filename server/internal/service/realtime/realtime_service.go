package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/infrastructure/eventbus"
	"github.com/nougght/monitoring-system/server/internal/model"
	"github.com/nougght/monitoring-system/server/internal/model/event"
	realtime_model "github.com/nougght/monitoring-system/server/internal/model/realtime"
	"github.com/nougght/monitoring-system/server/internal/service/overview"
	"github.com/nougght/monitoring-system/server/internal/transport/dto/mapper"
	dto "github.com/nougght/monitoring-system/server/internal/transport/dto/types"
	"github.com/nougght/monitoring-system/server/internal/transport/ws"
	"github.com/nougght/monitoring-system/shared/go/util"
)

type RealtimeService struct {
	cfg        *config.Config
	overview   *overview.OverviewService
	transactor model.Transactor
	hub        *ws.Hub
	upgrader   *websocket.Upgrader
	bus        *eventbus.EventBus
	handlers   map[string]event.EventHandler
}

func NewRealtimeService(cfg *config.Config,
	overviewService *overview.OverviewService,
	transactor model.Transactor,
	bus *eventbus.EventBus,
) (*RealtimeService, error) {
	s := &RealtimeService{
		cfg:        cfg,
		overview:   overviewService,
		transactor: transactor,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		bus: bus,
	}
	s.hub = ws.NewHub(s.HandleSubscription, s.HandleUnsubscription)
	s.handlers = map[string]event.EventHandler{
		event.AgentMetricsEventName: s.metricsSeriesEventHandler,
	}
	go s.hub.Run()
	return s, nil
}

func (s *RealtimeService) Run(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				overview, err := s.overview.GetOverview(ctx)
				if err != nil {
					log.Printf("failed to get overview: %s", err.Error())
					continue
				}
				raw, err := json.Marshal(mapper.OverviewToDTO(overview))
				if err != nil {
					log.Printf("failed to marshal json overview: %s", err.Error())
					continue
				}
				s.hub.SendMessage(event.FleetOverviewEventSubjectPrefix, &realtime_model.Message{
					Type:    realtime_model.MessageTypeFleetOverview,
					Payload: raw,
				})
			}

		}
	}()
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

	var name string
	switch {
	case strings.HasPrefix(subject, event.AgentMetricsEventSubjectPrefix):
		name = event.AgentMetricsEventName
	default:
		return
	}
	err := s.bus.Subscribe(subject, s.handlers[name], 1000, nil)
	if err != nil {
		log.Printf("failed to sub bus: %s", err.Error())
	}
}

func (s *RealtimeService) metricsSeriesEventHandler(ctx context.Context, e event.Event) {
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
	s.hub.SendMessage(fmt.Sprintf("%s.%s", event.AgentMetricsEventSubjectPrefix, metricsEvent.AgentID), &realtime_model.Message{
		Type:    realtime_model.MessageTypeSeries,
		AgentID: metricsEvent.AgentID,
		Payload: raw,
	})
}

func (s *RealtimeService) HandleUnsubscription(subject string) {

}

package enrollment_server

import (
	pb "github.com/nougght/monitoring-system/pkg/proto/gen/enrollment/v1"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertEnrollParamsFromProto(p *pb.EnrollRequest) *agent_model.EnrollParams {
	if p == nil {
		return nil
	}
	return &agent_model.EnrollParams{
		EnrollmentKey: p.EnrollmentKey,
		CsrDer:        p.CsrDer,
	}
}

func convertEnrollResultToProto(p *agent_model.EnrollResult) *pb.EnrollResponse {
	if p == nil {
		return nil
	}
	return &pb.EnrollResponse{
		CertificateDer: p.CertDer,
		CaChainDer:     p.CAChainDer,
		NotAfter:       timestamppb.New(p.NotAfter),
	}
}

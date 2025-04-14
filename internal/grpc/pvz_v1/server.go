package pvz_v1

import (
	"context"

	"avito-pvz-service/internal/services"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCServer struct {
	UnimplementedPVZServiceServer
	svc services.PVZService
}

func NewGRPCServer(svc services.PVZService) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func (s *GRPCServer) GetPVZList(ctx context.Context, req *GetPVZListRequest) (*GetPVZListResponse, error) {
	pvzs, err := s.svc.GetPVZList(nil, nil, 1, 1000000000)
	if err != nil {
		return nil, err
	}

	resp := &GetPVZListResponse{}
	for _, pvz := range pvzs {
		resp.Pvzs = append(resp.Pvzs, &PVZ{
			Id:               pvz.PVZ.ID.String(),
			RegistrationDate: timestamppb.New(pvz.PVZ.RegistrationDate),
			City:             pvz.PVZ.City,
		})
	}
	return resp, nil
}

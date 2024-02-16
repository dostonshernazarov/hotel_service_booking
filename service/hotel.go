package service

import (
	"context"
	"github.com/jmoiron/sqlx"
	pb "hotel_service_booking/genproto/hotel_proto"
	l "hotel_service_booking/pkg/logger"
	"hotel_service_booking/storage"
)

// HotelService ...
type HotelService struct {
	storage storage.IStorage
	logger  l.Logger
}

// NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger) *HotelService {
	return &HotelService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}

func (s *HotelService) Create(ctx context.Context, req *pb.Hotel) (*pb.Hotel, error) {
	return nil, nil
}

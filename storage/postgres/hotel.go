package postgres

import (
	"github.com/jmoiron/sqlx"
	pb "hotel_service_booking/genproto/hotel_proto"
)

type hotelRepo struct {
	db *sqlx.DB
}

// NewHotelRepo ...
func NewHotelRepo(db *sqlx.DB) *hotelRepo {
	return &hotelRepo{db: db}
}

func (r *hotelRepo) Create(user *pb.Hotel) (*pb.Hotel, error) {
	return nil, nil
}

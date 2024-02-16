package repo

import (
	pb "hotel_service_booking/genproto/hotel_proto"
)

// HotelStorageI ...
type HotelStorageI interface {
	Create(hotel *pb.Hotel) (*pb.Hotel, error)
}

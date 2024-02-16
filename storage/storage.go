package storage

import (
	"github.com/jmoiron/sqlx"
	"hotel_service_booking/storage/postgres"
	"hotel_service_booking/storage/repo"
)

// IStorage ...
type IStorage interface {
	Hotel() repo.HotelStorageI
}

type storagePg struct {
	db        *sqlx.DB
	hotelRepo repo.HotelStorageI
}

// NewStoragePg ...
func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:        db,
		hotelRepo: postgres.NewHotelRepo(db),
	}
}

func (s storagePg) Hotel() repo.HotelStorageI {
	return s.hotelRepo
}

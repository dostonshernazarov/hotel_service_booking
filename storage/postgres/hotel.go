package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	pb "hotel_service_booking/genproto/hotel_proto"
	"hotel_service_booking/pkg/logger"
	"time"
)

type hotelRepo struct {
	db *sqlx.DB
}

// NewHotelRepo ...
func NewHotelRepo(db *sqlx.DB) *hotelRepo {
	return &hotelRepo{db: db}
}

// Owner CRUD
func (r *hotelRepo) CreateOwner(user *pb.Owner) (*pb.Owner, error) {
	id := uuid.New()
	var full_name, email, password, birthday, imageUrl, phone, role, refreshToken sql.NullString
	var userRes pb.Owner

	query := `INSERT INTO owner (id, full_name, email ,password ,
                  birthday ,phone ,
                  image_url ,role, refresh_token) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id,
                                                  full_name, email, password, birthday, imageUrl,
                                                  phone, role, refreshToken, created_at, updated_at`

	err := r.db.QueryRow(query, id, sql.NullString{String: user.FullName, Valid: user.FullName != ""},
		sql.NullString{String: user.Email, Valid: user.Email != ""},
		sql.NullString{String: user.Password, Valid: user.Password != ""},
		sql.NullString{String: user.Birthday, Valid: user.Birthday != ""},
		sql.NullString{String: user.ImageUrl, Valid: user.ImageUrl != ""},
		sql.NullString{String: user.Phone, Valid: user.Phone != ""},
		sql.NullString{String: user.Role, Valid: user.Role != ""},
		sql.NullString{String: user.RefreshToken, Valid: user.RefreshToken != ""}).Scan(&userRes.Id,
		&full_name, &email,
		&password, &birthday, &imageUrl,
		&phone, &role, &refreshToken, &userRes.CreatedAt, &userRes.UpdatedAt)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if full_name.Valid {
		userRes.FullName = full_name.String
	}

	if email.Valid {
		userRes.Email = email.String
	}
	if password.Valid {
		userRes.Password = password.String
	}
	if birthday.Valid {
		userRes.Birthday = birthday.String
	}
	if imageUrl.Valid {
		userRes.ImageUrl = imageUrl.String
	}

	if phone.Valid {
		userRes.Phone = phone.String
	}
	if role.Valid {
		userRes.Role = role.String
	}
	if refreshToken.Valid {
		userRes.RefreshToken = refreshToken.String
	}
	return &userRes, nil

}

func (r *hotelRepo) GetUserByID(id string) (*pb.Owner, error) {
	var user pb.Owner
	var full_name, email, password, birthday, imageUrl, phone, role, refreshToken sql.NullString

	query := `SELECT id, full_name,
	email, password, birthday, image_url,
	phone, role, refresh_token, created_at, updated_at FROM owner WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, id).Scan(&user.Id, &full_name,
		&email, &password, &birthday, &imageUrl,
		&phone, &role, &refreshToken, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// Convert sql.NullString to regular strings if they are valid
	user.FullName = stringValue(full_name)
	user.Email = stringValue(email)
	user.Password = stringValue(password)
	user.Birthday = stringValue(birthday)
	user.ImageUrl = stringValue(imageUrl)
	user.Phone = stringValue(phone)
	user.Role = stringValue(role)
	user.RefreshToken = stringValue(refreshToken)

	return &user, nil
}

func (r *hotelRepo) GetOwnerByColumnItem(req *pb.GetOwnerByColumnItemReq) (*pb.Owner, error) {
	var user pb.Owner
	var full_name, email, password, birthday, imageUrl, phone, role, refreshToken sql.NullString

	col := "id, full_name, email, password, birthday, image_url,phone, role, refresh_token, created_at, updated_at"
	query := fmt.Sprintf("SELECT %s FROM users WHERE %s = $1 AND deleted_at IS NULL", col, req.Column)
	if err := r.db.QueryRow(
		query,
		req.Item,
	).Scan(&user.Id, &full_name, &email,
		&password, &birthday, &imageUrl,
		&phone, &role, &refreshToken, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	// Convert sql.NullString to regular strings if they are valid
	user.FullName = stringValue(full_name)
	user.Email = stringValue(email)
	user.Password = stringValue(password)
	user.Birthday = stringValue(birthday)
	user.ImageUrl = stringValue(imageUrl)
	user.Phone = stringValue(phone)
	user.Role = stringValue(role)
	user.RefreshToken = stringValue(refreshToken)

	return &user, nil
}

func (r *hotelRepo) GetAllOwners(req *pb.GetAllOwnerReq) ([]*pb.Owner, error) {
	var users []*pb.Owner

	query := `
		SELECT id, full_name, email, password, birthday, image_url,
		       phone, role, refresh_token, created_at, updated_at
		FROM owner
		WHERE deleted_at IS NULL LIMIT $1 OFFSET $2
	`

	offset := req.Limit * (req.Page - 1)
	rows, err := r.db.Query(query, req.Limit, offset)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user pb.Owner
		var full_name, email, password, birthday, imageUrl, phone, role, refreshToken sql.NullString

		err := rows.Scan(&user.Id, &full_name,
			&email, &password, &birthday, &imageUrl,
			&phone, &role, &refreshToken, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		user.FullName = stringValue(full_name)
		user.Email = stringValue(email)
		user.Password = stringValue(password)
		user.Birthday = stringValue(birthday)
		user.ImageUrl = stringValue(imageUrl)
		user.Phone = stringValue(phone)
		user.Role = stringValue(role)
		user.RefreshToken = stringValue(refreshToken)

		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		logger.Error(err)
		return nil, err
	}

	return users, nil
}

func (r *hotelRepo) UpdateUserByID(user *pb.Owner) (*pb.Owner, error) {
	//defer r.db.Close()
	var userReq pb.Owner
	var full_name, email, password, birthday, imageUrl, phone, role, refreshToken sql.NullString

	currentTime := time.Now()

	query := `UPDATE owner SET full_name = $1, email = $2,
                    password = $3, birthday = $4, image_url = $5,
                    phone = $6, role = $7,
                    refresh_token = $8, updated_at = $9 WHERE id = $10 RETURNING id,
                                        full_name, email,
                                        password, birthday, image_url,
                                        phone, role, refresh_token, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		user.FullName,
		user.Email, user.Password, user.Birthday,
		user.ImageUrl, user.Phone,
		user.Role, user.RefreshToken, currentTime, user.Id).Scan(&userReq.Id, &full_name,
		&email, &password, &birthday, &imageUrl,
		&phone, &role, &refreshToken, &userReq.CreatedAt, &userReq.UpdatedAt)
	if err != nil {
		logger.Error(err)
		return nil, err

	}

	// Convert sql.NullString to regular strings if they are valid
	userReq.FullName = stringValue(full_name)
	userReq.Email = stringValue(email)
	userReq.Password = stringValue(password)
	userReq.Birthday = stringValue(birthday)
	userReq.ImageUrl = stringValue(imageUrl)
	userReq.Phone = stringValue(phone)
	userReq.Role = stringValue(role)
	userReq.RefreshToken = stringValue(refreshToken)

	return &userReq, nil

}

func (r *hotelRepo) SoftDeleteUserByID(userid string) error {
	query := "UPDATE owner SET deleted_at = $1 WHERE id = $2"
	currentTime := time.Now()

	result, err := r.db.Exec(query, currentTime, userid)
	if err != nil {
		logger.Error(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error(err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no user found with the provided ID")
	}

	return nil

}

// Hotel CRUD
func (r *hotelRepo) CreateHotel(hotel *pb.Hotel) (*pb.Hotel, error) {
	var country, city, province, address, name, phone, email, license, image_url sql.NullString
	var hotelRes pb.Hotel

	//Query to add items in location
	queryLoc := `INSERT INTO location (id, country,
                     city,
                     province,
                     addres) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) returning city,
                     province, addres, created_at, updated_at`
	idLoc := uuid.New().String()
	err := r.db.QueryRow(queryLoc, idLoc, sql.NullString{String: hotel.Country, Valid: hotel.Country != ""},
		sql.NullString{String: hotel.City, Valid: hotel.City != ""},
		sql.NullString{String: hotel.Province, Valid: hotel.Province != ""},
		sql.NullString{String: hotel.Address, Valid: hotel.Address != ""}).Scan(&country, &city,
		&province, &address)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	queryHotel := `INSERT INTO hotel_info (id, name, phone, email,
					license, image_url, location_id) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING name,
					phone, email, license, image_url, created_at, updated_at`
	idHotel := uuid.New().String()
	err = r.db.QueryRow(queryHotel, idHotel, sql.NullString{String: hotel.Name, Valid: hotel.Name != ""},
		sql.NullString{String: hotel.Phone, Valid: hotel.Phone != ""},
		sql.NullString{String: hotel.Email, Valid: hotel.Email != ""},
		sql.NullString{String: hotel.License, Valid: hotel.License != ""},
		sql.NullString{String: hotel.ImageUrl, Valid: hotel.ImageUrl != ""},
		idLoc).Scan(&name, &phone, &email, &license,
		&image_url, &hotelRes.CreatedAt, &hotelRes.UpdatedAt)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	hotelRes.Id = idHotel
	hotelRes.Name = stringValue(name)
	hotelRes.Phone = stringValue(phone)
	hotelRes.Email = stringValue(email)
	hotelRes.License = stringValue(license)
	hotelRes.ImageUrl = stringValue(image_url)
	hotelRes.Country = stringValue(country)
	hotelRes.City = stringValue(city)
	hotelRes.Province = stringValue(province)
	hotelRes.Address = stringValue(address)

	return &hotelRes, nil
}

func (r *hotelRepo) GetHotelByID(hotelId string) (*pb.Hotel, error) {

}

// stringValue returns the string value of a sql.NullString, handling null values.
func stringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

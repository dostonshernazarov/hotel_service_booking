CREATE TABLE owner (id UUID primary key,
                    full_name varchar(255),
                    email varchar(255),
                    password varchar(255),
                    birthday varchar(255),
                    phone varchar(255),
                    image_url varchar(255),
                    role varchar(255),
                    refresh_token varchar(255),
                    created_at timestamp default current_timestamp,
                    updated_at timestamp default current_timestamp,
                    deleted_at timestamp);

CREATE TABLE hotel_info (id UUID primary key,
                         name varchar(255),
                         phone varchar(255),
                         email varchar(255),
                         license varchar(255),
                         image_url varchar(255),
                         country varchar(255),
                         city varchar(255),
                         province varchar(255),
                         address varchar(255),
                         owner_id UUID references owner(id) NOT NULL,
                         created_at timestamp default current_timestamp,
                         updated_at timestamp default current_timestamp,
                         deleted_at timestamp);

CREATE TABLE room (id UUID primary key,
                   hotel_id uuid references hotel_info(id) NOT NULL ,
                   price varchar(255),
                   description varchar(255),
                   holidays varchar(255),
                   free_days varchar(255),
                   number_of_rooms int,
                   discount varchar(255),
                   created_at timestamp default current_timestamp,
                   updated_at timestamp default current_timestamp,
                   deleted_at timestamp);




CREATE TABLE image_room (id UUID primary key,
                         room_id UUID references room(id),
                         image_url varchar(255));

CREATE TABLE review (id UUID primary key,
                     comment varchar(500),
                     stars int,
                     user_id UUID NOT NULL ,
                     room_id UUID references room(id));



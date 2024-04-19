CREATE EXTENSION  IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
   uuid   UUID  Primary Key,
   email varchar(255) NOT NULL Unique,
   username varchar(255) NOT NULL,
   phone varchar(255) NOT NULL Unique,
   dateOfBirth varchar(255) NOT null,
   pass_hash TEXT NOT NULL
);
Create index if not exists idx_email on users(email);
Create index if not exists idx_username on users(username);
Create index if not exists idx_phone on users(phone);


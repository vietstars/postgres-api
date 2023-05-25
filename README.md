# postgres-go-api

go mod init github.com/vietstars/postgres-api

go mod tidy

go get -u gorm.io/gorm

go get -u gorm.io/driver/postgres

go get github.com/go-playground/validator/v10

go get github.com/joho/godotenv

brew install postgresql 

postgres -V 

brew services stop postgresql@14

brew services start postgresql@14

createdb dev_master 

psql dev_master  

CREATE ROLE admin WITH LOGIN SUPERUSER PASSWORD 'master'

\l list of users

\du list of DB


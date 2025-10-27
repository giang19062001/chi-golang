# init
go mod init github.com/giang19062001/chi-golang 

# install
# go install = npm i -g   ( ko thêm dependenci vào go.mod )
# go get = npm i  (thêm dependenci vào go.mod )

go get github.com/joho/godotenv
go get github.com/go-chi/chi
go get github.com/go-chi/cors
go get github.com/google/uuid
# postgres sql
 go get github.com/lib/pq
# tool generate code Go từ SQL
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# quản lý migration database
go install github.com/pressly/goose/v3/cmd/goose@latest 
# build
go build

# run
./chi-golang.exe
 
# build & run
# ubuntu/mac
go build && ./chi-golang.exe
# win
go build; ./chi-golang.exe 


# delete
# xóa các thư viện ko được import
go mod tidy

# check version
vd: sqlc version
vd: goose -version


# query truy vấn bằng goose
cd \sql\schema> 
# chạy các query trong -- +goose Up
goose postgres postgres://postgres:123@localhost:5432/rssagg up
# chạy các query trong -- +goose Down
goose postgres postgres://postgres:123@localhost:5432/rssagg down

# tự động generate code trong go từ câu query
sqlc generate
mkdir myapi
cd myapi

go mod init myapi

go get -u github.com/gorilla/mux
go get -u github.com/go-sql-driver/mysql

go run main.go
from golang:latest
workdir /app
copy go.mod go.sum ./
run go mod download
copy . .
run go build -o main
expose 8020
run chmod 777 start-server.sh
cmd ["./start-server.sh"]
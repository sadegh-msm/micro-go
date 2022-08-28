package main

import (
	"database/sql"
	"micro-go/authontication-service/cmd/data"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

}

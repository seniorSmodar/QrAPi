package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Введите url базы")
	}

	return os.Getenv("MONGOURI")
}

func EnvPort() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Введите порт сервера")
	}

	return os.Getenv("PORT")
}

 

func EnvDuration() int {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Введите duration")
	}

	value := os.Getenv("CINEMADURATION")
	i, err := strconv.Atoi(value)
	if err != nil{
		return 0
	}
	log.Fatal(i)
	return i
}
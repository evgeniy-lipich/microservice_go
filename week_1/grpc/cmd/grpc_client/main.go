package main

import (
	"context"
	"log"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	desc "github.com/evgeniy-lipich/microservice_go/week_1/grpc/pkg/note_v1"
)

const (
	address = "localhost:50051"
	noteId  = 12
)

func main() {
	// установить соединение
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	// в контексте задаем ограничение по таймауту
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// создать клиент Note, передаем соединение
	client := desc.NewNoteV1Client(conn)

	// get запрос
	getResponse, err := client.Get(ctx, &desc.GetRequest{Id: noteId})
	if err != nil {
		log.Fatalf("failed to get note by id: %v", err)
	}
	log.Printf(color.RedString("Note info: \n"), color.GreenString("%+v", getResponse.GetNote()))
}

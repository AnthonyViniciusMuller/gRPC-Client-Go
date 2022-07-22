package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	proto "menssenger/client/protos"

	"google.golang.org/grpc"
)

func main() {
	fmt.Print("\033[H\033[2J")

	fmt.Print("Go Client \n\n\n")

	connection, err := grpc.Dial(":4444", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	client := proto.NewChatClient(connection)
	stream, err := client.GetMessages(context.Background(), &proto.Void{})
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	done := make(chan bool)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				fmt.Printf("cannot receive %v", err)
			}
			fmt.Println("\n" + resp.User + ": " + resp.Message + "\n")
			fmt.Print("> ")

		}
	}()

	var userName string

	fmt.Print("Qual é o seu usuário: ")
	fmt.Scanln(&userName)
	fmt.Print("> ")

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		message := scanner.Text()

		_, err := client.SendMessage(context.Background(), &proto.Message{Message: message, User: userName})
		if err != nil {
			fmt.Print(err)
		}
	}

	<-done
	log.Printf("finished")
}

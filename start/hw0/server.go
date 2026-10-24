package main

import (
	"fmt"
	"net"
)

type Server struct {
	Sk Socket
}

func (s* Server) CreateServer() error{
	listener, err := net.Listen(s.Sk.Protocol, s.Sk.Address)
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Printf("Сервер запущен на %s://%s", s.Sk.Protocol, s.Sk.Address)
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Ошибка при приеме соединения: %v", err)
			return err
		}
		go handleClient(connection)
	}
}

func handleClient (connection net.Conn) {
	defer connection.Close()
	_, err := connection.Write([]byte("OK\n"))
	if err != nil {
		fmt.Printf("Ошибка отправки клиенту: %v\n", err)
		return
	}
}

func main() {
    sk := Socket{
        Protocol: "tcp",
        Address:  "localhost:8080",
    }

    srv := Server{Sk: sk}
    err := srv.CreateServer()
    if err != nil {
        panic(err)
    }
}



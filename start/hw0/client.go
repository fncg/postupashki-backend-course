package main

import (
	"bytes"
	"fmt"
	"net"
)

type Client struct{
	Protocol string 
	Address string
}

func (c *Client) ConnectToServer()(net.Conn, error){
	connection, err := net.Dial(c.Protocol, c.Address)
	if err != nil {
		fmt.Printf("Ошибка подключения к серверу: %v", err)
		return nil, err
	}
	return connection, err
}

func Waiter(connection net.Conn){
	defer connection.Close()
	msg := []byte("OK\n")
	response := make([]byte, 32)
	for {
		n, err := connection.Read(response)
		if err != nil {
			fmt.Printf("Ошибка чтения с сервера")
			break
		}
		
		if bytes.Equal(msg, response[:n]){
			fmt.Println("Получен ответ: 'OK'")
			break
		}
	}
}

func main() {
	c := Client{
		Protocol: "tcp",
		Address:  "localhost:8080",
	}

	fmt.Println("Попытка подключения к серверу...")
	conn, err := c.ConnectToServer()
	if err != nil {
		return
	}
	defer conn.Close() 

	fmt.Println("Успешно подключено!")
}
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	expectedMessage = "OK\n"
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

func Waiter(connection net.Conn) (error) {
	defer connection.Close()
	reader := bufio.NewReader(connection)
	var response []string
	for {
		message, err := reader.ReadString('\n')
		if message != ""{
			response = append(response, message)
		}

		if err != nil{
			if errors.Is(err, io.EOF){
				break
			}
			return err
		}
	}
	if (len(response) != 1 || response[0] != expectedMessage){
		return errors.New("Невереное содержание ответа")
	}

	return nil
}

func main() {
	c := Client{
		Protocol: "tcp",
		Address:  "localhost:8080",
	}

	fmt.Println("Попытка подключения к серверу...")
	conn, err := c.ConnectToServer()
	if err != nil {
		os.Exit(1)
	}

	e := Waiter(conn)
	if e != nil {
		fmt.Printf("Ошибка ответа с сервера: %v\n", e)
		os.Exit(1)
	}
	defer conn.Close() 
	fmt.Println("Успешно подключено!")
}
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "Таймаут подключения")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		log.Fatalf("Не удалось подключиться к %s: %v", address, err)
	}
	defer conn.Close()
	fmt.Printf("Подключено к %s\n", address)

	done := make(chan struct{})

	// Перенаправление STDIN в соединение
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		_, err := io.Copy(conn, os.Stdin)
		if err != nil {
			log.Printf("Ошибка при чтении из STDIN: %v", err)
		}
	}()

	// Перенаправление соединения в STDOUT
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			log.Printf("Ошибка при чтении из соединения: %v", err)
		}
	}()

	<-done
	fmt.Println("\nСоединение закрыто.")
}

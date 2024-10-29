package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"time"
)

func main() {
	exactTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка получения точного времени: %v\n", err)
		os.Exit(1)
	}

	localTime := time.Now()

	fmt.Printf("Текущее локальное время - %s\n", localTime.Format("15-04-05"))
	fmt.Printf("Текущее время по NTP - %s\n", exactTime.Format("15-04-05"))
}

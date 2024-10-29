package main

import "fmt"

// Интерфейс команды
type Command interface {
	Execute()
}

// Получатель
type Light struct {
	IsOn bool
}

func (l *Light) On() {
	l.IsOn = true
	fmt.Println("Свет включен")
}

func (l *Light) Off() {
	l.IsOn = false
	fmt.Println("Свет выключен")
}

// Конкретные команды
type TurnOnCommand struct {
	light *Light
}

func (c *TurnOnCommand) Execute() {
	c.light.On()
}

type TurnOffCommand struct {
	light *Light
}

func (c *TurnOffCommand) Execute() {
	c.light.Off()
}

// Отправитель
type RemoteControl struct {
	command Command
}

func (r *RemoteControl) SetCommand(c Command) {
	r.command = c
}

func (r *RemoteControl) PressButton() {
	r.command.Execute()
}

func main() {
	light := &Light{}
	onCommand := &TurnOnCommand{light}
	offCommand := &TurnOffCommand{light}

	remote := &RemoteControl{}

	remote.SetCommand(onCommand)
	remote.PressButton()

	remote.SetCommand(offCommand)
	remote.PressButton()
}

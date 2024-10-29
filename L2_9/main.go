package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Go Shell! Type \\quit to exit.")

	for {
		// Отображаем приглашение
		fmt.Print("go-shell> ")

		// Читаем ввод от пользователя
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения ввода:", err)
			continue
		}

		// Убираем лишние пробелы и символы новой строки
		input = strings.TrimSpace(input)

		// Обрабатываем команду выхода
		if input == "\\quit" {
			fmt.Println("Завершение работы Go Shell.")
			break
		}

		// Разбиваем ввод на части для определения команд и аргументов
		args := strings.Split(input, " | ")

		// Если команда содержит пайп, обрабатываем пайплайн
		if len(args) > 1 {
			handlePipes(args)
		} else {
			// Выполняем одиночную команду
			handleCommand(strings.Fields(args[0]))
		}
	}
}

// Обработчик одиночной команды
func handleCommand(args []string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "cd":
		changeDirectory(args[1:])
	case "pwd":
		printWorkingDirectory()
	case "echo":
		echo(args[1:])
	case "kill":
		killProcess(args[1:])
	case "ps":
		printProcesses()
	default:
		executeExternalCommand(args)
	}
}

// Команда cd: смена директории
func changeDirectory(args []string) {
	if len(args) < 1 {
		fmt.Println("cd: требуется аргумент")
		return
	}
	err := os.Chdir(args[0])
	if err != nil {
		fmt.Println("cd:", err)
	}
}

// Команда pwd: вывод текущего пути
func printWorkingDirectory() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("pwd:", err)
		return
	}
	fmt.Println(dir)
}

// Команда echo: вывод аргументов
func echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

// Команда kill: завершение процесса
func killProcess(args []string) {
	if len(args) < 1 {
		fmt.Println("kill: требуется PID процесса")
		return
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("kill: неверный PID:", args[0])
		return
	}

	err = syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		fmt.Println("kill: ошибка при завершении процесса:", err)
	}
}

// Команда ps: вывод информации о процессах
func printProcesses() {
	cmd := exec.Command("ps", "-e")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("ps:", err)
	}
}

// Выполнение внешних команд
func executeExternalCommand(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Ошибка выполнения команды:", err)
	}
}

// Обработка команд с пайпами
func handlePipes(commands []string) {
	var cmds []*exec.Cmd

	// Создаем команды для пайпа
	for _, cmdStr := range commands {
		args := strings.Fields(cmdStr)
		if len(args) > 0 {
			cmds = append(cmds, exec.Command(args[0], args[1:]...))
		}
	}

	// Соединяем пайпы между собой
	for i := 0; i < len(cmds)-1; i++ {
		pipe, err := cmds[i].StdoutPipe()
		if err != nil {
			fmt.Println("Ошибка при создании пайпа:", err)
			return
		}
		cmds[i+1].Stdin = pipe
	}

	// Устанавливаем вывод последней команды на STDOUT
	cmds[len(cmds)-1].Stdout = os.Stdout

	// Запускаем команды в пайпе
	for _, cmd := range cmds {
		err := cmd.Start()
		if err != nil {
			fmt.Println("Ошибка запуска команды:", err)
			return
		}
	}

	// Ожидаем завершения всех команд
	for _, cmd := range cmds {
		err := cmd.Wait()
		if err != nil {
			fmt.Println("Ошибка выполнения команды:", err)
		}
	}
}

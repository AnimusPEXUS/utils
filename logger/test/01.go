package main

import (
	"bufio"
	"os"

	"github.com/AnimusPEXUS/utils/logger"
	// "time"
)

func main() {

	l := logger.New()

	l.AddOutput(os.Stdout)

	l.Text("test text")
	l.Info("test info")
	l.Warning("test warning")
	l.Error("test error")

	in_reader := bufio.NewReader(os.Stdin)

	in_reader.ReadLine()

}

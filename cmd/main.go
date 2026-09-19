package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/elitracy/void-exchange/pkg/gamestate"
	"github.com/elitracy/void-exchange/pkg/scenario"
)

func main() {
	var listScenariosFlag = flag.Bool("ls", false, "list scenarios")
	var scenarioFlag = flag.String("s", "", "the name of the scenario to run")

	flag.Parse()

	gs := gamestate.NewGameState(0)

	if *listScenariosFlag {
		path := "./scenarios"
		files, err := os.ReadDir(path)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Scenarios:")

		for _, file := range files {
			name := strings.Split(file.Name(), ".")
			fmt.Printf("%s", name[0])
		}

		fmt.Println()
	}

	if *scenarioFlag != "" {
		path := fmt.Sprintf("./scenarios/%s.json", *scenarioFlag)

		file, err := os.Open(path)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}

		cfg, err := scenario.Load(file)

		if err != nil {
			fmt.Println("Error:", err.Error())
			return
		}
		scenario.Build(cfg, gs)

	}

}

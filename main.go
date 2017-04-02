package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

type Shell struct {
	colibri       *Colibri
	currentFlower *Flower
	rl            *readline.Instance
}

func NewShell() *Shell {
	return &Shell{
		colibri: NewColibri(),
	}
}

func (shell *Shell) buildConfig() *readline.Config {
	return &readline.Config{
		Prompt: 	shell.buildPrompt(),
		AutoComplete:	shell.buildCompleter(),
	}
}

func (shell *Shell) buildPrompt() string {
	// let's create the fancy prompt
	promptSuffix := " "
	if shell.currentFlower != nil {
		promptSuffix += fmt.Sprintf("(%s~%s) ",
			color.CyanString(shell.currentFlower.Container.ShortID),
			color.YellowString(shell.currentFlower.Container.Name))
	}

	return fmt.Sprintf("%s%s", color.MagentaString("\\/°-"), promptSuffix)
}

func (shell *Shell) buildCompleter() *readline.PrefixCompleter {
	completer := readline.NewPrefixCompleter(
		readline.PcItem("fly-to",
			readline.PcItemDynamic(completeIdentifiers(shell.colibri))),
		readline.PcItem("ps"),
		readline.PcItem("clear"),
		readline.PcItem("exit"),
		readline.PcItem("help"),
	)

	// TODO add flower commands
	//completer.SetChildren()

	return completer
}

func completeIdentifiers(colibri *Colibri) func(string) []string {
	return func(line string) []string {
		colibri.Refresh()
		return append(colibri.ListNames(), colibri.ListShortIDs()...)
	}
}

func (shell *Shell) cmdFlyTo(identifier string) {
	flower := shell.colibri.GetFlower(identifier)

	switch flower {
	case nil:
		fmt.Fprintln(os.Stderr, color.RedString("Unknown container: is it a flower?"))
	default:
		// TODO parse flower

		shell.currentFlower = flower
		shell.Build()
	}
}

func (shell *Shell) cmdClear() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func (shell *Shell) cmdExit() {
	switch shell.currentFlower {
	case nil:
		fmt.Fprintln(os.Stdout, "Bye!")
		os.Exit(0)
	default:
		shell.currentFlower = nil
		shell.Build()
	}
}

func (shell *Shell) cmdHelp() {
	fmt.Fprintln(os.Stdout, "USAGE")
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintf(w, "%s\t\t\t%s\n", "fly-to CONTAINER", "Loads commands from CONTAINER (id or name).")
	fmt.Fprintf(w, "%s\t\t\t%s\n", "ps", "Shows containers which are exposing commands.")
	fmt.Fprintf(w, "%s\t\t\t%s\n", "clear", "Clears the screen.")
	fmt.Fprintf(w, "%s\t\t\t%s\n", "exit", "Quits the current flower or the application.")
	fmt.Fprintf(w, "%s\t\t\t%s\n", "help", "Shows this information.")
	w.Flush()
}

func (shell *Shell) Build() {
	switch shell.rl {
	case nil:
		rl, err := readline.NewEx(shell.buildConfig())
		if err != nil {
			panic(err)
		}

		shell.rl = rl
	default:
		shell.rl.SetConfig(shell.buildConfig())
	}
}

func (shell *Shell) Start() {
	defer shell.rl.Close()

	for {
		line, err := shell.rl.Readline()
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		args := strings.Fields(line)
		name := args[0]
		switch name {
		case  "fly-to":
			if len(args) > 1 {
				shell.cmdFlyTo(args[1])
			}
		case  "ps":
			println("List of flowers")
		case "clear":
			shell.cmdClear()
		case "exit":
			shell.cmdExit()
		case "help":
			shell.cmdHelp()
		}
	}
}

func main() {
	//shell := NewShell()
	//shell.Build()
	//shell.Start()
	colibri := NewColibri()
	colibri.Refresh()
	flower := colibri.GetFlower("examples_flower_1_1")
	_, err := flower.Parse()
	if err != nil {
		panic(err)
	}


}
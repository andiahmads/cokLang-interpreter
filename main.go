package main

import (
	"fmt"
	"go-intepreter/repl"
	"os"
	"os/user"
)

const LOGO = `
_________  ________   ____  __. .____       _____    _______    ________ 
\_   ___ \ \_____  \ |    |/ _| |    |     /  _  \   \      \  /  _____/ 
/    \  \/  /   |   \|      <   |    |    /  /_\  \  /   |   \/   \  ___ 
\     \____/    |    \    |  \  |    |___/    |    \/    |    \    \_\  \
 \______  /\_______  /____|__ \ |_______ \____|__  /\____|__  /\______  /
        \/         \/        \/         \/       \/         \/        \/ 
`

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", LOGO)

	fmt.Printf("hello %s! welcome to  REPL The COKLang v1.0.0\n", user.Username)
	fmt.Printf("Feel free to type commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

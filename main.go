package main

// Copyright (c) 2018, Arm Limited and affiliates.
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"flag"
	"fmt"
	"os"

	. "github.com/WigWagCo/maestro-shell/shell"
	"github.com/c-bata/go-prompt"
)

// func executor(t string) {
// 	if t == "bash" {
// 		cmd := exec.Command("bash")
// 		cmd.Stdin = os.Stdin
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		cmd.Run()
// 	}
// 	return
// }

// func completer(t prompt.Document) []prompt.Suggest {
// 	return []prompt.Suggest{
//         {Text: "bash"},
//         {}
// 	}
// }

var version = "0.0.1"
var defaultSock = "/tmp/maestroapi.sock"
var helpString = `
maestro shell        ver %s
--sock [socket]      Use the given socket instead of the default %s.
`

func main() {
	help := flag.Bool("h", false, "print help & options")
	sockSet := flag.String("s", defaultSock, "Use maestro socket [path]")
	flag.Parse()

	if *help {
		fmt.Printf(helpString, version, defaultSock)
		os.Exit(0)
	}

	client, err := NewUnixClient(*sockSet)

	if err != nil {
		fmt.Errorf("Failed to create connection to maestro: %s", err.Error())
		os.Exit(1)
	}

	SetDefaultClient(client)

	p := prompt.New(
		Executor,
		Completer,
	)
	p.Run()
}

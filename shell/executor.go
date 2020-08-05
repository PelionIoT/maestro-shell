package shell

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
	"errors"
	"fmt"
	"os"
	"strings"
)

// ConsoleOut Dump to console
func ConsoleOut(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	fmt.Printf("> %s\n", s)
}

// Successf prints success events to user
func Successf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func SuccessOut(format string, a ...interface{}) {
	s := Successf(format, a...)
	ConsoleOut("%s", s)
}

// Errorf prints errors to user
func Errorf(format string, a ...interface{}) string {
	s := fmt.Sprintf(format, a...)
	return fmt.Sprintf("[ERROR] %s", s)
}

func ErrorOut(format string, a ...interface{}) {
	s := Errorf(format, a...)
	ConsoleOut("%s", s)
}

// Eventf prints events to user
func Eventf(category string, format string, a ...interface{}) string {
	s := fmt.Sprintf(format, a...)
	return fmt.Sprintf("[event(%s)] %s", category, s)
}

func EventOut(cat string, format string, a ...interface{}) {
	s := Eventf(cat, format, a...)
	ConsoleOut("%s", s)
}

// DebugOut prints debug output to user
func DebugOut(format string, a ...interface{}) {
	if debug_on {
		s := fmt.Sprintf(format, a...)
		fmt.Printf("[DEBUG] %s\n", s)
	}
}

// Command is the function definition for any command
type Command func([]string) (out string, err error)

func cmdExit(args []string) (out string, err error) {
	os.Exit(0)
	return
}

var errors_no_client error
var errors_not_implemented error
var debug_on bool

func init() {
	errors_no_client = errors.New("Maestro could not connect.")
	errors_not_implemented = errors.New("Not implemented yet")
	debug_on = false
}

func cmdGetAlive(args []string) (out string, err error) {
	if defaultClient != nil {
		res, err2 := defaultClient.GetAlive()
		DebugOut("getAlive:%+v %+v", res, err2)
		if err2 == nil {
			out = Successf("Maestro Up. Uptime = %d.%ds\n", res.Uptime/1000000000, res.Uptime%1000000000)
		} else {
			err = err2
		}
	} else {
		err = errors_no_client
	}
	return
}

func cmdDebug(args []string) (out string, err error) {
	if len(args) > 1 {
		if args[1] == "on" {
			debug_on = true
		} else if args[1] == "off" {
			debug_on = false
		} else {
			err = errors.New("Must be on/off")
		}
	}
	s := "off"
	if debug_on {
		s = "on"
	}
	out = Successf("Debug is %s", s)
	return
}

func cmdLog(args []string) (out string, err error) {
	if len(args) > 1 {
		cmd, ok := logCommands[args[1]]
		if ok {
			out, err := cmd(args)
			if err != nil {
				//			DebugOut("here1")
				fmt.Println(Errorf("%s", err.Error()))
			} else {
				//			DebugOut("here")
				fmt.Println(out)
			}
		} else {
			fmt.Printf("%s\n", Errorf("no command: net %s", args[1]))
		}
	} else {
		fmt.Printf("%s\n", Errorf("net: not enough args"))
	}
	return
}

func cmdNet(args []string) (out string, err error) {
	if len(args) > 1 {
		cmd, ok := netCommands[args[1]]
		if ok {
			out, err := cmd(args)
			if err != nil {
				//			DebugOut("here1")
				fmt.Println(Errorf("%s", err.Error()))
			} else {
				//			DebugOut("here")
				fmt.Println(out)
			}
		} else {
			fmt.Printf("%s\n", Errorf("no command: net %s", args[1]))
		}
	} else {
		fmt.Printf("%s\n", Errorf("net: not enough args"))
	}
	return
}

func cmdJobs(args []string) (out string, err error) {
	if len(args) > 1 {
		cmd, ok := jobsCommands[args[1]]
		if ok {
			out, err := cmd(args)
			if err != nil {
				//			DebugOut("here1")
				fmt.Println(Errorf("%s", err.Error()))
			} else {
				//			DebugOut("here")
				fmt.Println(out)
			}
		} else {
			fmt.Printf("%s\n", Errorf("no command: jobs %s", args[1]))
		}
	} else {
		fmt.Printf("%s\n", Errorf("jobs: not enough args"))
	}
	return
}

func netConfigInterface(args []string) (out string, err error) {
	if defaultClient != nil {
		res, err2 := defaultClient.ConfigNetInterface(args)
		DebugOut("net ConfigInterface:%+v %+v", res, err2)
		if err2 == nil {
			out = Successf("%v", res)
		} else {
			err = err2
		}
	} else {
		err = errors_no_client
	}
	return
}

func logSet(args []string) (out string, err error) {
	if defaultClient != nil {
		res, err2 := defaultClient.SetLogging(args)
		DebugOut("log set:%+v %+v", res, err2)
		if err2 == nil {
			out = Successf("%v", res)
		} else {
			err = err2
		}
	} else {
		err = errors_no_client
	}
	return
}

func logGet(args []string) (out string, err error) {
	if defaultClient != nil {
		res, err2 := defaultClient.GetLogging()
		DebugOut("logging:%+v %+v", res, err2)
		if err2 == nil {
			out = Successf("%v", res)
		} else {
			err = err2
		}
	} else {
		err = errors_no_client
	}
	return
}

func logDelete(args []string) (out string, err error) {
	if defaultClient != nil {
		res, err2 := defaultClient.DeleteLogging(args)
		DebugOut("log delete:%+v %+v", res, err2)
		if err2 == nil {
			out = Successf("%v", res)
		} else {
			err = err2
		}
	} else {
		err = errors_no_client
	}
	return
}

func dnsAdd(args []string) (out string, err error) {
	if defaultClient != nil {
		res, err2 := defaultClient.AddDNS(args)
		DebugOut("dns add:%+v %+v", res, err2)
		if err2 == nil {
			out = Successf("%v", res)
		} else {
			err = err2
		}
	} else {
		err = errors_no_client
	}
	return
}

func dnsGet(args []string) (out string, err error) {
	if defaultClient != nil {
		res, err2 := defaultClient.GetDNS()
		DebugOut("dns get:%+v %+v", res, err2)
		if err2 == nil {
			out = Successf("%v", res)
		} else {
			err = err2
		}
	} else {
		err = errors_no_client
	}
	return
}

func dnsDelete(args []string) (out string, err error) {
	if defaultClient != nil {
		res, err2 := defaultClient.DeleteDNS(args)
		DebugOut("dns delete:%+v %+v", res, err2)
		if err2 == nil {
			out = Successf("%v", res)
		} else {
			err = err2
		}
	} else {
		err = errors_no_client
	}
	return
}

func netGetInterfaces(args []string) (out string, err error) {
	if defaultClient != nil {
		res, err2 := defaultClient.GetNetInterfaces()
		DebugOut("net getInterfaces:%+v %+v", res, err2)
		if err2 == nil {
			out = Successf("%v", res)
		} else {
			err = err2
		}
	} else {
		err = errors_no_client
	}
	return
}

func netEvents(args []string) (out string, err error) {
	if defaultClient != nil {
		res, err2 := defaultClient.SubscribeToNetEvents()
		DebugOut("net SubscribeToNetEvents: %+v %+v", res, err2)
		if err2 == nil {
			out = Successf("%v", res)
		} else {
			err = err2
		}
	} else {
		err = errors_no_client
	}
	return
}

func jobsGet(args []string) (out string, err error) {
	if defaultClient != nil {
		res, err2 := defaultClient.GetAllJobStatus()
		DebugOut("net JobsGetAll:%+v %+v", res, err2)
		if err2 == nil {
			out = Successf("%v", res)
		} else {
			err = err2
		}
	} else {
		err = errors_no_client
	}
	return
}

func notImplemented(args []string) (out string, err error) {
	if defaultClient != nil {
		err = errors_not_implemented
		//		fmt.Printf("%s\n", Errorf("Not implemented yet: ", args[1]))
	} else {
		err = errors_no_client
	}
	return
}

var commandMap = map[string]Command{
	"exit":  cmdExit,
	"alive": cmdGetAlive,
	"log":   cmdLog,
	"net":   cmdNet,
	"jobs":  cmdJobs,
	"debug": cmdDebug,
	"help":  GetCommandsHelpString,
}

var logCommands = map[string]Command{
	"get":    logGet,
	"set":    logSet,
	"delete": logDelete,
	"help":   GetLogSubcommandsHelpString,
}

var netCommands = map[string]Command{
	"get-interfaces":   netGetInterfaces,
	"events":           netEvents,
	"config-interface": netConfigInterface,
	"get-dns":          dnsGet,
	"add-dns":          dnsAdd,
	"delete-dns":       dnsDelete,
	"help":             GetNetSubcommandsHelpString,
}

var jobsCommands = map[string]Command{
	"get":      jobsGet,
	"start":    notImplemented, // jobsStart,
	"stop":     notImplemented, // jobsStop,
	"register": notImplemented, // jobsRegister,
}

func Executor(t string) {
	argz := strings.Split(t, " ")
	if len(argz) > 0 {
		cmd, ok := commandMap[argz[0]]
		if ok {
			out, err := cmd(argz)
			if err != nil {
				//				DebugOut("here1")
				fmt.Println(Errorf("%s", err.Error()))
			} else {
				//				DebugOut("here2")
				fmt.Println(out)
			}
		} else {
			fmt.Printf("%s\n", Errorf("no command: %s", argz[0]))
		}
	}
	//	fmt.Printf("Unhandled <%s>\n", t)
}

var defaultClient *MaestroClient

func SetDefaultClient(client *MaestroClient) {
	defaultClient = client
}

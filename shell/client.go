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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/armPelionEdge/maestroSpecs"
)

type MaestroClient struct {
	// len non zero based on transport we use
	unixPath string
	url      string
	httpc    http.Client
	tr       *http.Transport // used to cancel requests
	// see: https://stackoverflow.com/questions/29197685/how-to-close-abort-a-golang-http-client-post-prematurely
	connected bool

	netEventsSubscribeID string
	//	netEventsIntervalSeconds time.Duration
	netEventsListenerRunning bool
}

const (
	defaultHttpTimeoutSeconds = 30 // seconds
)

func (self *MaestroClient) post() {

}

func (self *MaestroClient) get(uri string) (resp *http.Response, err error) {
	resp, err = self.httpc.Get("http://unix" + uri)
	if err != nil {
		return
	}
	DebugOut("http resp:%+v", resp)
	if resp == nil {
		err = errors.New("nil response.")
	}
	return
}

func (self *MaestroClient) put(uri string, body []byte) (resp *http.Response, err error) {
	req, err2 := http.NewRequest(http.MethodPut, "http://unix"+uri, bytes.NewReader(body))
	if err2 != nil {
		err = err2
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err = self.httpc.Do(req)
	if err != nil {
		return
	}
	DebugOut("http resp:%+v", resp)
	if resp == nil {
		err = errors.New("nil response.")
	}
	return
}

func NewUnixClient(path string) (ret *MaestroClient, err error) {
	ret = new(MaestroClient)
	// ret.netEventsIntervalSeconds = time.Duration(defaultNetEventsListenTimeoutSeconds) * time.Second
	DebugOut("creating client on UNIX sock: %s", path)

	ret.tr = &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", path)
		},
	}
	ret.httpc = http.Client{
		Transport: ret.tr,
	}

	return
}

type AliveResponse struct {
	Ok     bool
	Uptime int64
}

func (self *MaestroClient) GetAlive() (alive *AliveResponse, err error) {
	resp, err := self.get("/alive")
	if err == nil {
		DebugOut("resp.Body = %+v", resp.Body)
		body, err2 := ioutil.ReadAll(resp.Body)
		DebugOut("resp.Body body = %+v", body)
		DebugOut("resp.Body body = %s", string(body))
		if err2 == nil {
			alive = &AliveResponse{}
			json.Unmarshal(body, alive)
		} else {
			DebugOut("Error on ReadAll %s", err2.Error())
			err = err2
		}
	}
	return
}

func FormatJsonEasyRead(out bytes.Buffer, rawjson []byte) (outs string, err error) {
	var data interface{}

	var mapfromjson func(m map[string]interface{}, level int)
	var fromjson func(d interface{}, level int)
	fromjson = func(d interface{}, level int) {
		space := strings.Repeat(" ", level*4)
		switch val := d.(type) {
		case []interface{}:
			for i, u := range val {
				out.WriteString(fmt.Sprintf("[%d]: ", i))
				fromjson(u, level)
			}
		case map[string]interface{}:
			m, ok := d.(map[string]interface{})
			if ok {
				out.WriteString("{\n")
				mapfromjson(m, level+1)
				out.WriteString(fmt.Sprintf("%s}\n", space))
			} else {
				err = errors.New(fmt.Sprintf("generic JSON decode went to unknown type. %s", reflect.TypeOf(d)))
			}
		case nil:
			out.WriteString("null")
		case bool:
			if val {
				out.WriteString("true")
			} else {
				out.WriteString("false")
			}
		case string:
			out.WriteString(fmt.Sprintf("\"%s\"", d.(string)))
		case float64:
			out.WriteString(fmt.Sprintf("%f", d.(float64)))
		case int:
			out.WriteString(fmt.Sprintf("%d", d.(int)))
		case int64:
			out.WriteString(fmt.Sprintf("%d", d.(int64)))
		default:
			out.WriteString(fmt.Sprintf("<can't handle type. %s>", reflect.TypeOf(d)))
		}
	}

	mapfromjson = func(m map[string]interface{}, level int) {
		space := strings.Repeat(" ", level*4)
		for k, v := range m {
			switch vv := v.(type) {
			// case string:
			// 	out.WriteString(fmt.Sprintf("%s%s: \"%s\"\n", space, k, v.(string)))
			// 	//				fmt.Println(k, "is string", vv)
			// case float64:
			// 	out.WriteString(fmt.Sprintf("%s%s: ", space, k, v.(float64)))
			// 	//				fmt.Println(k, "is float64", vv)
			// case int:
			// 	out.WriteString(fmt.Sprintf("%s%s: %d\n", space, k, v.(int)))
			// case int64:
			// 	out.WriteString(fmt.Sprintf("%s%s: %d\n", space, k, v.(int64)))
			// 	//				fmt.Println(k, "is float64", vv)
			case map[string]interface{}:
				out.WriteString(fmt.Sprintf("%s%s: {\n", space, k))
				mapfromjson(v.(map[string]interface{}), level+1)
				out.WriteString(fmt.Sprintf("%s}\n", space))
			case []interface{}:
				//				fmt.Println(k, "is an array:")
				out.WriteString(fmt.Sprintf("%s%s: []", space, k))
				for i, u := range vv {
					out.WriteString(fmt.Sprintf("[%d]: ", i))
					fromjson(u, level+1)
					out.WriteString("\n")
					//					fmt.Println(i, u)
				}
			default:
				out.WriteString(fmt.Sprintf("%s%s: ", space, k))
				fromjson(v, level)
				out.WriteString("\n")
				//				fmt.Println(k, " - can't handle type. %s", reflect.TypeOf(v))
			}
		}

	}

	err = json.Unmarshal(rawjson, &data)
	if err == nil {
		out.WriteString("\n")
		fromjson(data, 0)
	}

	outs = out.String()

	return
}

func (self *MaestroClient) ConfigNetInterface(args []string) (string, error) {
	var netIfConfig maestroSpecs.NetIfConfigPayload

	// check for addition args beyond "net config-interface"
	if len(args)-2 <= 0 {
		return "Incorrect number of opts:", errors.New("Missing interface options")
	}

	for _, opt := range args[2:] {
		val := strings.Split(opt, "=")
		if len(val) < 2 {
			return "Invalid option", fmt.Errorf("Invalid option: %s", val)
		}
		DebugOut("opt=%s, arg=%s", val[0], val[1])
		//TODO: netIfConfig.AliasAddrV4
		//TODO: netIfConfig.WiFiSettings
		//TODO: netIfConfig.IEEE8021x
		//TODO: netIfConfig.Routes
		//TODO: netIfConfig.TestHttpsRouteOut
		//TODO: netIfConfig.TestICMPv4EchoOut
		switch strings.ToLower(val[0]) {
		case "ifname":
			netIfConfig.IfName = val[1]
		case "ifindex":
			i, err := strconv.Atoi(val[1])
			if err != nil {
				return "Invalid argument", err
			}
			netIfConfig.IfIndex = i
		case "dhcpv4enabled":
			b, err := strconv.ParseBool(val[1])
			if err != nil {
				return "Invalid argument", err
			}
			netIfConfig.DhcpV4Enabled = b
		case "ipv4addr":
			netIfConfig.IPv4Addr = val[1]
		case "ipv4mask":
			i, err := strconv.Atoi(val[1])
			if err != nil {
				return "Invalid argument", err
			}
			netIfConfig.IPv4Mask = i
		case "ipv4bcast":
			netIfConfig.IPv4BCast = val[1]
		case "ipv6addr":
			netIfConfig.IPv6Addr = val[1]
		case "hwaddr":
			netIfConfig.IPv6Addr = val[1]
		case "replaceaddress":
			netIfConfig.ReplaceAddress = val[1]
		case "clearaddresses":
			b, err := strconv.ParseBool(val[1])
			if err != nil {
				return "Invalid argument", err
			}
			netIfConfig.ClearAddresses = b
		case "down":
			b, err := strconv.ParseBool(val[1])
			if err != nil {
				return "Invalid argument", err
			}
			netIfConfig.Down = b
		case "defaultgateway":
			netIfConfig.DefaultGateway = val[1]
		case "fallbackdefaultgateway":
			netIfConfig.FallbackDefaultGateway = val[1]
		case "routepriority":
			i, err := strconv.Atoi(val[1])
			if err != nil {
				return "Invalid argument", err
			}
			netIfConfig.RoutePriority = i
		case "aux":
			b, err := strconv.ParseBool(val[1])
			if err != nil {
				return "Invalid argument", err
			}
			netIfConfig.Aux = b
		case "nameserveroverrides":
			netIfConfig.NameserverOverrides = val[1]
		case "dhcpdisableclearaddresses":
			b, err := strconv.ParseBool(val[1])
			if err != nil {
				return "Invalid argument", err
			}
			netIfConfig.DhcpDisableClearAddresses = b
		case "dhcpsteptimeout":
			i, err := strconv.Atoi(val[1])
			if err != nil {
				return "Invalid argument", err
			}
			netIfConfig.DhcpStepTimeout = i
		case "existing":
			netIfConfig.Existing = val[1]
		}
	}

	if netIfConfig.IfName == "" {
		return "Missing IfName", errors.New("Missing IfName")

	}

	var configs = []maestroSpecs.NetIfConfigPayload{netIfConfig}
	bytes, err := json.Marshal(configs)
	if err != nil {
		return "Failed to encode to JSON", err
	}

	resp, err2 := self.put("/net/interfaces", bytes)

	return resp.Status, err2
}

func (self *MaestroClient) GetNetInterfaces() (out string, err error) {
	resp, err := self.get("/net/interfaces")
	var buf bytes.Buffer
	if err == nil {
		DebugOut("resp.Body = %+v", resp.Body)
		body, err2 := ioutil.ReadAll(resp.Body)
		DebugOut("resp.Body body = %+v", body)
		DebugOut("resp.Body body = %s", string(body))
		if err2 == nil {
			buf.WriteString("interfaces:")
			out, err = FormatJsonEasyRead(buf, body)
			//			out = string()
			//			json.Unmarshal(body, alive)
		} else {
			DebugOut("Error on ReadAll %s", err2.Error())
			err = err2
		}
	}
	return
}

// SubscribeNetEventsResponse is the response from a /net/events call
type SubscribeNetEventsResponse struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}

// to be ran as a go routine
func (client *MaestroClient) netEventListener() {
	for {
		if len(client.netEventsSubscribeID) < 1 {
			DebugOut("network events listener stopping.")
			break
		}
		resp, err := client.get(fmt.Sprintf("/net/events/%s", client.netEventsSubscribeID))
		// var buf bytes.Buffer
		if err == nil {
			DebugOut("resp.Body = %+v", resp.Body)
			body, err2 := ioutil.ReadAll(resp.Body)
			DebugOut("resp.Body body = %+v", body)
			DebugOut("resp.Body body = %s", string(body))
			if resp.StatusCode != 200 && resp.StatusCode != 204 {
				Errorf("failed to get network events (%d): %s - Stopping listener.", resp.StatusCode, resp.Status)
				client.netEventsSubscribeID = ""
				break
			} else {
				if err2 == nil {
					var buf bytes.Buffer
					buf.WriteString("JSON:")
					out, err := FormatJsonEasyRead(buf, body)
					if err == nil {
						EventOut("network", "%s", out)
					} else {
						Errorf("Could not parse network events: %s", err.Error())
					}
				} else {
					DebugOut("Error on ReadAll %s", err2.Error())
					err = err2
				}
			}
		} else {
			Errorf("Error polling net events: %s", err.Error())
		}
	}
	client.netEventsListenerRunning = false
}

// SubscribeToNetEvents shell will subscribe to network events
func (client *MaestroClient) SubscribeToNetEvents() (out string, err error) {
	resp, err := client.get("/net/events")
	// var buf bytes.Buffer
	if err == nil {
		DebugOut("resp.Body = %+v", resp.Body)
		body, err2 := ioutil.ReadAll(resp.Body)
		DebugOut("resp.Body body = %+v", body)
		DebugOut("resp.Body body = %s", string(body))
		if resp.StatusCode != 200 {
			err = fmt.Errorf("failed to subscribe to network events (%d): %s", resp.StatusCode, resp.Status)
		} else {
			if err2 == nil {
				evresp := &SubscribeNetEventsResponse{}
				json.Unmarshal(body, evresp)
				if len(evresp.Error) < 1 && len(evresp.ID) > 0 {
					client.netEventsSubscribeID = evresp.ID
					if !client.netEventsListenerRunning {
						go client.netEventListener()
					}
				} else {
					err = fmt.Errorf("failed to subscribe to network events: %s", evresp.Error)
				}
				// buf.WriteString("interfaces:")
				// out, err = FormatJsonEasyRead(buf, body)
				//			out = string()
				//			json.Unmarshal(body, alive)
			} else {
				DebugOut("Error on ReadAll %s", err2.Error())
				err = err2
			}

		}
	}
	return
}

// jobs

func (client *MaestroClient) GetAllJobStatus() (out string, err error) {
	resp, err := client.get("/jobs")
	var buf bytes.Buffer
	if err == nil {
		DebugOut("resp.Body = %+v", resp.Body)
		body, err2 := ioutil.ReadAll(resp.Body)
		DebugOut("resp.Body body = %+v", body)
		DebugOut("resp.Body body = %s", string(body))
		if err2 == nil {
			buf.WriteString("interfaces:")
			out, err = FormatJsonEasyRead(buf, body)
			//			out = string()
			//			json.Unmarshal(body, alive)
		} else {
			DebugOut("Error on ReadAll %s", err2.Error())
			err = err2
		}
	}
	return
}

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
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

const (
	prompt_mode_normal   = 0
	prompt_mode_ifconfig = 1
)

var prompt_mode = prompt_mode_normal

var commands = []prompt.Suggest{
	{Text: "exit", Description: "Exit shell"},
	{Text: "alive", Description: "Check if maestro is running & get up time"},
	{Text: "debug", Description: "Turn on / off debug print outs"},
	{Text: "net", Description: "Query or change network interfaces"},
	{Text: "log", Description: "Query or change logging parameters"},
	{Text: "jobs", Description: "Query or change job configs"},
	{Text: "help", Description: "Print available commands."},
}

var logSubCommands = []prompt.Suggest{
	{Text: "get", Description: "Show configurations for all logging targets"},
	{Text: "set", Description: "Set configurations for a logging target"},
	{Text: "delete", Description: "Delete a configuration for a logging target"},
}

var netSubcommands = []prompt.Suggest{
	//{Text: "get-up-interfaces", Description: "Show configuration interfaces which are up."},
	//{Text: "get-enabled-interfaces", Description: "Show configuration interfaces which are enabled."},
	{Text: "get-interfaces", Description: "Show configurations for all interfaces"},
	//{Text: "renew-dhcp", Description: "Renew DHCP lease for a specific interface"},
	//{Text: "release-dhcp", Description: "Release DHCP lease for a specific interface"},
	//{Text: "ifdown", Description: "Shutdown an interface"},
	//{Text: "ifup", Description: "Bring up an interface"},
	{Text: "events", Description: "Listen for network events [interval-seconds]"},
	{Text: "config-interface", Description: "Enter config for an interface"},
	{Text: "get-dns", Description: "Show all domain name servers"},
	{Text: "add-dns", Description: "Add a new domain name server"},
	{Text: "delete-dns", Description: "Delete an existing domain name server"},
}

func GetCommandsHelpString([]string) (ret string, err error) {
	buffer := bytes.NewBufferString("Commands:\n")
	for _, cmd := range commands {
		buffer.WriteString(fmt.Sprintf("%-15s- %s\n", cmd.Text, cmd.Description))
	}
	buffer.WriteString("--\n")
	ret = string(buffer.Bytes())
	return
}

func GetNetSubcommandsHelpString([]string) (ret string, err error) {
	buffer := bytes.NewBufferString("Net Subcommands:\n")
	for _, cmd := range netSubcommands {
		buffer.WriteString(fmt.Sprintf("%-17s- %s\n", cmd.Text, cmd.Description))
	}
	buffer.WriteString("--\n")
	buffer.WriteString("Specify options as <opt>=<arg>, like IfName=eth0")
	ret = string(buffer.Bytes())
	return
}

func GetLogSubcommandsHelpString([]string) (ret string, err error) {
	buffer := bytes.NewBufferString("Log Subcommands:\n")
	for _, cmd := range logSubCommands {
		buffer.WriteString(fmt.Sprintf("%-17s- %s\n", cmd.Text, cmd.Description))
	}
	buffer.WriteString("--\n")
	buffer.WriteString("Specify options as <opt>=<arg>, like target=id")
	ret = string(buffer.Bytes())
	return
}

func Completer(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")
	//	w := d.GetWordBeforeCursor()

	// If PIPE is in text before the cursor, returns empty suggestions.
	for i := range args {
		if args[i] == "|" {
			return []prompt.Suggest{}
		}
	}

	// If word before the cursor starts with "-", returns CLI flag options.
	// if strings.HasPrefix(w, "-") {
	// 	return optionCompleter(args, strings.HasPrefix(w, "--"))
	// }

	// Return suggestions for option
	// if suggests, found := completeOptionArguments(d); found {
	// 	return suggests
	// }

	return argumentsCompleter(args) // excludeOptions(args))
}

func init() {
	fileListCache = map[string][]prompt.Suggest{}
}

func getPreviousOption(d prompt.Document) (cmd, option string, found bool) {
	args := strings.Split(d.TextBeforeCursor(), " ")
	l := len(args)
	if l >= 2 {
		option = args[l-2]
	}
	if strings.HasPrefix(option, "-") {
		return args[0], option, true
	}
	return "", "", false
}

func completeOptionArguments(d prompt.Document) ([]prompt.Suggest, bool) {
	cmd, option, found := getPreviousOption(d)
	if !found {
		return []prompt.Suggest{}, false
	}
	switch cmd {
	case "get", "describe", "create", "delete", "replace", "patch",
		"edit", "apply", "expose", "rolling-update", "rollout",
		"label", "annotate", "scale", "convert", "autoscale":
		switch option {
		case "-f", "--filename":
			return fileCompleter(d), true
		}
	}
	return []prompt.Suggest{}, false
}

/* file list */

var fileListCache map[string][]prompt.Suggest

func fileCompleter(d prompt.Document) []prompt.Suggest {
	path := d.GetWordBeforeCursor()
	if strings.HasPrefix(path, "./") {
		path = path[2:]
	}
	dir := filepath.Dir(path)
	if cached, ok := fileListCache[dir]; ok {
		return cached
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Print("[ERROR] catch error " + err.Error())
		return []prompt.Suggest{}
	}
	suggests := make([]prompt.Suggest, 0, len(files))
	for _, f := range files {
		if !f.IsDir() &&
			!strings.HasSuffix(f.Name(), ".yml") &&
			!strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}
		suggests = append(suggests, prompt.Suggest{Text: filepath.Join(dir, f.Name())})
	}
	return prompt.FilterHasPrefix(suggests, path, false)
}

func argumentsCompleter(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(commands, args[0], true)
	}

	first := args[0]
	second := args[1]
	switch first {
	case "net":
		if len(args) == 2 {
			return prompt.FilterHasPrefix(netSubcommands, second, true)
		}

		if len(args) >= 3 {
			last := args[len(args)-1]
			switch second {
			case "add-dns", "delete-dns":
				dns_set_args := []prompt.Suggest{
					{Text: "<server>", Description: "Domain name server address"},
				}
				return prompt.FilterHasPrefix(dns_set_args, last, true)
			case "config-interface":
				iface_args := []prompt.Suggest{
					{Text: "IfName", Description: "Interface name, like eth0"},
					{Text: "DhcpV4Enabled", Description: "true or false"},
					{Text: "IPv4Addr", Description: "IPv4 Address"},
					{Text: "IPv4Mask", Description: "IPv4 Netmask integer (CIDR format)"},
					{Text: "IPv4BCast", Description: "IPv4 Broadcast Address"},
					//{Text: "AliasAddrV4", Description: "NOT IMPLEMENTED"},
					{Text: "IPv6Addr", Description: "IPv6 Address"},
					{Text: "HwAddr", Description: "MAC Address or similar"},
					//{Text: "WiFiSettings", Description: "NOT IMPLEMENTED"},
					//{Text: "IEEE8021x", Description: "NOT IMPLEMENTED"},
					{Text: "ReplaceAddress", Description: "Address to delete before setting the new address"},
					{Text: "ClearAddresses", Description: "true or false.  if true, remove all existing addresses before setting the new address"},
					{Text: "Down", Description: "true or false.  if true, the interface is disabled"},
					{Text: "DefaultGateway", Description: "Default route associated with this interface"},
					//{Text: "FallbackDefaultGateway", Description: "NOT IMPLEMENTED"},
					{Text: "RoutePriority", Description: "Interface priority as a default route, ranked across all interfaces.  range 0-9, 0=first priority, 9=last"},
					{Text: "Aux", Description: "true or false"},
					{Text: "NameserverOverrides", Description: "Override DNS"},
					//{Text: "Routes", Description: "NOT IMPLEMENTED"},
					//{Text: "TestHttpsRouteOut", Description: "NOT IMPLEMENTED"},
					//{Text: "TestICMPv4EchoOut", Description: "NOT IMPLEMENTED"},
					{Text: "DhcpDisableClearAddresses", Description: "Don't allow DHCP to clear all addresses"},
					{Text: "DhcpStepTimeout", Description: "Max seconds to wait for DHCP address"},
					{Text: "Existing", Description: "override=replace any data in the db, replace=remove any data in the db"},
					{Text: "Type", Description: "Type of the connection, such as wifi or lte"},
					{Text: "SerialDevice", Description: "Path to the LTE modem serial device"},
					{Text: "APN", Description: "LTE modem access point name"},
				}
				return prompt.FilterHasPrefix(iface_args, last, true)
			}
		}

	case "log":
		if len(args) == 2 {
			return prompt.FilterHasPrefix(logSubCommands, second, true)
		}

		if len(args) >= 3 {
			last := args[len(args)-1]
			switch second {
			case "set", "delete":
				log_set_args := []prompt.Suggest{
					{Text: "target", Description: "Log filter target"},
					{Text: "levels", Description: "Log filter level"},
					{Text: "tag", Description: "Log filter tag"},
					{Text: "pre", Description: "Log pre filter"},
					{Text: "post", Description: "Log post filter"},
					{Text: "post-fmt-pre-msg", Description: "Log post format pre message"},
				}
				return prompt.FilterHasPrefix(log_set_args, last, true)
			}
		}

	case "debug":
		if len(args) == 2 {
			subcommands := []prompt.Suggest{
				{Text: "on", Description: "Turn on debugging"},
				{Text: "off", Description: "Turn off debugging"},
			}
			return prompt.FilterHasPrefix(subcommands, second, true)
		}
	case "jobs":
		if len(args) == 2 {
			subcommands := []prompt.Suggest{
				{Text: "get", Description: "Show all running jobs."},
				{Text: "stop", Description: "Stop one or more jobs by unique name"},
				{Text: "start", Description: "Start one or more jobs by unique name"},
				{Text: "register", Description: "Register (define) a new job using a JSON string"},
			}
			return prompt.FilterHasPrefix(subcommands, second, true)
		}
	case "get":
	// 	second := args[1]
	// 	if len(args) == 2 {
	// 		subcommands := []prompt.Suggest{
	// 			{Text: "componentstatuses"},
	// 			{Text: "configmaps"},
	// 			{Text: "daemonsets"},
	// 			{Text: "deployments"},
	// 			{Text: "endpoints"},
	// 			{Text: "events"},
	// 			{Text: "horizontalpodautoscalers"},
	// 			{Text: "ingresses"},
	// 			{Text: "jobs"},
	// 			{Text: "limitranges"},
	// 			{Text: "namespaces"},
	// 			{Text: "networkpolicies"},
	// 			{Text: "nodes"},
	// 			{Text: "persistentvolumeclaims"},
	// 			{Text: "persistentvolumes"},
	// 			{Text: "pod"},
	// 			{Text: "podsecuritypolicies"},
	// 			{Text: "podtemplates"},
	// 			{Text: "replicasets"},
	// 			{Text: "replicationcontrollers"},
	// 			{Text: "resourcequotas"},
	// 			{Text: "secrets"},
	// 			{Text: "serviceaccounts"},
	// 			{Text: "services"},
	// 			{Text: "statefulsets"},
	// 			{Text: "storageclasses"},
	// 			{Text: "thirdpartyresources"},
	// 			// aliases
	// 			{Text: "cs"},
	// 			{Text: "cm"},
	// 			{Text: "ds"},
	// 			{Text: "deploy"},
	// 			{Text: "ep"},
	// 			{Text: "hpa"},
	// 			{Text: "ing"},
	// 			{Text: "limits"},
	// 			{Text: "ns"},
	// 			{Text: "no"},
	// 			{Text: "pvc"},
	// 			{Text: "pv"},
	// 			{Text: "po"},
	// 			{Text: "psp"},
	// 			{Text: "rs"},
	// 			{Text: "rc"},
	// 			{Text: "quota"},
	// 			{Text: "sa"},
	// 			{Text: "svc"},
	// 		}
	// 		return prompt.FilterHasPrefix(subcommands, second, true)
	// 	}

	// 	third := args[2]
	// 	if len(args) == 3 {
	// 		switch second {
	// 		case "componentstatuses", "cs":
	// 			return prompt.FilterContains(getComponentStatusCompletions(), third, true)
	// 		case "configmaps", "cm":
	// 			return prompt.FilterContains(getConfigMapSuggestions(), third, true)
	// 		case "daemonsets", "ds":
	// 			return prompt.FilterContains(getDaemonSetSuggestions(), third, true)
	// 		case "deploy", "deployments":
	// 			return prompt.FilterContains(getDeploymentSuggestions(), third, true)
	// 		case "endpoints", "ep":
	// 			return prompt.FilterContains(getEndpointsSuggestions(), third, true)
	// 		case "ingresses", "ing":
	// 			return prompt.FilterContains(getIngressSuggestions(), third, true)
	// 		case "limitranges", "limits":
	// 			return prompt.FilterContains(getLimitRangeSuggestions(), third, true)
	// 		case "namespaces", "ns":
	// 			return prompt.FilterContains(getNameSpaceSuggestions(), third, true)
	// 		case "no", "nodes":
	// 			return prompt.FilterContains(getNodeSuggestions(), third, true)
	// 		case "po", "pod", "pods":
	// 			return prompt.FilterContains(getPodSuggestions(), third, true)
	// 		case "persistentvolumeclaims", "pvc":
	// 			return prompt.FilterContains(getPersistentVolumeClaimSuggestions(), third, true)
	// 		case "persistentvolumes", "pv":
	// 			return prompt.FilterContains(getPersistentVolumeSuggestions(), third, true)
	// 		case "podsecuritypolicies", "psp":
	// 			return prompt.FilterContains(getPodSecurityPolicySuggestions(), third, true)
	// 		case "podtemplates":
	// 			return prompt.FilterContains(getPodTemplateSuggestions(), third, true)
	// 		case "replicasets", "rs":
	// 			return prompt.FilterContains(getReplicaSetSuggestions(), third, true)
	// 		case "replicationcontrollers", "rc":
	// 			return prompt.FilterContains(getReplicationControllerSuggestions(), third, true)
	// 		case "resourcequotas", "quota":
	// 			return prompt.FilterContains(getResourceQuotasSuggestions(), third, true)
	// 		case "secrets":
	// 			return prompt.FilterContains(getSecretSuggestions(), third, true)
	// 		case "sa", "serviceaccounts":
	// 			return prompt.FilterContains(getServiceAccountSuggestions(), third, true)
	// 		case "svc", "services":
	// 			return prompt.FilterContains(getServiceSuggestions(), third, true)
	// 		}
	// 	}
	// case "describe":
	// 	second := args[1]
	// 	if len(args) == 2 {
	// 		return prompt.FilterHasPrefix(resourceTypes, second, true)
	// 	}

	// 	third := args[2]
	// 	if len(args) == 3 {
	// 		switch second {
	// 		case "componentstatuses", "cs":
	// 			return prompt.FilterContains(getComponentStatusCompletions(), third, true)
	// 		case "configmaps", "cm":
	// 			return prompt.FilterContains(getConfigMapSuggestions(), third, true)
	// 		case "daemonsets", "ds":
	// 			return prompt.FilterContains(getDaemonSetSuggestions(), third, true)
	// 		case "deploy", "deployments":
	// 			return prompt.FilterContains(getDeploymentSuggestions(), third, true)
	// 		case "endpoints", "ep":
	// 			return prompt.FilterContains(getEndpointsSuggestions(), third, true)
	// 		case "ingresses", "ing":
	// 			return prompt.FilterContains(getIngressSuggestions(), third, true)
	// 		case "limitranges", "limits":
	// 			return prompt.FilterContains(getLimitRangeSuggestions(), third, true)
	// 		case "namespaces", "ns":
	// 			return prompt.FilterContains(getNameSpaceSuggestions(), third, true)
	// 		case "no", "nodes":
	// 			return prompt.FilterContains(getNodeSuggestions(), third, true)
	// 		case "po", "pod", "pods":
	// 			return prompt.FilterContains(getPodSuggestions(), third, true)
	// 		case "persistentvolumeclaims", "pvc":
	// 			return prompt.FilterContains(getPersistentVolumeClaimSuggestions(), third, true)
	// 		case "persistentvolumes", "pv":
	// 			return prompt.FilterContains(getPersistentVolumeSuggestions(), third, true)
	// 		case "podsecuritypolicies", "psp":
	// 			return prompt.FilterContains(getPodSecurityPolicySuggestions(), third, true)
	// 		case "podtemplates":
	// 			return prompt.FilterContains(getPodTemplateSuggestions(), third, true)
	// 		case "replicasets", "rs":
	// 			return prompt.FilterContains(getReplicaSetSuggestions(), third, true)
	// 		case "replicationcontrollers", "rc":
	// 			return prompt.FilterContains(getReplicationControllerSuggestions(), third, true)
	// 		case "resourcequotas", "quota":
	// 			return prompt.FilterContains(getResourceQuotasSuggestions(), third, true)
	// 		case "secrets":
	// 			return prompt.FilterContains(getSecretSuggestions(), third, true)
	// 		case "sa", "serviceaccounts":
	// 			return prompt.FilterContains(getServiceAccountSuggestions(), third, true)
	// 		case "svc", "services":
	// 			return prompt.FilterContains(getServiceSuggestions(), third, true)
	// 		}
	// 	}
	// case "create":
	// 	subcommands := []prompt.Suggest{
	// 		{Text: "configmap", Description: "Create a configmap from a local file, directory or literal value"},
	// 		{Text: "deployment", Description: "Create a deployment with the specified name."},
	// 		{Text: "namespace", Description: "Create a namespace with the specified name"},
	// 		{Text: "quota", Description: "Create a quota with the specified name."},
	// 		{Text: "secret", Description: "Create a secret using specified subcommand"},
	// 		{Text: "service", Description: "Create a service using specified subcommand."},
	// 		{Text: "serviceaccount", Description: "Create a service account with the specified name"},
	// 	}
	// 	if len(args) == 2 {
	// 		return prompt.FilterHasPrefix(subcommands, args[1], true)
	// 	}
	// case "delete":
	// 	second := args[1]
	// 	if len(args) == 2 {
	// 		return prompt.FilterHasPrefix(resourceTypes, second, true)
	// 	}

	// 	third := args[2]
	// 	if len(args) == 3 {
	// 		switch second {
	// 		case "componentstatuses", "cs":
	// 			return prompt.FilterContains(getComponentStatusCompletions(), third, true)
	// 		case "configmaps", "cm":
	// 			return prompt.FilterContains(getConfigMapSuggestions(), third, true)
	// 		case "daemonsets", "ds":
	// 			return prompt.FilterContains(getDaemonSetSuggestions(), third, true)
	// 		case "deploy", "deployments":
	// 			return prompt.FilterContains(getDeploymentSuggestions(), third, true)
	// 		case "endpoints", "ep":
	// 			return prompt.FilterContains(getEndpointsSuggestions(), third, true)
	// 		case "ingresses", "ing":
	// 			return prompt.FilterContains(getIngressSuggestions(), third, true)
	// 		case "limitranges", "limits":
	// 			return prompt.FilterContains(getLimitRangeSuggestions(), third, true)
	// 		case "namespaces", "ns":
	// 			return prompt.FilterContains(getNameSpaceSuggestions(), third, true)
	// 		case "no", "nodes":
	// 			return prompt.FilterContains(getNodeSuggestions(), third, true)
	// 		case "po", "pod", "pods":
	// 			return prompt.FilterContains(getPodSuggestions(), third, true)
	// 		case "persistentvolumeclaims", "pvc":
	// 			return prompt.FilterContains(getPersistentVolumeClaimSuggestions(), third, true)
	// 		case "persistentvolumes", "pv":
	// 			return prompt.FilterContains(getPersistentVolumeSuggestions(), third, true)
	// 		case "podsecuritypolicies", "psp":
	// 			return prompt.FilterContains(getPodSecurityPolicySuggestions(), third, true)
	// 		case "podtemplates":
	// 			return prompt.FilterContains(getPodTemplateSuggestions(), third, true)
	// 		case "replicasets", "rs":
	// 			return prompt.FilterContains(getReplicaSetSuggestions(), third, true)
	// 		case "replicationcontrollers", "rc":
	// 			return prompt.FilterContains(getReplicationControllerSuggestions(), third, true)
	// 		case "resourcequotas", "quota":
	// 			return prompt.FilterContains(getResourceQuotasSuggestions(), third, true)
	// 		case "secrets":
	// 			return prompt.FilterContains(getSecretSuggestions(), third, true)
	// 		case "sa", "serviceaccounts":
	// 			return prompt.FilterContains(getServiceAccountSuggestions(), third, true)
	// 		case "svc", "services":
	// 			return prompt.FilterContains(getServiceSuggestions(), third, true)
	// 		}
	// 	}
	// case "edit":
	// 	if len(args) == 2 {
	// 		return prompt.FilterHasPrefix(resourceTypes, args[1], true)
	// 	}

	// 	if len(args) == 3 {
	// 		third := args[2]
	// 		switch args[1] {
	// 		case "componentstatuses", "cs":
	// 			return prompt.FilterContains(getComponentStatusCompletions(), third, true)
	// 		case "configmaps", "cm":
	// 			return prompt.FilterContains(getConfigMapSuggestions(), third, true)
	// 		case "daemonsets", "ds":
	// 			return prompt.FilterContains(getDaemonSetSuggestions(), third, true)
	// 		case "deploy", "deployments":
	// 			return prompt.FilterContains(getDeploymentSuggestions(), third, true)
	// 		case "endpoints", "ep":
	// 			return prompt.FilterContains(getEndpointsSuggestions(), third, true)
	// 		case "ingresses", "ing":
	// 			return prompt.FilterContains(getIngressSuggestions(), third, true)
	// 		case "limitranges", "limits":
	// 			return prompt.FilterContains(getLimitRangeSuggestions(), third, true)
	// 		case "namespaces", "ns":
	// 			return prompt.FilterContains(getNameSpaceSuggestions(), third, true)
	// 		case "no", "nodes":
	// 			return prompt.FilterContains(getNodeSuggestions(), third, true)
	// 		case "po", "pod", "pods":
	// 			return prompt.FilterContains(getPodSuggestions(), third, true)
	// 		case "persistentvolumeclaims", "pvc":
	// 			return prompt.FilterContains(getPersistentVolumeClaimSuggestions(), third, true)
	// 		case "persistentvolumes", "pv":
	// 			return prompt.FilterContains(getPersistentVolumeSuggestions(), third, true)
	// 		case "podsecuritypolicies", "psp":
	// 			return prompt.FilterContains(getPodSecurityPolicySuggestions(), third, true)
	// 		case "podtemplates":
	// 			return prompt.FilterContains(getPodTemplateSuggestions(), third, true)
	// 		case "replicasets", "rs":
	// 			return prompt.FilterContains(getReplicaSetSuggestions(), third, true)
	// 		case "replicationcontrollers", "rc":
	// 			return prompt.FilterContains(getReplicationControllerSuggestions(), third, true)
	// 		case "resourcequotas", "quota":
	// 			return prompt.FilterContains(getResourceQuotasSuggestions(), third, true)
	// 		case "secrets":
	// 			return prompt.FilterContains(getSecretSuggestions(), third, true)
	// 		case "sa", "serviceaccounts":
	// 			return prompt.FilterContains(getServiceAccountSuggestions(), third, true)
	// 		case "svc", "services":
	// 			return prompt.FilterContains(getServiceSuggestions(), third, true)
	// 		}
	// 	}

	// case "namespace":
	// 	if len(args) == 2 {
	// 		return prompt.FilterContains(getNameSpaceSuggestions(), args[1], true)
	// 	}
	// case "logs":
	// 	if len(args) == 2 {
	// 		return prompt.FilterContains(getPodSuggestions(), args[1], true)
	// 	}
	// case "rolling-update", "rollingupdate":
	// 	if len(args) == 2 {
	// 		return prompt.FilterContains(getReplicationControllerSuggestions(), args[1], true)
	// 	} else if len(args) == 3 {
	// 		return prompt.FilterContains(getReplicationControllerSuggestions(), args[2], true)
	// 	}
	// case "scale", "resize":
	// 	if len(args) == 2 {
	// 		// Deployment, ReplicaSet, Replication Controller, or Job.
	// 		r := getDeploymentSuggestions()
	// 		r = append(r, getReplicaSetSuggestions()...)
	// 		r = append(r, getReplicationControllerSuggestions()...)
	// 		return prompt.FilterContains(r, args[1], true)
	// 	}
	// case "cordon":
	// 	fallthrough
	// case "drain":
	// 	fallthrough
	// case "uncordon":
	// 	if len(args) == 2 {
	// 		return prompt.FilterHasPrefix(getNodeSuggestions(), args[1], true)
	// 	}
	// case "attach":
	// 	if len(args) == 2 {
	// 		return prompt.FilterContains(getPodSuggestions(), args[1], true)
	// 	}
	// case "exec":
	// 	if len(args) == 2 {
	// 		return prompt.FilterContains(getPodSuggestions(), args[1], true)
	// 	}
	// case "port-forward":
	// 	if len(args) == 2 {
	// 		return prompt.FilterContains(getPodSuggestions(), args[1], true)
	// 	}
	// 	if len(args) == 3 {
	// 		return prompt.FilterHasPrefix(getPortsFromPodName(args[1]), args[2], true)
	// 	}
	// case "rollout":
	// 	subCommands := []prompt.Suggest{
	// 		{Text: "history", Description: "view rollout history"},
	// 		{Text: "pause", Description: "Mark the provided resource as paused"},
	// 		{Text: "resume", Description: "Resume a paused resource"},
	// 		{Text: "undo", Description: "undoes a previous rollout"},
	// 	}
	// 	if len(args) == 2 {
	// 		return prompt.FilterHasPrefix(subCommands, args[1], true)
	// 	}
	// case "annotate":
	// case "config":
	// 	subCommands := []prompt.Suggest{
	// 		{Text: "current-context", Description: "Displays the current-context"},
	// 		{Text: "delete-cluster", Description: "Delete the specified cluster from the kubeconfig"},
	// 		{Text: "delete-context", Description: "Delete the specified context from the kubeconfig"},
	// 		{Text: "get-clusters", Description: "Display clusters defined in the kubeconfig"},
	// 		{Text: "get-contexts", Description: "Describe one or many contexts"},
	// 		{Text: "set", Description: "Sets an individual value in a kubeconfig file"},
	// 		{Text: "set-cluster", Description: "Sets a cluster entry in kubeconfig"},
	// 		{Text: "set-context", Description: "Sets a context entry in kubeconfig"},
	// 		{Text: "set-credentials", Description: "Sets a user entry in kubeconfig"},
	// 		{Text: "unset", Description: "Unsets an individual value in a kubeconfig file"},
	// 		{Text: "use-context", Description: "Sets the current-context in a kubeconfig file"},
	// 		{Text: "view", Description: "Display merged kubeconfig settings or a specified kubeconfig file"},
	// 	}
	// 	if len(args) == 2 {
	// 		return prompt.FilterHasPrefix(subCommands, args[1], true)
	// 	}
	// 	if len(args) == 3 {
	// 		third := args[2]
	// 		switch args[1] {
	// 		case "use-context":
	// 			return prompt.FilterContains(getContextSuggestions(), third, true)
	// 		}
	// 	}
	// case "cluster-info":
	// 	subCommands := []prompt.Suggest{
	// 		{Text: "dump", Description: "Dump lots of relevant info for debugging and diagnosis"},
	// 	}
	// 	if len(args) == 2 {
	// 		return prompt.FilterHasPrefix(subCommands, args[1], true)
	// 	}
	// case "explain":
	// 	return prompt.FilterHasPrefix(resourceTypes, args[1], true)
	default:
		return []prompt.Suggest{}
	}
	return []prompt.Suggest{}
}

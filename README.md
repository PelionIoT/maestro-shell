# maestro-shell

An interactive shell for controlling maestro locally on deviceOS


### Commands

#### `net events`
Opens a connection to maestro and continually waits for network events, printing them out as they occur

Usage: `net events`

Example
```
> net events
[OK] 

> [event(network)] JSON:
[0]: {
    data: {
        data: {
            type: "interface-state-down"
            interface: {
                id: "enxf4f951f22db3"
                index: 10.000000
                address: ""
                addressV6: ""
                linkstate: "LOWER_DOWN"
            }
        }
    }
}

> [event(network)] JSON:
[0]: {
    data: {
        data: {
            type: "interface-state-up"
            interface: {
                id: "enxf4f951f22db3"
                index: 10.000000
                address: "10.10.102.218"
                addressV6: ""
                linkstate: "LOWER_UP"
            }
        }
    }
}
```

#### `net get-interfaces`
Dumps information on all managed interfaces.

Usage: `net get-interfaces`

Example
```
> net get-interfaces
[OK] interfaces:
[0]: {
    StoredIfconfig: {
        replace_addr: ""
        down: false
        aux: false
        nameserver_overrides: ""
        test_echo_route_out: ""
        hw_addr: ""
        wifi: null
        test_https_route_out: ""
        alias_ipv4: null
        clear_addresses: false
        fallback_default_gateway: ""
        existing: ""
        ipv4_bcast: ""
. . .
```

#### `net config-interface`
Configures a single existing managed interface.

Usage: `net config-interface <opt1>=<arg1> <opt2>=<arg2> ...`

Example
```
> net config-interface ifname=eth1 dhcpv4enabled=false ipv4addr=192.168.1.1 ipv4mask=24 clearaddresses=true
[OK] 200 OK
```

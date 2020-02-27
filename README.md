# maestro-shell

An interactive shell for controlling maestro locally on deviceOS


### Commands

#### `log delete`
Delete a maestro log filter.

Example
```
> log delete target=/var/log/maestro/maestro.log levels=info

[DEBUG] opt=target, arg=/var/log/maestro/maestro.log
[DEBUG] opt=levels, arg=info
Log Sending: {"Target":"/var/log/maestro/maestro.log","Levels":"info","Tag":"","Pre":"","Post":"","PostFmtPreMsg":""}
[DEBUG] http resp:&{Status:200 OK StatusCode:200 Proto:HTTP/1.1 ProtoMajor:1 ProtoMinor:1 Header:map[Content-Length:[0] Date:[Thu, 27 Feb 2020 21:10:30 GMT]] Body:{} ContentLength:0 TransferEncoding:[] Close:false Uncompressed:false Trailer:map[] Request:0xc000148300 TLS:<nil>}
[DEBUG] log delete:200 OK <nil>
[OK] 200 OK
```

#### `log get`
Get a list of the log filter targets.

Example
```
> log get
[OK] Log Recieved.
[0]: {
    Format: {
        Time: ""
        Level: ""
        Tag: ""
        Origin: ""
    }
    Filters: [][0]: {
            Target: "default"
            Levels: ""
            Tag: ""
            Pre: ""
            Post: ""
            PostFmtPreMsg: ""
        }

    Rotate: {
        MaxFiles: 0.000000
        RotateOnStart: false
        MaxFileSize: 0.000000
        MaxTotalSize: 0.000000
    }
    Name: "default"
    FormatLevel: ""
    FormatOrigin: ""
    FormatPreMsg: ""
    TTY: ""
    FormatPre: ""
    FormatTime: ""
    FormatPost: ""
    File: ""
    ExampleFileOpts: {
        Mode: 0.000000
        Flags: 0.000000
        Max_files: 0.000000
        Max_file_size: 0.000000
        Max_total_size: 0.000000
        Rotate_on_start: false
    }
    Delim: ""
    FormatTag: ""
    Flag_json_escape_strings: false
}
[1]: {
    TTY: ""
    Delim: ""
    FormatTag: ""
    FormatOrigin: ""
    File: "/var/log/maestro/maestro.log"
    FormatTime: ""
    FormatPost: ""
    FormatPreMsg: ""
    Format: {
        Time: ""
        Level: ""
        Tag: ""
        Origin: ""
    }
    Filters: [][0]: {
            PostFmtPreMsg: ""
            Target: "/var/log/maestro/maestro.log"
            Levels: "warn"
            Tag: ""
            Pre: ""
            Post: ""
        }

        [1]: {
            Target: "/var/log/maestro/maestro.log"
            Levels: "error"
            Tag: ""
            Pre: ""
            Post: ""
            PostFmtPreMsg: ""
        }

    Rotate: {
        MaxFiles: 0.000000
        RotateOnStart: false
        MaxFileSize: 0.000000
        MaxTotalSize: 0.000000
    }
    ExampleFileOpts: {
        Max_files: 0.000000
        Max_file_size: 0.000000
        Max_total_size: 0.000000
        Rotate_on_start: false
        Mode: 0.000000
        Flags: 0.000000
    }
    FormatPre: ""
    FormatLevel: ""
    Flag_json_escape_strings: false
    Name: "/var/log/maestro/maestro.log"
}

```

#### `log set`
Set a log filter target.

`NOTE:` This function will not update an existing filter. You currently need to `log delete` an existing filter then call `log set` to add it back with changes.

Example
```
> log set target=/var/log/maestro/maestro.log levels=info

[DEBUG] opt=target, arg=/var/log/maestro/maestro.log
[DEBUG] opt=levels, arg=info
Log Sending: {"Target":"/var/log/maestro/maestro.log","Levels":"info","Tag":"","Pre":"","Post":"","PostFmtPreMsg":""}
[DEBUG] http resp:&{Status:200 OK StatusCode:200 Proto:HTTP/1.1 ProtoMajor:1 ProtoMinor:1 Header:map[Content-Length:[0] Date:[Thu, 27 Feb 2020 21:10:01 GMT]] Body:{} ContentLength:0 TransferEncoding:[] Close:false Uncompressed:false Trailer:map[] Request:0xc000148100 TLS:<nil>}
[DEBUG] log set:200 OK <nil>
[OK] 200 OK
```

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

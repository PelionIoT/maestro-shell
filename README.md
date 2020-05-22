# maestro-shell

An interactive shell for controlling maestro locally on Pelion Edge.

## Supported Operating Systems

* Ubuntu Core 16
* Yocto Linux

## Build
Maestro-shell is an application written in Go. To build maestro-shell, run the following commands:

```bash
./build-deps.sh
go build -o maestro-shell main.go
```

## Run

## Commands

### `alive`
Check if maestro is running and show uptime.

#### Examples

Maestro is up and running.
```
> alive
Maestro Up. Uptime = 1041.434485848s
```

Maestro is down.  To recover, start maestro.
```
> alive
[ERROR] Get http://unix/alive: dial unix /tmp/maestroapi.sock: connect: connection refused
```

Insufficient permissions to communicate with Maestro.  To recover, exit maestro-shell and restart with `sudo`.
```
> alive
[ERROR] Get http://unix/alive: dial unix /tmp/maestroapi.sock: connect: permission denied
```

### `log`
View and modify maestro log output.

#### `log get`
Get a list of the log output targets.

Example
```
> log get

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

#### `log set`
Add a log target filter.

`NOTE:` This function will not update an existing filter. You need to `log delete` an existing filter then call `log set` to add a new filter with the desired configuration.

This example adds a filter for `info` level messages which would allow the target file `/var/log/maestro/maestro.log` to start receiving `info` messages.
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


# Developing Maestro-Shell

### Compiling
1. Since Maestro-Shell is the front end for Maestro you need to first follow all the current build instruction to build [Maestro](https://github.com/armPelionEdge/maestro).
1. Once that environment is setup you need to call `vagrant ssh` to get into the build system from the `Maestro` folder.
1. After you logged into vagrant: `cd $MAESTRO_SRC/../maestro-shell`
1. To do the pre-build step run: `./build-deps.sh`
1. When that is done to build Maestro-Shell run `go build`
1. You have built the system

### Testing
1. Open two separate `vagarant ssh` shells from the Maestro folder on your main system.
1. In the 1st shell run `sudo maestro` to start the Maestro instance you will be sending commands to.
1. From the 2nd shell `cd $MAESTRO_SRC/../maestro-shell`
1. In that folder run `sudo maestro-shell.sh`
1. You are now ready to send commands from Maestro-Shell to Maestro

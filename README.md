# IGMP Leave Forge

A simple tool to create IGMPv2 leave messages.
## Usage
```
n: number of leave message, -1 for infinite loop, default to 3
src: source address
group: group address
```

example

```
go build
./igmp_leave_forge -n -1 -src 192.168.4.4 -group 224.3.4.7
```

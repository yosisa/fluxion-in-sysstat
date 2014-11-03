# fluxion-in-sysstat
[![wercker status](https://app.wercker.com/status/c3235e06012ecb2e4c7642e33a8ccefa/s/master "wercker status")](https://app.wercker.com/project/bykey/c3235e06012ecb2e4c7642e33a8ccefa)

fluxion-in-sysstat is an input plugin for [Fluxion](https://github.com/yosisa/fluxion), which collects and emits system stats periodically. It supports following stats:

* load average
* memory usage
* disk usage
* process info (cpu usage and memory usage)
* process count by name

NOTE: Currently, this plugin only works on linux

## Installation
```
go get github.com/yosisa/fluxion-in-sysstat
```

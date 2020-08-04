# ppinger

Inspect the status of the network and notify you if there is an error.

## Getting Started

```shell
% git clone https://github.com/nakajima0301/ppinger.git
% go run main.go
```

## Config

`config/default.csv`

```csv
name,hostname
target1,hostname1
target2,hostname2
target3,hostname3
```

e.g.

```csv
name,hostname
Example,www.example.com
Gateway,192.168.0.1
Me,127.0.0.1
```

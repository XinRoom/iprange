# iprange

这是一个通过给定IP范围字符串生成IP集合的工具.

This is a tool for generating IP SETS from a given IP range string.

## Build

```
git clone https://github.com/XinRoom/iprange
cd iprange
go build cmd/iprange.go
```

## Usage

```
.\iprange.exe
Gen Ip Set.
Usage: iprange.exe ipStr [ipStr/file] ...
IP format can :
        1.1.1.1
        1.1.1.1-2
        1.1.1-2.0-1
        1.1.1.1/30
        2001::59:63
        2001::59:63-89
        ...
in addition: Support multiple parameters, file and commas
```

## Feature

### 1. CidrMode 1.1.1.1/30

```
.\iprange.exe 1.1.1.1/30     
1.1.1.0
1.1.1.1
1.1.1.2
1.1.1.3
```

### 2. WideMode 1.1.1.6-1.1.1.8

```
.\iprange.exe 1.1.1.6-1.1.1.8                  
1.1.1.6
1.1.1.7
1.1.1.8
```

### 3. NarrowMode 1.1-2.1-2.1

```
.\iprange.exe 1.1-2.1-2.1
1.1.1.1
1.1.2.1
1.2.1.1
1.2.2.1
```

### 4. IPv6

```
.\iprange.exe 2002:2::1-9
2002:2::1
2002:2::2
2002:2::3
2002:2::4
2002:2::5
2002:2::6
2002:2::7
2002:2::8
2002:2::9
```

### 5. File

```
.\iprange.exe .\ips.txt
1.1.2.3
1.1.2.4
1.1.2.5
```

### 6. Stdin

```
echo 1.1.1.2-5 | .\iprange.exe
1.1.1.2
1.1.1.3
1.1.1.4
1.1.1.5
```

## Use as a library

import: `import "github.com/Xinroom/iprange"`

Simple to use:

```go
package main

import "github.com/XinRoom/iprange"
import "fmt"

func main() {
	ipSet, err := iprange.GenIpSet("1.1.1.1/30")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Sprintln("1.1.1.1/30 GenIpSet is %s", ipSet)
}
```

Iterator (Save memory):

```go
package main

import "github.com/XinRoom/iprange"
import "fmt"

func main() {
	it, err := iprange.NewIter("1.1.1.1/30")
	if err != nil {
		fmt.Println(err)
		return
	}
	for itn := it.Next(); it.HasNext(); itn = it.Next() {
		fmt.Println(itn)
	}
}
```
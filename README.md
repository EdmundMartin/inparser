## Inparser
A simple .ini parser specifically written to parse pgbouncer.ini files. A couple of 
edge cases are not properly covered due to them being outside our use case.

### Example Usage
```go
package main

import (
	"fmt"

	"github.com/EdmundMartin/inparser"
)

func main() {
	res, _ := inparser.ParseIni("pgbouncer.ini")
	sect := res.GetSection("databases")
	for _, prop := range sect.Properties {
		fmt.Println(prop.Key)
		fmt.Println(prop.Mapping)
	}
}

```
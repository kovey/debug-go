# debug library by golang
### Usage
    go get -u github.com/kovey/debug-go
### Examples
```golang
    package main

    import (
        "github.com/kovey/debug-go/debug"
    )

    func main() {
        debug.SetLevel(debug.Debug_Info)
        debug.Info("test [%s]", kovey)
        debug.Dbug("test [%s]", kovey)
        debug.Warn("test [%s]", kovey)
        debug.Erro("test [%s]", kovey)
    }

```

# go-pushnotifier
Go library for [pushnotifier.de](https://pushnotifier.de).

## Installation

```
go get github.com/temal-/go-pushnotifier
```

## Usage

Lets export some environment variables:

```
export PUSHNOTIFIER_PACKAGE="com.foo.bar"
export PUSHNOTIFIER_TOKEN="ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
export PUSHNOTIFIER_USERNAME="user"
export PUSHNOTIFIER_PASSWORD="secretpassword"
export PUSHNOTIFIER_DEBUG="false"
```

With the exception of `PUSHNOTIFIER_DEBUG` (which defaults to `false`) all of
them are required!

```
package main$
$
import ($
>-"fmt"$
$
>-"github.com/temal-/go-pushnotifier"$
)$
$
func main() {$
>-foo := pushnotifier.NewClientFromEnv()$
>-foo.Login()$
>-bar, _ := foo.ListDevices()$
>-fmt.Println(bar)$
>-baz, _ := foo.SendText([]string{bar[0].Id}, "test")$
>-fmt.Println(baz)$
}$
```

Alternatively you can also create a `NewClient` and supply all [needed parameters](https://github.com/temal-/go-pushnotifier/blob/master/pushnotifier.go#L61)
on your own.

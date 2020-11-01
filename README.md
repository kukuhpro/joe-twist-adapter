<h1 align="center">Joe Bot - Twist Adapter</h1>
<p align="center">Connecting joe with the Twist chat application. https://github.com/go-joe/joe</p>

---

This repository contains a module for the [Joe Bot library][joe].

## Getting Started

This library is packaged as [Go module][go-modules]. You can get it via:

```
go get github.com/kukuhpro/joe-twist-adapter
```

### Example usage

In order to connect your bot to slack you can simply pass it as module when
creating a new bot:

```go
package main

import (
	"os"

	"github.com/go-joe/joe"
	"github.com/kukuhpro/joe-twist-adapter"
)

func main() {
	b := joe.New("example-bot",
		twist.Adapter(),
		…
    )
	
	b.Respond("ping", Pong)

	err := b.Run()
	if err != nil {
		b.Logger.Fatal(err.Error())
	}
}
```

## Contributing

The current implementation is rather minimal and there are many more features
that could be implemented on the slack adapter so you are highly encouraged to
contribute. If you want to hack on this repository, please read the short
[CONTRIBUTING.md](CONTRIBUTING.md) guide first.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available,
see the [tags on this repository][tags]. 

## Authors

- **Kukuh Prabowo** - [kukuhpro](https://github.com/kukuhpro)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

[joe]: https://github.com/go-joe/joe
[go-modules]: https://github.com/golang/go/wiki/Modules

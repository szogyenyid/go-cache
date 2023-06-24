# go-cache

go-cache is a concurrent key-value store library for efficient in-memory caching with expiration in GoLang.

## Features

- Concurrent key-value store for efficient in-memory caching
- Supports expiration of key-value pairs
- Thread-safe operations with built-in synchronization
- Simple API for putting, getting, and deleting key-value pairs

## Installation

```shell
go get github.com/szogyenyid/go-cache
```

## Usage

```go
package main

import (
	"github.com/szogyenyid/go-cache"
	"time"
)

func main() {
	kvStore := cache.NewKeyValueStore()

	kvStore.Put("key1", "value1", time.Second*5)
	kvStore.Put("key2", "value2", time.Second*10)

	time.Sleep(time.Second * 7)

	value, found := kvStore.Get("key1")
	if found {
		println("Value:", value)
	} else {
		println("Key not found or expired.")
	}

	value, found = kvStore.Get("key2")
	if found {
		println("Value:", value)
	} else {
		println("Key not found or expired.")
	}
}
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any bug fixes, improvements, or new features.

## Upgrading

Phocus follows [semantic versioning](https://semver.org/), which means breaking changes may occur between major releases.

## License

Phocus is licensed under [MIT License](LICENSE).
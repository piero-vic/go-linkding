# go-linkding

Go client library for the [Linkding](https://github.com/sissbruecker/linkding) API.

## Installation

```bash
go get -u github.com/piero-vic/go-linkding
```

## Getting Started

### List Bookmarks

```go
package main

import (
	"fmt"

	"github.com/piero-vic/go-linkding"
)

func main() {
	client := linkding.NewClient("https://linkding.example.org", "secret-token")

	params := linkding.ListBookmarksParams{
		Limit: 15,
	}

	response, err := client.ListBookmarks(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(response.Results)
}
```

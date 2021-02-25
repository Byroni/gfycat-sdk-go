# The unofficial Gfycat SDK for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/byroni/gfycat-sdk-go)](https://goreportcard.com/report/github.com/byroni/gfycat-sdk-go)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-58%25-brightgreen.svg?longCache=true&style=flat)</a>
---
# STILL IN ACTIVE DEVELOPMENT
This is the unofficial Gfycat SDK for Golang. See official Gfycat API documentation [here](https://developers.gfycat.com/api/).

## Getting Started

### Installing

Run `go get github.com/byroni/gfycat-sdk-go`

### Configuring client credentials and authenticating

```go
import (
  "github.com/byroni/gfycat-sdk-go"
)

func main() {
    config := gfycat.ClientConfig{
    ClientID:     "", // Gfycat API client ID
    ClientSecret: "", // Gfycat API client secret
    }
    
    // Create a Gfycat client
    client := gfycat.New(config)
    
    // Authenticate
    if err := client.Authenticate(); err != nil {
      log.Fatal(err)
    }
    
    // Get the access token
    fmt.Printf("%s", client.AccessToken)
}
```
That's it! Now you can use any of the SDK's methods.  

You can request Gfycat client credentials [here](https://developers.gfycat.com/).

## Usage

All the methods are super easy to use:

### Searching for a Gfycat

```go
searchResponse, err := client.Search("michael scott")
```

### Get a Gfycat by ID

```go
gfycatResponse, err := client.GetGfycat("bountifulexaltedflycatcher")
```

## License

[MIT](https://github.com/Byroni/gfycat-sdk-go/raw/main/LICENSE.MIT)
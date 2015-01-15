# go-httpreplay

Marshal and unmarshal request/response, so that it can be persisted or replayed any time.

## Usage

```go
package main

import (
    "io/ioutil"
    "github.com/ddliu/go-httpreplay"
)

func main() {
    resp, err := http.Get("http://google.com/")

    // marshal and cache the response
    data, err := httpreplay.Marshal(resp)
    ioutil.WriteFile("/tmp/cache.data", data, 0777)

    // unmarshal the response
    data, err := ioutil.ReadFile("/tmp/cache.data")
    resp, err := httpreplay.Unmarshal(data)
}
```
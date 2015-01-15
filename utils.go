// Copyright 2015 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package httpreplay

import (
    "bytes"
    "io"
    "io/ioutil"
)

func newBytesReadCloser(data []byte) io.ReadCloser {
    return ioutil.NopCloser(bytes.NewReader(data))
}

func cloneAndCloseBody(body io.ReadCloser) (io.ReadCloser, error) {
    defer body.Close()

    data, err := ioutil.ReadAll(body)
    if err != nil {
        return nil, err
    }

    return newBytesReadCloser(data), nil
}
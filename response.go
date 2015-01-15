// Copyright 2015 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package httpreplay

import (
    "io/ioutil"
    "net/http"
    "encoding/json"
)

type RawResponse struct {
    Status     string // e.g. "200 OK"
    StatusCode int    // e.g. 200
    Proto      string // e.g. "HTTP/1.0"
    ProtoMajor int    // e.g. 1
    ProtoMinor int    // e.g. 0
    Header map[string][]string
    Body []byte
    ContentLength int64
    TransferEncoding []string
    Close bool
    Trailer map[string][]string
    Request *RawRequest
    // TLS *tls.ConnectionState
}

func toRawResponse(res *http.Response) (*RawResponse, error) {
    resBody := res.Body
    defer resBody.Close()
    data, err := ioutil.ReadAll(resBody)

    if err != nil {
        return nil, err
    }

    res.Body = newBytesReadCloser(data)

    var rawReq *RawRequest
    if res.Request != nil {
        rawReq, err = toRawRequest(res.Request)
        if err != nil {
            return nil, err
        }
    }

    raw := &RawResponse {
        Status: res.Status,
        StatusCode: res.StatusCode,
        Proto: res.Proto,
        ProtoMajor: res.ProtoMajor,
        ProtoMinor: res.ProtoMinor,
        Header: res.Header,
        Body: data,
        ContentLength: res.ContentLength,
        TransferEncoding: res.TransferEncoding,
        Close: res.Close,
        Trailer: res.Trailer,
        Request: rawReq,
    }

    return raw, nil
}

func fromRawResponse(raw *RawResponse) (*http.Response, error) {
    var req *http.Request
    var err error
    if raw.Request != nil {
        req, err = fromRawRequest(raw.Request)
        if err != nil {
            return nil, err
        }
    } 
    res := &http.Response {
        Status: raw.Status,
        StatusCode: raw.StatusCode,
        Proto: raw.Proto,
        ProtoMajor: raw.ProtoMajor,
        ProtoMinor: raw.ProtoMinor,
        Header: raw.Header,
        Body: newBytesReadCloser(raw.Body),
        ContentLength: raw.ContentLength,
        TransferEncoding: raw.TransferEncoding,
        Close: raw.Close,
        Trailer: raw.Trailer,
        Request: req,
    }

    return res, nil
}

func MarshalResponse(res *http.Response) ([]byte, error) {
    raw, err := toRawResponse(res)
    if err != nil {
        return nil, err
    }

    return json.Marshal(raw)
}

func UnmarshalResponse(data []byte) (*http.Response, error) {
    var raw RawResponse

    err := json.Unmarshal(data, &raw)
    if err != nil {
        return nil, err
    }

    return fromRawResponse(&raw)
}
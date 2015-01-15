// Copyright 2015 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package httpreplay

import (
    "net/http"
    "net/url"
    "io/ioutil"
    "encoding/json"
)

type RawRequest struct {
    Method string
    URL string
    Proto      string // "HTTP/1.0"
    ProtoMajor int    // 1
    ProtoMinor int    // 0
    Header map[string][]string
    Body []byte
    ContentLength int64
    TransferEncoding []string
    Close bool
    Host string
    Form map[string][]string
    PostForm map[string][]string

    // Ignored...
    // MultipartForm *multipart.Form

    Trailer map[string][]string
    RemoteAddr string
    RequestURI string

    // TLS *tls.ConnectionState    
}

func toRawRequest(req *http.Request) (*RawRequest, error) {
    var data []byte

    if req.Body != nil {
        reqBody := req.Body
        defer reqBody.Close()
        data, err := ioutil.ReadAll(reqBody)

        if err != nil {
            return nil, err
        }

        req.Body = newBytesReadCloser(data)
    }


    raw := &RawRequest {
        Method: req.Method,
        URL: req.URL.String(),
        Proto: req.Proto,
        ProtoMajor: req.ProtoMajor,
        ProtoMinor: req.ProtoMinor,
        Header: req.Header,
        Body: data,
        ContentLength: req.ContentLength,
        TransferEncoding: req.TransferEncoding,
        Close: req.Close,
        Host: req.Host,
        Form: req.Form,
        PostForm: req.PostForm,
        Trailer: req.Trailer,
        RemoteAddr: req.RemoteAddr,
        RequestURI: req.RequestURI,
    }

    return raw, nil
}

func fromRawRequest(raw *RawRequest) (*http.Request, error) {
    urlinfo, err := url.Parse(raw.URL)
    if err != nil {
        return nil, err
    }

    req := &http.Request {
        Method: raw.Method,
        URL: urlinfo,
        Proto: raw.Proto,
        ProtoMajor: raw.ProtoMajor,
        ProtoMinor: raw.ProtoMinor,
        Header: raw.Header,
        Body: newBytesReadCloser(raw.Body),
        ContentLength: raw.ContentLength,
        TransferEncoding: raw.TransferEncoding,
        Close: raw.Close,
        Host: raw.Host,
        Form: raw.Form,
        PostForm: raw.PostForm,
        Trailer: raw.Trailer,
        RemoteAddr: raw.RemoteAddr,
        RequestURI: raw.RequestURI,
    }

    return req, nil
}



func MarshalRequest(req *http.Request) ([]byte, error) {
    raw, err := toRawRequest(req)
    if err != nil {
        return nil, err
    }

    return json.Marshal(raw)
}

func UnmarshalRequest(data []byte) (*http.Request, error) {
    var raw RawRequest

    err := json.Unmarshal(data, &raw)
    if err != nil {
        return nil, err
    }

    return fromRawRequest(&raw)
}
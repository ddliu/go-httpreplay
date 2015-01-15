// Copyright 2015 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package httpreplay

import (
    "testing"
    "net/http"
    "io/ioutil"
    "strings"
)

func TestResponse(t *testing.T) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", "http://httpbin.org/get", nil)
    if err != nil {
        t.Error(err)
    }

    req.Header.Set("x-request-with", "go-httpreplay")

    enReq, err := MarshalRequest(req)
    if err != nil {
        t.Error(err)
    }

    reqClone, err := UnmarshalRequest(enReq)
    if err != nil {
        t.Error(err)
    }

    if reqClone.Header.Get("x-request-with") != "go-httpreplay" {
        t.Error()
    }

    resp1, err := client.Do(req)
    if err != nil {
        t.Error(err)
    }

    resp2, err := client.Do(reqClone)
    if err != nil {
        t.Error(err)
    }

    deResp1, err := MarshalResponse(resp1)
    if err != nil {
        t.Error(err)
    }

    resp1Clone, err := UnmarshalResponse(deResp1)
    if err != nil {
        t.Error(err)
    }

    defer resp1.Body.Close()
    defer resp2.Body.Close()
    defer resp1Clone.Body.Close()
    
    body1, err := ioutil.ReadAll(resp1.Body)
    if err != nil {
        t.Error(err)
    }

    body2, err := ioutil.ReadAll(resp2.Body)
    if err != nil {
        t.Error(err)
    }

    body1Clone, err := ioutil.ReadAll(resp1Clone.Body)
    if err != nil {
        t.Error(err)
    }

    if string(body1) != string(body1Clone) {
        t.Error()
        println("======================")
        println(string(body1))
        println("======================")
        println(string(body1Clone))
    }

    if !strings.Contains(string(body2), "go-httpreplay") {
        t.Error()
        println(string(body2))
    }
}
package main

// refer
// https://en.wikipedia.org/wiki/Real_Time_Streaming_Protocol
// https://tools.ietf.org/html/rfc2326

// go lang test
// https://tour.golang.org/flowcontrol/1

import (
    "errors"
    "fmt"
    "strings"
)

type request struct {  
    // Request_Line
    method      string
    media       string
    version     string

    cseq        string

    // general_header

    // Cache_Control
    cache_control   string
    cache_directive   string
    cache_request_directive string
    cache_response_directive    string
    // connection
    // date
    // via

    // Request Header
    accept          string
    accept_encoding string
    accept_language string
    authorization   string
    from   string
    if_Modified_Since   string
    ranges       string
    referer     string
    user_agent  string

    // entity_header
    allow       string
    content_base    string
    content_encoding string
    content_language string
    content_length string
    content_location string
    content_type string
    expires     string
    last_modified   string
    extension_header    string

    // CRLF

    // message_body
}

func (req *request) init() {
    req.method = ""
    req.media = ""
    req.version = ""
}

func (req *request) parse_request_line(line string) (string, string, string, error) {
    metas := strings.SplitN(line, " ", 3)
    if metas != nil && len(metas) == 3 {
        method := metas[0]
        media := metas[1]
        versions := strings.SplitN(metas[2], "/", 2)

        if versions != nil && len(versions) == 2 {
            version := versions[1]
            return method, media, version, nil
        } else {
            return "", "", "", errors.New("version haven't '/'")
        }
    }
    return "", "", "", errors.New("line format wrong")
}

func (req *request) parse_header(line string) (error) {
    headers := strings.SplitN(line, ": ", 2)
    if len(headers) == 2 {
        key := strings.ToLower(headers[0])
        switch key {
        case "cseq":
            req.cseq = headers[1]
            break
        case "user-agent":
            req.user_agent = headers[1]
            break
        default:
            break
        }
    } else {
        return errors.New("header line format is invalid")
    }

    return nil
}

func parse_request(message string) (request, error) {
    req := request{}
    req.init()

    lines := strings.Split(message, "\r\n")
    if len(lines) == 0 {
        return req, errors.New("message havn't CRLF")
    }

    has_method := false

    for i := 0; i < len(lines); i++ {
        if len(lines[i]) == 0 {
            continue
        }

        if has_method {
            err := req.parse_header(lines[i])
            if err != nil {
                fmt.Println(err)
                return req, nil
            }
        } else {
            has_method = true

            method, media, version, err := req.parse_request_line(lines[i])
            if err != nil {
                fmt.Println(err)
                return req, nil
            }
            req.method = method
            req.media = media
            req.version = version
        }
    }

    return req, nil
}
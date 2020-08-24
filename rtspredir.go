package main

// refer
// https://en.wikipedia.org/wiki/Real_Time_Streaming_Protocol
// https://tools.ietf.org/html/rfc2326

// go lang test
// https://tour.golang.org/flowcontrol/1

import (
    "bufio"
    "bytes"
    "errors"
    "fmt"
    "net"
    "time"
)

const (
    CONN_TYPE = "tcp"
    CONN_HOST = "0.0.0.0"
    CONN_PORT = "554"
)

// TODO: refactor CRLF split
func read_request_message(conn net.Conn) (string, error) {
    var buffer bytes.Buffer

    reader := bufio.NewReader(conn)

    for {
        data, err := reader.ReadBytes('\n')
        if err != nil {
            fmt.Println(err)
            return "", errors.New("read fail")
        }
        data = bytes.TrimRight(data, "\r\n")
        line := string(data[:])

        if len(line) == 0 {
            break
        }
        buffer.WriteString(line + "\r\n")
    }

    return buffer.String(), nil
}

func on_rtsp_client(conn net.Conn) {
    defer conn.Close()

    req_idx := 1
    res_idx := 1

    fmt.Printf("welcome new client: %s\n", conn.RemoteAddr().String())

    for {
        message, err := read_request_message(conn)
        if err != nil {
            fmt.Println(err)
            break
        }

        req, err := parse_request_message(message)
        if err != nil {
            fmt.Println(err)
            break
        }

        fmt.Printf("\n<<<<<<<<<<<<<<<<\n")
        fmt.Printf("[%03d] request:\n%s", req_idx, message)
        fmt.Printf("\n<<<<<<<<<<<<<<<<\n")

        response, err := take_response(req)

        fmt.Printf("\n>>>>>>>>>>>>>>>>\n")
        fmt.Printf("[%03d] response:\n%s", res_idx, response)
        fmt.Printf("\n>>>>>>>>>>>>>>>>\n")

        conn.Write([]byte(string(response)))
        if err != nil {
            fmt.Println(err)
            break
        }

        req_idx++
        res_idx++
    }

    fmt.Printf("client %s closed\n", conn.RemoteAddr().String())
}

func launch_server(protocol string, host string, port string, user string, pwd string) int {
    switch protocol {
    case "tcp", "tcp4", "tcp6":
    default:
    return -1
    }

    tcp_server, err := net.Listen(protocol, ":" + port)
    if err != nil {
        fmt.Println("error listening:", err.Error())
        return -2
    }
    defer tcp_server.Close()
    fmt.Printf("listening on %s://%s:%s\n", protocol, host, port)

    for {
        conn, err := tcp_server.Accept()
        if err != nil {
            fmt.Println("error accepting: ", err.Error())
            break
        }

        go on_rtsp_client(conn)
    }

    return -3
}

func main() {
    ret := launch_server(CONN_TYPE, CONN_HOST, CONN_PORT, "", "")

    if ret == 0 {
        for {
            time.Sleep(time.Duration(2) * time.Second)
        }
    }
}
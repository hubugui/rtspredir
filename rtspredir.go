package main

// refer
// https://en.wikipedia.org/wiki/Real_Time_Streaming_Protocol
// https://tools.ietf.org/html/rfc2326

// go lang test
// https://tour.golang.org/flowcontrol/1

import (
    "bufio"
    "errors"
    "fmt"
    "net"
    "time"
    "strings"
)

const (
    CONN_TYPE = "tcp"
    CONN_HOST = "0.0.0.0"
    CONN_PORT = "554"
)

func parseRequestFirstLine(line string) (string, string, string, error) {
    metas := strings.SplitN(line, " ", 3)
    if len(metas) == 3 {
        method := metas[0]
        media := metas[1]
        versions := strings.SplitN(metas[2], "/", 2)

        if len(versions) == 2 {
            version := versions[1]
            return method, media, version, nil
        } else {
            return "", "", "", errors.New("version haven't '/'")
        }
    }
    return "", "", "", errors.New("first line content wrong")
}

func on_new_client(conn net.Conn) {
    fmt.Printf("%s connected\n", conn.RemoteAddr().String())

    for {
        netData, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil {
            fmt.Println(err)
            break
        }

        request := strings.TrimSpace(string(netData))
        method, media, version, err := parseRequestFirstLine(request)
        if err != nil {
            fmt.Println(err)
            continue
        }
        fmt.Printf("method:%s, %s, %s\n", method, media, version)

        switch method {
        case "DESCRIBE":
        case "ANNOUNCE":
        case "GET_PARAMETER":
        case "OPTIONS":
        case "PAUSE":
        case "PLAY":
        case "RECORD":
        case "SETUP":
        case "SET_PARAMETER":
        case "TEARDOWN":
        case "EXIT":
            break
        default:
            fmt.Printf("unknown method: %s.\n", method)
        }

        result := "bye\n"
        conn.Write([]byte(string(result)))
    }
    fmt.Printf("Serving %s closed\n", conn.RemoteAddr().String())
    conn.Close()
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

        go on_new_client(conn)
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
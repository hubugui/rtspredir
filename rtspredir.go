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
    "os"
    "strings"
)

const (
    CONN_HOST = "0.0.0.0"
    CONN_PORT = "554"
    CONN_TYPE = "tcp"
)

func main() {
    l, err := net.Listen("tcp4", ":" + CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()

    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            break
        }

        go handleRequest(conn)
    }
}

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

func handleRequest(conn net.Conn) {
	fmt.Print("Serving %s connected\n", conn.RemoteAddr().String())
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
        fmt.Print("%s, %s, %s\n", method, media, version)

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
            fmt.Print("unknown method: %s.\n", method)
        }

        result := "bye\n"
        conn.Write([]byte(string(result)))
    }
    conn.Close()
    fmt.Print("Serving %s closed\n", conn.RemoteAddr().String())
}
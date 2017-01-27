package main

import (
    "os"
    "fmt"
    "flag"

    "github.com/danbovey/Chappy/server"
    "github.com/danbovey/Chappy/secret"
)

func main() {
    if len(os.Args) < 2 {
        help()
        return
    }

    switch os.Args[1] {
    case "start":
        server.Start()
    case "secret":
        secret.Generate()
    case "help":
        help()
    default:
        flag.PrintDefaults()
    }
}

func help() {
    fmt.Println("Chappy - the simplest way to deploy websites using GitHub webhooks.\n")
    fmt.Println("Commands: start, secret\n")
    fmt.Println("Options:")
    fmt.Println("-ip\t\t\tIP that Chappy should serve webhooks on")
    fmt.Println("-port\t\t\tPort that Chappy should serve webhooks on")
    fmt.Println("-baseurl\t\tURL prefix to serve from (ip:port/BASEURL/:webhook-id")
    fmt.Println("-hotreload\t\tWatch .chappy.json file for changes and reload them automatically")
    fmt.Println("-projectsfilepath\tPath to the json file containing webhooks config")
    fmt.Println("-secure\t\t\tUse HTTPS instead of HTTP")
    fmt.Println("-cert\t\t\tPath to the HTTPS certificate pem file")
    fmt.Println("-key\t\t\tPath to the HTTPS certificate private key pem file")
}

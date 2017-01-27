package server

import (
    "fmt"
    "log"
    "os"
    "flag"
    "os/exec"
    "path/filepath"
    "io/ioutil"
    "net/http"
    "net/url"
    "encoding/json"
    "strings"
    "github.com/gorilla/mux"
    "github.com/urfave/negroni"

    "github.com/danbovey/Chappy/project"

    fsnotify "gopkg.in/fsnotify.v1"
)

var (
    startCommand = flag.NewFlagSet("start", flag.ContinueOnError)
    
    ip                  = startCommand.String("ip", "0.0.0.0", "IP that Chappy should serve webhooks on")
    port                = startCommand.Int("port", 9000, "Port that Chappy should serve webhooks on")
    baseURL             = startCommand.String("baseurl", "", "URL prefix to serve from (ip:port/BASEURL/:webhook-id")
    hotReload           = startCommand.Bool("hotreload", false, "Watch .chappy.json file for changes and reload them automatically")
    projectsFilePath    = startCommand.String("projectsfilepath", "projects.json", "Path to the json file containing webhooks config")
    secure              = startCommand.Bool("secure", false, "Use HTTPS instead of HTTP")
    cert                = startCommand.String("cert", "cert.pem", "Path to the HTTPS certificate pem file")
    key                 = startCommand.String("key", "key.pem", "Path to the HTTPS certificate private key pem file")

    projects project.Projects
    watcher *fsnotify.Watcher
)

func Start() {
    startCommand.Parse(os.Args[2:])

    loadProjects()

    if projects != nil {
        if *hotReload {
            // Set up file watcher
            var err error
            watcher, err = fsnotify.NewWatcher()
            if err != nil {
                log.Fatal("Error creating file watcher instance", err)
            }

            defer watcher.Close()

            go watchForFileChange()

            err = watcher.Add(*projectsFilePath)
            if err != nil {
                log.Fatal("error adding projects file to the watcher", err)
            }
        }

        // Allow the router to continue running by recovering from errors
        negroniRecovery := &negroni.Recovery{
            PrintStack: true,
            StackAll:   false,
            StackSize:  1024 * 8,
        }
        n := negroni.New(negroniRecovery)

        router := mux.NewRouter()

        var webhookURL string
        if *baseURL == "" {
            webhookURL = "/{name}"
        } else {
            webhookURL = "/" + *baseURL + "/{name}"
        }

        router.HandleFunc(webhookURL, hookHandler)
        n.UseHandler(router)

        if *secure {
            log.Printf("🤖\tChappy is serving webhooks on https://%s:%d%s", *ip, *port, *baseURL)
            log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", *ip, *port), *cert, *key, n))
        } else {
            log.Printf("🤖\tChappy is serving webhooks on http://%s:%d%s", *ip, *port, *baseURL)
            log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *ip, *port), n))
        }
    }
}

func loadProjects() {
    newProjects := make([]project.Project, 0)

    log.Printf("🔄\tLoading projects from %s\n", *projectsFilePath)

    raw, err := ioutil.ReadFile(*projectsFilePath)
    if err != nil {
        fmt.Println("Couldn't load projects from file!\n", err)
    } else {
        json.Unmarshal(raw, &newProjects)

        seenNames := make(map[string]bool)
        for _, project := range newProjects {
            if seenNames[project.Name] == true {
                log.Printf("Error: The project %s has already been loaded!\nPlease check your projects file for duplicate names!", project.Name)
                return
            }
            seenNames[project.Name] = true
            log.Printf("✅\tLoaded: %s\n", project.Name)
        }
        log.Printf("\n")

        projects = newProjects
    }
}

func hookHandler(w http.ResponseWriter, r *http.Request) {
    name := mux.Vars(r)["name"]

    if matchedProject := projects.Match(name); matchedProject != nil {
        log.Printf("Received event for '%s'\n", name)

        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            log.Printf("error reading the request body. %+v\n", err)
        }

        // parse headers
        headers := valuesToMap(r.Header)

        // parse query variables
        query := valuesToMap(r.URL.Query())

        // parse body
        var payload map[string]interface{}

        contentType := r.Header.Get("Content-Type")
        if strings.Contains(contentType, "json") {
            decoder := json.NewDecoder(strings.NewReader(string(body)))
            decoder.UseNumber()

            err := decoder.Decode(&payload)

            if err != nil {
                log.Printf("Error: Failed to parse JSON payload %+v\n", err)
            }
        } else if strings.Contains(contentType, "form") {
            fd, err := url.ParseQuery(string(body))
            if err != nil {
                log.Printf("Error: Failed to parse form payload %+v\n", err)
            } else {
                payload = valuesToMap(fd)
            }
        }

        // Check GitHub Event is correct
        githubEvent := r.Header.Get("X-GitHub-Event")
        if githubEvent == "" {
            log.Print("Error: Received request without an X-GitHub-Event header")
            return
        }
        if githubEvent != "push" {
            log.Printf("Received request for '%s' event, skipping and waiting for a 'push' event", githubEvent)
            return
        }

        // Check GitHub Signature is valid
        if !matchedProject.CheckHubSignature(body, r.Header.Get("X-Hub-Signature")) {
            return
        }

        // Check the branch name
        if payload["ref"] != "refs/heads/" + matchedProject.Branch {
            log.Print("Not the branch we're listening for")
            return
        }

        log.Printf("Deployment for '%s' triggered\n", matchedProject.Name)

        go runDeployScript(matchedProject, &headers, &query, &payload, &body)

        return
    } else {
        log.Printf("Error: Received event for '%s' but no project found", name)
    }
}

func runDeployScript(project *project.Project, headers, query, payload *map[string]interface{}, body *[]byte) (string, error) {
    cmd := exec.Command(project.Script)
    cmd.Dir = filepath.Dir(project.Script)
    cmd.Args = project.ExtractCommandArguments(headers, query, payload)

    if project.Script != cmd.Path {
        log.Printf("Executing %s (%s) with arguments %q using %s as cwd\n", project.Script, cmd.Path, cmd.Args[1:], cmd.Dir)
    } else {
        log.Printf("Executing %s with arguments %q using %s as cwd\n", project.Script, cmd.Args[1:], cmd.Dir)
    }

    out, err := cmd.Output()

    log.Printf("Command output: %s\n", out)

    if err != nil {
        log.Printf("Error occurred: %+v\n", err)
    }

    log.Printf("Finished handling '%s'\n", project.Name)

    return string(out), err
}


func watchForFileChange() {
    for {
        select {
        case event := <-(*watcher).Events:
            if event.Op&fsnotify.Write == fsnotify.Write {
                log.Println("Projects file modified")

                loadProjects()
            }
        case err := <-(*watcher).Errors:
            log.Println("Watcher error:", err)
        }
    }
}

// valuesToMap converts map[string][]string to a map[string]string object
func valuesToMap(values map[string][]string) map[string]interface{} {
    ret := make(map[string]interface{})

    for key, value := range values {
        if len(value) > 0 {
            ret[key] = value[0]
        }
    }

    return ret
}

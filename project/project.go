package project

import (
    "fmt"
    "log"
    "strconv"
    "strings"
    "reflect"
    "crypto/hmac"
    "crypto/sha1"
    "encoding/hex"
)

type Project struct {
    Name    string  // Name of the project
    Branch  string  // The repo branch to handle and deploy
    Script  string  // Script to run to deploy
    Secret  string  // Secret for the GitHub webhook
}

// ExtractCommandArguments creates a list of arguments from the GitHub
// payload that are ready to be used with exec.Command()
func (p *Project) ExtractCommandArguments(headers, query, payload *map[string]interface{}) []string {
    payloadToExtract := []string{ "head_commit.id", "pusher.name", "pusher.email" }

    var args = make([]string, 0)
    args = append(args, p.Script)

    for i := range payloadToExtract {
        if arg, ok := ExtractParameterAsString(payloadToExtract[i], *payload); ok {
            args = append(args, arg)
        } else {
            args = append(args, "")
            fmt.Sprintf("Error: Couldn't retrieve argument for %+v", payloadToExtract[i])
        }
    }

    return args
}

func (p *Project) CheckHubSignature(payload []byte, signature string) bool {
    if signature == "" {
        log.Printf("Error: Received request without an X-Hub-Signature")
        return false
    }

    // Calculate and verify the SHA1 signature of the given payload
    if strings.HasPrefix(signature, "sha1=") {
        signature = signature[5:]
    }

    mac := hmac.New(sha1.New, []byte(p.Secret))
    _, err := mac.Write(payload)
    if err != nil {
        log.Printf("Failed to test X-Hub-Signature", err)
        return false
    }
    expectedMAC := hex.EncodeToString(mac.Sum(nil))

    if !hmac.Equal([]byte(signature), []byte(expectedMAC)) {
        log.Print("Error: Received request with an invalid X-Hub-Signature")
        return false
    }

    return true
}

type Projects []Project

// Match iterates through Hooks and returns first one that matches the given ID,
// if no hook matches the given ID, nil is returned
func (p *Projects) Match(name string) *Project {
    for i := range *p {
        if (*p)[i].Name == name {
            return &(*p)[i]
        }
    }

    return nil
}

// ExtractParameterAsString extracts value from interface{} as string based on the passed string
func ExtractParameterAsString(s string, params interface{}) (string, bool) {
    if pValue, ok := GetParameter(s, params); ok {
        return fmt.Sprintf("%v", pValue), true
    }
    return "", false
}

// GetParameter extracts interface{} value based on the passed string
func GetParameter(s string, params interface{}) (interface{}, bool) {
    if params == nil {
        return nil, false
    }

    if paramsValue := reflect.ValueOf(params); paramsValue.Kind() == reflect.Slice {
        if paramsValueSliceLength := paramsValue.Len(); paramsValueSliceLength > 0 {

            if p := strings.SplitN(s, ".", 2); len(p) > 1 {
                index, err := strconv.ParseUint(p[0], 10, 64)

                if err != nil || paramsValueSliceLength <= int(index) {
                    return nil, false
                }

                return GetParameter(p[1], params.([]interface{})[index])
            }

            index, err := strconv.ParseUint(s, 10, 64)

            if err != nil || paramsValueSliceLength <= int(index) {
                return nil, false
            }

            return params.([]interface{})[index], true
        }

        return nil, false
    }

    if p := strings.SplitN(s, ".", 2); len(p) > 1 {
        if paramsValue := reflect.ValueOf(params); paramsValue.Kind() == reflect.Map {
            if pValue, ok := params.(map[string]interface{})[p[0]]; ok {
                return GetParameter(p[1], pValue)
            }
        } else {
            return nil, false
        }
    } else {
        if pValue, ok := params.(map[string]interface{})[p[0]]; ok {
            return pValue, true
        }
    }

    return nil, false
}

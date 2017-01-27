package secret

import (
    "math/rand"
    "fmt"
)

func Generate() {
    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

    b := make([]rune, 32)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }

    fmt.Println("Generated secret: " + string(b))
}

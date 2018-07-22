package main
import "fmt"
func main(){
    fmt.Println("This is for Raspberry Pi")
    fmt.Println("Please build following:")
    fmt.Println("GOOS=linux GOARCH=arm go build server/main.go")
}

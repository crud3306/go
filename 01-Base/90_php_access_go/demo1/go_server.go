package main
  import (
      "bufio"
      "fmt"
      "io"
      "os"
      "strings"
  )
 
  func main() {
 
      inputReader := bufio.NewReader(os.Stdin)
      
      for {
          s, err := inputReader.ReadString('\n')
          if err != nil && err == io.EOF {
              break
          }
          s = strings.TrimSpace(s)
 
          if s != "" {
            fmt.Printf("hello ：%s \n", s)
            // fmt.Println("hello："+s)

          } else {
              fmt.Println("get empty \n")
          }
      }
  }
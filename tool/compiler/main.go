package main

import (
  "fmt"
  "github.com/Ankr-network/ankr-chain/tool/compiler/compile"
)

func main() {
  err := compile.CompileCmd.Execute()
  if err != nil {
    fmt.Println(err)
  }
}

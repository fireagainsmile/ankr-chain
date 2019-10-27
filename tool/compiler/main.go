package main

import (
  "fmt"

  "github.com/Ankr-network/ankr-chain/tool/compiler/root"
)

func main() {
  err := root.RootCmd.Execute()
  if err != nil {
    fmt.Println(err)
  }
}

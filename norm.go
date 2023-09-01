/*
* [norm]
* Developer: YeJun, Jung
* E-Mail: yejun614@naver.com
*
* [Description]
* macOS는 폴더명과 파일명이 터미널에서 한국어 자소, 자모 분리되어 출력되는
* 현상이 발생한다. 이는 유니코드 관련 이슈로 이를 해결하기 위해서는
* 출력 결과를 NFC형태로 unicode normalize 해주어야 정상적으로 출력된다.
*
* 해당 소스코드는 터미널 출력을 입력받아 NFC 처리후 다시 출력하는
* 프로그램으로 macOS에서 한국어 분리현상에 대응할 수 있게 해준다.
*
* [Build & Install]
* go build norm.go
* mv ./norm /usr/local/bin
*
* [Changelogs]
* 2023-02-25
*  - Create norm program with golang
*/

package main

import (
  "os"
  "io"
  "fmt"
  "flag"
  "bufio"
  "strings"
  "os/exec"
  "github.com/fatih/color"
  "golang.org/x/text/unicode/norm"
)

var (
  DebugMode bool
)

func main() {
  // flags
  flag.Usage = help
  flag.BoolVar(&DebugMode, "debug", false, "Turn on the debugging mode")
  shell := flag.String("sh", "bash -c", "shell")
  flag.Parse()

  if args := flag.Args(); len(args) > 0 {
    // run mode
    run(*shell, args)

  } else {
    // scan mode
    scan(os.Stdin)
  }
}

func run(shell string, args []string) {
  shellArgs := strings.Fields(shell)
  cmdArgs := append(shellArgs[1:], strings.Join(args, " "))

  cmd := exec.Command(shellArgs[0], cmdArgs...);
  out, err := cmd.StdoutPipe()

  if err != nil {
    panic(err)
  }

  cmd.Start()
  scan(out)
  cmd.Wait()
}

func scan(r io.Reader) {
  scanner := bufio.NewScanner(r)

  for {
    scanner.Scan()

    if str := scanner.Bytes(); len(str) != 0 {
      if DebugMode {
        buf := norm.NFC.Bytes(str)

        color.Set(color.FgWhite)
        fmt.Print("[DEBUG]")

        color.Set(color.FgCyan)
        fmt.Print(" INPUT  <<< ")
        for index, el := range str {
          if index >= len(buf) || el != buf[index] {
            color.Set(color.FgRed)
          } else {
            color.Set(color.FgCyan)
          }
          fmt.Printf("%d ", el)
        }
        fmt.Println("")

        color.Set(color.FgMagenta)
        fmt.Print("        OUTPUT >>> ")
        for index, el := range buf {
          if index >= len(str) || el != str[index] {
            color.Set(color.FgRed)
          } else {
            color.Set(color.FgMagenta)
          }
          fmt.Printf("%d ", el)
        }
        fmt.Println("")
      } else {
        os.Stdout.Write(norm.NFC.Bytes(str))
        os.Stdout.Write([]byte{ 10 }) // new line
      }
    } else {
      break
    }
  }
}

func help() {
  fmt.Fprintf(os.Stderr, "Unicode Normarlize (v.0.1)\n\n")

  fmt.Fprintf(os.Stderr, "[Usage]\n")
  flag.PrintDefaults()

  fmt.Fprintf(os.Stderr, "\n[Example]\n")
  fmt.Fprintf(os.Stderr, " > norm ls ~                    # bash -c \"ls ~\"\n")
  fmt.Fprintf(os.Stderr, " > norm -sh \"fish -c\" ls ~      # fish -c \"ls ~\"\n")
  fmt.Fprintf(os.Stderr, "\n")
  fmt.Fprintf(os.Stderr, " > alias ls=\"norm ls\"\n")
  fmt.Fprintf(os.Stderr, " > ls ~                         # norm ls ~\n")
  fmt.Fprintf(os.Stderr, " > ls -al ~                     # norm ls -al ~\n")
  fmt.Fprintf(os.Stderr, "\n")
  fmt.Fprintf(os.Stderr, " > ls ~ | norm\n")
}

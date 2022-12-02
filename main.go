package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	var keyPath, password, filePath string
	for idx, param := range os.Args {
		if param == "/f" {
			keyPath = os.Args[idx+1]
		}
		if param == "/p" {
			password = os.Args[idx+1]
		}
		if idx == (len(os.Args) - 1) {
			filePath = param
		}
	}

	if keyPath == "" {
		log.Fatal("keyPath not set")
	}

	if password == "" {
		log.Fatal("password not set")
	}

	if filePath == "" {
		log.Fatal("exe file path not set")
	}

	for idx, hash := range []string{"sha1", "sha256"} {
		args := []string{
			"/c",
			"start",
			"/wait",
			"/unix",
			"/usr/local/bin/osslsigncode", // Path for osslsigncode on OSX/Linux
			"-in",
			filePath,
			"-out",
			filePath + ".signed",
			"-pkcs12",
			keyPath,
			"-pass",
			password}

		if idx == 1 {
			args = append(args, "-nest") // add instead of replace
		}
		args = append(args, "-h", hash)

		cmd := exec.Command("cmd", args...)

		// call command
		log.Println("osslsigncode", strings.Join(args, " "))
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		// wait for result, max time = 5s
		i := 0
		for i < 10 {
			if _, err := os.Stat(filePath + ".signed"); err == nil {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}

		// check result
		if _, err := os.Stat(filePath + ".signed"); err != nil {
			log.Println("Signing", hash, "fail")
		} else {
			log.Println("Signing", hash, "success")
			if err := os.Remove(filePath); err != nil {
				log.Fatal(err)
			}

			if err := os.Rename(filePath+".signed", filePath); err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}

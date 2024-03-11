package util

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"pi-hosting/config"
)

func GetPublicIPAddress(config *config.Config) (string, error) {
	resp, err := http.Get(config.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println("IP address:", string(ip))
	return string(ip), nil
}

func CheckServiceStatus(service string) (string, error) {
	out, err := exec.Command("systemctl", "check", service).CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return "", err
		}
	}
	return string(out), nil
}

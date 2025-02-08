package utils

import (
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
)

func GetConfigPath() string {
	config := os.Getenv("XDG_CONFIG_HOME")
	if IsFileExists(config) {
		return config
	}
	config = os.Getenv("HOME") + "/.config"
	if IsFileExists(config) {
		return config
	}
	return ""
}

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func GenerateMachineId() string {
	length := 65
	hexChars := "0123456789abcdef"
	hex := ""
	for i := 0; i < length; i++ {
		index := rand.Intn(len(hexChars))
		hex += string(hexChars[index])
	}
	return hex
}

func Mkdir(path string) error {
	return os.MkdirAll(path, 0o775)
}

func LoadFileJson(p string, v interface{}) error {
	bin, err := LoadFile(p)
	if err != nil {
		return err
	}
	return json.Unmarshal(bin, v)
}

func LoadFile(p string) ([]byte, error) {
	return os.ReadFile(p)
}

func MustJsonBytes(data interface{}) []byte {
	bin, _ := json.Marshal(data)
	return bin
}

func SaveFile(p string, data interface{}) error {
	dir := filepath.Dir(p)
	_ = Mkdir(dir)

	var bin []byte

	switch t := data.(type) {
	case []byte:
		bin = t
	case string:
		bin = []byte(t)
	case io.Reader:
		f, _ := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o664)
		_, err := io.Copy(f, t)
		return err
	default:
		bin = MustJsonBytes(data)
	}

	return os.WriteFile(p, bin, 0o664)
}

func SetCurrentOSName() string {
	os := runtime.GOOS
	switch os {
	case "darwin":
		os = "Darwin"
	case "linux":
		os = "Linux"
	case "windows":
		os = "Windows"
	}
	return os
}

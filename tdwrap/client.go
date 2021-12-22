package tdwrap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/sgorblex/tgsend/utils"
	"golang.org/x/net/proxy"
)

var dataDir string
var clientFile string
var sessionFile string

func init() {
	var err error
	dataDir, err = getDataDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding data directory: %v\n", err)
		os.Exit(1)
	}
	clientFile = path.Join(dataDir, "client.json")
	sessionFile = path.Join(dataDir, "session.json")
}

func getDataDir() (string, error) {
	dataDir, err := utils.UserDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "tgsend"), nil
}

type clientInfo struct {
	ID   int    `json:"id"`
	Hash string `json:"hash"`
}

func SaveClientCreds() error {
	dataDir, err := getDataDir()
	if err != nil {
		return err
	}

	fmt.Println("You can create Client ID and hash via https://my.telegram.org/")
	fmt.Print("Client ID: ")
	var info clientInfo
	_, err = fmt.Scan(&info.ID)
	if err != nil {
		return err
	}
	fmt.Print("Client hash: ")
	_, err = fmt.Scan(&info.Hash)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return errors.New("client dir creation failed")
	}
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(clientFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func getClientCreds(dataDir string) (int, string, error) {
	var info clientInfo
	data, err := ioutil.ReadFile(filepath.Join(dataDir, "client.json"))
	if err != nil {
		return -1, "", err
	}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return -1, "", err
	}

	return info.ID, info.Hash, nil
}

func getOptions(dataDir string) (telegram.Options, error) {
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return telegram.Options{}, errors.New("session dir creation failed")
	}

	return telegram.Options{SessionStorage: &session.FileStorage{Path: sessionFile}, Resolver: dcs.Plain(dcs.PlainOptions{Dial: proxy.Dial})}, nil
}

func GetClient() (*telegram.Client, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return nil, err
	}

	id, hash, err := getClientCreds(dataDir)
	if err != nil {
		return nil, err
	}
	opts, err := getOptions(dataDir)
	if err != nil {
		return nil, err
	}
	return telegram.NewClient(id, hash, opts), nil
}

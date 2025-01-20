package app

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rivo/tview"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Auth struct {
	app  *tview.Application
	athm *AuthMemory
}

func NewAuth(app *tview.Application, athm *AuthMemory) *Auth {
	return &Auth{
		app:  app,
		athm: athm,
	}
}

func openBrowser(link string) error {
	if link == "" {
		return errors.New("no url provided")
	}
	link = strings.Trim(link, `"'`)
	var cmd string
	var args []string

	switch os := runtime.GOOS; os {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", link}
	case "darwin":
		cmd = "open"
		args = []string{link}
	default:
		cmd = "xdg-open"
		args = []string{link}
	}

	return exec.Command(cmd, args...).Start()
}

func (a *Auth) CheckIsLogged() (string, error) {
	fmt.Println("Checking IsLogged...")
	resp, err := http.Get("http://localhost:8989/auth-polling?code=" + a.athm.GetAuthId())
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("error reading response: %w", err)
		}

		fmt.Println(line)

		if line != "" && line != "data: false\n" {
			var decodedMap map[string]string
			err := json.Unmarshal([]byte(line), &decodedMap)
			if err != nil {
				return "", fmt.Errorf("failed to decode response: %w", err)
			}
			fmt.Println("User is logged in:", decodedMap)
			return decodedMap["data"], nil
		}

		// Handle server-sent "data: false"
		if line == "data: false\n" {

		}

		time.Sleep(2 * time.Second)
	}
}

func (a *Auth) getGithubOauthUrl() (string, error) {
	resp, err := http.Get("http://localhost:8989/github-url")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (a *Auth) logWithGithub(auth chan AuthState) error {
	link, err := a.getGithubOauthUrl()
	if err != nil {
		return err
	}
	fmt.Println("Log link:", link)
	parsedURL, err := url.Parse(link)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return err
	}

	queryParams := parsedURL.Query()
	fmt.Println(queryParams)
	query := queryParams.Get("termiserv")
	fmt.Println("------", query)
	a.athm.SetAuthId(query)

	err = openBrowser(link)
	if err != nil {
		return err
	}

	auth <- Islogged
	return nil
}

func (a *Auth) RenderAuth(auth chan AuthState) {
	button := tview.NewButton("Log With Github").SetSelectedFunc(func() {
		err := a.logWithGithub(auth)
		if err != nil {
			log.Println("Error logging with Github:", err)
			auth <- AuthFailed
			return
		}
	})
	a.app.EnableMouse(true)
	if err := a.app.SetRoot(button, false).Run(); err != nil {
		panic(err)
	}
}

func (a *Auth) UnsetAppRoot() {
	a.app.Stop()
}

func (a *Auth) IsAuth() bool {
	if a.athm.GetAuthId() == "" {
		return false
	}
	return true
}

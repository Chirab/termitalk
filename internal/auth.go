package app

import (
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

func (a *Auth) logWithGithub() error {
	link, err := a.getGithubOauthUrl()
	if err != nil {
		return err
	}
	parsedURL, err := url.Parse(link)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return err
	}

	queryParams := parsedURL.Query()
	query := queryParams.Get("code")
	a.athm.SetAuthId(query)

	err = openBrowser(link)
	if err != nil {
		return err
	}
	return nil
}

func (a *Auth) RenderAuth() {
	button := tview.NewButton("Log With Github").SetSelectedFunc(func() {
		err := a.logWithGithub()
		if err != nil {
			log.Println("Error logging with Github:", err)
			return
		}
	})
	a.app.EnableMouse(true)
	if err := a.app.SetRoot(button, false).Run(); err != nil {
		panic(err)
	}
}

func (a *Auth) IsAuth() bool {

	fmt.Println(a.athm.GetAuthId())
	return false
}

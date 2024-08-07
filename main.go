package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Proxy struct {
	Host string `json:"host"`
	Name string `json:"name"`
	User *User  `json:"user"`
}

type Config struct {
	Proxies []*Proxy
}

func checkAddr(addr string, p *Proxy) error {
	transport := &http.Transport{}
	if p.Host != "" {
		u := &url.URL{Host: p.Host}
		if p.User != nil {
			u.User = url.UserPassword(p.User.Username, p.User.Password)
		}
		transport = &http.Transport{
			Proxy: http.ProxyURL(u),
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   3 * time.Second,
	}

	_, err := client.Get(addr)
	return err
}

func changeStatus(address string, p *Proxy, circle *canvas.Circle) {
	if err := checkAddr(address, p); err != nil {
		circle.StrokeColor = color.RGBA{R: 255, A: 255}
	} else {
		circle.StrokeColor = color.RGBA{G: 255, A: 255}
	}
	circle.Show()
}

func checkAddress(address string, proxies []*Proxy, circles []*canvas.Circle) {
	if !strings.HasPrefix(address, "https://") && !strings.HasPrefix(address, "http://") {
		address = "https://" + address
	}
	for i, p := range proxies {
		go changeStatus(address, p, circles[i])
	}
}

func main() {
	app := app.New()
	w := app.NewWindow("Proxy Checker")

	w.Resize(fyne.NewSize(330, 330))
	w.SetFixedSize(true)

	data, err := os.ReadFile("config.json")
	if err != nil {
		return
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg.Proxies)
	if err != nil {
		fmt.Errorf("incorrect config: %w", err)
	}

	labels := make([]*widget.Label, 0, len(cfg.Proxies))
	for _, p := range cfg.Proxies {
		var label *widget.Label
		if p.Host == "" {
			label = widget.NewLabel(p.Name)
		} else {
			label = widget.NewLabel(fmt.Sprintf("%s (%s)", p.Name, p.Host))
		}
		labels = append(labels, label)
	}

	circles := make([]*canvas.Circle, 0, len(cfg.Proxies))

	for range cfg.Proxies {
		c := canvas.NewCircle(color.White)
		c.StrokeColor = color.White
		c.StrokeWidth = 5
		circles = append(circles, c)
	}

	addressInput := widget.NewEntry()
	addressInput.SetPlaceHolder("https://...")

	check := widget.NewButton("Check...", func() {
		address := addressInput.Text
		go checkAddress(address, cfg.Proxies, circles)
	})
	addressInput.OnSubmitted = func(address string) {
		go checkAddress(address, cfg.Proxies, circles)
	}
	addressInput.OnChanged = func(_ string) {
		for i := range circles {
			i := i
			go func() {
				circles[i].StrokeColor = color.White
			}()
		}
	}

	vBox := container.NewVBox(addressInput)
	for i := range cfg.Proxies {
		vBox.Add(container.NewHBox(labels[i], layout.NewSpacer(), circles[i]))
	}

	vBox.Add(check)
	w.SetContent(vBox)
	w.ShowAndRun()
}

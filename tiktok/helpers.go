package tiktok

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func generateRandomNumber() string {
	max := 1999999999999999999
	min := 1000000000000000000
	return strconv.Itoa(min + rand.Intn(max-min))
}

func replaceUnicode(URL string) string {
	return strings.ReplaceAll(URL, "\u0026", "&")
}

func saveTiktok(path string, resp *http.Response) error {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func newClient(proxy string, jar *cookiejar.Jar) *http.Client {
	client := &http.Client{Jar: jar}
	if proxy != "" {
		if !strings.Contains(proxy, "://") {
			proxy = "http://" + proxy
		}
		if proxyURL, err := url.Parse(proxy); err == nil {
			client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		} else {
			fmt.Println(err)
			fmt.Println("Not Using Proxy")
		}
	}
	return client
}

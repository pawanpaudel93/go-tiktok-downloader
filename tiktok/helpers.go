package tiktok

import (
	"io"
	"math/rand"
	"net/http"
	"os"
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

func saveTiktok(filepath string, resp *http.Response) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

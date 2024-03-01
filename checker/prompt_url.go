package checker

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	ImagePromptRef string = "https://docs.midjourney.com/docs/image-prompts"

	ErrInvalidProxyUrl            string = "Invalid proxy url."
	ErrInternalError              string = "Internal error."
	ErrInvalidImageUrl            string = "Invalid image url."
	ErrInvalidImageContentType    string = "Invalid image content type, file should end in .png, .gif, .webp, .jpg, or .jpeg. You can get more details at " + ImagePromptRef
	ErrInvalidImagePromptPosition string = "Invalid image prompt position, image prompt should go at the front of a prompt. You can get more details at " + ImagePromptRef
	ErrInvalidPromptParts         string = "Invalid prompt parts, prompts must have two images or one image and text to work. You can get more details at " + ImagePromptRef
)

func CheckImageUrl(prompt string, urls []string, proxyUrl string) error {
	if len(urls) == 0 {
		return nil
	}
	if len(urls) > 0 {
		parts := regexp.MustCompile(`[ ,]`).Split(prompt, -1)
		if err := CheckImageUrlPosition(parts, urls); err != nil {
			return err
		}
		// image url is the only content in prompt
		if len(parts) == 1 {
			return errors.New(ErrInvalidPromptParts)
		}
	}
	if err := CheckImageUrlValid(urls, proxyUrl); err != nil {
		return err
	}
	return nil
}

func CheckImageUrlPosition(parts, urls []string) error {
	// image prompt should be at the front of a prompt
	// ref: https://docs.midjourney.com/docs/image-prompts
	cnt := len(urls)
	i := 0
	for ; i < len(parts); i++ {
		if parts[i] == "" {
			continue
		}
		if !strings.HasPrefix(parts[i], "http") {
			break
		}
		cnt--
	}
	// v6 allows image urls after --sref
	for ; i < len(parts); i++ {
		if parts[i] != "--sref" {
			continue
		}
		j := i + 1
		for ; j < len(parts); j++ {
			if parts[j] == "" {
				continue
			}
			if !strings.HasPrefix(parts[j], "http") {
				break
			}
			cnt--
		}
		break
	}
	if cnt != 0 {
		return errors.New(ErrInvalidImagePromptPosition)
	}
	return nil
}

func CheckImageUrlValid(urls []string, proxyUrl string) error {
	var client *http.Client
	skipHeadRequest := false
	if proxyUrl == "" {
		// 未设置代理地址时跳过图片链接校验
		skipHeadRequest = true
	} else {
		proxy, err := url.Parse(proxyUrl)
		if err != nil {
			return errors.New(ErrInvalidProxyUrl)
		}
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}
	for _, url := range urls {
		// v6 has relative weights such as urlA::2 urlB::3 urlC::5
		// not consider of regular image urls with :: in it
		url = strings.Split(url, "::")[0]
		// file should end in .png, .gif, .webp, .jpg, or .jpeg.
		if !strings.HasSuffix(url, ".jpg") && !strings.HasSuffix(url, ".jpeg") && !strings.HasSuffix(url, ".png") && !strings.HasSuffix(url, ".gif") && !strings.HasSuffix(url, ".webp") {
			return errors.New(fmt.Sprintf("%s url: %s", ErrInvalidImageContentType, url))
		}
		if skipHeadRequest {
			continue
		}

		req, err := http.NewRequest("HEAD", url, nil)
		if err != nil {
			return errors.New(ErrInternalError)
		}
		resp, err := client.Do(req)
		if err != nil {
			return errors.New(ErrInternalError)
		}
		if resp.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("%s url: %s, head code: %d", ErrInvalidImageUrl, url, resp.StatusCode))
		}
		// contentType := resp.Header.Get("Content-Type")
		// if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" && contentType != "image/webp" {
		// 	return errors.New(fmt.Sprintf("%s: %s, content type: %s", ErrInvalidImageContentType, url, contentType))
		// }
	}
	return nil
}

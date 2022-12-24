package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	defaultCallbackTimeout = 30 * time.Second
)

func PostDo(req *http.Request, rep interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body err:%s", err.Error())
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("http.status=%s", resp.Status)
	}

	if err := json.Unmarshal(b, rep); err != nil {
		return fmt.Errorf("http.status=%s Unmarshal failed cause: %s", resp.Status, err.Error())
	}

	return nil
}

func Post(ctx context.Context, url string, reader io.Reader) error {
	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), defaultCallbackTimeout)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, reader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}

	}()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("http.status=%s", resp.Status)
	}
	return nil
}

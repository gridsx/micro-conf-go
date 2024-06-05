package conn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/grdisx/micro-conf-go/utils"
	"io"
	"net/http"
	"strings"
)

const addr = `http://%s/api/app/reg`

func (r *RegAppPayload) toPayload() io.Reader {
	d, _ := json.Marshal(r)
	return bytes.NewReader(d)
}

// RegApp TODO add sign related logic
func RegApp(payload *RegAppPayload, metaServers, token string) error {
	chosenServer := utils.ChooseAddr(metaServers)
	url := fmt.Sprintf(addr, chosenServer)
	resp, err := http.Post(url, "application/json", payload.toPayload())
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return readErr
		}
		httpResult := new(HttpResult)
		if err := json.Unmarshal(body, httpResult); err != nil {
			return err
		}
		if strings.EqualFold(httpResult.Code, "0") {
			return nil
		} else {
			return errors.New(httpResult.Msg)
		}
	} else {
		return errors.New(fmt.Sprintf("error register， code: %d", resp.StatusCode))
	}
}

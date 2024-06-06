package conn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/grdisx/micro-conf-go/auth"
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
	cli := auth.DefaultClient(payload.AppId, token)
	resp, err := cli.Post(url, payload.toPayload(), http.Header{"Content-Type": []string{"application/json"}})
	if err != nil {
		return err
	}
	httpResult := new(HttpResult)
	if err := json.Unmarshal([]byte(resp), httpResult); err != nil {
		return err
	}
	if strings.EqualFold(httpResult.Code, "0") {
		return nil
	} else {
		return errors.New(httpResult.Msg)
	}
}

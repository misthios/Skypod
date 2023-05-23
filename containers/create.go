package containers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/misthios/skypod"
)

/*
Create a container from json
returns error
*/
func Create(ctx context.Context,spec string) error{
	client, err := skypod.GetClient(ctx)
	if err !=nil {return err}

	body := strings.NewReader(spec)

	response, err := client.Request(ctx,http.MethodPost,"/containers/create",body,nil,nil)
	defer response.Body.Close()
	if err !=nil {return err}

	if response.StatusCode != 201 {
    body,err := ioutil.ReadAll(response.Body)
    if err !=nil {return err}

    cause,err := skypod.HandleApiError(body)
    if err !=nil {return err}

		return fmt.Errorf("Failed to create container: %s",cause)
	}
	return nil
}

package containers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		body, err := io.ReadAll(response.Body)
		if err !=nil {return err}

		var apierror skypod.ApiError
		err = json.Unmarshal(body,&apierror)
		if err !=nil {return err}
		return fmt.Errorf("Failed to create container: %s",apierror.Cause)
	}
	return nil
}

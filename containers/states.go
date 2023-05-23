package containers

import (
  "context"
  "fmt"
  "io/ioutil"
  "net/http"
  "github.com/misthios/skypod"
)


/*
shared code for changing states
returns error
*/
func changestate(ctx context.Context, name string,state string) error{
  client, err := skypod.GetClient(ctx)
  if err !=nil {return err}

  response, err := client.Request(ctx,http.MethodPost,"/containers/" + name+ "/" + state,nil,nil,nil)
  defer response.Body.Close()
  if err !=nil {return err}

  switch response.StatusCode{
  case 204:
    return nil
  case 304:
    return fmt.Errorf("Container is already in state %s",state)
  case 404:
    return fmt.Errorf("Container not found")
  case 500:
    body,err := ioutil.ReadAll(response.Body)
    if err !=nil {return err}

    cause,err := skypod.HandleApiError(body)
    if err !=nil {return err}

		return fmt.Errorf("Failed to change container state: %s",cause)
  }
  return nil 
 }

 func Start(ctx context.Context,name string) error{
   return changestate(ctx,name,"start")
 }

 func Stop(ctx context.Context,name string) error{
   return changestate(ctx,name,"stop")
 }

 func Pause(ctx context.Context,name string) error{
   return changestate(ctx,name,"pause")
 }
 func Unpause(ctx context.Context,name string) error{

   return changestate(ctx,name,"unpause")
 }







package spec

import (
	"encoding/json"
	"fmt"
)


type Spec struct{
	Name	 string `json:"name"`
	Image	 string `json:"image"`
}

/*
Convert a Spec to a json string
returns the string or error
*/
func ToJson(s *Spec) (string,error){
	sjson,err := json.Marshal(s)
	if err != nil {return "",fmt.Errorf("Failed to marshal spec: %w",err)}

	return string(sjson),nil
}


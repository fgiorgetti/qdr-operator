package entities

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Connection represents a Dispatch Router connection
type Connection struct {
	EntityCommon
	Active bool `json:"active"`
	AdminStatus AdminStatusType `json:"adminStatus,string"`
	OperStatus OperStatusType `json:"operStatus,string"`
	Container string `json:"container"`
	Opened bool `json:"opened"`
	Host string `json:"host"`
	Direction DirectionType `json:"dir,string"`
	Role string `json:"role"`
	IsAuthenticated bool `json:"isAuthenticated"`
	IsEncrypted bool `json:"isEncrypted"`
	Sasl string `json:"sasl"`
	User string `json:"user"`
	Ssl bool `json:"ssl"`
	SslProto string `json:"sslProto"`
	SslCipher string `json:"sslCipher"`
	SslSsf int `json:"sslSsf"`
	Tenant string `json:"tenant"`
	Properties map[string]string `json:"properties"`
}

func (Connection) GetEntityId() string {
	return "connection"
}

type AdminStatusType int
const (
	AdminStatusEnabled AdminStatusType = iota
	AdminStatusDeleted
)

func (a *AdminStatusType) UnmarshalJSON(b []byte) error {
	var s string

	if len(b) == 0 {
		return nil
	}
	if b[0] != '"' {
		b = []byte(strconv.Quote(string(b)))
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "enabled":
		*a = AdminStatusEnabled
	case "deleted":
		*a = AdminStatusDeleted
	}
	return nil
}

func (a AdminStatusType) MarshalJSON() ([]byte, error) {
	var s string
	switch a {
	case AdminStatusEnabled:
		s = "enabled"
	case AdminStatusDeleted:
		s = "deleted"
	}
	return json.Marshal(s)
}

type OperStatusType int
const (
	OperStatusUp OperStatusType = iota
	OperStatusClosing
)

func (o *OperStatusType) UnmarshalJSON(b []byte) error {
	var s string

	if len(b) == 0 {
		return nil
	}
	if b[0] != '"' {
		b = []byte(strconv.Quote(string(b)))
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "up":
		*o = OperStatusUp
	case "closing":
		*o = OperStatusClosing
	}
	return nil
}

func (o OperStatusType) MarshalJSON() ([]byte, error) {
	var s string
	switch o {
	case OperStatusUp:
		s = "up"
	case OperStatusClosing:
		s = "closing"
	}
	return json.Marshal(s)
}

type DirectionType int
const (
	DirectionTypeIn DirectionType = iota
	DirectionTypeOut
)

func (d *DirectionType) UnmarshalJSON(b []byte) error {
	var s string

	if len(b) == 0 {
		return nil
	}
	if b[0] != '"' {
		b = []byte(strconv.Quote(string(b)))
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "in":
		*d = DirectionTypeIn
	case "out":
		*d = DirectionTypeOut
	}
	return nil
}

func (d DirectionType) MarshalJSON() ([]byte, error) {
	var s string
	switch d {
	case DirectionTypeIn:
		s = "in"
	case DirectionTypeOut:
		s = "out"
	}
	return json.Marshal(s)
}

//const input string = `[
// {
//   "opened": true,
//   "adminStatus": "enabled",
//   "container": "6a32521a-4850-44dd-b9cf-4cfe76b58f36",
//   "name": "connection/127.0.0.1:34594",
//   "operStatus": "up",
//   "ssl": false,
//   "host": "127.0.0.1:34594",
//   "isEncrypted": false,
//   "role": "normal",
//   "identity": "4",
//   "isAuthenticated": false,
//   "active": true,
//   "sslSsf": 0,
//   "type": "org.apache.qpid.dispatch.connection",
//   "properties": {},
//   "dir": "in",
//   "user": "anonymous"
// }
//]`
//
//func main() {
//	var connections2 []Connection
//	json.Unmarshal([]byte(input), &connections2)
//	for i, v := range connections2 {
//		fmt.Printf("Connection[%d]\n", i)
//		fmt.Printf("\tOpened: %v\n", v.Opened)
//		fmt.Printf("\tAdminStatus: %v\n", v.AdminStatus)
//		fmt.Printf("\tContainer: %v\n", v.Container)
//		fmt.Printf("\tName: %v\n", v.Name)
//		fmt.Printf("\tOperStatus: %v\n", v.OperStatus)
//		fmt.Printf("\tSsl: %v\n", v.Ssl)
//		fmt.Printf("\tHost: %v\n", v.Host)
//		fmt.Printf("\tIsEncrypted: %v\n", v.IsEncrypted)
//		fmt.Printf("\tRole: %v\n", v.Role)
//		fmt.Printf("\tIdentity: %v\n", v.Identity)
//		fmt.Printf("\tIsAuthenticated: %v\n", v.IsAuthenticated)
//		fmt.Printf("\tActive: %v\n", v.Active)
//		fmt.Printf("\tsslSsf: %v\n", v.SslSsf)
//		fmt.Printf("\tType: %v\n", v.Type)
//		fmt.Printf("\tProperties: %v\n", v.Properties)
//		fmt.Printf("\tDir: %v\n", v.Direction)
//		fmt.Printf("\tUser: %v\n", v.User)
//	}
//
//}

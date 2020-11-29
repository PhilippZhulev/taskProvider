package handler

import (
	base64 "encoding/base64"
	"encoding/json"
	"log"
	"strings"

	"gitlab.com/taskProvider/services/user/internal/app/store/sqlstore"
	pb "gitlab.com/taskProvider/services/user/proto/user"
)

// HandleUserGet ...
func (init Init) HandleUserGet(in *pb.UserGetRequest, store *sqlstore.Store) ([]string, error) {

	type Clams struct {
		UUID string `json:"uuid"`
		ID string `json:"user_id"`
		Exp int `json:"exp"`
	}

	return func() ([]string, error) {
		jwtBase64 := strings.Split(in.GetData(), ".")[1]

		jwtJSON, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(jwtBase64)
		if err != nil {
			return nil, err;
		}

		c := Clams{}
		err = json.Unmarshal(jwtJSON, &c)
		if err != nil {
			return nil, err;
		}

		u, err := store.User().FindByKey("uuid", c.UUID)
		if err != nil {
			log.Println(err)
			return nil, err;
		}

		return []string{
			u.Login, 
			u.Name, 
			u.Email,
			c.ID,
		}, nil
	}()
}
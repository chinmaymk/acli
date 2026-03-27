package bitbucket

import "encoding/json"

type CurrentUser struct {
	UUID        string `json:"uuid"`
	Nickname    string `json:"nickname"`
	DisplayName string `json:"display_name"`
	AccountID   string `json:"account_id"`
}

func (c *Client) GetCurrentUser() (*CurrentUser, error) {
	data, err := c.get("/user")
	if err != nil {
		return nil, err
	}
	var user CurrentUser
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

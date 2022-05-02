package cache

import "strconv"

const (
	adminSsoPre   = "asso:" // +id
	adminTokenPre = "atk:"  // +token
)

func AdminSsoKey(id int) string { return adminSsoPre + strconv.Itoa(id) }

func AdminTokenKey(token string) string { return adminTokenPre + token }

type AdminToken struct {
	ID       int      `json:"i"`
	Username string   `json:"u"`
	IsSuper  bool     `json:"s,omitempty"`
	Actions  []string `json:"a,omitempty"`
}

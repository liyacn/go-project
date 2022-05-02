package cache

import "strconv"

const (
	adminSsoPre         = "admin:sso:"          // +id
	adminTokenPre       = "admin:token:"        // +token
	adminRoleActionsPre = "admin:role:actions:" // +role_id, value:json([]action)
)

func AdminSsoKey(id int) string { return adminSsoPre + strconv.Itoa(id) }

func AdminTokenKey(token string) string { return adminTokenPre + token }

type AdminToken struct {
	ID       int    `json:"i"`
	Username string `json:"u"`
	IsSuper  bool   `json:"s,omitempty"`
	RoleID   int    `json:"r,omitempty"`
}

func AdminRoleActionsKey(id int) string { return adminRoleActionsPre + strconv.Itoa(id) }

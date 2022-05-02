package proto

import (
	"project/model/entity"
)

type SystemActionSyncData struct {
	Delete []string
	Create []*entity.SystemAction
}

type SystemActionMenuListResp struct {
	List []*SystemActionGroup `json:"list"`
}

type SystemActionGroup struct {
	*entity.SystemAction
	Sub []*entity.SystemAction `json:"sub"`
}

type SystemActionUpdateArgs struct {
	KeyName string `json:"key_name" binding:"required"`
	Title   string `json:"title" binding:"max=16"`
	Sort    int    `json:"sort" binding:"min=0,max=9999"`
}

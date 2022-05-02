package handler

import (
	"github.com/gin-gonic/gin"
	"project/cms/internal/proto"
	"project/model/entity"
	"project/pkg/logger"
	"project/pkg/types"
	"project/pkg/web/errcode"
	"sort"
	"strings"
)

func (h *Handler) SystemActionSync(g *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		aks, err := h.service.SystemActionKeyNames(c, false)
		if err != nil {
			logger.FromContext(c).Error("service.SystemActionKeyNames error", err)
			c.JSON(errcode.FromError(err))
			return
		}
		actionMap := make(map[string]struct{})
		groupMap := make(map[string]struct{})
		for _, v := range aks {
			if strings.HasPrefix(v, "/") {
				actionMap[v] = struct{}{}
			} else {
				groupMap[v] = struct{}{}
			}
		}
		groupSet := make(map[string]struct{})
		routes := g.Routes()
		actionCreate := make([]*entity.SystemAction, 0, len(routes))
		sortNum := 1
		for _, v := range routes {
			paths := strings.Split(v.Path[1:], "/")
			if len(paths) < 2 {
				continue
			}
			groupSet[paths[0]] = struct{}{}
			if _, ok := actionMap[v.Path]; ok { // 存在则跳过
				delete(actionMap, v.Path)
			} else { // 不存在则插入
				actionCreate = append(actionCreate, &entity.SystemAction{
					KeyName: v.Path,
					Sort:    sortNum,
				})
			}
			sortNum++
		}
		sortNum = 1
		for key := range groupSet {
			if _, ok := groupMap[key]; ok { // 存在则跳过
				delete(groupMap, key)
			} else { // 不存在则插入
				actionCreate = append(actionCreate, &entity.SystemAction{
					KeyName: key,
					Sort:    sortNum,
				})
			}
			sortNum++
		}
		actionDelete := append(types.Keys(actionMap), types.Keys(groupMap)...) // 多余的删除
		err = h.service.SystemActionSave(c, &proto.SystemActionSyncData{
			Delete: actionDelete,
			Create: actionCreate,
		})
		if err != nil {
			logger.FromContext(c).Error("service.SystemActionSave error", err)
			c.JSON(errcode.FromError(err))
			return
		}
		h.SystemActionList(c)
	}
}

func (h *Handler) SystemActionList(c *gin.Context) {
	data, err := h.service.SystemActionList(c)
	if err != nil {
		logger.FromContext(c).Error("service.SystemActionList error", err)
		c.JSON(errcode.FromError(err))
		return
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Sort < data[j].Sort
	})
	var list []*proto.SystemActionGroup
	routes := make(map[string][]*entity.SystemAction)
	for _, v := range data {
		if strings.HasPrefix(v.KeyName, "/") {
			group := v.KeyName[1 : strings.Index(v.KeyName[1:], "/")+1]
			routes[group] = append(routes[group], v)
		} else {
			list = append(list, &proto.SystemActionGroup{SystemAction: v})
		}
	}
	for _, v := range list {
		v.Sub = routes[v.KeyName]
	}
	c.JSON(OK, &proto.SystemActionMenuListResp{List: list})
}

func (h *Handler) SystemActionUpdate(c *gin.Context) {
	var r proto.SystemActionUpdateArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.FromError(err))
		return
	}
	err := h.service.SystemActionUpdate(c, &entity.SystemAction{
		KeyName: r.KeyName,
		Title:   r.Title,
		Sort:    r.Sort,
	})
	if err != nil {
		logger.FromContext(c).Error("service.SystemActionUpdate error", &r, err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, nil)
}

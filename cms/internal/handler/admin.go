package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"github.com/gin-gonic/gin"
	"project/cms/internal/proto"
	"project/model/cache"
	"project/model/entity"
	"project/pkg/logger"
	"project/pkg/random"
	"project/pkg/types"
	"project/pkg/web/errcode"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) captchaSign(exp, code string) string {
	mac := hmac.New(sha1.New, []byte(h.captcha))
	mac.Write([]byte(exp))
	mac.Write([]byte(strings.ToUpper(code)))
	return base32.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (h *Handler) Captcha(c *gin.Context) {
	code := random.Strings(4)
	bin := h.drawer.Draw(code)
	exp := strconv.FormatInt(time.Now().Unix()+65, 10)
	key := exp + "." + h.captchaSign(exp, code)
	c.JSON(OK, &proto.CaptchaResp{
		SessionKey:  key,
		Base64Image: bin,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var r proto.LoginArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.ParamsInvalid.Response())
		return
	}

	exp, sign, ok := strings.Cut(r.SessionKey, ".")
	if !ok {
		c.JSON(errcode.ParamsInvalid.Response())
		return
	}
	if e, _ := strconv.ParseInt(exp, 10, 64); e < time.Now().Unix() {
		c.JSON(errcode.CaptchaExpired.Response())
		return
	}
	if h.captchaSign(exp, r.Captcha) != sign {
		c.JSON(errcode.CaptchaWrong.Response())
		return
	}

	user, err := h.service.AdminUserFindByName(c, r.Username)
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserFindByName error", r.Username, err)
		c.JSON(errcode.FromError(err))
		return
	}
	if user.ID == 0 || !user.CheckPassword(r.Password) {
		c.JSON(errcode.UserOrPwdWrong.Response())
		return
	}
	if user.Status != entity.StatusEnabled {
		c.JSON(errcode.AccountDisabled.Response())
		return
	}
	pwdExp := user.PwdExp - time.Now().Unix()
	if pwdExp < 0 {
		c.JSON(errcode.PasswordExpired.Response())
		return
	}

	store := &cache.AdminToken{
		ID:       user.ID,
		Username: user.Username,
	}
	resp := &proto.LoginResp{Username: user.Username}
	if user.IsSuper() {
		store.IsSuper = true
		resp.IsSuper = true
	} else {
		role, err := h.service.AdminRoleFindByID(c, user.RoleID)
		if err != nil {
			logger.FromContext(c).Error("service.AdminRoleFindByID error", user.RoleID, err)
			c.JSON(errcode.FromError(err))
			return
		}
		store.RoleID = user.RoleID
		resp.Actions = role.Actions
		resp.Menus = role.Menus
	}
	token, err := h.service.AdminTokenSet(c, store)
	if err != nil {
		logger.FromContext(c).Error("service.AdminTokenSet error", user, err)
		c.JSON(errcode.FromError(err))
		return
	}
	resp.Token = token

	if day := (pwdExp + 43200) / 86400; day < 2 {
		resp.PwdTip = "密码即将过期，请立即修改！"
	} else if day < 16 {
		resp.PwdTip = "密码将于" + strconv.FormatInt(day, 10) + "天内过期，请尽快修改！"
	}
	c.JSON(OK, resp)
}

func (h *Handler) Logout(c *gin.Context) {
	auth := getAuth(c)
	err := h.service.AdminUserLogout(c, auth.ID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserLogout error", err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, nil)
}

func (h *Handler) Password(c *gin.Context) {
	var r proto.PasswordArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.ParamsInvalid.Response())
		return
	}
	auth := getAuth(c)
	user, err := h.service.AdminUserFindByID(c, auth.ID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserFindByID error", auth.ID, err)
		c.JSON(errcode.FromError(err))
		return
	}
	if user.CheckPassword(r.Password) {
		c.JSON(errcode.NeedModified.WithMsg("新密码不能与原密码相同"))
		return
	}
	err = h.service.AdminUserUpdate(c, &entity.AdminUser{
		ID:       auth.ID,
		Password: r.Password,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserUpdate error", auth.ID, err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, nil)
}

func (h *Handler) AdminRoleList(c *gin.Context) {
	var r proto.ListArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.ParamsInvalid.Response())
		return
	}
	total, list, err := h.service.AdminRoleList(c, &r)
	if err != nil {
		logger.FromContext(c).Error("service.AdminRoleList error", &r, err)
		c.JSON(errcode.FromError(err))
		return
	}
	items := make([]*proto.AdminRoleItem, 0, len(list))
	for _, v := range list {
		items = append(items, &proto.AdminRoleItem{
			ID:      v.ID,
			Name:    v.Name,
			Actions: v.Actions,
			Menus:   v.Menus,
		})
	}
	c.JSON(OK, &proto.AdminRoleListResp{
		Total: total,
		List:  items,
	})
}

func (h *Handler) AdminRoleOption(c *gin.Context) {
	data, err := h.service.AdminRoles(c)
	if err != nil {
		logger.FromContext(c).Error("service.AdminRoles error", err)
		c.JSON(errcode.FromError(err))
		return
	}
	opt := make([]*proto.OptionItem, 0, len(data))
	for _, d := range data {
		opt = append(opt, &proto.OptionItem{
			ID:   d.ID,
			Name: d.Name,
		})
	}
	c.JSON(OK, &proto.OptionResp{List: opt})
}

func (h *Handler) AdminRoleSave(c *gin.Context) {
	var r proto.AdminRoleSaveArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.ParamsInvalid.Response())
		return
	}
	allActions, err := h.service.SystemActionKeyNames(c, true)
	if err != nil {
		logger.FromContext(c).Error("service.SystemActionKeyNames error", err)
		c.JSON(errcode.FromError(err))
		return
	}
	actionMap := types.SliceToMap(allActions)
	var actions []string
	for _, v := range r.Actions {
		if _, ok := actionMap[v]; ok {
			actions = append(actions, v)
			delete(actionMap, v)
		}
	}
	menus := types.Deduplicate(r.Menus)

	if r.ID != 0 {
		role, err := h.service.AdminRoleFindByID(c, r.ID)
		if err != nil {
			logger.FromContext(c).Error("service.AdminRoleFindByID error", r.ID, err)
			c.JSON(errcode.FromError(err))
			return
		}
		if role.ID == 0 {
			c.JSON(errcode.AssociatedIdInvalid.Response())
			return
		}
		if r.Name == role.Name &&
			types.CheckEqualSet(actions, role.Actions) &&
			types.CheckEqualSet(menus, role.Menus) {
			c.JSON(OK, nil)
			return
		}
	}
	err = h.service.AdminRoleSave(c, &entity.AdminRole{
		ID:      r.ID,
		Name:    r.Name,
		Actions: actions,
		Menus:   menus,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AdminRoleSave error", err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, nil)
}

func (h *Handler) AdminUserList(c *gin.Context) {
	var r proto.AdminUserListArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.ParamsInvalid.Response())
		return
	}
	total, list, err := h.service.AdminUserList(c, &r)
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserList error", &r, err)
		c.JSON(errcode.FromError(err))
		return
	}
	items := make([]*proto.AdminUserItem, 0, len(list))
	for _, v := range list {
		items = append(items, &proto.AdminUserItem{
			ID:       v.ID,
			Username: v.Username,
			RoleID:   v.RoleID,
			IsSuper:  v.IsSuper(),
			Status:   v.Status,
		})
	}
	c.JSON(OK, &proto.AdminUserListResp{
		Total: total,
		List:  items,
	})
}

func (h *Handler) AdminUserCreate(c *gin.Context) {
	var r proto.AdminUserCreateArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.ParamsInvalid.Response())
		return
	}
	role, err := h.service.AdminRoleFindByID(c, r.RoleID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminRoleFindByID error", r.RoleID, err)
		c.JSON(errcode.FromError(err))
		return
	}
	if role.ID == 0 {
		c.JSON(errcode.AssociatedIdInvalid.Response())
		return
	}
	ok, err := h.service.AdminUserCreate(c, &entity.AdminUser{
		Username: r.Username,
		Password: r.Password,
		RoleID:   r.RoleID,
		Status:   entity.StatusEnabled,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserCreate error", r.RoleID, err)
		c.JSON(errcode.FromError(err))
		return
	}
	if !ok {
		c.JSON(errcode.UniqueKeyExist.WithMsg("用户名已存在"))
		return
	}
	c.JSON(OK, nil)
}

func (h *Handler) AdminUserPassword(c *gin.Context) {
	var r proto.AdminUserPasswordArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.ParamsInvalid.Response())
		return
	}
	user, err := h.service.AdminUserFindByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserFindByID error", r.ID, err)
		c.JSON(errcode.FromError(err))
		return
	}
	if user.ID == 0 {
		c.JSON(errcode.AssociatedIdInvalid.Response())
		return
	}
	if user.IsSuper() {
		c.JSON(errcode.OperationDeny.Response())
		return
	}
	err = h.service.AdminUserUpdate(c, &entity.AdminUser{
		ID:       r.ID,
		Password: r.Password,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserUpdate error", &r, err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, nil)
}

func (h *Handler) AdminUserRole(c *gin.Context) {
	var r proto.AdminUserRoleArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.ParamsInvalid.Response())
		return
	}
	user, err := h.service.AdminUserFindByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserFindByID error", r.ID, err)
		c.JSON(errcode.FromError(err))
		return
	}
	if user.ID == 0 {
		c.JSON(errcode.AssociatedIdInvalid.Response())
		return
	}
	if user.IsSuper() {
		c.JSON(errcode.OperationDeny.Response())
		return
	}
	if user.RoleID == r.RoleID {
		c.JSON(OK, nil)
		return
	}
	role, err := h.service.AdminRoleFindByID(c, r.RoleID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminRoleFindByID error", r.RoleID, err)
		c.JSON(errcode.FromError(err))
		return
	}
	if role.ID == 0 {
		c.JSON(errcode.AssociatedIdInvalid.Response())
		return
	}
	err = h.service.AdminUserUpdate(c, &entity.AdminUser{
		ID:     r.ID,
		RoleID: r.RoleID,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserUpdate error", &r, err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, nil)
}

func (h *Handler) AdminUserStatus(c *gin.Context) {
	var r proto.SwitchStatusArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.ParamsInvalid.Response())
		return
	}
	user, err := h.service.AdminUserFindByID(c, r.ID)
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserFindByID error", r.ID, err)
		c.JSON(errcode.FromError(err))
		return
	}
	if user.ID == 0 {
		c.JSON(errcode.AssociatedIdInvalid.Response())
		return
	}
	if user.IsSuper() {
		c.JSON(errcode.OperationDeny.Response())
		return
	}
	if user.Status == r.Status {
		c.JSON(OK, nil)
		return
	}
	err = h.service.AdminUserUpdate(c, &entity.AdminUser{
		ID:     r.ID,
		Status: r.Status,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AdminUserUpdate error", &r, err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, nil)
}

package auth

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	kitcache "github.com/icowan/kit-cache"
	"github.com/jtblin/go-ldap-client"
	"github.com/pkg/errors"

	"github.com/kplcloud/kplcloud/src/encode"
	kpljwt "github.com/kplcloud/kplcloud/src/jwt"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util"
)

type Middleware func(Service) Service

type Service interface {
	// Login 登陆
	Login(ctx context.Context, username, password string) (rs string, sessionTimeout int64, err error)
	// Register 注册用户
	Register(ctx context.Context, username, email, password, mobile, remark string) (err error)
	// AuthLoginGithub github 授权登陆跳转
	//AuthLoginGithub(w http.ResponseWriter, r *http.Request)
	//// AuthLoginGithubCallback github 授权登陆回调
	//AuthLoginGithubCallback(w http.ResponseWriter, r *http.Request)
	//// AuthLoginType 是否启用第三方授权登陆
	//AuthLoginType(ctx context.Context) string
}

type service struct {
	logger         log.Logger
	repository     repository.Repository
	appKey         string
	sessionTimeout int64
	cache          kitcache.Service
	traceId        string
}

func (s *service) Register(ctx context.Context, username, email, password, mobile, remark string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId))

	_, err = s.repository.SysUser().FindByEmail(ctx, username)
	if err == nil {
		_ = level.Warn(logger).Log("repository.SysUser", "FindByEmail")
		return encode.ErrAuthRegisterExists.Error()
	}
	if !gorm.IsRecordNotFoundError(err) {
		_ = level.Error(logger).Log("repository.SysUser", "FindByEmail", "err", err.Error())
		return encode.ErrAuthRegisterExists.Error()
	}
	var loginName string
	loginName = strings.Split(email, "@")[0]
	passwordHashed := util.EncodePassword(password, s.appKey)
	exp := time.Now().AddDate(1, 0, 0)
	err = s.repository.SysUser().Save(ctx, &types.SysUser{
		Username:  username,
		Mobile:    mobile,
		LoginName: strings.ToLower(loginName),
		Email:     email,
		Password:  passwordHashed,
		Locked:    false,
		Remark:    remark,
		ExpiresAt: &exp,
	})
	if err != nil {
		_ = level.Error(logger).Log("repository.SysUser", "Save", "err", err.Error())
		return encode.ErrAuthRegisterExists.Error()
	}

	return
}

func (s *service) Login(ctx context.Context, username, password string) (rs string, sessionTimeout int64, err error) {
	sysUser, err := s.repository.SysUser().FindByEmail(ctx, username)
	if err != nil {
		// 用户名或密码错误
		err = encode.ErrAccountLogin.Error()
		return
	}
	passwordHashed := util.EncodePassword(password, s.appKey)
	if !strings.EqualFold(sysUser.Password, passwordHashed) {
		// 用户名或密码错误
		err = encode.ErrAccountLogin.Error()
		return
	}
	if sysUser.Locked {
		// 用户已锁定
		err = encode.ErrAccountLocked.Error()
		return
	}
	sessionTimeout = s.sessionTimeout
	rs, err = s.jwtToken(ctx, sysUser)
	go func(sysUser *types.SysUser) {
		t := time.Now()
		sysUser.LastLogin = &t
		if e := s.repository.SysUser().Save(ctx, sysUser); e != nil {
			_ = level.Error(s.logger).Log("repository.SysUse", "Save", "err", e.Error())
		}
	}(&sysUser)
	return
}

/**
 * @Title ldap登陆
 */
func (s *service) ldapLogin(email, password string) (ok bool, user map[string]string, err error) {
	client := &ldap.LDAPClient{
		//Base:         c.config.GetString("ldap", "ldap_base"),
		//Host:         c.config.GetString("ldap", "ldap_host"),
		//Port:         c.config.GetInt("ldap", "ldap_port"),
		//UseSSL:       c.config.GetBool("ldap", "ldap_useSSL"),
		//BindDN:       c.config.GetString("ldap", "ldap_bindDN"),
		//BindPassword: c.config.GetString("ldap", "ldap_bind_password"),
		//UserFilter:   c.config.GetString("ldap", "ldap_user_filter"),
		//GroupFilter:  c.config.GetString("ldap", "ldap_group_filter"),
		//Attributes:   c.config.GetStrings("ldap", "ldap_attr"),
	}

	defer client.Close()

	return client.Authenticate(email, password)
}

func (s *service) jwtToken(ctx context.Context, sysUser types.SysUser) (tk string, err error) {
	timeout := time.Duration(s.sessionTimeout) * time.Second
	expAt := time.Now().Add(timeout).Unix()

	// 创建声明
	claims := kpljwt.ArithmeticCustomClaims{
		UserId:  sysUser.Id,
		IsAdmin: false,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAt,
			Issuer:    "system",
		},
	}

	//创建token，指定加密算法为HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//生成token
	tk, err = token.SignedString([]byte(kpljwt.GetJwtKey()))
	if err != nil {
		return tk, nil
	}

	var namespaces []string
	//var groups []int64
	var roleIds []int64
	var clusters []string
	var permissions []types.SysPermission

	for _, ns := range sysUser.SysNamespaces {
		namespaces = append(namespaces, ns.Name)
	}
	for _, role := range sysUser.SysRoles {
		roleIds = append(roleIds, role.Id)
		// TODO: 去重
		permissions = append(permissions, role.SysPermissions...)
	}
	for _, v := range sysUser.Clusters {
		clusters = append(clusters, v.Name)
	}

	if err = s.cache.Set(ctx, fmt.Sprintf("user:%d:info", sysUser.Id), sysUser, timeout); err != nil {
		err = encode.ErrAuthLogin.Wrap(errors.Wrap(err, "info"))
		return tk, err
	}
	if err = s.cache.Set(ctx, fmt.Sprintf("user:%d:permissions", sysUser.Id), permissions, timeout); err != nil {
		err = encode.ErrAuthLogin.Wrap(errors.Wrap(err, "permissions"))
		return tk, err
	}
	if err = s.cache.Set(ctx, fmt.Sprintf("user:%d:clusters", sysUser.Id), clusters, timeout); err != nil {
		err = encode.ErrAuthLogin.Wrap(errors.Wrap(err, "clusters"))
		return tk, err
	}
	if err = s.cache.Set(ctx, fmt.Sprintf("user:%d:namespaces", sysUser.Id), roleIds, timeout); err != nil {
		err = encode.ErrAuthLogin.Wrap(errors.Wrap(err, "namespaces"))
		return tk, err
	}
	if err = s.cache.Set(ctx, fmt.Sprintf("login:%d:token", sysUser.Id), tk, timeout); err != nil {
		err = encode.ErrAuthLogin.Wrap(errors.Wrap(err, "token"))
		return tk, err
	}

	return
}

func New(logger log.Logger, traceId string, store repository.Repository, cacheSvc kitcache.Service, appKey string, sessionTimeout int64) Service {
	return &service{
		logger:         logger,
		repository:     store,
		appKey:         appKey,
		sessionTimeout: sessionTimeout,
		cache:          cacheSvc,
		traceId:        traceId,
	}
}

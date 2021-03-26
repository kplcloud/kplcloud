package auth

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/go-github/v26/github"
	"github.com/icowan/config"
	"github.com/jtblin/go-ldap-client"
	kplcasbin "github.com/kplcloud/kplcloud/src/casbin"
	kpljwt "github.com/kplcloud/kplcloud/src/jwt"
	"github.com/kplcloud/kplcloud/src/middleware"
	"github.com/kplcloud/kplcloud/src/repository"
	"github.com/kplcloud/kplcloud/src/repository/types"
	"github.com/kplcloud/kplcloud/src/util/encode"
	"golang.org/x/oauth2"
	oauthgithub "golang.org/x/oauth2/github"
	"gopkg.in/guregu/null.v3"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidArgument            = errors.New("invalid argument")
	ErrUserOrPassword             = errors.New("邮箱或密码错误.")
	ErrUserStateFail              = errors.New("账号受限无法登陆.")
	ErrAuthLoginDefaultNamespace  = errors.New("默认空间不存在,请在app.cfg配置文件设置默认空间.")
	ErrAuthLoginDefaultRoleID     = errors.New("默认角色不存在,请在app.cfg配置文件设置默认角色ID.")
	ErrAuthLoginGitHubGetUser     = errors.New("获取Github用户邮箱及名称失败.")
	ErrAuthLoginGitHubPublicEmail = errors.New("请您在您的Github配置您的Github公共邮箱，否则无法进行授权。在 https://github.com/settings/profile 选择 public email 后重新进行授权")
)

const (
	LoginTypeLDAP = "ldap"
	UserStateFail = 2
)

type Service interface {
	// 登陆
	Login(ctx context.Context, email, password string) (rs string, err error)

	// 解析Token
	ParseToken(ctx context.Context, token string) (map[string]interface{}, error)

	// github 授权登陆跳转
	AuthLoginGithub(w http.ResponseWriter, r *http.Request)

	// github 授权登陆回调
	AuthLoginGithubCallback(w http.ResponseWriter, r *http.Request)

	// 是否启用第三方授权登陆
	AuthLoginType(ctx context.Context) string
}

type service struct {
	logger     log.Logger
	config     *config.Config
	casbin     kplcasbin.Casbin
	repository repository.Repository
}

func (c *service) AuthLoginType(ctx context.Context) string {
	return c.config.GetString("server", "auth_login")
}

func (c *service) AuthLoginGithub(w http.ResponseWriter, r *http.Request) {
	githubOauthConfig := c.auth2Config()

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)
	u := githubOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (c *service) auth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.config.GetString("server", "client_id"),
		ClientSecret: c.config.GetString("server", "client_secret"),
		Scopes:       []string{"SCOPE1", "SCOPE2", "user:email"},
		Endpoint:     oauthgithub.Endpoint,
	}
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	// todo 用jwt生成 然后 jwt 解析出来
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func (c *service) AuthLoginGithubCallback(w http.ResponseWriter, r *http.Request) {
	var resp authResponse

	ctx := context.Background()
	// state := r.URL.Query().Get("state") // todo 它需要验证一下可以考虑使用jwt生成  先用cookie 简单处理一下吧...
	if httpProxy := c.config.GetString(config.SectionServer, "http_proxy"); httpProxy != "" {
		_ = level.Debug(c.logger).Log("use-proxy", httpProxy)
		dialer := &net.Dialer{
			Timeout:   time.Duration(5 * int64(time.Second)),
			KeepAlive: time.Duration(5 * int64(time.Second)),
		}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
			Transport: &http.Transport{
				Proxy: func(_ *http.Request) (*url.URL, error) {
					return url.Parse(httpProxy)
				},
				DialContext: dialer.DialContext,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
		})
	}

	githubOauthConfig := c.auth2Config()

	if r.URL.Query().Get("error") != "" {
		resp.Err = errors.New(r.URL.Query().Get("error") + ": " + r.URL.Query().Get("error_description"))
		_ = encodeLoginResponse(ctx, w, resp)
		return
	}

	token, err := githubOauthConfig.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		if strings.Contains(err.Error(), "server response missing access_token") {
			http.Redirect(w, r, c.config.GetString("server", "domain")+"/#/user/login", http.StatusPermanentRedirect)
		}
		_ = level.Error(c.logger).Log("githubOauthConfig", "Exchange", "err", err.Error())
		resp.Err = err
		_ = encodeLoginResponse(ctx, w, resp)
		return
	}

	if token == nil || !token.Valid() {
		_ = level.Error(c.logger).Log("token", "nil", "or", "token.valid is false")
		resp.Err = errors.New("token is nil or token.valid is false")
		_ = encodeLoginResponse(ctx, w, resp)
		return
	}

	client := github.NewClient(githubOauthConfig.Client(ctx, token))
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		_ = level.Error(c.logger).Log("client.users", "Get", "err", err.Error())
		resp.Err = err
		_ = encodeLoginResponse(ctx, w, resp)
		return
	}

	if user == nil {
		resp.Err = ErrAuthLoginGitHubGetUser
		_ = encodeLoginResponse(ctx, w, resp)
		return
	}

	if user.GetEmail() == "" {
		resp.Err = ErrAuthLoginGitHubPublicEmail
		_ = encodeLoginResponse(ctx, w, resp)
		return
	}

	username := user.GetName()
	if username == "" {
		username = user.GetLogin()
	}

	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		_ = level.Warn(c.logger).Log("invalid", "oauth github state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rs, member, namespaces, err := c.AuthLogin(user.GetEmail(), username)
	if err != nil {
		resp = authResponse{Err: err}
		_ = encodeLoginResponse(ctx, w, resp)
		return
	}

	//_ = c.casbin.GetEnforcer().LoadPolicy()

	params := url.Values{}
	params.Add("token", rs)
	params.Add("email", member.Email)
	params.Add("username", member.Username)
	//params.Add("namespaces", strings.Join(namespaces, ","))
	params.Add("namespace", namespaces[0])

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    rs,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   7200})

	http.Redirect(w, r, c.config.GetString("server", "domain")+"/#/user/login?"+params.Encode(), http.StatusPermanentRedirect)

}

func (c *service) AuthLogin(email, username string) (rs string, member *types.Member, nss []string, err error) {
	member, err = c.repository.Member().Find(email)

	if member == nil || err != nil {
		ns, err := c.repository.Namespace().Find(c.config.GetString("server", "default_namespace"))
		if err != nil {
			_ = level.Error(c.logger).Log("namespaceRepository", "Find", "err", err.Error())
			return "", nil, nil, ErrAuthLoginDefaultNamespace
		}
		role, err := c.repository.Role().FindById(int64(c.config.GetInt("server", "default_role_id")))
		if err != nil {
			_ = level.Error(c.logger).Log("roleRepository", "FindById", "err", err.Error())
			return "", nil, nil, ErrAuthLoginDefaultRoleID
		}
		var namespaces []types.Namespace
		var roles []types.Role
		namespaces = append(namespaces, *ns)
		roles = append(roles, *role)

		member = &types.Member{
			Username: username,
			Email:    email,
			Password: null.StringFrom(encode.EncodePassword(encode.GetRandomString(32),
				c.config.GetString("server", "app_key"))),
			Namespaces: namespaces,
			Roles:      roles,
		}
		if err = c.repository.Member().CreateMember(member); err != nil {
			_ = level.Error(c.logger).Log("member", "create", "err", err.Error())
			return "", nil, nil, err
		}
		for _, role := range roles {
			if _, err = c.casbin.GetEnforcer().AddGroupingPolicySafe(strconv.Itoa(int(member.ID)), strconv.Itoa(int(role.ID))); err != nil {
				_ = level.Warn(c.logger).Log("GetEnforcer", "AddGroupingPolicySafe", "err", err.Error())
			}
		}
	}

	if member.State == UserStateFail {
		_ = level.Error(c.logger).Log("email", "login", "email", email, "state", member.State, "err", "user state is fail")
		return rs, member, nss, ErrUserStateFail
	}

	var groups []int64

	for _, ns := range member.Namespaces {
		nss = append(nss, ns.Name)
	}

	for _, group := range member.Groups {
		groups = append(groups, group.ID)
	}

	rs, err = c.sign(email, member.ID, nss, groups, member.Roles)
	rs = "Bearer " + rs
	return
}

func (c *service) ParseToken(ctx context.Context, token string) (map[string]interface{}, error) {
	token = strings.Split(token, "Bearer ")[1]

	var clustom kpljwt.ArithmeticCustomClaims
	tk, err := jwt.ParseWithClaims(token, &clustom, kpljwt.JwtKeyFunc)
	if err != nil {
		return nil, err
	}
	claim, ok := tk.Claims.(*kpljwt.ArithmeticCustomClaims)
	if !ok {
		_ = level.Error(c.logger).Log("tk", "Claims", "err", ok)
		return nil, middleware.ErrorASD
	}

	m, _ := c.repository.Member().FindById(claim.UserId)

	return map[string]interface{}{
		"namespaces": claim.Namespaces,
		"email":      m.Email,
		"username":   m.Username,
	}, nil
}

func (c *service) Login(ctx context.Context, email, password string) (rs string, err error) {
	var (
		ok   bool
		user map[string]string
		info *types.Member
	)
	if c.config.GetString("server", "login_type") == LoginTypeLDAP {
		// LDAP登陆
		ok, user, err = c.ldapLogin(email, password)
		if err != nil || !ok {
			_ = level.Error(c.logger).Log("ldap", "login", "err", err, "ok", ok)
			return rs, ErrUserOrPassword
		}
		info, err = c.repository.Member().Find(email)
		if err != nil {
			ns, err := c.repository.Namespace().Find(c.config.GetString("server", "default_namespace"))
			if err != nil {
				_ = level.Error(c.logger).Log("namespaceRepository", "Find", "err", err.Error())
				return "", ErrAuthLoginDefaultNamespace
			}
			role, err := c.repository.Role().FindById(int64(c.config.GetInt("server", "default_role_id")))
			if err != nil {
				_ = level.Error(c.logger).Log("roleRepository", "FindById", "err", err.Error())
				return "", ErrAuthLoginDefaultRoleID
			}
			var namespaces []types.Namespace
			var roles []types.Role
			namespaces = append(namespaces, *ns)
			roles = append(roles, *role)
			info = &types.Member{
				Username: user["name"],
				Email:    email,
				State:    1,
				Password: null.StringFrom(encode.EncodePassword(encode.GetRandomString(32),
					c.config.GetString("server", "app_key"))),
				Namespaces: namespaces,
				Roles:      roles,
			}
			if err = c.repository.Member().CreateMember(info); err != nil {
				_ = level.Error(c.logger).Log("member", "create", "err", err.Error())
				return "", err
			}
			for _, role := range roles {
				if _, err = c.casbin.GetEnforcer().AddGroupingPolicySafe(strconv.Itoa(int(info.ID)), strconv.Itoa(int(role.ID))); err != nil {
					_ = level.Warn(c.logger).Log("GetEnforcer", "AddGroupingPolicySafe", "err", err.Error())
				}
			}

		}
	} else {
		// 邮箱登陆
		info, err = c.emailLogin(email, password)
		if err != nil {
			_ = level.Error(c.logger).Log("email", "login", "err", err)
			return rs, ErrUserOrPassword
		}
	}

	if info == nil || info.State == UserStateFail {
		_ = level.Error(c.logger).Log("email", "login", "email", email, "err", err)
		return rs, ErrUserStateFail
	}

	var namespaces []string
	var groups []int64

	for _, ns := range info.Namespaces {
		namespaces = append(namespaces, ns.Name)
	}

	for _, group := range info.Groups {
		groups = append(groups, group.ID)
	}

	rs, err = c.sign(email, info.ID, namespaces, groups, info.Roles)
	rs = "Bearer " + rs
	return
}

/**
 * @Title ldap登陆
 */
func (c *service) ldapLogin(email, password string) (ok bool, user map[string]string, err error) {
	client := &ldap.LDAPClient{
		Base:         c.config.GetString("ldap", "ldap_base"),
		Host:         c.config.GetString("ldap", "ldap_host"),
		Port:         c.config.GetInt("ldap", "ldap_port"),
		UseSSL:       c.config.GetBool("ldap", "ldap_useSSL"),
		BindDN:       c.config.GetString("ldap", "ldap_bindDN"),
		BindPassword: c.config.GetString("ldap", "ldap_bind_password"),
		UserFilter:   c.config.GetString("ldap", "ldap_user_filter"),
		GroupFilter:  c.config.GetString("ldap", "ldap_group_filter"),
		Attributes:   c.config.GetStrings("ldap", "ldap_attr"),
	}

	defer client.Close()

	return client.Authenticate(email, password)
}

/**
 * @Title 邮箱登陆
 */
func (c *service) emailLogin(email, password string) (*types.Member, error) {
	passwordHashed := encode.EncodePassword(password, c.config.GetString("server", "app_key"))
	info, err := c.repository.Member().Login(email, passwordHashed)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (c *service) sign(email string, uid int64, namespaces []string, groups []int64, roles []types.Role) (string, error) {
	sessionTimeout, err := c.config.Int64("server", "session_timeout")
	if err != nil {
		sessionTimeout = 3600
	}
	expAt := time.Now().Add(time.Duration(sessionTimeout) * time.Second).Unix()

	_ = level.Debug(c.logger).Log("expAt", expAt)

	var isTrue bool
	var roleIds []int64
	for _, v := range roles {
		roleIds = append(roleIds, v.ID)
		if v.Level <= int(types.LevelOps) {
			isTrue = true
		}
	}

	// 创建声明
	claims := kpljwt.ArithmeticCustomClaims{
		UserId:     uid,
		Name:       email,
		Namespaces: namespaces,
		Groups:     groups,
		IsAdmin:    isTrue,
		RoleIds:    roleIds,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAt,
			Issuer:    "system",
		},
	}

	//创建token，指定加密算法为HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//生成token
	return token.SignedString([]byte(kpljwt.GetJwtKey()))
}

func NewService(logger log.Logger, cf *config.Config, kplcasbin kplcasbin.Casbin, store repository.Repository) Service {
	return &service{
		logger:     logger,
		config:     cf,
		casbin:     kplcasbin,
		repository: store,
	}
}

package oauth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lexkong/log"
	"gopkg.in/oauth2.v4"
	"gopkg.in/oauth2.v4/errors"
	"gopkg.in/oauth2.v4/generates"
	"gopkg.in/oauth2.v4/manage"
	"gopkg.in/oauth2.v4/models"
	"gopkg.in/oauth2.v4/server"
	"gopkg.in/oauth2.v4/store"
)

var (
	authCodeExp     = time.Hour * 3
	accessTokenExp  = time.Hour * 24
	refreshTokenExp = time.Hour * 48
	jwtKey          = "oauth"
)

type OauthServerModel struct {
	Server      *server.Server
	ClientStore *store.ClientStore
}

var OauthServer *OauthServerModel

func (*OauthServerModel) Init() {
	manager, clientStore := GetManager()
	srv := server.NewDefaultServer(manager)

	srv.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.Refreshing)
	srv.SetAllowGetAccessRequest(false)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Info("Internal Error:" + err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Info("Response Error:" + re.Error.Error())
	})

	// UserAuthorizationHandler get user id from request authorization
	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		// return "", errors.ErrAccessDenied
		// userID = r.Context().Value("userID").(string)
		return
	})

	// get client info (clientID and clientSecret)
	srv.SetClientInfoHandler(func(r *http.Request) (clientID, clientSecret string, err error) {
		clientID = r.FormValue("client_id")
		clientSecret = r.FormValue("client_secret")
		log.Info("client info: " + clientID + clientSecret)
		return
	})

	// srv.SetAccessTokenExpHandler(func(w http.ResponseWriter, r *http.Request) (time.Duration, error) {
	// })

	// ClientAuthorizedHandler check the client allows to use this authorization grant type
	srv.SetClientAuthorizedHandler(func(clientID string, grant oauth2.GrantType) (allowed bool, err error) {
		return true, nil
	})

	OauthServer = &OauthServerModel{
		Server:      srv,
		ClientStore: clientStore,
	}
}

func GetManager() (oauth2.Manager, *store.ClientStore) {
	manager := manage.NewDefaultManager()

	manager.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    accessTokenExp,
		RefreshTokenExp:   refreshTokenExp,
		IsGenerateRefresh: true,
	})

	manager.SetAuthorizeCodeExp(authCodeExp)

	// token store
	manager.MustTokenStorage(getTokenStore())
	// manager.MapTokenStorage(tokenStore)
	// token generate
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte(jwtKey), jwt.SigningMethodHS512))

	// client store
	clientStore := getClientStore()
	manager.MapClientStorage(clientStore)

	return manager, clientStore
}

func getTokenStore() (oauth2.TokenStore, error) {
	return store.NewFileTokenStore("store.db")
}

func getClientStore() *store.ClientStore {
	clientStore := store.NewClientStore()

	clientStore.Set("test", &models.Client{
		ID:     "test",
		Secret: "2",
		Domain: "http://localhost:9094",
	})

	return clientStore
	// return store.NewClientStore()
}

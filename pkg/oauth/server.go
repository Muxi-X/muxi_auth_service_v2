package oauth

import (
	"net/http"
	"time"

	store "github.com/Shadowmaple/oauth2-mysql-store"
	"github.com/dgrijalva/jwt-go"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"gopkg.in/oauth2.v4"
	"gopkg.in/oauth2.v4/errors"
	"gopkg.in/oauth2.v4/generates"
	"gopkg.in/oauth2.v4/manage"
	"gopkg.in/oauth2.v4/server"
)

const (
	authCodeExp     = time.Hour * 3
	accessTokenExp  = time.Hour * 24
	refreshTokenExp = time.Hour * 48

	tokenStoreTableName  = "oauth2_token"
	clientStoreTableName = "oauth2_client"
)

type OauthServerModel struct {
	Server      *server.Server
	ClientStore *store.ClientStore
}

var OauthServer *OauthServerModel

func (*OauthServerModel) Init() {
	clientStore := getClientStore()
	manager := getManager()
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)

	serverConfig(srv)

	OauthServer = &OauthServerModel{
		Server:      srv,
		ClientStore: clientStore,
	}
}

func serverConfig(srv *server.Server) {
	srv.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.Refreshing)
	srv.SetAllowGetAccessRequest(false)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Infof("Internal Error: %s\n", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Infof("Response Error: %s\n", re.Error.Error())
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

	// ClientAuthorizedHandler check the client allows to use this authorization grant type
	srv.SetClientAuthorizedHandler(func(clientID string, grant oauth2.GrantType) (allowed bool, err error) {
		return true, nil
	})
}

func getManager() *manage.Manager {
	manager := manage.NewDefaultManager()

	manager.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    accessTokenExp,
		RefreshTokenExp:   refreshTokenExp,
		IsGenerateRefresh: true,
	})

	manager.SetAuthorizeCodeExp(authCodeExp)

	// token store
	manager.MapTokenStorage(getTokenStore())
	// token generate
	jwtKey := viper.GetString("jwt_secret")
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte(jwtKey), jwt.SigningMethodHS512))

	// client store
	clientStore := getClientStore()
	manager.MapClientStorage(clientStore)

	return manager
}

func getTokenStore() oauth2.TokenStore {
	return store.NewTokenStore(&store.TokenConfig{
		BasicConfig: getDBBasicConfig(tokenStoreTableName),
		GcDisabled:  false,
		GcInterval:  time.Minute * 30,
	})
}

func getClientStore() *store.ClientStore {
	return store.NewClientStore(&store.ClientConfig{
		BasicConfig: getDBBasicConfig(clientStoreTableName),
	})

	// clientStore.Create(&models.Client{
	// 	ID:     "test",
	// 	Secret: "2",
	// 	Domain: "http://localhost:9094",
	// })

	// return clientStore
}

func getDBBasicConfig(table string) store.BasicConfig {
	return store.BasicConfig{
		Addr:     viper.GetString("db.addr"),
		UserName: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		Database: viper.GetString("db.name"),
		Table:    table,
	}
}

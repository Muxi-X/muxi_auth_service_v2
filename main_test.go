package main


import (
	"testing"
	"github.com/spf13/pflag"
	"github.com/Muxi-X/muxi_auth_service_v2/router"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/config"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/constvar"
	"github.com/gin-gonic/gin"
	"os"
)

func TestMain(m *testing.M) {
	pflag.Parse()

	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	model.DB.Init()
	defer model.DB.Close()

	constvar.TestRouter = gin.Default()
	router.Load(constvar.TestRouter)

	os.Exit(m.Run())
}
package main

import (
	"context"
	"fmt"
	"kenshop/goken/server/httpserver"
	"kenshop/goken/server/httpserver/middlewares/jwt"

	errors "kenshop/pkg/error"
	"kenshop/pkg/log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func test1() {

	// t, err := motel.NewJaegerTraceProvider(context.Background(), "192.168.199.128:4318")
	// if err != nil {
	// 	panic(err)
	// }
	// defer t.Shutdown(context.Background())
	// tracer := otel.Tracer("ken")
	// _, span := tracer.Start(context.Background(), "ken-span")
	// defer span.End()
	//otelLogger := log.MustNewOtelLogger(log.WithFormat(log.ConsoleFormat), log.WithOutputPaths("./a.log"))
	//log.Errorf("w1%s", "w")
	l := log.MustNewOtelLogger(log.WithErrOutPaths("stderr"))
	l.Sugar().Errorf("w1%s", "w")
}

func test2() {
	s := httpserver.NewServer(httpserver.WithTransLocale("zh"))
	s.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"w": "w",
		})
	})
	s.Start(context.Background())
}

func test3() {
	f := func(ctx *gin.Context) (interface{}, error) {
		pwd := ctx.DefaultQuery("pwd", "")
		usname := ctx.DefaultQuery("usname", "")
		if pwd != "" && usname != "" {
			return usname, nil
		}

		return nil, errors.New("用户或密码错误")
	}

	r, err := jwt.NewGinJWTMiddleware(jwt.WithAuthenticator(f), jwt.WithSecureKey("kensame"))
	if err != nil {
		panic(err)
	}
	s := httpserver.NewServer()
	r.Timeout = 3 * time.Second

	s.GET("/", r.RefreshHandler, r.MiddlewareFunc(), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"w": "w",
		})
	})

	r.Authenticator = f

	s.GET("/login", r.LoginHandler)

	s.Start(context.Background())
}

func test_error1() error {
	err := test_error2()
	if err != nil {
		return errors.Wrap(err, "this is func external test_error1\n")
	}
	return nil
}

func test_error2() error {
	err := errors.New("this is func internal test_error2\n")
	return err
}

func test_caller() errors.StackTrace {
	return errors.Callers().StackTrace()
}

type VT struct {
	A int     `mapstructure:"a"`
	B string  `mapstructure:"b"`
	C float32 `mapstructure:"c"`
}

func testViper(prefix ...string) {
	v := viper.New()
	s, _ := os.Getwd()
	v.AddConfigPath(s)
	v.SetConfigFile("a.yaml")

	// os.Setenv("A", "1234")
	// os.Setenv("WW_A", "2234")
	v.AutomaticEnv()
	if len(prefix) > 0 {
		v.SetEnvPrefix(strings.Replace(strings.ToUpper(prefix[0]), "-", "_", -1))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	a := v.GetString("b")
	fmt.Println(a)
	v.SetConfigFile("b.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	a = v.GetString("b")
	fmt.Println(a)
	vt := &VT{}
	v.Unmarshal(&vt)
	fmt.Println(vt)

	v.WatchConfig()
	t := 0
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		v.Unmarshal(&vt)
		fmt.Println(vt, t)
		t++
	})
}

func main() {
	// err := test_error1()
	// fmt.Println("begin")
	// fmt.Printf("%+v\n", err)
	//s := test_caller()
	//fmt.Printf("%v", s)
	testViper()
	grpc.NewClient("w", grpc.WithResolvers())
	time.Sleep(1 * time.Minute)
}

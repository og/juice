package main

import (
	"context"
	"github.com/michaeljs1990/sqlitestore"
	ogjson "github.com/og/json"
	jhttpttp "github.com/og/juice/http"
	vd "github.com/og/juice/validator"
	ge "github.com/og/x/error"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var sessionStore *sqlitestore.SqliteStore
func init() {
	var err error
	sessionStore, err = sqlitestore.NewSqliteStore(
		"./test/session_sqllite",
		"sessions",
		"/",
		3600*24,
		[]byte("production environment must write secretKey"),
		)
	if err != nil {
		panic(err)
	}
}
type ReqHome struct {
	Name string `query:"name"` // query json form
	Age uint `query:"age"`
}

func (v ReqHome) VD(r *vd.Rule) {
	r.String(v.Name, vd.StringSpec{
		Name:              "姓名",
		MinRuneLen: 2,
		MaxRuneLen:10,
	})
	r.Uint(v.Age, vd.IntSpec{
		Name:           "年龄",
		Min: vd.Int(18),
		Max: vd.Int(80),
	})
}
type ReqUserDetail struct {
	ID string `param:"id"`
}
func main() {
	router := jhttpttp.NewRouter(jhttpttp.RouterOption{
		OnCatchError: func(c *jhttpttp.Context, errInterface interface{}) error {
			log.Print(errInterface)
			switch errInterface.(type) {
			case error:
				err := errInterface.(error)
				reject, isReject := ge.ErrorToReject(err)
				if isReject {
					return c.Bytes(ogjson.Bytes(reject.Response))
				} else {
					return c.Bytes([]byte("server error!"))
				}
			default:
				return c.Bytes([]byte("server error."))
			}
		},
	})
	requestLogMiddleware := func(c *jhttpttp.Context, next jhttpttp.Next) error {
		log.Print("Request:", c.R.Method, " ", c.R.URL)
		err := next() ; if err != nil {panic(err)}
		log.Print("Response:", c.R.Method, " ", c.R.URL)
		return nil
	}
	router.Use(requestLogMiddleware)
	router.HandleFunc(jhttpttp.GET, "/", func(c *jhttpttp.Context) (reject error) {
		time.Sleep(6*time.Second)
		/* 绑定请求 */{
			req := ReqHome{}
			reject = c.BindRequest(&req) ;if reject != nil {return}
			report := vd.NewCN().Check(req)
			if report.Fail {
				return ge.NewReject(report.Message, false)
			}
		}
		/* 读写 session */{
			sess := c.Session("juice_session", sessionStore)
			// err := sess.SetString("time", time.Now().String()) ; if err !=nil { return err}
			var timeStr string
			var hasTime bool
			timeStr, hasTime, reject = sess.GetString("time") ; if reject != nil {return}
			if !hasTime {timeStr = ""}
			return c.Bytes([]byte("time:" + timeStr))
		}
	})
	router.HandleFunc(jhttpttp.GET,"/user/{id}", func(c *jhttpttp.Context) (reject error) {
		req := ReqUserDetail{}
		reject = c.BindRequest(&req) ; if reject != nil {return}
		id, reject := c.Param("id") ; if reject != nil {return}
		return c.Bytes([]byte("BindRequestID:" + req.ID + " ParamID:" + id))
	})
	{
		g := router.Group()
		g.Use(func(c *jhttpttp.Context, next jhttpttp.Next) error {
			token := c.R.Header.Get("token")
			if token != "abc" {
				return c.Bytes([]byte("token 错误"))
			}
			return next()
		})
		g.HandleFunc(jhttpttp.GET, "/user", func(c *jhttpttp.Context) error {
			return c.Bytes([]byte("some data"))
		})
	}
	serve := http.Server{
		Addr: ":1219",
		Handler: router,
	}
	go func(){
		log.Print("Listen http://127.0.0.1" + serve.Addr)
		err := serve.ListenAndServe()
		if err != nil {
			log.Print(err)
		}
	}()
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	log.Print("Shuting down server...")
	if err := serve.Shutdown(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
	go func() {
		<-exit
	}()

}
package main

import (
	"database/sql"
	"log"
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/casbin/casbin"
	sqladapter "github.com/Blank-Xu/sql-adapter"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"time"
)

var cas *casbin.Enforcer

func finalizer(db *sql.DB) {
	err := db.Close()
	if err != nil {
		panic(err)
	}
}

func mysqlinit() *sqladapter.Adapter {
	// connect to the database first.
	db, err := sql.Open("mysql", "root:842b1255668c@tcp(192.168.3.30:3306)/abc")
	if err != nil {
		panic(err)
	}

	if err = db.Ping();err!=nil{
		panic(err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 10)

	// need to control by user, not the package
	runtime.SetFinalizer(db, finalizer)

	// Initialize an adapter and use it in a Casbin enforcer:
	// The adapter will use the Sqlite3 table name "casbin_rule_test",
	// the default table name is "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	a, err := sqladapter.NewAdapter(db, "mysql", "casbin_rule")
	if err != nil {
		panic(err)
	}

	log.Println("sql adapter new success!")

	return a
}

func main() {
	var err error

	cas, err = casbin.NewEnforcer("rbac_model.conf", mysqlinit())
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Load the policy from DB.
	if err = cas.LoadPolicy(); err != nil {
		log.Println("LoadPolicy failed, err: ", err)
		return
	}

	// Check the permission.
	has, err := cas.Enforce("alice", "data1", "read")
	if err != nil {
		log.Println("Enforce failed, err: ", err)
	}

	if !has {
		log.Println("do not have permission")
	}

	//ok, err := cas.AddPolicy("alice", "data1", "read")
	//if err != nil {
	//	log.Println("AddPolicy failed, err: ", err)
	//}

	//ok, err = cas.AddPolicy("user","/","GET")
	//if err != nil {
	//	log.Println("AddPolicy failed, err: ", err)
	//}
	//
	//if !ok {
	//	log.Println("add policy not ok")
	//}

	ok, err := cas.RemovePolicy("user","/","GET")
	if err != nil {
		log.Println("RemovePolicy failed, err: ", err)
	}

	if !ok {
		log.Println("do not have permission")
	}

	// Save the policy back to DB.
	if err = cas.SavePolicy(); err != nil {
		log.Println("SavePolicy failed, err: ", err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//e.Use(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	//	return func(context echo.Context) error {
	//
	//	}
	//})
	
	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	ok, err := cas.Enforce("user", c.Path(), c.Request().Method)
	if err != nil {
		c.Error(err)
		return err
	}
	if ok == true {
		return c.String(http.StatusOK, "Hello, World!")
	}
	return c.String(http.StatusNonAuthoritativeInfo, "deny")
}
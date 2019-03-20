// +build go1.12

package main

import (
	"database/sql"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/bcicen/grmon/agent"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/panjf2000/ants"
	"github.com/pborman/uuid"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// gin不能设置goroutine池，大并发时会出问题。
	d := gin.Default()
	fmt.Println(uuid.New())
	grmon.Start()

	pool, err := ants.NewPool(10000)
	if err != nil {
		panic(err)
	}
	defer pool.Release()

	db, err := sql.Open("mysql", "root:1314@/runoob?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(150)
	db.SetMaxIdleConns(10)

	gdb, err := gorm.Open("mysql", db)
	if err != nil {
		panic(err)
	}

	defer gdb.Close()

	d.POST("/addr", func(ctx *gin.Context) {
		var person Users
		if ctx.ShouldBindJSON(&person) == nil {
			fmt.Println(person.Id)
			fmt.Println(person.Name)

			person.Id = uuid.New()
			fmt.Println(gdb.NewRecord(person))
			gdb.Create(&person)
			ctx.String(http.StatusOK, "Success")
		} else {
			ctx.String(http.StatusBadRequest, "fail")
		}
	})

	err = d.Run(":8080")
	if err != nil {
		panic(err)
	}
}

type Users struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

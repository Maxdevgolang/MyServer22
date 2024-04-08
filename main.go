package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

var db *sql.DB

func main() {
	con, err := sql.Open("sqlite3", "myBase.db")
	if err != nil {
		panic("Не загружена база данных, дальнейшая работа невозможна")
	} else {
		db = con
	}
	r := gin.Default()
	r.LoadHTMLGlob("tmpl/*.html")
	r.GET("/f", MainPage)
	r.POST("/fgos", FindGos)
	r.POST("/fgar", FindGar)
	r.Run(":80")
}
func MainPage(ctx *gin.Context) {
	ctx.HTML(200, "main.html", nil)
}

func FindGar(ctx *gin.Context) { //Поиск гаражного по гос
	inp := ctx.PostForm("inp")

	inp = strings.ToUpper(inp)
	inp = strings.ReplaceAll(inp, " ", "")
	log.Println("Запрос поиска гаражного номера", inp)
	gar, gv := getGarDB(inp)
	ctx.HTML(200, "main.html", gin.H{"Gos": inp, "Gar": gar, "Gv": gv})

}

func FindGos(ctx *gin.Context) {
	inp := ctx.PostForm("inp")
	log.Println("Запрос поиска гос номера", inp)
	inp = strings.ToUpper(inp)
	inp = strings.ReplaceAll(inp, " ", "")
	gos, gv := getGosDB(inp)
	ctx.HTML(200, "main.html", gin.H{"Gos": gos, "Gar": inp, "Gv": gv})
}

func getGarDB(gos string) (string, string) {
	var tmpGar, tmpGv string
	raw, err := db.Query("SELECT Gar,Gv FROM TSBase WHERE Gos=?", gos)
	if err != nil {
		log.Println("Не верный запрос", err)
	}
	for raw.Next() {
		if raw.Scan() == nil {
			return tmpGar, tmpGv
		} else {
			raw.Scan(&tmpGar, &tmpGv)
		}
	}
	if tmpGar != "" {
		return tmpGar, tmpGv
	}
	return "Не найдено", "0"
}

func getGosDB(gar string) (string, string) {
	var tmpGos, tmpGv string
	raw, err := db.Query("SELECT Gos,Gv FROM TSBase WHERE Gar=?", gar)
	if err != nil {
		log.Println("Не верный запрос", err)
	}
	for raw.Next() {
		if raw.Scan() == nil {
			return tmpGos, tmpGv
		} else {
			raw.Scan(&tmpGos, &tmpGv)
		}
	}
	if tmpGos != "" {
		return tmpGos, tmpGv
	}
	return "Не найдено", "0"
}

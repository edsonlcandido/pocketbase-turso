package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func init() {
	dbx.BuilderFuncMap["libsql"] = dbx.BuilderFuncMap["sqlite3"]
}

func main() {
	//
	// "libsql://nome-do-db.turso.io?authToken=seu-token-aqui"
	tursoUrl := "libsql://nome-do-db.turso.io?authToken=seu-token-aqui"

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DBConnect: func(dbPath string) (*dbx.DB, error) {
			fmt.Printf("Verificando conexão para: %s\n", dbPath)

			if strings.HasSuffix(dbPath, "data.db") {
				fmt.Println("--- CONECTANDO AO TURSO (NUVEM) ---")
				return dbx.Open("libsql", tursoUrl)
			}

			fmt.Println("--- CONECTANDO AO BANCO LOCAL (LOGS) ---")
			return core.DefaultDBConnect(dbPath)
		},
	})

	// Servir arquivos estáticos de pb_public (com fallback para index.html)
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		publicDir := filepath.Join(filepath.Dir(app.DataDir()), "pb_public")
		se.Router.GET("/{path...}", apis.Static(os.DirFS(publicDir), true))
		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

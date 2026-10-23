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
	_ "golang.org/x/crypto/x509roots/fallback"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func init() {
	dbx.BuilderFuncMap["libsql"] = dbx.BuilderFuncMap["sqlite3"]
}

func loadTursoURLs(getenv func(string) string) (string, string, error) {
	tursoURL := getenv("URL_LIBSQL_TURSO")
	if tursoURL == "" {
		return "", "", fmt.Errorf("URL_LIBSQL_TURSO não definida")
	}

	auxTursoURL := getenv("URL_LIBSQL_TURSO_AUX")
	if auxTursoURL == "" {
		return "", "", fmt.Errorf("URL_LIBSQL_TURSO_AUX não definida")
	}

	return tursoURL, auxTursoURL, nil
}

func main() {
	tursoURL, auxTursoURL, err := loadTursoURLs(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DBConnect: func(dbPath string) (*dbx.DB, error) {
			fmt.Printf("Verificando conexão para: %s\n", dbPath)

			switch {
			case strings.HasSuffix(dbPath, "data.db"):
				fmt.Println("--- CONECTANDO data.db AO TURSO (NUVEM) ---")
				return dbx.Open("libsql", tursoURL)
			case strings.HasSuffix(dbPath, "auxiliary.db"):
				fmt.Println("--- CONECTANDO auxiliary.db AO TURSO (NUVEM) ---")
				return dbx.Open("libsql", auxTursoURL)
			default:
				fmt.Println("--- CONECTANDO AO BANCO LOCAL (LOGS) ---")
				return core.DefaultDBConnect(dbPath)
			}
		},
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		publicDir := filepath.Join(filepath.Dir(app.DataDir()), "pb_public")
		se.Router.GET("/{path...}", apis.Static(os.DirFS(publicDir), true))
		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

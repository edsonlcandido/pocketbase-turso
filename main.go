package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"ganigeorgiev.com/pocketbase/dbx"
	"ganigeorgiev.com/pocketbase/pocketbase"
	"ganigeorgiev.com/pocketbase/pocketbase/apis"
	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"ganigeorgiev.com/pocketbase/pocketbase/plugins/jsvm"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "golang.org/x/crypto/x509roots/fallback"
)

func init() {
	dbx.BuilderFuncMap["libsql"] = dbx.BuilderFuncMap["sqlite3"]
}

// loadTursoURLs lê a URL do banco principal e, opcionalmente, a do auxiliar.
//
// A URL auxiliar é OPCIONAL de propósito: quando ela não vem, o auxiliary.db
// (onde fica a tabela _logs) fica no disco local. Nunca se deve apontar os
// dois para o mesmo banco — data.db e auxiliary.db têm as mesmas tabelas de
// sistema (_collections, _records) e uma colide com a outra.
func loadTursoURLs(getenv func(string) string) (string, string, error) {
	tursoURL := getenv("URL_LIBSQL_TURSO")
	if tursoURL == "" {
		return "", "", fmt.Errorf("URL_LIBSQL_TURSO não definida")
	}

	auxTursoURL := getenv("URL_LIBSQL_TURSO_AUX")

	return tursoURL, auxTursoURL, nil
}

func main() {
	tursoURL, auxTursoURL, err := loadTursoURLs(os.Getenv)
	if err != nil {
		log.Fatal(err)
	}

	auxIsRemote := auxTursoURL != ""
	if !auxIsRemote {
		fmt.Println("URL_LIBSQL_TURSO_AUX não definida: auxiliary.db (logs) fica no disco local.")
	}

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DBConnect: func(dbPath string) (*dbx.DB, error) {
			fmt.Printf("Verificando conexão para: %s\n", dbPath)

			switch {
			case strings.HasSuffix(dbPath, "data.db"):
				fmt.Println("--- CONECTANDO data.db AO TURSO (NUVEM) ---")
				return dbx.Open("libsql", tursoURL)
			case strings.HasSuffix(dbPath, "auxiliary.db") && auxIsRemote:
				fmt.Println("--- CONECTANDO auxiliary.db AO TURSO (NUVEM) ---")
				return dbx.Open("libsql", auxTursoURL)
			default:
				fmt.Println("--- CONECTANDO AO BANCO LOCAL (LOGS) ---")
				return core.DefaultDBConnect(dbPath)
			}
		},
	})

	// ---- JSVM: carrega pb_hooks/*.js e pb_migrations/*.js ----
	//
	// Desde o PB v0.39 o JSVM virou plugin: num build Go customizado ele só
	// é carregado se for registrado explicitamente (é o que faz o binário
	// oficial, compilado a partir de examples/base). Sem esta chamada, os
	// arquivos .js são IGNORADOS SILENCIOSAMENTE — o app sobe, mas o schema
	// nunca é criado.
	jsvm.MustRegister(app, jsvm.Config{
		MigrationsDir: filepath.Join(filepath.Dir(app.DataDir()), "pb_migrations"),
		HooksDir:      filepath.Join(filepath.Dir(app.DataDir()), "pb_hooks"),
		HooksWatch:    false, // em container não faz sentido vigiar arquivos
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

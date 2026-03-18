package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func init() {
	// Faz o driver do Turso ser reconhecido como SQLite pelo PocketBase
	dbx.BuilderFuncMap["libsql"] = dbx.BuilderFuncMap["sqlite3"]
}

func main() {
	// SUBSTITUA PELA SUA URL DO TURSO
	// Exemplo: "libsql://nome-do-db.turso.io?authToken=seu-token-aqui"
	tursoUrl := "libsql://nome-do-db.turso.io?authToken=seu-token-aqui"

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DBConnect: func(dbPath string) (*dbx.DB, error) {
			// Log para debug: vamos ver o que o PocketBase está tentando abrir
			fmt.Printf("Verificando conexão para: %s\n", dbPath)

			// Verifica se o arquivo que ele quer abrir é o banco de dados principal
			// Usamos HasSuffix para funcionar em Windows, Linux ou caminhos relativos
			if strings.HasSuffix(dbPath, "data.db") {
				fmt.Println("--- CONECTANDO AO TURSO (NUVEM) ---")
				return dbx.Open("libsql", tursoUrl)
			}

			// Para logs e backups, ele usa o SQLite local (pb_data/auxiliary.db)
			fmt.Println("--- CONECTANDO AO BANCO LOCAL (LOGS) ---")
			return core.DefaultDBConnect(dbPath)
		},
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
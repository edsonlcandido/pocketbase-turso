# PocketBase + Turso (libSQL)

Este projeto é uma implementação customizada do [PocketBase](https://pocketbase.io/) **v0.40.6** que utiliza o [Turso (libSQL)](https://turso.tech/) como banco de dados principal, em vez do SQLite local padrão.

Isso permite que você tenha uma instância do PocketBase rodando localmente ou em containers (Edge/Serverless) enquanto seus dados permanecem sincronizados em um banco de dados distribuído na nuvem.

## 📋 Pré-requisitos

*   [Go 1.27](https://go.dev/dl/) ou superior instalado.
*   Uma conta no [Turso](https://turso.tech/) e **um banco de dados** criado (o auxiliar é opcional — veja abaixo).
*   URL e Token do seu banco de dados Turso.

## 🚀 Como começar

1.  **Clone o repositório:**
    ```bash
    git clone <url-do-seu-repositorio>
    cd <nome-da-pasta>
    ```

2.  **Crie um banco no Turso:**
    - um banco principal para o `data.db`
    - (opcional) um banco auxiliar para o `auxiliary.db`

3.  **Configure suas credenciais:**
    Defina as variáveis de ambiente com as URLs de conexão do Turso (incluindo o token):
    ```bash
    export URL_LIBSQL_TURSO="libsql://seu-db.turso.io?authToken=seu-token-aqui"
    # opcional — se não definir, o auxiliary.db (logs) fica no disco local:
    export URL_LIBSQL_TURSO_AUX="libsql://seu-db-aux.turso.io?authToken=seu-token-aux-aqui"
    ```

    | Variável | Obrigatória | Descrição |
    |---|---|---|
    | `URL_LIBSQL_TURSO` | ✅ Sim | URL do banco principal (`data.db`) |
    | `URL_LIBSQL_TURSO_AUX` | ⚠️ Opcional | URL do banco auxiliar (`auxiliary.db`, onde fica a tabela `_logs`). **Se for definida, tem que apontar para um segundo database separado** — `data.db` e `auxiliary.db` têm as mesmas tabelas de sistema (`_collections`, `_records`) e uma colide com a outra. Sem ela, o auxiliar fica no disco local. |

    > ⚠️ O `.env` **não** é lido pelo binário: a URL é lida do ambiente em runtime (`os.Getenv`). Use `docker run --env-file .env`, `env_file:` no compose ou `EnvironmentFile=` no systemd. O `.env` está no `.gitignore` — a URL carrega o token dentro da query string (`?authToken=`).

4.  **Instale as dependências:**
    ```bash
    go mod tidy
    ```

5.  **Execute em modo de desenvolvimento:**
    ```bash
    go run main.go serve
    ```

## 🐳 Imagem Docker (ghcr)

A cada push na `main` o CI publica a imagem em `ghcr.io/edsonlcandido/pocketbase-turso`:

```bash
docker run -d --name pocketbase \
  -p 8090:8090 \
  -v pb_data:/app/pb_data \
  -e URL_LIBSQL_TURSO="libsql://seu-db.turso.io?authToken=seu-token-aqui" \
  ghcr.io/edsonlcandido/pocketbase-turso:main
```

Tags disponíveis: `main` (último push), `v0.40.6` e `0.40` (a partir de tags `v*.*.*`).

> A URL pode ser de **qualquer** servidor que fale o protocolo hrana v2 do libSQL — Turso, um `sqld` auto-hospedado ou um proxy (formato `https://host/dev/<namespace>/?authToken=...`). O binário só precisa de `https://` (ou `libsql://`) + o token na query string.

## 🛠️ Como Compilar (Build)

O Go permite gerar executáveis para diferentes sistemas operacionais, independente de qual sistema você está usando para desenvolver.

### Para Windows (64 bits)
Se você estiver no Linux/Mac e quiser gerar o `.exe`:
```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o pocketbase.exe
```

### Para Linux (64 bits)
Útil para deploy em VPS, Docker ou serviços como Fly.io:
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pocketbase
```

### Para macOS
*   **Intel:** `GOOS=darwin GOARCH=amd64 go build -o pocketbase_mac`
*   **Apple Silicon (M1/M2/M3):** `GOOS=darwin GOARCH=arm64 go build -o pocketbase_mac`

## 📂 Estrutura de Pastas

*   `main.go`: Contém a lógica de inicialização e o desvio da conexão para o Turso.
*   `pb_hooks/` e `pb_migrations/`: Hooks e migrations em JS (`*.js`). São carregados pelo plugin JSVM registrado em `main.go` — sem esse registro os arquivos são ignorados silenciosamente.
*   `/pb_data`: Pasta criada automaticamente pelo PocketBase.
    *   `data.db`: Ficará vazio/mínimo, pois os dados estão no Turso.
    *   `auxiliary.db`: Vai para o Turso via `URL_LIBSQL_TURSO_AUX`. **Se apontar para o Turso, tem que ser um database separado** do `data.db`, para evitar colisão de tabelas. Sem a variável, fica no disco local.
    *   `/storage`: Guarda os arquivos de upload (fotos, docs). **Nota:** Arquivos de upload não vão para o Turso, ficam nesta pasta.

## ⚠️ Observações Importantes

1.  **Latência:** Como o banco de dados está na nuvem (Turso) e não no disco local, você pode notar um pequeno delay nas operações do painel administrativo.
2.  **Isolamento dos bancos:** se `URL_LIBSQL_TURSO_AUX` estiver definida, `data.db` e `auxiliary.db` precisam apontar para databases distintos.
3.  **JSVM (PB v0.39+):** desde a v0.39 o JSVM virou *plugin* (`plugins/jsvm`). Num build Go customizado ele só é carregado se for registrado explicitamente — é o que o `main.go` deste repo faz. **Sem esse registro, `pb_hooks/*.js` e `pb_migrations/*.js` são ignorados silenciosamente** (o app sobe, mas o schema nunca é criado).
4.  **Segurança:** a URL de conexão carrega o token do Turso dentro da query string (`?authToken=`). Nunca versionar o `.env` — ele já está no `.gitignore` e no `.dockerignore`. No Docker, defina a URL como **variável de ambiente de runtime**; se for build arg, o token fica gravado na layer da imagem.
5.  **CGO:** Utilizamos `CGO_ENABLED=0` durante o build para garantir que o executável seja estático e funcione em qualquer servidor sem depender de bibliotecas do sistema.
6.  **Uploads:** os arquivos enviados ficam em `/app/pb_data/storage` (disco). Monte um volume nesse caminho, senão eles são perdidos a cada novo deploy.

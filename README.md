# PocketBase + Turso (libSQL)

Este projeto é uma implementação customizada do [PocketBase](https://pocketbase.io/) que utiliza o [Turso (libSQL)](https://turso.tech/) como banco de dados principal, em vez do SQLite local padrão.

Isso permite que você tenha uma instância do PocketBase rodando localmente ou em containers (Edge/Serverless) enquanto seus dados permanecem sincronizados em um banco de dados distribuído na nuvem.

## 📋 Pré-requisitos

*   [Go 1.21](https://go.dev/dl/) ou superior instalado.
*   Uma conta no [Turso](https://turso.tech/) e um banco de dados criado.
*   URL e Token do seu banco de dados Turso.

## 🚀 Como começar

1.  **Clone o repositório:**
    ```bash
    git clone <url-do-seu-repositorio>
    cd <nome-da-pasta>
    ```

2.  **Configure suas credenciais:**
    Abra o arquivo `main.go` e substitua a variável `tursoUrl` pela sua URL de conexão completa (incluindo o token):
    ```go
    tursoUrl := "libsql://seu-db.turso.io?authToken=seu-token-aqui"
    ```

3.  **Instale as dependências:**
    ```bash
    go mod tidy
    ```

4.  **Execute em modo de desenvolvimento:**
    ```bash
    go run main.go serve
    ```

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
*   `/pb_data`: Pasta criada automaticamente pelo PocketBase.
    *   `data.db`: Ficará vazio/mínimo, pois os dados estão no Turso.
    *   `auxiliary.db`: Continua sendo usado localmente para logs de requisições.
    *   `/storage`: Guarda os arquivos de upload (fotos, docs). **Nota:** Arquivos de upload não vão para o Turso, ficam nesta pasta.

## ⚠️ Observações Importantes

1.  **Latência:** Como o banco de dados está na nuvem (Turso) e não no disco local, você pode notar um pequeno delay nas operações do painel administrativo.
2.  **Segurança:** Evite subir seu `main.go` para repositórios públicos com o Token do Turso exposto. Recomenda-se usar variáveis de ambiente em produção.
3.  **CGO:** Utilizamos `CGO_ENABLED=0` durante o build para garantir que o executável seja estático e funcione em qualquer servidor sem depender de bibliotecas do sistema.
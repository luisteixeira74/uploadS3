# S3 Parallel Uploader (Go)

Este projeto é uma ferramenta simples e eficiente escrita em Go para fazer upload de múltiplos arquivos para um bucket S3 da AWS, com suporte a **concorrência**, **controle de erro** e **retry automático** em caso de falha.

## 🚀 Funcionalidades

- Upload concorrente com limite configurável (100 por padrão).
- Retry automático para arquivos que falharem no upload.
- Leitura contínua de arquivos de um diretório local (`./tmp`).
- Uso de goroutines, canais e `sync.WaitGroup` para controle de execução.

## 🧠 Como funciona

1. O programa abre o diretório `./tmp` e lê arquivos continuamente.
2. A cada arquivo encontrado:
   - Ele adiciona uma goroutine para o upload.
   - Usa um canal `uploadControl` para limitar o número de uploads simultâneos.
3. Se um upload falhar:
   - O nome do arquivo é enviado para um canal `errorFileUpload`.
   - Uma goroutine especial escuta este canal e tenta fazer o reenvio automaticamente.
4. O programa aguarda todos os uploads (incluindo os de retry) finalizarem antes de encerrar.

## 📁 Estrutura esperada

```bash
.
├── main.go
├── .env
└── tmp/
    ├── file1.jpg
    ├── file2.jpg
    └── ...
```

## 🔐 Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com suas credenciais da AWS:

```env
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
S3_BUCKET_NAME=your-bucket-name
```

## ⚙️ Execução

### Pré-requisitos

- Go 1.20+ instalado
- Credenciais AWS válidas com permissão para `s3:PutObject`
- Um bucket S3 já criado

### Rodar o projeto

```bash
go run main.go
```

## 🔄 Retry Automático

Arquivos que falham no upload são reenviados automaticamente pelo sistema de retry baseado em canal. Não há limite de tentativas atualmente (pode ser adicionado facilmente).

## 📌 Notas Técnicas

- `uploadControl` é um canal com buffer que limita a concorrência a 100 uploads simultâneos.
- `errorFileUpload` é um canal com buffer que captura os arquivos com falha para retry.
- `sync.WaitGroup` garante que todos os uploads (incluindo retries) sejam finalizados antes de sair.

## 🧪 Testes

Este projeto não possui testes automatizados neste momento, mas há pontos claros onde `uploadFile` poderia ser testado com mocks da AWS.

## 📚 Aprendizados

Este projeto é ideal para praticar:

- Concorrência com goroutines e canais em Go
- Integração com AWS SDK
- Tratamento de erros e retry resiliente

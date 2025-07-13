Como rodar o projeto em ambiente dev.

Antes de começar, certifique-se de ter instalado:
	Docker
	Docker Compose
	Go

1. Clone o repositório

2. Configure o arquivo .env e certifique-se que está no caminho cmd/auction/.env.com o conteúdo:
	AUCTION_DURATION_MINUTES=10
	AUCTION_INTERVAL=20s
	BATCH_INSERT_INTERVAL=20s
	MAX_BATCH_SIZE=4
	
	MONGO_INITDB_ROOT_USERNAME=admin
	MONGO_INITDB_ROOT_PASSWORD=admin
	MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
	MONGODB_DB=auctions

3. Suba os containers com Docker Compose:
	docker-compose up --build

4. Acesse a aplicação em:
	http://localhost:8080

5. Para rodar os testes automatizados:
	go test ./internal/infra/database/auction -v
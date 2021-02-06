![Imersão Full Stack && Full Cycle](https://events-fullcycle.s3.amazonaws.com/events-fullcycle/static/site/img/grupo_4417.png)

Participe gratuitamente: https://imersao.fullcycle.com.br/

## Microsserviço CodePix

Esse microsserviço tem o objetivo de ser um hub de transações entre os bancos que simularemos durante o projeto.

## Como executar

Utilizamos Docker para que todos os serviços que utilizaremos fiquem disponíveis.

- Faça o clone do projeto
- Tendo o docker instalado em sua máquina apenas execute:
  `docker-compose up -d`

### Como executar a aplicação

- Acesse o container da aplicação executando: `docker exec -it codepix_app bash`
- Rode `go run cmd/codepix/main.go`

### Serviços utilizados ao executar o docker-compose

- Aplicação principal
- Postgres
- PgAdmin
- Apache Kafka
- Criador dos tópicos a serem utilizados pelo Kafka
- Confluent control center
- ZooKeeper
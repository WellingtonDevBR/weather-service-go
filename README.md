# Weather Service

Este projeto implementa um sistema em Go que recebe um CEP, identifica a cidade e retorna o clima atual (temperatura em graus Celsius, Fahrenheit e Kelvin). O sistema é implantado no Google Cloud Run.

## Funcionalidades

- Recebe um CEP válido de 8 dígitos.
- Pesquisa a localização com a API viaCEP.
- Obtém a temperatura atual com a WeatherAPI.
- Converte e retorna as temperaturas em Celsius, Fahrenheit e Kelvin.
- Responde adequadamente nos seguintes cenários:
  - Sucesso: HTTP 200
  - CEP inválido: HTTP 422
  - CEP não encontrado: HTTP 404

## Requisitos

- Go 1.16+
- Docker
- Conta no Google Cloud

## Uso

### Executar Localmente

1. Clone o repositório:

```sh
git clone https://github.com/WellingtonDevBR/weather-service-go
cd weather-service-go
```

2. Instale as dependências:

```sh
go mod tidy
```

3. Compile a aplicação:

```sh
go build -o weather-service-go main.go
```

4. Execute a aplicação:

```sh
./weather-service
```

5. Faça uma requisição para a aplicação:

```sh
curl http://localhost:8080/weather/05144085
```

### Executar com Docker

```sh
docker build -t weather-service .
docker run -p 8080:8080 weather-service
curl http://localhost:8080/weather/05144085
```

### Executar Docker no Google Cloud

```sh
curl https://weather-service-q2pgawiciq-uc.a.run.app/weather/05144085
```
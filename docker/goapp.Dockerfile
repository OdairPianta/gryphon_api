# Imagem base
FROM golang

# Adiciona informações da pessoa que mantem a imagem
LABEL maintainer="Odair Pianta <odair@spotec.app>"

# diretoria um diretorio de trabalho
WORKDIR /app/src/gryphon_api

# aponta a variavel gopath do go para o diretorio app
ENV GOPATH=/app

# define a variavel de ambiente TZ para America/Sao_Paulo
ENV TZ="America/Sao_Paulo"

# copia os arquivos do projeto para o workdir do container
COPY . /app/src/gryphon_api

# execulta o main.go e baixa as dependencias do projeto
RUN go build main.go

# Comando para rodar o executavel
ENTRYPOINT ["./main"]

# expose port 8083
EXPOSE 8083
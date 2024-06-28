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

# install google-chrome
RUN apt update && apt upgrade
RUN apt install software-properties-common apt-transport-https ca-certificates curl -y
RUN curl -fSsL https://dl.google.com/linux/linux_signing_key.pub | gpg --dearmor | tee /usr/share/keyrings/google-chrome.gpg >> /dev/null
RUN echo deb [arch=amd64 signed-by=/usr/share/keyrings/google-chrome.gpg] http://dl.google.com/linux/chrome/deb/ stable main | tee /etc/apt/sources.list.d/google-chrome.list
RUN apt update
RUN apt install google-chrome-stable -y

# copia os arquivos do projeto para o workdir do container
COPY . /app/src/gryphon_api

# execulta o main.go e baixa as dependencias do projeto
RUN go build main.go

# Comando para rodar o executavel
ENTRYPOINT ["./main"]

# expose port 8083
EXPOSE 8083
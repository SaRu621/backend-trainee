FROM golang:latest
 
WORKDIR /app

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go build -o main .
 
CMD ["./main"] 

FROM golang:1.8

RUN go get golang.org/x/net/context
RUN git clone https://github.com/Va5kez/Tarea1.git
RUN go get googlemaps.github.io/maps

EXPOSE 8080

CMD cd /Tarea1 && go run main.go

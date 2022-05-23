FROM golang:latest 
WORKDIR /app 
ADD ./ /app/ 
RUN go build ./main.go
EXPOSE 9600/udp
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["sh"]

## TradeSimulator

### Installation
```
clone epository
cd [workspace]/TradeSimulator
go get -u
go build main.go
go run main.go
```
### Usage
- WebPage: `http://localhost:8000/home.html`
 
### Use Docker to Install

```
cd [workspace]/TradeSimulator
docker build -t tradesimulator  .
docker run -d -p 8000:8000 --name tradesimulator tradesimulator
```



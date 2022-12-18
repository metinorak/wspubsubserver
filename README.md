# Websocket PubSub server written in Go
* This project uses [gorilla/websocket](https://github.com/gorilla/websocket) to implement a simple PubSub server. And for managing the connections, subscriptions and publishing messages, it uses [metinorak/varto](https://github.com/metinorak/varto) package and its gorilla ws wrapper package [metinorak/wspubsub](https://github.com/metinorak/wspubsub).

## Usage
* To run the server, you need to have Go installed on your machine. Then you can run the following command to start the server:
```bash
go run main.go
```

* You can specify the port number to run the server on by setting the `PORT` environment variable.
* To connect to the server, you can use any websocket client.

## Example
* You can use the following code to connect to the server and subscribe to a topic:
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onopen = () => {
  ws.send(JSON.stringify({
    action: 'SUBSCRIBE',
    topic: 'test',
  }));
};
ws.onmessage = (msg) => {
  console.log(msg);
};
```

* You can publish a message to a topic by sending a message to the server:
```javascript
ws.send(JSON.stringify({
  action: 'PUBLISH',
  topic: 'test',
  message: 'Hello World!',
}));
```

* You can unsubscribe from a topic by sending a message to the server:
```javascript
ws.send(JSON.stringify({
  action: 'UNSUBSCRIBE',
  topic: 'test',
}));
```

* You can broadcast a message to all the clients by sending a message to the server:
```javascript
ws.send(JSON.stringify({
  action: 'BROADCASTALL',
  message: 'Hello World!',
}));
```
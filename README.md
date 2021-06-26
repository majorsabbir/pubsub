# Pub/sub Implementation of gRPC

## Description
At the time of writting this document, I'm planing a write a program for realtime client-server-client communication. Codes will be written in golang, use redis to store data and use pubsub mechanism of it. Frontend client may be written in Javascript. Javascript client pass a message along with channel name to gRPC server. This message will published on provided redis channel. Another gRPC server will listening on that particular channel and streaming to a client. May be it will be a chat application.

I tried to learn golang last year(2020). But after few days due to work load I don't get time to continue my learning. Last Month (May) I continue my learning. And this week I learn gRPC protocol and implementation in golang. And this is my first application written in golang from scratch.

## API Operation with Evans
Connect to gRPC server with evans on default port. For Evans quick command check below. 

#### Publish an event
```
call PublishEvent
``` 
It will ask a ```channel``` name and followed by ```message```. If there is no subscriber on that particular ```channel``` message will not be published, get restroyed and return response of ```PublishEventResponse``` with ```publishEvent```. In case of any subscriber you will get ```subscriber_count``` also on response.

Response Example:
```
{
  "publishEvent": {
    "channel": "test",
    "msg": "test message one"
  }
}
```

#### Listen an event
```
call ListenEvent
```
You will be asked for a ```channel``` name and get stream response of ```ListenEventResponse``` with ```event```.

Response Example:
```
{
  "event": {
    "channel": "test",
    "msg": "test message one"
  }
}
{
  "event": {
    "channel": "test",
    "msg": "test message two"
  }
}
```

### Connecting gRPC server with Evans
+ ```evans -p 50051 -r -t --cacert ssl/ca.crt --servername localhost``` - when ```tls``` enabled
+ ```evans -p 50051 -r``` - when ```tls``` disabled
+ ```show service``` - list of all services with associated RPC, Request and Response type
+ ```call [rpcMethod]``` - call RPC
+ ```ctrl``` + ```d``` - exit
# Pub/sub Implementation of gRPC

## Description
At the time of writting this document, I'm planing a write a program for realtime client-server-client communication. Codes will be written in golang, use redis to store data and use pubsub mechanism of it. Frontend client may be written in Javascript. Javascript client pass a message along with channel name to gRPC server. This message will published on provided redis channel. Another gRPC server will listening on that particular channel and streaming to a client. May be it will be a chat application.

I tried to learn golang last year(2020). But after few days due to work load I don't get time to continue my learning. Last Month (May) I continue my learning. And this week I learn gRPC protocol and implementation in golang. And this is my first application written in golang from scratch.
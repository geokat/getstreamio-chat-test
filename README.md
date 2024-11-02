# A [Stream Chat](https://getstream.io/chat/) API Test

A demo Golang/JS app that creates a channel and an arbitrary number
of users, and sends a stream of messages from each of them.

To run, create a `dev.ini` file (see [dev.ini.example](dev.ini.example))
and run:
```shell
./run.sh
```

Open the dev console in the browser and visit
http://localhost:8000/client1.html

This will make the app create a channel and user `u1` and start
sending a stream of messages into the channel, one per second. Watch
the dev console for incoming messages in the channel.

Open more `client{n}.html` files in different browser windows/tabs
and watch the created users join the channel and send messages.


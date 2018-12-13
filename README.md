# /
## Launch

```bash
git clone https://github.com/reinerRubin/froppyshima.git
cd froppyshima/
sudo docker-compose build
sudo docker-compose up
```

Open http://127.0.0.1:3000/ & enjoy "froppyshima the game"!

## Cover letter
I have tried to make a complete "working product" with UI and a back-end part. It is quite different from the task. But if I understand properly you can do something based on it, because the main purpose of the task is to show my approach to a code. Currently an user can start a new game, play it and load an old one by a game-code.

### Decisions were made
#### Communication
I chose to use web-sockets, because a game logic was unclear at the beginning. Web-sockets allow to send events and can be very useful in such applications, but with current flow they are a little bit excessive.

#### Data model
On the first glance a matrix seems a good way to represent such games. But you must be able to answer to questions like "how many lives has this ship?", "where is my L shaped ship", "which ship was hit"? So you have to use some external catalog of ships with their mapping to a matrix. Furthermore, if a field is big enough a matrix should be sparse. But I have tried to avoid matrix representation at all. Because it seems more like representation than an actual model of data. But such approach requires a better implementation than I have. I have written down some details below (see "Things that are terrible").


### Implementation details
#### Front-end part
I have decided not to use any extra libraries (react, vue) in a such small task and stick with vanilla js. Web development is not my passion so I have tried to keep things simple and "bruteforced" a solution. There is a mix between DOM and data, terrible big functions and so on. So I have tried to add some "style" to compensate this.

#### Backend part
##### Client package
Currently the client (client.go) seems to use too many goroutines and maybe too complicated. It is a result of an attempt to separate client levels.
* client.go passes data between the connection client and the web client with minimum modifications
* clientconnection.go works with connections and raw data
* clientweb.go is a controller for "web" calls. It parses arguments and passes them to the game client
* clientgame.go is a client separated from all "web" context. It calls the game object and does some corresponding actions


##### Engine package
* finding intersections between objects must be optimized. Current implementation is too expensive
* some "ideas" like object skirt are not fully implemented and used ad hoc for the task
* putting ships at random positions must be also optimized (probably using adaptive strategies)

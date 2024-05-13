
# Groupie Tracker

## Prerequisites
You need to use :  
* Golang  
* Git for the code management  
  


## Installation
First you need to clone the repository  
```bash
git clone https://ytrack.learn.ynov.com/git/aerwan/HANGMAN-Web.git
```  
and after 
```bash
git init
```
And you need to install the Spotify, oauth2, unicode, transform package
```bash
go get "github.com/gorilla/sessions"
go get "github.com/gorilla/websocket"
go get "github.com/mattn/go-sqlite3"
go get "github.com/zmb3/spotify"
go get "golang.org/x/crypto"
go get "golang.org/x/oauth2"
go get "golang.org/x/text"
go get "github.com/gorilla/securecookie"
go get "golang.org/x/net"
go get "golang.org/x/sys"
go get "golang.org/x/text"
go get "google.golang.org/appengine" 
go get "google.golang.org/protobuf"

```
To read the database you need : 
https://sqlitebrowser.org/dl/

And finally, you need to install the __*SQLite Viewer*__ extension in your code editor 



## Start


#### To launch the game, enter the command :  
1) To access the web page you had to write in the terminal  :
```bash
go run main.go
```
2) Then open a web page and write (Preferably run on Google Chrome): 
```bash
http://localhost:8080/
```

## How To play

You'll have access to 3 multiplayer games, with a choice of Blind Test, Deaf Test or Petit Bac. You can join your friend by entering her game session code.

### Blind Test 

You have to listen to an audio file and give the title of the music, the fastest wins points.

### Deaf Test

You'll have to read the lyrics of a piece of music and name the title, the fastest wins points.

### Petit Bac
You have to fill in the requested information as quickly as possible, respecting the requested letter, and the game creator will have the honour of correcting the handles.

## Version 
go 1.21.0

## Author
By Erwan AGESNE, Brandon LUTULA, Aurélien BRANCO
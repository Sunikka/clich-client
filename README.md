# CLI Chat application in Go

A chat application in Go, where your terminal serves as the interface! The project is mainly for fun and learning. It's also a work in progress. 

![alt text](chat_view.png)

## Working on (Core functionalities):
  - Refactoring the code and returning functionality after moving to multirepo architecture. 
  
## Implemented
 -  The UI's for login screen and chat view
 -  Websocket communication between client and server
 -  Login logic works
 -  Messaging almost working (client sends and receives messages, but currently can't be rendered into the UI)

## Working on (Core functionalities):
  - Refactoring the code and returning functionality after moving to multirepo architecture. Currently trying to fix the issue with message handler not returning a bubbletea message to the main model.

## Tech stack:
- Go for frontend and backend
- Charmbracelet/bubbletea CLI frontend
- Websockets + Protobuf communication (Protobuf to be implemented)

![alt text](login.png)
# CLI Chat application in Go

A chat application in Go, where your terminal serves as the interface! The project is mainly for fun and learning. It's also a work in progress.

![alt text](chat_view.png)
## Implemented

- Bubbletea frontend UI base (Used the example from their github for starting the project out, so I can test my backend which I am focusing on first. I'm planning to customize the UI a lot later on)
- Basic websocket server
- Some basic functionality implemented:
  - Messages of the connected clients are rendered into the UI

## What I plan working on next (Core functionalities):
  - Make the register button functional in the login screen (the backend endpoint is already implemented)
  - Contacts/Friends system
  - Chat rooms

## Tech stack:
- Both frontend and backend are written in go
- Charmbracelet/bubbletea CLI frontend
- Websockets + Protobuf communication (Protobuf to be implemented)

![alt text](login.png)

# CLI Chat application in Go

A chat application in Go, where your terminal serves as the interface! The project is mainly for fun and learning. It's also a work in progress.

![alt text](chat_view.png)
## Implemented
- Some basic functionality implemented:
  - Messages of the connected clients are rendered into the UI
  - Working login system (even though far from complete)
- Custom themes (Fetched from ~/.config/clich/default_theme.yml)
  - Planning to implement a theme switcher inside the app, but currently the colors can be switched by adjusting the values inside the default_theme YAML file.



## What I plan working on next (Core functionalities):
  - Contacts/Friends system
  - Chat rooms

## Tech stack:
- Both frontend and backend are written in go
- Charmbracelet/bubbletea CLI frontend
- Websockets + Protobuf communication (Protobuf to be implemented)

![alt text](login.png)

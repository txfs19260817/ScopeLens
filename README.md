# ScopeLens
ScopeLens is a web platform that allows trainers to share awesome teams with others around the world. 

[Features](https://scopelens.team/#/about)

## Built with
- [Gin-gonic](https://github.com/gin-gonic/gin) - Gin is a web framework written in Go (Golang).
- [Vuetify](https://github.com/vuetifyjs/vuetify) - Vuetify is a Vue UI Library with beautifully handcrafted Components using the Material Design specification.
- [MongoDB](https://github.com/mongodb/mongo) & [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) - MongoDB is a document-based database. MongoDB is one of the leading NoSQL databases.

## Docker
The project now provides Dockerfile which lets you build Scopelens Server under `./server` with a single command: `docker build -t <TAG> --build-arg PORT=<PORT> . `
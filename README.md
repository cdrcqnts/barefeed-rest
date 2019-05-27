# Barefeed - REST API
A most simple feed reader web app for podcasts. No registration. No cookies.

This is the backend REST API application using [gin-gonic](https://github.com/gin-gonic/gin).
For the corresponding frontend application, see [this repository](https://github.com/cdrcqnts/barefeed).

### How it works
A random unique key is generated when the user adds a feed.
The key is mapped as a parameter to the url. All further feeds added by the user are linked with the same key.
Barefeed therefore does not require any personal user data. 
No cookies are stored, the data is queried exclusively on the basis of the URL parameter key.

### Run
```
// Serves the project at localhost:8080
go run main.go 
```

#### License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/cdrcqnts/barefeed/blob/master/LICENSE) file for details.
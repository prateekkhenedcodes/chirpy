# Chirpy 

A HTTP webserver that manages chirps (kind of tweets), built with Golang, managing data with PostgreSQL database and with a json parser.

## Requirements

`go version go1.24.1` or greater. 

## Installation and setup 

- clone the project with `git clone https://github.com/prateekkhenedcodes/chirpy.git`
- Install project dependencies with `git mod tidy` in the root of the project

## Running the server

Run the server with `go run .` in the root of the project or compile and execute the binary with `go build && ./chirpy`

## Endpoints

```go 
    mux.Handle("/app/", apicfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", ReadinessHandler)
	mux.HandleFunc("GET /admin/metrics", apicfg.CountHandler)
	mux.HandleFunc("POST /admin/reset", apicfg.ResetHandler)
	mux.HandleFunc("POST /api/users", apicfg.CreatUserHandler)
	mux.HandleFunc("POST /api/chirps", apicfg.ChirpHandler)
	mux.HandleFunc("GET /api/chirps", apicfg.ChirpsGetHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apicfg.ChirpGetHandler)
	mux.HandleFunc("POST /api/login", apicfg.LoginHandler)
	mux.HandleFunc("POST /api/refresh", apicfg.RefreshTokenHandler)
	mux.HandleFunc("POST /api/revoke", apicfg.RevokeTokenHandler)
	mux.HandleFunc("PUT /api/users", apicfg.UpdateUsersHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apicfg.ChirpDeleteHandler)
	mux.HandleFunc("POST /api/polka/webhooks", apicfg.WebhookHandler)
```

## Technicals

- JWT
- CRUD operations
- JSON 

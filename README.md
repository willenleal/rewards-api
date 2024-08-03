## Rewards API

### Running the app
#### Option 1: Docker and Docker Compose
Run `docker compose build && docker compose up`
#### Option 2: Makefile
Ensure Make and Go 1.22 or higher are installed on your sytem and run `make run`

#### Notes: 
1. The app will run on port 3000
2. The app is built on top of htt/net with as few 3rth party libraries as possible
3. Most of the api implementation can be found on api/impl.go file
4. I extended but did not modify the openapi spec to also codegen validation tags




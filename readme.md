# gRPC shop

### run in local:
cd colis/grpcs/auth
env CHADMIN_DB_HOST=127.0.0.1:27017 env CHADMIN_DB_NAME=cuahang env CHADMIN_DB_USER=cuahang env CHADMIN_DB_PASS=cuahang1234@ env PORT=32002 go run shop.go 

### run in docker:
docker build -t tidusant/colis-grpc-shop . && docker run -p 32002:8901 --env CLUSTERIP=127.0.0.1 --name colis-grpc-auth tidusant/colis-grpc-shop  
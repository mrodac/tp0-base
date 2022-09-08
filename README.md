git checkout ej1

### Ejercicio N°2:

./docker-compose-dev.sh 3 > docker-compose-dev.yaml

make docker-image

docker-compose -f docker-compose-dev.yaml up

sed -i 's/12345/54321/g' config/*

docker-compose -f docker-compose-dev.yaml up



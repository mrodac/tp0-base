git checkout ej8

### Ejercicio N°8:

./docker-compose-dev.sh 5 > docker-compose-dev.yaml 

make docker-compose-up && make docker-compose-logs &

make docker-compose-down

wc -l .data/winners 

### Protocolo:

Se agregan los mensajes nuevos:

Petición (TOTAL_WINNER_QUERY, NULL)
->
Respuesta (TOTAL_WINNER_RESPONSE, TotalWinnersResponse {totalWinners: uint32, pending: byte })

Luego de cada intercambio de mensaje se cierra la conexión. Para enviar un nuevo mensaje se el cliente debe volver a conectarse.

### Sincronización

Se usan dos enteros en memoria compartida (multiprocessing.Value) para compartir el estado del total de ganadores y peticiones pendientes.
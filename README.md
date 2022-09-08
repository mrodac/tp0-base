git checkout ej7

### Ejercicio N°7:

./docker-compose-dev.sh 5 > docker-compose-dev.yaml 

make docker-compose-up && make docker-compose-logs &

make docker-compose-down


### Sincronización

Se utiliza un pool de procesos de tamaño configurable por achivo. El servidor al recibir un mensaje de un cliente lo agrega a una cola del paquete multiprocessing para que uno de los procesos hijos lo tome y maneje.
Cuando el proceso padre recibe la señal SIGTERM, la reenvía a los hijos para que salgan de la espera sobre la cola y puedan ser joineados.
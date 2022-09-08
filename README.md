git checkout ej5

### Ejercicio N°5:

make docker-compose-up && make docker-compose-logs

make docker-compose-down


### Protocolo:

El protocolo consisten en envío de mensajes, enviando primero el tipo de mensaje en un byte y luego el contenido serializado.

Como solo hay un intercambio los mensajes son:

Petición
(WINNER_QUERY, ContestantsMessage {contestants: Contestant[] })  
->
Respuesta
(WINNER_RESPONSE, WinnersMessage {winners: uint32[] })


### Serialización:

Para la serialización se utiliza: 
- logitud fija para enteros
- byte delimitador para strings
- logitud + valor para colecciones

 
import socket
import logging
import signal
        
class SignalException(Exception):
	pass

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.active = True

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communication
        finishes, servers starts to accept new connections again
        """
        
        signal.signal(signal.SIGTERM, self.sig_handler)

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while self.active:
        	client_sock = self.__accept_new_connection()
        	if client_sock: 
        		self.__handle_client_connection(client_sock)
        

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            msg = client_sock.recv(1024).rstrip().decode('utf-8')
            logging.info(
                'Message received from connection {}. Msg: {}'
                .format(client_sock.getpeername(), msg))
            client_sock.send("Your Message has been received: {}\n".format(msg).encode('utf-8'))
        except OSError:
            logging.info("Error while reading socket {}".format(client_sock))
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info("Proceed to accept new connections")
        try:
        	c, addr = self._server_socket.accept()
        	logging.info('Got connection from {}'.format(addr))
        	return c
        except SignalException:
            logging.info("Got SIGTERM. Shutting down socket...")
            return None
            
    def sig_handler(self, sig, frame):
        self.active = False
        self._server_socket.shutdown(socket.SHUT_RDWR)
        raise SignalException()


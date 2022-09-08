import socket
import logging
import signal
import os

from . import client_request
        
class SignalException(Exception):
    pass

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
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
        
        while self.active:
            try:
                client_sock, addr = self._server_socket.accept()
                logging.info('Got connection from {}'.format(addr))
                
                request = client_request.ClientRequest(client_sock)
                request.handle()
            except SignalException:
                logging.info('Accept interrupted by signal')
                break
    
    
    def sig_handler(self, sig, frame):
        self.active = False
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()
        raise SignalException()


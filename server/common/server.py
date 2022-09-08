import socket
import logging
import signal
import os
from multiprocessing import Queue, Process, Value, active_children

from . import client_request
        
class SignalException(Exception):
	pass

def child_signal_handler(sig, frame):
    raise SignalException() 

def do_work(workQueue, totalWinners, processingCount):
    logging.info("Starting child process {}".format(os.getpid()))
    signal.signal(signal.SIGTERM, child_signal_handler)

    while True:
        try:
            request = workQueue.get()
            request.handle(totalWinners, processingCount)
        except SignalException:
            logging.debug("Process {} got SIGTERM".format(os.getpid()))
            break
        except ValueError:
            logging.debug("Process {} queue closed".format(os.getpid()))
            break

    logging.debug("Process {} stopped".format(os.getpid()))

class Server:
    def __init__(self, port, server_child_processes):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        self._server_socket.bind(('', port))
        self._server_socket.listen()

        self.workerCount = server_child_processes
        self.workQueue = Queue()


    def run(self):
        """
        Server loop

        Server that accept a new connections and delegates
        its handling to a process pool
        """
        
        signal.signal(signal.SIGTERM, self.sig_handler)
        
        totalWinners = Value('i', 0)
        processingCount = Value('i', 0)
        for i in range(self.workerCount):
            process = Process(target=do_work, args=(self.workQueue, totalWinners, processingCount))
            process.start()
        
        while True:
            try:
                client_sock, addr = self._server_socket.accept()
                logging.info('Got connection from {}'.format(addr))
                
                request = client_request.ClientRequest(client_sock)
                self.workQueue.put(request)
            except SignalException:
                logging.info('Accept interrupted by signal')
                break
        
        for process in active_children():
            process.join()
    
    
    def sig_handler(self, sig, frame):
        logging.debug("Parent process got SIGTERM")

        self.workQueue.close()

        for process in active_children():
            process.terminate()

        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()
        raise SignalException()


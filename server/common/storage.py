
import signal
import os
from multiprocessing import Process, Queue
import logging

from .utils import persist_winners, SignalException

def child_signal_handler(sig, frame):
    raise SignalException() 

def writer(dataQueue):
    logging.info("Starting writer process {}".format(os.getpid()))
    signal.signal(signal.SIGTERM, child_signal_handler)

    while True:
        try:
            winners = dataQueue.get()
            persist_winners(winners)
        except SignalException:
            logging.debug("Writer Process {} got SIGTERM".format(os.getpid()))
            break
        except ValueError:
            logging.debug("Writer Process {} queue closed".format(os.getpid()))
            break
    logging.debug("Writer Process {} stopped".format(os.getpid()))

class Storage:

    process = None
    dataQueue = None
    
    @staticmethod
    def start():
        Storage.dataQueue = Queue()
        Storage.process = Process(target=writer, args=(Storage.dataQueue,))
        Storage.process.start()
        
    @staticmethod
    def stop():
        Storage.dataQueue.close()
        Storage.process.terminate()
        Storage.process.join()
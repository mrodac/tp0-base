import socket
import logging
from multiprocessing import Queue, Value

from . import messages
from . import utils

class ClientRequest:
    def __init__(self, conn: socket.socket):
        self.conn = conn

    def handle(self, dataQueue:Queue, totalWinners: Value, processingCount: Value):
        msgType, msg = messages.read_message(self.conn)
        logging.debug('handling {} {}'.format(msgType, msg))

        match msgType:
            case messages.MessageType.WINNER_QUERY:
                self.handle_winners_query(msg.contestants, dataQueue, totalWinners, processingCount)
            case messages.MessageType.TOTAL_WINNER_QUERY:
                self.handle_total_winners_query(totalWinners, processingCount)
            case _:
                logging.error("Got unkwown msgType {}".format(msgType))


    def handle_winners_query(self, contestants, dataQueue:Queue, totalWinners: Value, processingCount: Value):
        with processingCount.get_lock():
            processingCount.value += 1
        
        msgType = messages.MessageType.WINNER_RESPONSE
        msg = messages.WinnersMessage()
        
        winners = list(filter(utils.is_winner, contestants))
        dataQueue.put(winners)
        msg.winners = list(map(lambda x: x.document, winners))

        with totalWinners.get_lock():
            totalWinners.value += len(msg.winners)

        with processingCount.get_lock():
            processingCount.value -= 1
        
        logging.debug('Sending response {} {}'.format(msgType, msg))
        messages.write_message(msgType, msg, self.conn)


    def handle_total_winners_query(self, totalWinners: Value, processingCount: Value):
        
        msg = messages.TotalWinnersResponse
        with processingCount.get_lock():
            msg.pending = processingCount.value

        with totalWinners.get_lock():
            msg.totalWinners  = totalWinners.value
        
        msgType = messages.MessageType.TOTAL_WINNER_RESPONSE
        
        logging.debug('Sending response {} {}'.format(msgType, msg))
        messages.write_message(msgType, msg, self.conn)        

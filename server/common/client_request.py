import socket
import logging

from . import messages
from . import utils
from .storage import Storage

class ClientRequest:
    def __init__(self, conn: socket.socket):
        self.conn = conn

    def handle(self, dataQueue):
        msgType, msg = messages.read_message(self.conn)

        logging.debug('handling {} {}'.format(msgType, msg))

        match msgType:
            case messages.MessageType.WINNER_QUERY:
                self.handle_winners_query(msg.contestants, dataQueue)
            case _:
                logging.error("Got unkwown msgType {}".format(msgType))


    def handle_winners_query(self, contestants, dataQueue):
        msgType = messages.MessageType.WINNER_RESPONSE
        msg = messages.WinnersMessage()

        winners = list(filter(utils.is_winner, contestants))
        dataQueue.put(winners)
        msg.winners = list(map(lambda x: x.document, winners))
        
        logging.debug('Sending response {} {}'.format(msgType, msg))
        messages.write_message(msgType, msg, self.conn)
        
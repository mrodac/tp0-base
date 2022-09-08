import socket
import logging

from . import messages
from . import utils

class ClientRequest:
    def __init__(self, conn: socket.socket):
        self.conn = conn

    def handle(self):
        msgType, msg = messages.read_message(self.conn)

        logging.debug('handling {} {}'.format(msgType, msg))

        match msgType:
            case messages.MessageType.WINNER_QUERY:
                self.handle_winners_query(msg.contestants)
            case _:
                logging.error("Got unkwown msgType {}".format(msgType))


    def handle_winners_query(self, contestants):
        msgType = messages.MessageType.WINNER_RESPONSE
        msg = messages.WinnersMessage()
        msg.winners = utils.winners_query(contestants)
        
        logging.debug('Sending response {} {}'.format(msgType, msg))
        messages.write_message(msgType, msg, self.conn)
        
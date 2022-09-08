import socket
import logging
from enum import IntEnum

from . import utils


class MessageType(IntEnum):
    WINNER_QUERY = 0
    WINNER_RESPONSE = 1

def read(conn, size):
    data = b""
    remaining = size
    while remaining != 0:
        data += conn.recv(remaining)
        remaining = size - len(data)
    return data

def read_delim(conn: socket.socket):
    data = b""
    
    c = read(conn, 1)
    while c != b'\0':
        data += c
        c = read(conn, 1)
    return data


class ContestantsMessage:
    def __init__(self):
        self.contestants = []
    
    def __repr__(self):
        return 'ContestantsMessage [' + str(len(self.contestants)) + ']'

class WinnersMessage:
    def __init__(self):
        self.winners = []
    
    def __repr__(self):
        return 'WinnersMessage [' + str(len(self.winners)) + ']'


def read_contestant(conn: socket.socket):
    b_document = read(conn, 4)
    document = int.from_bytes(b_document, byteorder='big')
    first_name = read_delim(conn).decode('utf-8')
    last_name = read_delim(conn).decode('utf-8')
    birthdate = read_delim(conn).decode('utf-8')

    return utils.Contestant(first_name, last_name, document, birthdate)
    

def read_winner_query(conn: socket.socket):
    b_size = read(conn, 4)
    items = int.from_bytes(b_size, byteorder='big')

    msg = ContestantsMessage()

    for i in range(items):
        msg.contestants.append(read_contestant(conn))

    return msg


def read_message(conn: socket.socket):
    ordinal = read(conn, 1)[0]
    msgType = MessageType(ordinal)
    logging.debug("Got msgType {}".format(msgType))
    
    match msgType:
        case MessageType.WINNER_QUERY:
            msg = read_winner_query(conn)
            return msgType, msg
        case _:
            logging.error("Got unkwown msgType {}".format(ordinal))


def write_winner_response(msg, conn: socket.socket):
    conn.sendall(len(msg.winners).to_bytes(4, 'big'))

    for winner in msg.winners:
        conn.sendall(winner.to_bytes(4, 'big'))


def write_message(msgType, msg, conn: socket.socket):
    ordinal = int(msgType)
    conn.sendall(ordinal.to_bytes(1, 'big'))

    match msgType:
        case MessageType.WINNER_RESPONSE:
            msg = write_winner_response(msg, conn)
        case _:
            logging.error("Got unkwown msgType {}".format(ordinal))



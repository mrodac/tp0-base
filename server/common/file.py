STORAGE = "./winners"

class File:
    _file = None

    @staticmethod
    def open():
       File._file = open(STORAGE, 'a+')

    @staticmethod
    def write(line):
       File._file.write(line)

    @staticmethod
    def close():
       File._file.close()
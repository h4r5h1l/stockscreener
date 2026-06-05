from ibapi.client import EClient
from ibapi.wrapper import EWrapper
import time
import threading
from datetime import datetime

class IBFundamentals(EWrapper, EClient):
    def __init__(self):
        EClient.__init__(self, self)

    def currentTime(self, time: int):
        print("Current Time:", datetime.fromtimestamp(time).strftime('%Y-%m-%d %H:%M:%S'))
        
    
app = IBFundamentals()
app.connect("localhost", 4001, clientId=1)
print("Connected to TWS API: ", app.isConnected())
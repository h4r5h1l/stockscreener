import sys
import time
from ibapi.client import EClient
from ibapi.wrapper import EWrapper
from ibapi.contract import Contract

class XMLStreamer(EWrapper, EClient):
    def __init__(self, conid_list):
        EClient.__init__(self, self)
        self.conid_list = conid_list
        self.idx = 0
        self.req_to_conid = {}

    def nextValidId(self, orderId: int):
        # Socket connection established handshake; begin processing the queue
        self.request_next()

    def request_next(self):
        # If we have reached the end of the provided ConID targets, exit cleanly
        if self.idx >= len(self.conid_list):
            self.disconnect()
            sys.exit(0)
            
        current_conid = int(self.conid_list[self.idx])
        req_id = 30000 + self.idx
        self.req_to_conid[req_id] = current_conid
        
        contract = Contract()
        contract.conId = current_conid
        contract.exchange = "SMART" # Best routing practice for fundamental snapshot requests
        
        # Request the data safely
        self.reqFundamentalData(reqId=req_id, contract=contract, reportType="ReportSnapshot", fundamentalDataOptions=[])
        
        # Enforce IBKR Pacing Limitations (keep requests under 50/sec)
        time.sleep(0.04)

    def fundamentalData(self, reqId: int, data: str):
        conid = self.req_to_conid.get(reqId)
        if conid:
            # Package stdout cleanly with token indicators for the Go scanner
            sys.stdout.write(f"START_CONID:{conid}\n")
            sys.stdout.write(data)
            sys.stdout.write("\nEND_XML_BLOCK\n")
            sys.stdout.flush()
        
        # Move the pointer forward and execute next queue step
        self.idx += 1
        self.request_next()

    def error(self, reqId: int, errorCode: int, errorString: str, contract=None):
        # Ignore passive background systemic notice codes
        if errorCode in [2104, 2106, 2158]:
            return
            
        # Code 430: "No fundamental data found" (Common for tiny cap / tracking shares)
        if errorCode == 430 and reqId in self.req_to_conid:
            conid = self.req_to_conid.get(reqId)
            sys.stdout.write(f"START_CONID:{conid}\nNOT_FOUND\nEND_XML_BLOCK\n")
            sys.stdout.flush()
            
            self.idx += 1
            self.request_next()
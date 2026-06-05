import requests
import urllib3

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

BASE_URL = "https://localhost:5010/v1/api"

def post_wrapper(path:str, payload: dict):
    url = f"{BASE_URL}{path}"
    resp = requests.post(url, json=payload, verify=False)
    
    try:
        return resp.json()
    except Exception as e:
        print("Non-JSON response:", resp.text)
        return None
    
def get_wrapper(path:str, params: dict):
    url = f"{BASE_URL}{path}"
    resp = requests.get(url, params=params, verify=False)
    
    try:
        return resp.json()
    except Exception as e:
        print("Non-JSON response:", resp.text)
        return None

def auth_status():
    return get_wrapper("/iserver/auth/status", {})
    
def search_symbol(symbol):
    return post_wrapper("/iserver/secdef/search", {"symbol": symbol})

def get_conids_by_exch(exch):
    return get_wrapper("/trsrv/all-conids", {"exchange": exch})

def get_all_equities():
    return [
        {
            "ticker": item.get("ticker"),
            "conid": item.get("conid"),
            "exchange": exch,
        }
        for exch in ["NYSE", "NASDAQ", "AMEX", "ARCA", "BATS", "IEX"]
        for item in get_conids_by_exch(exch)
    ]

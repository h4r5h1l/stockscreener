import requests
import urllib3

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

BASE_URL = "https://localhost:5010/v1/api"

def search_symbol(symbol):
    url = f"{BASE_URL}/iserver/secdef/search"
    resp = requests.post(url, json={"symbol": symbol}, verify=False)
    return resp.json()

print(search_symbol("AAPL"))

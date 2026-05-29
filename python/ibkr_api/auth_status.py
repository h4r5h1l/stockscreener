import requests
import urllib3

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

BASE = "https://localhost:5010/v1/api"

def auth_status():
    r = requests.get(f"{BASE}/iserver/auth/status", verify=False)
    print(r.json())
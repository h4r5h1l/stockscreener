import os
import sys
import json

# Forces Python to look inside the ibpython folder for its modules
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

# Absolute package pathing relative to ibpython workspace root
from web_api.w_requests import get_all_equities
from tws_api.get_fundamentals import XMLStreamer

def run_universe_fetch():
    """Fetches broad stock definitions from the Client Portal API and passes back minified JSON."""
    try:
        data = get_all_equities()
        if not data:
            sys.stderr.write("Error: Client Portal Web API returned an empty universe list.\n")
            sys.exit(1)
            
        # Write clean JSON back to the Go execution thread
        sys.stdout.write(json.dumps(data))
        sys.stdout.flush()
    except Exception as e:
        sys.stderr.write(f"Exception encountered inside web_api universe processing: {str(e)}\n")
        sys.exit(1)

def run_fundamental_stream(conids):
    """Binds to the TWS API background socket loops to stream data blocks."""
    if not conids:
        sys.stderr.write("Error: No target ConID arguments provided to stream.\n")
        sys.exit(1)
        
    app = XMLStreamer(conids)
    
    # Connect directly to your running instance of IB Gateway
    # Port 4001 is standard for IB Gateway Client connections (change to 7497 if using desktop TWS application instead)
    app.connect("127.0.0.1", 4001, clientId=1)
    app.run()

if __name__ == "__main__":
    if len(sys.argv) < 2:
        sys.stderr.write("Usage: python main.py [fetch_universe | stream_fundamentals] [conids...]\n")
        sys.exit(1)

    sub_command = sys.argv[1]

    if sub_command == "fetch_universe":
        run_universe_fetch()
        
    elif sub_command == "stream_fundamentals":
        # Pull everything after the second index positional string argument as target tokens
        target_arguments = sys.argv[2:]
        run_fundamental_stream(target_arguments)
        
    else:
        sys.stderr.write(f"Unknown system runtime execution argument flag: {sub_command}\n")
        sys.exit(1)
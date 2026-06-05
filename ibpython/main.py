from web_api.w_requests import auth_status, get_all_equities, search_symbol
def main():
    print("Running IBKR Python entrypoint...!")
    auth_status()
    
    data = get_all_equities()
    for item in data[:1]:
        print(search_symbol(item["ticker"]))


if __name__ == "__main__":
    main()

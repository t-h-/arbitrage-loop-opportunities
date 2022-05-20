import requests

def calc_arbitrage_loop(max_path_length=3):
    exchange_rates, currencies = get_exchange_rates()
    beginning_cur = currencies.pop()
    res = {}
    backtrack(currencies, (beginning_cur,), max_path_length, res, exchange_rates)
    return res
    
def backtrack(currencies, path, max_path_length, res, exchange_rates):
    if len(path) > max_path_length:
        return
    if len(path) >= 3:
        arbitrage = calc_path_arbitrage(path, exchange_rates)
        if arbitrage > 1:
            res[path] = arbitrage
    for i, elem in enumerate(currencies):
        unvisited_currencies = currencies.copy()
        del unvisited_currencies[i]
        next_path = path + (elem,)
        backtrack(unvisited_currencies, next_path, max_path_length, res, exchange_rates)
        

def calc_path_arbitrage(path, exchange_rates):
    rotated_path = path[1:] + (path[0],)
    
    arbitrage = 1
    for c1, c2 in zip(path, rotated_path):
        cur_pair_key = (c1, c2)
        arbitrage *= exchange_rates[cur_pair_key]
    return arbitrage
    
    
def get_exchange_rates():
    resp = requests.get("https://api.swissborg.io/v1/challenge/rates")
    pairs = resp.json()
    pairs = {(k.split("-")[0], k.split("-")[1]): float(v) for k, v in pairs.items()}
    currencies = list({cur[0] for cur in pairs.keys()})
    return pairs, currencies


if __name__ == "__main__":
    res = calc_arbitrage_loop(max_path_length=4)
    print(res)
# Arbitrage Loops 
This implements exact solutions for calculating arbitrage loops utilizing some sort of fully exhaustive backtracking in Python and Go. The Python version was to verify the approach. As I am currently getting into Golang and wanted to know more about its powerful concurrency model, I took some additional time to actually implement a concurrent version of the algorithm. The code is not fully hardened and tested.

## Complexity
Given `n` the number of currencies and `k` the maximum number of hops, the complexity is `O((n-1)!/(n-k-1)!)` = `O(n!/(n-k)!)`. Solutions are circular, i.e. a path "USD-BTC-EUR-USD" provides the same arbitrage as "BTC-EUR-USD-BTC".

## Parallelism
With a large number of potential currencies to handle, the runtime of this algorithm explodes. Fortunately, the problem is easily parallelizable as kicking of calculations with a given `path` and `unvisitedCurrencies` will independently yield unique solutions. Therefore, we can distribute the computation by generating the first `n` levels of `paths` and according `unvisitedCurrencies`, then kicking off the computations on worker nodes and then combining the returned results. 

Technically, this could be achieved by wrapping the code in a (REST, gRPC, async messaging) API and querying the workers with the generated parameters. Alternatively, if we already have some system for distributed processing in place, we can utilize that (I am thinking of something like Celery in the Python universe or big guns like Hadoop in the Java world).

If this is still not fast enough, we can try to further reduce runtime by not calculating exact but only approximated solutions:
- Randomly subsample from the large set of currencies and calculate a solution for each subset. This will of course produce some overlap in solutions, but with the right subsampling technique and parameters it could still be efficient enough. After calculating the solutions of the subsets, reduce them to a unique set of solutions.
- Maybe we can use some heuristics for pruning the `unvisitedCurrencies` array before doing the recursion. For that, we could keep track of how often a given currency pair or currency was seen in a solution. After the algorithm has gathered "sufficient" data, we can sort and prune the `unvisitedCurrencies` array accordingly. This is more of a last resort solution though :)

## Run
- Go: `go run main.go` from within the `arbitrageloop-go` directory should do the trick. Parameters given in the main() function should be self-explanatory.
    - there are some tests implemented for debugging and checking purposes during implementation, no assertions are made. Tests can be run with `go test -v`
    - I implemented a method to create synthetic currency and exchange rate data to test the implementation on a larger scale.
- Python: `python3 arbitrage-loop.py` from within the `arbitrageloop-py` directory. Package `requests` needed.
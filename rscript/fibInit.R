fib <- function(n) {
  if (n <= 0) {
    stop("n must be a positive integer")
  }
  if (n == 1) {
    return(0)
  } else if (n == 2) {
    return(1)
  } else {
    return(fib(n-1) + fib(n-2))
  }
}
source("./rscript/http.R")
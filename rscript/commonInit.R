library(hiplotlib)
source("./rscript/http.R")

hiFunc <- function(inputFile = "", confFile = "", outputFilePrefix = "", tool = "", module = "") {
  err <- ""
  glo_env_vars <- ls(envir = .GlobalEnv)
  pkgs <- (.packages())
  log <- ""
  opt <<- list(
    module = module,
    tool = tool,
    inputFile = inputFile,
    confFile = confFile,
    outputFilePrefix = outputFilePrefix
  )

  uuid <- system(sprintf("grep uuid %s | cut -f 2 -d ':' | sed 's;[\", ];;g'",
    confFile), inter = T)
  cat(sprintf("[%s] %s %s/%s %s\n", Sys.time(),
    ifelse(length(uuid) == 0, "Anonymous running", paste0(uuid, " running")),
    module, tool, dirname(outputFilePrefix)))
  tb <- ""

  script_dir <<- sprintf("../%s/", module)
  print(script_dir)
  options(hiplotlib.script_dir=script_dir)
  wd <- getwd()
  dir.create(file.path(dirname(opt$outputFilePrefix), "log"))
  err <- tryCatch(
    suppressWarnings(run_hiplot()),
    error = function(e) {
       print(e)
       print(conditionMessage(e))
       traceback()
      return(as.character(e))
    }, warning = function(w) {
      traceback()
      return(as.character(w))
    }
  )
  setwd(wd)
  callback(taskID, "common", FALSE)

  glo_env_vars2 <- ls(envir = .GlobalEnv)
  pkgs2 <- (.packages())
  pkgs2 <- pkgs2[!pkgs2 %in% pkgs]
  sapply(pkgs2, function(x) {
    detach(sprintf("package:%s", x), character.only = TRUE)
  })
  rm(list = glo_env_vars2[!glo_env_vars2 %in% glo_env_vars], envir = .GlobalEnv)
  gc()
  list(err = err)
}

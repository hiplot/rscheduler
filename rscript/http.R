callback = function(taskID, taskName, kill = FALSE){
        httr::GET(paste("http://localhost:8080/completed?",
        "taskID=", taskID,
        "&taskName=", taskName,
        "&kill=", kill,
        sep=""))
    }
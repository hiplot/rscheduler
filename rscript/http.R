callback = function(taskID, taskName, kill = 0){
        httr::GET(paste("http://localhost:8080/completed?",
        "taskID=", taskID,
        "&taskName=", taskName,
        "&kill=", kill,
        sep=""))
    }
callback = function(taskID, taskName){
        httr::GET(paste("http://localhost:8080/completed?taskID=",taskID, "&taskName=", taskName, sep=""))
    }
callback = function(taskID){
        httr::GET("http://localhost:8080/completed?taskID=123&taskName=add")
    }
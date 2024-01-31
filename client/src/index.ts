import * as http from "http"

http.get("http://localhost:8080/count", (res) => {
    const { statusCode } = res;
    
    let error;
    if (statusCode !== 200) {
        error = new Error("Request Failed.\n" + 
                         `Status Code: ${statusCode}`);
    }

    if (error) {
        console.error(error.message);
        res.resume();
        return;
    }

    res.setEncoding("utf-8");

    let rawData = "";

    res.on("data", (chunk) => { rawData += chunk });

    res.on("end", () => {
        console.log(rawData);
    }).on("error", (e) => {
        console.log(`Got error: ${e.message}`);
    });
})

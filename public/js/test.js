$(function() {

    ws = new WebSocket("ws://aaronstgeorge.co/mpws");
    // WebSocket connection established
    ws.onopen = function(e) {};

    ws.onmessage = function(event) {
        obj = JSON.parse(event.data);
        if (obj.Content == "PING") {
            document.getElementById(obj.Origin.Id).style.color = "green";
            ws.send("PONG");
        }
        if (obj.Content == "PONG") {
            document.getElementById(obj.Origin.Id).style.color = "green";
        }
        if (obj.Content == "CLOSE") {
            document.getElementById(obj.Origin.Id).style.color = "black";
        }
    };
});

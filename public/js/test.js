$(function() {

    ws = new WebSocket("ws://aaronstgeorge.co/ws");
    // WebSocket connection established
    ws.onopen = function(e) {
        alert("Connection established");
    };
});

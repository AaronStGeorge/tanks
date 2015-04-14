$(function() {

    var conn;
    var msg = $("#msg");
    var log = $("#log");

    function appendLog(msg) {
        var d = log[0];
        var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
        msg.appendTo(log);
        if (doScroll) {
            d.scrollTop = d.scrollHeight - d.clientHeight;
        }
    }

    $("#form").submit(function() {
        if (!conn) {
            return false;
        }
        if (!msg.val()) {
            return false;
        }
        conn.send(msg.val());
        msg.val("");
        return false;
    });

    conn = new WebSocket("ws://aaronstgeorge.co/ws");
    conn.onclose = function(evt) {
        appendLog($("<div><b>Connection closed.</b></div>"));
    };
    conn.onmessage = function(evt) {
        appendLog($("<div/>").text(evt.data));
    };
});

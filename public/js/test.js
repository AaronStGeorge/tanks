var me;
var ws;

$(function() {

    ws = new WebSocket("ws://aaronstgeorge.co/mpws");
    // WebSocket connection established

    ws.onmessage = function(event) {

        obj = JSON.parse(event.data);

        if (obj.Content == "INIT") {
            me = obj.Origin;
            sessionStorage.setItem("me", JSON.stringify(me));
            ws.send(JSON.stringify({
                Origin: me,
                PubTo: me.Twitter,
                Content: "PING"
            }));
        } else if (obj.Content == "PING") {
            document.getElementById(obj.Origin.Id).style.color = "green";
            ws.send(JSON.stringify({
                Origin: me,
                PubTo: obj.Origin.PhoneNumber,
                Content: "PONG"
            }));
        } else if (obj.Content == "PONG") {
            document.getElementById(obj.Origin.Id).style.color = "green";
        } else if (obj.Content == "CLOSE") {
            document.getElementById(obj.Origin.Id).style.color = "black";
        } else if (obj.Content == "ASK") {
            var r = confirm("Play a game with " + obj.Origin.UserName + "?");
            if (r === true) {
                // pass CONFIRM to other player
                ws.send(JSON.stringify({
                    Origin: me,
                    PubTo: obj.Origin.PhoneNumber,
                    Content: "CONFIRM"
                }));
                // store opponent
                localStorage.setItem("lastname", "Smith");
                // redirect to game play
                window.location.href = "/play";

            } else {
                ws.send(JSON.stringify({
                    Origin: me,
                    PubTo: obj.Origin.PhoneNumber,
                    Content: "DENY"
                }));
            }
        } else if (obj.Content == "CONFIRM") {
            // save opponent in session
            sessionStorage.setItem("opponent", JSON.stringify(obj.Origin));
            // redirect to game play
            window.location.href = "/play";
        } else if (obj.Content == "DENY") {
            alert(obj.UserName + " is busy");
        } else {
            alert(event.data);
        }
    };
});

function myJsFunc(PhoneNumber) {
    ws.send(JSON.stringify({
        Origin: me,
        PubTo: PhoneNumber,
        Content: "ASK"
    }));
}

var me,
    opponent,
    ws,
    width,
    height;

var pixelsPerTick = 10;


$(function() {

    /* TODO: uncoment these lines
       // remove old values
       sessionStorage.removeItem("me");
       sessionStorage.removeItem("opponent");
      */

    ws = new WebSocket("ws://aaronstgeorge.co/gpws");

    ws.onmessage = function(event) {};
});

width = $(window).width();
height = $(window).height();

// who am I playing against?
me = JSON.parse(sessionStorage.getItem("me"));
opponent = JSON.parse(sessionStorage.getItem("opponent"));

me.coords = {
    "cx": width / 3,
    "cy": height / 2
};
opponent.coords = {
    "cx": 2 * width / 3,
    "cy": height / 2
};


var dataset = [me, opponent];

var svg = d3.select("body").append("svg")
    .attr("width", width)
    .attr("height", height);


svg.selectAll("circle")
    .data(dataset)
    .enter().append("circle")
    .attr("id", function(d, i) {
        return d.PhoneNumber;
    })
    .attr("r", 100)
    .attr("cx", function(d) {
        return d.coords.cx;
    })
    .attr("cy", function(d) {
        return d.coords.cy;
    })
    .style("fill", function(d, i) {
        if (i === 0) {
            return "blue";
        } else {
            return "red";
        }
    })
    .style("stroke", "black")
    .style("stroke-width", "3");


function move() {
    svg.select("#" + me.PhoneNumber)
        .transition()
        .ease("linear")
        .duration(300)
        .attr("cx", function(d) {
            return d.coords.cx;
        })
        .attr("cy", function(d) {
            return d.coords.cy;
        });
}


var map = {
    37: false, // left
    38: false, // up
    39: false, // right
    40: false, // down
};

$(document).keydown(function(e) {
    e.preventDefault(); // prevent the default action (scroll / move caret)
    if (e.keyCode in map) {
        map[e.keyCode] = true;
        if (map[37] && map[38] && !map[39] && !map[40]) { // up and left
            me.coords.cx -= pixelsPerTick;
            me.coords.cy -= pixelsPerTick;
            move();
        } else if (!map[37] && map[38] && map[39] && !map[40]) { // up and right
            me.coords.cx += pixelsPerTick;
            me.coords.cy -= pixelsPerTick;
            move();
        } else if (map[37] && !map[38] && !map[39] && map[40]) { // down and left
            me.coords.cx -= pixelsPerTick;
            me.coords.cy += pixelsPerTick;
            move();
        } else if (!map[37] && !map[38] && map[39] && map[40]) { // down and right
            me.coords.cx += pixelsPerTick;
            me.coords.cy += pixelsPerTick;
            move();
        } else if (map[37] && !map[38] && !map[39] && !map[40]) { // left
            me.coords.cx -= pixelsPerTick;
            move();
        } else if (!map[37] && map[38] && !map[39] && !map[40]) { // up
            me.coords.cy -= pixelsPerTick;
            move();
        } else if (!map[37] && !map[38] && map[39] && !map[40]) { // right
            me.coords.cx += pixelsPerTick;
            move();
        } else if (!map[37] && !map[38] && !map[39] && map[40]) { // down 
            me.coords.cy += pixelsPerTick;
            move();
        }
    }
}).keyup(function(e) {
    if (e.keyCode in map) {
        map[e.keyCode] = false;
    }
});

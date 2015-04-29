var me,
    opponent,
    ws,
    width,
    height;

var pixelsPerTick = 10;



$(function() {

    // remove old values
    sessionStorage.removeItem("me");
    sessionStorage.removeItem("opponent");

    ws = new WebSocket("ws://aaronstgeorge.co/gpws");

    ws.onmessage = function(event) {
        obj = JSON.parse(event.data);

        if ('x' in obj.Content && 'y' in obj.Content) {
            moveOpponent(obj.Content.x, obj.Content.y);
        } else {
            alert(event.data);
        }
    };
});

width = $(window).width();
height = $(window).height();

// who am I playing against?
me = JSON.parse(sessionStorage.getItem("me"));
opponent = JSON.parse(sessionStorage.getItem("opponent"));

if (me.Id > opponent.Id) {
    meCoords = {
        "x": width / 3,
        "y": height / 2,
        "color": "blue",
        "name": me.UserName
    };
    opponentCoords = {
        "x": 2 * width / 3,
        "y": height / 2,
        "color": "red",
        "name": opponent.UserName
    };

} else {
    meCoords = {
        "x": 2 * width / 3,
        "y": height / 2,
        "color": "red",
        "name": me.UserName
    };
    opponentCoords = {
        "x": width / 3,
        "y": height / 2,
        "color": "blue",
        "name": opponent.UserName
    };
}

var dataset = [meCoords, opponentCoords];

// height of text in pixels
var textpx = 40;

function adjustYCoordForText(yCoord) {
    return yCoord + textpx / 4.0;
}

var svg = d3.select("body").append("svg")
    .attr("width", width)
    .attr("height", height);

svg.selectAll("circle")
    .data(dataset)
    .enter().append("circle")
    .attr("id", function(d, i) {
        if (i === 0) {
            return me.PhoneNumber;
        } else {
            return opponent.PhoneNumber;
        }
    })
    .attr("r", 100)
    .attr("cx", function(d) {
        return d.x;
    })
    .attr("cy", function(d) {
        return d.y;
    })
    .style("fill", function(d) {
        return d.color;
    })
    .style("stroke", "black")
    .style("stroke-width", "3");

//Add the SVG Text Element to the svgContainer
var text = svg.selectAll("text")
    .data(dataset)
    .enter()
    .append("text");

//Add SVG Text Element Attributes
var textLabels = text
    .attr("x", function(d) {
        return d.x;
    })
    .attr("y", function(d) {
        return adjustYCoordForText(d.y);
    })
    .text(function(d) {
        return d.name;
    })
    .attr("id", function(d, i) {
        if (i === 0) {
            return "text" + me.PhoneNumber;
        } else {
            return "text" + opponent.PhoneNumber;
        }
    })
    .attr("font-size", textpx + "px")
    .attr("font-family", "sans-serif")
    .attr("text-anchor", "middle")
    .attr("fill", "black");

function move() {
    svg.select("#" + me.PhoneNumber)
        .transition()
        .ease("linear")
        .duration(300)
        .attr("cx", function(d) {
            return d.x;
        })
        .attr("cy", function(d) {
            return d.y;
        });
    svg.select("#text" + me.PhoneNumber)
        .transition()
        .ease("linear")
        .duration(300)
        .attr("x", function(d) {
            return d.x;
        })
        .attr("y", function(d) {
            return adjustYCoordForText(d.y);
        });
    ws.send(JSON.stringify({
        Origin: me,
        PubTo: opponent.PhoneNumber,
        Content: {
            "x": meCoords.x,
            "y": meCoords.y
        }
    }));
}

function moveOpponent(x, y) {
    opponentCoords.x = x;
    opponentCoords.y = y;
    svg.select("#" + opponent.PhoneNumber)
        .transition()
        .ease("linear")
        .duration(300)
        .attr("cx", function(d) {
            return d.x;
        })
        .attr("cy", function(d) {
            return d.y;
        });
    svg.select("#text" + opponent.PhoneNumber)
        .transition()
        .ease("linear")
        .duration(300)
        .attr("x", function(d) {
            return d.x;
        })
        .attr("y", function(d) {
            return adjustYCoordForText(d.y);
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
            meCoords.x -= pixelsPerTick;
            meCoords.y -= pixelsPerTick;
            move();
        } else if (!map[37] && map[38] && map[39] && !map[40]) { // up and right
            meCoords.x += pixelsPerTick;
            meCoords.y -= pixelsPerTick;
            move();
        } else if (map[37] && !map[38] && !map[39] && map[40]) { // down and left
            meCoords.x -= pixelsPerTick;
            meCoords.y += pixelsPerTick;
            move();
        } else if (!map[37] && !map[38] && map[39] && map[40]) { // down and right
            meCoords.x += pixelsPerTick;
            meCoords.y += pixelsPerTick;
            move();
        } else if (map[37] && !map[38] && !map[39] && !map[40]) { // left
            meCoords.x -= pixelsPerTick;
            move();
        } else if (!map[37] && map[38] && !map[39] && !map[40]) { // up
            meCoords.y -= pixelsPerTick;
            move();
        } else if (!map[37] && !map[38] && map[39] && !map[40]) { // right
            meCoords.x += pixelsPerTick;
            move();
        } else if (!map[37] && !map[38] && !map[39] && map[40]) { // down 
            meCoords.y += pixelsPerTick;
            move();
        }
    }
}).keyup(function(e) {
    if (e.keyCode in map) {
        map[e.keyCode] = false;
    }
});

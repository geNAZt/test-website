var servers = {};
var sorted = [];
var conn = null;
var skip = 0;
var chart;
var offset = new Date().getTimezoneOffset() * 60;
var graphSettings = [];
var colorArray = ["#5ca1c8", "#e42570", "#ced8b5", "#8774de", "#b8a0d0", "#7a645d", "#4f5e73", "#46be6d", "#961f7f", "#8e7867", "#6d53a7", "#992c7e", "#5c274e", "#e67c4b", "#d2ad7e", "#ed06f", "#878ea7", "#c2ef88", "#4b5aca", "#b974c4", "#85f6bb", "#2feb6a", "#7a650a", "#70a1aa", "#762128", "#8e8cf", "#d8d81e", "#14533f", "#f4e9e1", "#70f317", "#755ce0", "#1b5aab", "#73d3fd", "#6f931e", "#2bfeea", "#3c5a53", "#e05e81", "#267118", "#26608a", "#351810", "#d0cb25", "#78b849", "#ffef26", "#6437bf", "#8133bb", "#354453", "#2ecaa9", "#cf6416", "#5def3d", "#1a6281", "#47532c", "#12ce13", "#55f153", "#6c8ff4", "#e32548", "#724925", "#cbfe76", "#dd1b04", "#7d1b14", "#21e130", "#60233e", "#2bb540", "#dd63be", "#63b267"];
var currentColorI = 0;
var time = 2 * 24 * 60;
var favicons = {};
var sendTimeout = -1;

//Better to construct options first and then pass it as a parameter
var options = {
    animationEnabled: true,
    title: {
        text: "Last 2 Days",
        fontSize: 12
    },
    zoomEnabled: true,
    axisX: {
        labelFontSize: 12
    },
    axisY: {
        labelFontSize: 12
    },
    legend: {
        fontSize: 12,
        cursor: "pointer"
    },
    backgroundColor: "#ebebeb",
    data: []
};

var wsFuncs = {
    servers: function (data) {
        servers = {};

        data.forEach(function (value) {
            value["Color"] = getColor();
            servers[value["Id"]] = value;
            conn.send("pings:" + value["Id"]);
        });

        sortServers( true );

        $('#page-selection').bootpag({
            total: Math.ceil(sorted.length / 5)
        }).on("page", function (event, num) {
            skip = 5 * (num - 1);
            sortServers( true );
        });
    },
    pings: function(data) {
        servers[data["Id"]]["Players"] = data["Players"];
        rerenderChart();
    },
    updatePlayer: function (data) {
        servers[data["Id"]]["Online"] = data["Online"];
        servers[data["Id"]]["Ping"] = data["Ping"];
        servers[data["Id"]]["Ping24"] = data["Ping24"];
        servers[data["Id"]]["Record"] = data["Record"];
        servers[data["Id"]]["Average"] = data["Average"];
        servers[data["Id"]]["Players"][data["Time"]] = data["Online"];

        newData = {};
        lowDate = data["Time"] - ( time * 60 );

        for (var key in servers[data["Id"]]["Players"]) {
            if (servers[data["Id"]]["Players"].hasOwnProperty(key) && key >= lowDate) {
                newData[key] = servers[data["Id"]]["Players"][key];
            }
        }

        servers[data["Id"]]["Players"] = newData;

        rerenderChart();
        sortServers( false );
    },
    maxPlayer: function(data) {
        servers[data["Id"]]["MaxPlayers"] = data["MaxPlayers"];
        sortServers( false );
    },
    favicon: function(data) {
        if (favicons[data["Server"]] !== undefined) {
            favicons[data["Server"]].attr("src", data["Icon"]);
        }
    }
};

function getColor() {
    if (currentColorI > colorArray.length - 1) {
        currentColorI = 0;
    }

    color = colorArray[currentColorI];
    currentColorI++;

    return color;
}

function rerenderChart() {
    options["data"] = [];

    // Convert them into the array
    serversCopy = servers;
    for (var key in serversCopy) {
        if (serversCopy.hasOwnProperty(key) && graphSettings.indexOf(serversCopy[key]["Name"]) == -1) {
            options["data"].push(generateData(serversCopy[key]));
        }
    }

    $("#chartContainer").CanvasJSChart().render();
}

function generateData(server) {
    data = {
        type: "line",
        xValueType: "dateTime",
        showInLegend: true,
        lineThickness: 2,
        color: server["Color"],
        indexLabelFontSize: 12,
        markerSize: 0,
        toolTipContent: "<b>" + server["Name"] + "</b> ({x})<br/>Players: {y}",
        name: server["Name"],
        dataPoints: []
    };

    if ( server["Players"] == undefined ) {
        return data;
    }

    length = Object.keys(server["Players"]).length;

    skip = 0;

    if ( length > 3000 ) {
        skip = ( length - 3000 ) / 3000;
    }

    counter = 0;
    currentTime = new Date().getTime();

    serversCopy = servers;
    for (var key in server["Players"]) {
        if (skip > counter) {
            counter++;
            continue;
        }

        counter = 0;

        if (server["Players"].hasOwnProperty(key)) {
            valueTime = ( parseInt(key) - offset ) * 1000;
            if (currentTime - valueTime > time * 60 * 1000) {
                continue;
            }

            data.dataPoints.push({
                x: ( parseInt(key) - offset ) * 1000,
                y: server["Players"][key]
            });
        }
    }

    return data;
}

function renderTable( renderAnimatedFavicons ) {
    table = $("<table />");
    table.addClass("table table-hover");
    table.append(createTH());

    if (renderAnimatedFavicons) {
        favicons = {};
    }

    counter = 0;
    rendered = 0;
    sorted.forEach(function (value) {
        if (counter < skip) {
            counter++;
            return;
        }

        if (rendered == 5) {
            return;
        }

        table.append(createTR(value, renderAnimatedFavicons));
        rendered++;
    });

    serverTableContainer = $('#server-table');
    serverTableContainer.children("table").remove();
    serverTableContainer.append(table);
}

function createTH() {
    return "<thead><tr><th></th><th>ID</th><th>Server Name</th><th>Minecraft IP</th><th>Website</th><th>Players</th><th>Record</th><th>Average (24h)</th><th>Ping</th></tr></thead>";
}

function createTR(server, renderAnimatedFavicons) {
    tr = $('<tr />');

    // Check
    checkTd = $('<td />');
    $('<input />', {
        type: 'checkbox',
        id: 'check_' + server["Id"],
        value: server["Name"],
        checked: graphSettings.indexOf(server["Name"]) == -1
    }).click(function () {
        if (graphSettings.indexOf(server["Name"]) == -1) {
            graphSettings.push(server["Name"]);
        } else {
            var index = graphSettings.indexOf(server["Name"]);
            if (index > -1) {
                graphSettings.splice(index, 1);
            }
        }

        rerenderChart();
    }).appendTo(checkTd);
    tr.append(checkTd);

    // ID
    idTd = $('<td />');
    idTd.text(server["Id"]);
    tr.append(idTd);

    // Name
    favicon = $('<img />');
    favicon.attr('src', server["Favicon"]);

    if ( renderAnimatedFavicons && conn != null ) {
        favicons[server["Name"]] = favicon;

        window.setTimeout(function() {
            conn.send("animated:" + server["Name"]);
        }, 500);
    }

    nameTd = $('<td />');
    nameTd.append(favicon);
    nameTd.append(server["Name"]);
    tr.append(nameTd);

    // IP
    ipTd = $('<td />');
    ipTd.text(server["IP"]);
    tr.append(ipTd);

    // Website
    websiteTd = $('<td />');
    link = $('<a />');
    link.attr('href', server['Website']);
    link.text(server['Website']);
    websiteTd.append(link);
    tr.append(websiteTd);

    // Players
    playersTd = $('<td />');
    if (server['Ping24'] !== undefined) {
        spanDir = $('<span />');
        if (server['Ping24'] > server['Online']) {
            spanDir.addClass("glyphicon glyphicon-arrow-down");
        } else {
            spanDir.addClass("glyphicon glyphicon-arrow-up");
        }
        playersTd.append(spanDir);
    }

    playersTd.append(server["Online"] + " / " + server["MaxPlayers"]);
    tr.append(playersTd);

    // Record
    recordTd = $('<td />');
    recordTd.text(server["Record"] + " Players");
    tr.append(recordTd);

    // Average
    averageTd = $('<td />');
    averageTd.text(server["Average"] + " Players");
    tr.append(averageTd);

    // Ping
    pingMS = Math.ceil(server["Ping"] / 1000000);
    pingTd = $('<td />');
    pingTd.text(pingMS + " ms");
    tr.append(pingTd);

    return tr;
}

function sortServers( renderAnimatedFavicons ) {
    serversCopy = servers;
    newSorted = [];

    // Convert them into the array
    for (var key in serversCopy) {
        if (serversCopy.hasOwnProperty(key)) {
            newSorted.push(serversCopy[key]);
        }
    }

    // Sort it
    newSorted.sort(function (a, b) {
        return b["Online"] - a["Online"];
    });

    sorted = newSorted;

    // Rerender it
    renderTable( renderAnimatedFavicons );
}

function connectToWebSocket(host) {
    // Connect to the WebSocket
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + host + "/ws");

        conn.onclose = function (evt) {
            console.log("Connection closed")
        };

        conn.onmessage = function (evt) {
            jsonData = JSON.parse(evt.data);
            if (jsonData["Ident"] !== undefined && wsFuncs[jsonData["Ident"]] !== undefined && jsonData["Value"] !== undefined) {
                wsFuncs[jsonData["Ident"]](jsonData["Value"]);
            }
        };
    }
}

$(document).ready(function () {
    $("#chartContainer").CanvasJSChart(options);

    var height = window.innerHeight - 49;
    $('#main').css('min-height', height);

    $("#slider").slider({
        range: "min",
        value: 2,
        min: 2,
        max: 60,
        slide: function( event, ui ) {
            time = ui.value * 24 * 60;
            options.title.text = "Last " + ui.value + " Days";

            if ( sendTimeout == -1 ) {
                sendTimeout = window.setTimeout(function() {
                    conn.send("range:" + time / (24*60));
                    sendTimeout = -1;
                }, 1000)
            }
        }
    });

    connectToWebSocket(host);
});
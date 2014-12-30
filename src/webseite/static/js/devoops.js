var servers = {};
var sorted = [];
var conn;
var skip = 0;
var chart;
var offset = new Date().getTimezoneOffset() * 60;


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
	colorSet: "colorset",
	backgroundColor: "#ebebeb",
	data: [ ]
};

var wsFuncs = {
	servers: function (data) {
		servers = {};

		data.forEach(function (value) {
			servers[value["Name"]] = value;
			options["data"].push(generateData(value));
		});

		$("#chartContainer").CanvasJSChart().render();
		sortServers();

		$('#page-selection').bootpag({
			total: Math.floor(sorted.length / 5) + 1
		}).on("page", function (event, /* page number here */ num) {
			skip = 5 * (num - 1);
			sortServers();
		});
	},
	updatePlayer: function (data) {
		servers[data["Name"]]["Online"] = data["Online"];
		servers[data["Name"]]["MaxPlayers"] = data["MaxPlayers"];
		servers[data["Name"]]["Ping"] = data["Ping"];
		servers[data["Name"]]["Players"].push({
			Time: data["Time"],
			Online: data["Online"]
		});

		options["data"] = [];

		// Convert them into the array
		serversCopy = servers;
		for (var key in serversCopy) {
			if (serversCopy.hasOwnProperty(key)) {
				options["data"].push(generateData(serversCopy[key]));
				newSorted.push(serversCopy[key]);
			}
		}

		$("#chartContainer").CanvasJSChart().render();
		sortServers();
	}
};

function generateData( server ) {
	data = {
		type: "line",
		xValueType: "dateTime",
		showInLegend: true,
		lineThickness: 2,
		indexLabelFontSize: 12,
		markerSize: 0,
		toolTipContent: "<b>" + server["Name"] + "</b> ({x})<br/>Players: {y}",
		name: server["Name"],
		dataPoints: [ ]
	};

	server["Players"].forEach(function(ping) {
		data.dataPoints.push({
			x: ( ping["Time"] - offset ) * 1000,
			y: ping["Online"]
		});
	});

	return data;
}

function renderTable() {
	table = $("<table />");
	table.addClass("table table-hover");
	table.append(createTH());

	counter = 0;
	rendered = 0;
	sorted.forEach(function(value) {
		if ( counter < skip ) {
			counter++;
			return;
		}

		if ( rendered == 5 ) {
			return;
		}

		table.append(createTR(value));
		rendered++;
	});

	serverTableContainer = $('#server-table');
	serverTableContainer.children("table").remove();
	serverTableContainer.append(table);
}

function createTH() {
	return "<thead><tr><th>#</th><th>Server Name</th><th>Minecraft IP</th><th>Website</th><th>Players</th><th>Record</th><th>Average (24h)</th><th>Ping</th></tr></thead>";
}

function createTR( server ) {
	tr = $('<tr />');

	// ID
	idTd = $('<td />');
	idTd.text( server["Id"] );
	tr.append(idTd);

	// Name
	favicon = $('<img />');
	favicon.attr('src', server["Favicon"]);
	nameTd = $('<td />');
	nameTd.append( favicon );
	nameTd.append( server["Name"] );
	tr.append(nameTd);

	// IP
	ipTd = $('<td />');
	ipTd.text( server["IP"] );
	tr.append(ipTd);

	// Website
	websiteTd = $('<td />');
	link = $('<a />');
	link.attr('href', server['Website']);
	link.text( server['Website'] );
	websiteTd.append( link );
	tr.append( websiteTd );

	// Players
	playersTd = $('<td />');
	playersTd.text( server["Online"] + " / " + server["MaxPlayers"] );
	tr.append( playersTd );

	// Record
	recordTd = $('<td />');
	recordTd.text( server["Record"] + " Players" );
	tr.append( recordTd );

	// Average
	averageTd = $('<td />');
	averageTd.text( server["Average"] + " Players" );
	tr.append( averageTd );

	// Ping
	pingMS = Math.ceil( server["Ping"] / 1000000 );
	pingTd = $('<td />');
	pingTd.text( pingMS + " ms" );
	tr.append( pingTd );

	return tr;
}

function sortServers() {
	serversCopy = servers;
	newSorted = [];

	// Convert them into the array
	for (var key in serversCopy) {
		if (serversCopy.hasOwnProperty(key)) {
			newSorted.push(serversCopy[key]);
		}
	}

	// Sort it
	newSorted.sort(function(a, b) {
		return b["Online"] - a["Online"];
	});

	sorted = newSorted;

	// Rerender it
	renderTable();
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

	var colorArray = ["#5ca1c8", "#e42570", "#ced8b5", "#8774de", "#b8a0d0", "#7a645d", "#4f5e73", "#46be6d", "#961f7f", "#8e7867", "#6d53a7", "#992c7e", "#5c274e", "#e67c4b", "#d2ad7e", "#ed06f", "#878ea7", "#c2ef88", "#4b5aca", "#b974c4", "#85f6bb", "#2feb6a", "#7a650a", "#70a1aa", "#762128", "#8e8cf", "#d8d81e", "#14533f", "#f4e9e1", "#70f317", "#755ce0", "#1b5aab", "#73d3fd", "#6f931e", "#2bfeea", "#3c5a53", "#e05e81", "#267118", "#26608a", "#351810", "#d0cb25", "#78b849", "#ffef26", "#6437bf", "#8133bb", "#354453", "#2ecaa9", "#cf6416", "#5def3d", "#1a6281", "#47532c", "#12ce13", "#55f153", "#6c8ff4", "#e32548", "#724925", "#cbfe76", "#dd1b04", "#7d1b14", "#21e130", "#60233e", "#2bb540", "#dd63be", "#63b267"];
	CanvasJS.addColorSet("colorset", colorArray);

	connectToWebSocket(host);
});
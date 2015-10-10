/**
 * Created by Fabian on 21.06.15.
 */
var offset = new Date().getTimezoneOffset() * 60;
var colorArray = [
    "#5ca1c8",
    "#e42570",
    "#ced8b5",
    "#8774de",
    "#b8a0d0",
    "#7a645d",
    "#4f5e73",
    "#46be6d",
    "#961f7f",
    "#8e7867",
    "#6d53a7",
    "#992c7e",
    "#5c274e",
    "#e67c4b",
    "#d2ad7e",
    "#ed06f",
    "#878ea7",
    "#c2ef88",
    "#4b5aca",
    "#b974c4",
    "#85f6bb",
    "#2feb6a",
    "#7a650a",
    "#70a1aa",
    "#762128",
    "#8e8cf",
    "#d8d81e",
    "#14533f",
    "#f4e9e1",
    "#70f317",
    "#755ce0",
    "#1b5aab",
    "#73d3fd",
    "#6f931e",
    "#2bfeea",
    "#3c5a53",
    "#e05e81",
    "#267118",
    "#26608a",
    "#351810",
    "#d0cb25",
    "#78b849",
    "#ffef26",
    "#6437bf",
    "#8133bb",
    "#354453",
    "#2ecaa9",
    "#cf6416",
    "#5def3d",
    "#1a6281",
    "#47532c",
    "#12ce13",
    "#55f153",
    "#6c8ff4",
    "#e32548",
    "#724925",
    "#cbfe76",
    "#dd1b04",
    "#7d1b14",
    "#21e130",
    "#60233e",
    "#2bb540",
    "#dd63be",
    "#63b267"
];

define(["app/timeSlider"], function (timeSlider) {
    return function Chart(ui) {
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
            },
            currentColorI = 0;

        // Buildup Chart
        $("#chartContainer").CanvasJSChart(options);

        // Private functions
        function getColor() {
            if (currentColorI > colorArray.length - 1) {
                currentColorI = 0;
            }

            var color = colorArray[currentColorI];
            currentColorI++;

            return color;
        }

        function generateData(server) {
            var data = {
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

            var length = Object.keys(server["Players"]).length,
                skipData = 0,
                counter = 0,
                currentTime = new Date().getTime();

            if ( length > 3000 ) {
                skipData = ( length - 3000 ) / 3000;
            }

            for (var key in server["Players"]) {
                if (skipData > counter) {
                    counter++;
                    continue;
                }

                counter = 0;

                if (server["Players"].hasOwnProperty(key)) {
                    var valueTime = ( parseInt(key) - offset ) * 1000;
                    if (currentTime - valueTime > timeSlider.getTime() * 60 * 1000) {
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

        this.onNewServer = function (server) {
            server["Color"] = getColor();
        };

        this.resetColors = function () {
            currentColorI = 0;
        };

        this.onChangeTime = function (time) {
            options.title.text = "Last " + time / ( 24 * 60 ) + " Days";
        };

        this.render = function () {
            options["data"] = [];

            // Convert them into the array
            var serversCopy = ui.getServers();
            for (var key in serversCopy) {
                if (serversCopy.hasOwnProperty(key)) {
                    options["data"].push(generateData(serversCopy[key]));
                }
            }

            console.log(options);

            $("#chartContainer").CanvasJSChart().render();
        };

        // Setup listeners
        ui.on("newServer", this.onNewServer);
        timeSlider.on("changeTime", this.onChangeTime);
    };
});
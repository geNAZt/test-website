/**
 * Created by Fabian on 20.06.15.
 *
 * All methods which are needed to communicate with the Backend over a websocket
 */
define(function () {
    return function Socket() {
        var socketConnection = null,
            self = this,
            wsFuncs = {};

        /**
         * Connect to a WebSocket server (/ws endpoint on given host)
         * @param host to which we want to connect
         */
        this.connectToWebSocket = function (host) {
            // Connect to the WebSocket
            if (window["WebSocket"]) {
                socketConnection = new WebSocket("ws://" + host + "/ws");

                // When the Connection closes log and reconnect
                socketConnection.onclose = function () {
                    console.log("Connection closed");

                    // Reconnect after 5 seconds
                    window.setTimeout(function () {
                        self.connectToWebSocket(host);
                    }, 5000);
                };

                socketConnection.onmessage = function (evt) {
                    jsonData = JSON.parse(evt.data);
                    if (jsonData["Ident"] !== undefined && wsFuncs[jsonData["Ident"]] !== undefined && jsonData["Value"] !== undefined) {
                        wsFuncs[jsonData["Ident"]](jsonData["Value"]);
                    }
                };
            } else {
                alert("This Website needs a Browser which supports WebSockets");
            }
        };

        /**
         * Register a new Function which should be called when Data is received with that ident
         * @param ident for which this callback should be called
         * @param cb which has one parameter for the data
         */
        this.registerFunction = function (ident, cb) {
            wsFuncs[ident] = cb;
        };

        /**
         * Tell the server which ping data it should send
         *
         * @param pings
         */
        this.sendPingIDs = function (pings) {
            var pingIds = "";

            for (var key in pings) {
                if (pings.hasOwnProperty(key)) {
                    pingIds += ":" + pings[key]["Id"];
                }
            }

            if (pingIds != "") {
                socketConnection.send("pings" + pingIds);
            }
        };

        /**
         * Tell the server how much data we want to get
         *
         * @param time
         */
        this.sendRange = function (time) {
            socketConnection.send("range:" + time / (24 * 60));
        };

        this.setView = function(view) {
            socketConnection.send("setview:" + view);
        }
    }
});
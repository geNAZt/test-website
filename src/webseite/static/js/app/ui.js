/**
 * Created by Fabian on 21.06.15.
 *
 * Chart and table rendering
 */
define(["app/chart", "app/timeSlider", "app/servertable", "app/viewselector", "lib/heir", "lib/eventEmitter"], function(Chart, timeSlider, ServerTable, ViewSelector, heir, eventEmitter) {
    function UI() {
        var socket = null,
            currentServers = {},
            allServers = {},
            viewSelector = new ViewSelector( this),
            serverTable = new ServerTable( this ),
            chart = new Chart( this ),
            self = this;

        /**
         * Bind the Websocket to this UI. It registeres all needed idents
         *
         * @param newSocket
         */
        this.bindSocket = function(newSocket) {
            socket = newSocket;
            socket.registerFunction("servers", this.onNewServers);
            socket.registerFunction("pings", this.onPings);
            socket.registerFunction("views", this.onViews);
            socket.registerFunction("updatePlayer", this.onPlayerUpdate);
            socket.registerFunction("maxPlayer", this.onMaxPlayerUpdate);
            socket.registerFunction("uptime", this.onUptimeUpdate);
        };

        this.getSocket = function() {
            return socket;
        };

        /**
         * Get current displayed servers
         *
         * @returns {{}}
         */
        this.getServers = function() {
            return currentServers;
        };

        this.onMaxPlayerUpdate = function(data) {
            allServers[data["Id"]]["MaxPlayers"] = data["MaxPlayers"];

            if (currentServers.hasOwnProperty(data["Id"])) {
                currentServers[data["Id"]] = allServers[data["Id"]];
                serverTable.render(false);
            }
        };

        this.onUptimeUpdate = function(data) {
            allServers[data["Id"]]["Uptime"] = data["Uptime"];
            allServers[data["Id"]]["UptimeLast"] = data["UptimeLast"];

            if (currentServers.hasOwnProperty(data["Id"])) {
                currentServers[data["Id"]] = allServers[data["Id"]];
                serverTable.render(false);
            }
        };

        this.onPlayerUpdate = function(data) {
            allServers[data["Id"]]["Online"] = data["Online"];
            allServers[data["Id"]]["Average"] = data["Average"];
            allServers[data["Id"]]["Record"] = data["Record"];
            allServers[data["Id"]]["Ping24"] = data["Ping24"];
            allServers[data["Id"]]["Players"][data["Time"]] = data["Online"];

            if (currentServers.hasOwnProperty(data["Id"])) {
                currentServers[data["Id"]] = allServers[data["Id"]];
                serverTable.render(false);
                chart.render();
            }
        };

        this.onPings = function(data) {
            data.forEach(function(value) {
                allServers[value["Id"]]["Players"] = value["Players"];
            });

            chart.render();
            timeSlider.resetDisabledState();
        };

        this.onViews = function(data) {
            self.emit("newViews", data);
        };

        /**
         * Called when the servers send completely new tracked Servers
         *
         * @param data
         */
        this.onNewServers = function(data) {
            // Reset the Chart
            chart.resetColors();

            // Reset the timeSlider
            timeSlider.resetTime();

            // Reset internal variables
            allServers = {};
            favicons = {};

            // Iterate over all new Servers
            data.forEach(function (value) {
                self.emit("newServer", value);
                allServers[value["Id"]] = value;
            });

            // Take over all servers
            currentServers = allServers;
            self.emit("serversInited");

            // Let the chart render
            chart.render();
            socket.sendPingIDs( currentServers );
            serverTable.render(true);
        };

        /**
         * Called when the user changed the range of time it wants to see
         *
         * @param time
         */
        this.onChangeTime = function(time) {
            socket.sendRange( time );
            socket.sendPingIDs( currentServers );
        };

        this.onHideServer = function(serverId) {
            var newCurrentServers = {};

            for(var key in currentServers) {
                if(currentServers.hasOwnProperty(key) && key != serverId) {
                    newCurrentServers[key] = currentServers[key];
                }
            }

            currentServers = newCurrentServers;
            chart.render();
        };

        this.onShowServer = function(serverID) {
            currentServers[serverID] = allServers[serverID];
            chart.render();
        };

        serverTable.on("hideServer", this.onHideServer);
        serverTable.on("showServer", this.onShowServer);
        timeSlider.on("changeTime", this.onChangeTime);
    }

    heir.inherit(UI, eventEmitter);
    return UI;
});

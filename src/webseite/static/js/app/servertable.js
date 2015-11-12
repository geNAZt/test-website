/**
 * Created by FFass on 10.10.15.
 */
define(["lib/heir", "lib/eventEmitter"], function (heir, eventEmitter) {
    function ServerTable(ui) {
        var self = this,
            animatedTimeout = 0,
            sorted = [],
            skip = 0;

        /** Private methods **/
        function createTH() {
            return "<thead><tr><th></th><th>ID</th><th>Server Name</th><th>Minecraft IP</th><th>Website</th><th>Players</th><th>Record</th><th>Average (24h)</th><th>Uptime (1 month)</th></tr></thead>";
        }

        function createTR(server, renderAnimatedFavicons) {
            var tr = $('<tr />');

            // Check
            var checkTd = $('<td />');

            $('<input />', {
                type: 'checkbox',
                id: 'check_' + server["Id"],
                value: server["Name"],
                checked: ui.getServers().hasOwnProperty(server["Id"])
            }).click(function () {
                if (ui.getServers().hasOwnProperty(server["Id"])) {
                    self.emit("hideServer", server["Id"]);
                } else {
                    self.emit("showServer", server["Id"]);
                }
            }).appendTo(checkTd);
            tr.append(checkTd);

            // ID
            var idTd = $('<td />');
            idTd.text(server["Id"]);
            tr.append(idTd);

            // Name
            var favicon = $('<img />');
            favicon.attr('src', server["Favicon"]);

            if (renderAnimatedFavicons && ui.getSocket() != null) {
                favicons[server["Id"]] = favicon;

                if (animatedTimeout == -1) {
                    animatedTimeout = window.setTimeout(function () {
                        ui.getSocket().sendAnimationRequest(ui.getServers());
                        animatedTimeout = -1;
                    }, 500);
                }
            }

            var nameTd = $('<td />');
            nameTd.append(favicon);
            nameTd.append(server["Name"]);
            tr.append(nameTd);

            // IP
            var ipTd = $('<td />');
            ipTd.text(server["IP"]);
            tr.append(ipTd);

            // Website
            var websiteTd = $('<td />');
            var link = $('<a />');
            link.attr('href', server['Website']);
            link.text(server['Website']);
            websiteTd.append(link);
            tr.append(websiteTd);

            // Players
            var playersTd = $('<td />');
            if (server['Ping24'] !== undefined) {
                var spanDir = $('<span />');
                if (server['Ping24'] > server['Online']) {
                    spanDir.addClass("glyphicon glyphicon-arrow-down");
                    spanDir.attr( "data-original-title", ( server['Ping24'] - server['Online'] ) + " users lost in 24 hours" );
                } else {
                    spanDir.addClass("glyphicon glyphicon-arrow-up");
                    spanDir.attr( "data-original-title", ( server['Onlint'] - server['Ping24'] ) + " users gained in 24 hours" );
                }

                playersTd.append(spanDir);
                spanDir.tooltip();
            }

            playersTd.append(server["Online"] + " / " + server["MaxPlayers"]);
            tr.append(playersTd);

            // Record
            var recordTd = $('<td />');
            recordTd.text(server["Record"] + " Players");
            tr.append(recordTd);

            // Average
            var averageTd = $('<td />');
            averageTd.text(server["Average"] + " Players");
            tr.append(averageTd);

            // Uptime
            var uptimeTd = $('<td />');
            uptimeTd.text(server["Uptime"] + " %");
            tr.append(uptimeTd);

            return tr;
        }

        function sortServers() {
            var serversCopy = ui.getServers();
            var newSorted = [];

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
        }

        /**
         * Render the table which holds all Servers
         *
         * @param renderAnimatedFavicons true if you want to render the favicons new (it will cause the images to reload)
         */
        this.render = function (renderAnimatedFavicons) {
            sortServers();

            var table = $("<table />");
            table.addClass("table table-hover");
            table.append(createTH());

            if (renderAnimatedFavicons) {
                favicons = {};
            }

            var counter = 0, rendered = 0;
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

            var serverTableContainer = $('#server-table');
            serverTableContainer.children("table").remove();
            serverTableContainer.append(table);
        };

        this.onServersInited = function () {
            var size = 0,
                serverCopy = ui.getServers();

            for (var key in serverCopy) {
                if (serverCopy.hasOwnProperty(key)) {
                    size++;
                }
            }

            $('#page-selection').bootpag({
                total: Math.ceil(size / 5)
            }).on("page", function (event, num) {
                skip = 5 * (num - 1);
                self.render(false);
            });
        };

        ui.on("serversInited", this.onServersInited);
    }

    // Make this a eventEmitter
    heir.inherit(ServerTable, eventEmitter);
    return ServerTable;
});
/**
 * Created by Fabian on 21.06.15.
 */
define(["app/socket", "app/ui"], function(Socket, UI) {
    // Create new Socket connection
    var socket = new Socket();
    socket.connectToWebSocket("minecrafttracker.net");

    // Create and bind UI
    var ui = new UI();
    ui.bindSocket( socket );
});
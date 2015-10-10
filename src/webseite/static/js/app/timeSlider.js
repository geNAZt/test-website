/**
 * Created by Fabian on 21.06.15.
 */
define(["lib/heir", "lib/eventEmitter"], function(heir, eventEmitter) {
    function TimeSlider() {
        var currentTime = 2 * 24 * 60,  // Current seconds selected on the slider
            sendTimeout = -1,           // Time to wait until we send the event down the road (to prevent spamming)
            self = this,
            obj = $("#slider");

        obj.slider({
            range: "min",
            value: 2,
            min: 2,
            max: 365,
            slide: function( event, ui ) {
                currentTime = ui.value * 24 * 60;

                if ( sendTimeout == -1 ) {
                    sendTimeout = window.setTimeout(function() {
                        self.emit("changeTime", currentTime);

                        sendTimeout = -1;
                        obj.slider("option", "disabled", true);
                    }, 1000)
                }
            }
        });

        this.resetTime = function() {
            currentTime = 2 * 24 * 60;
            obj.slider("value", 2);
        };

        this.resetDisabledState = function() {
            obj.slider("option", "disabled", false);
        };

        this.getTime = function() {
            return currentTime;
        };
    }

    // Make this object a EventEmitter
    heir.inherit(TimeSlider, eventEmitter);
    return new TimeSlider();
});
/**
 * Created by FFass on 10.10.15.
 */
define(function() {
    return function ViewSelector(ui) {
        this.onNewViews = function(data) {
            var first = null,
                selectData = [];

            for (var key in data) {
                if (data.hasOwnProperty(key)) {
                    if ( first == null ) {
                        first = key;
                    }

                    selectData.push({
                        id: data[key],
                        text: key
                    })
                }
            }

            $("#views").select2({
                data: selectData,
                initSelection: function(element, callback) {
                    callback({id: data[element.val()], text: element.val()});
                }
            }).on('change', function(val) {
                ui.getSocket().setView(val.val);
            });

            $("#views").select2("val", first);

            $(".select2-arrow").html("<i class=\"fa fa-angle-down pull-right\" style=\"margin-top: 0;\"></i>");
        };

        ui.on("newViews", this.onNewViews);
    };
});
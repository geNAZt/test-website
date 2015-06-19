{{ template "layout.tpl" . }}

{{ define "js" }}
    <!-- All functions for this theme + document.ready processing -->
    <script src="{{ "js/jquery.bootpag.min.js" | asset }}"></script>
    <script src="{{ "js/jquery.canvasjs.min.js" | asset }}"></script>
    <script src="{{ "plugins/select2/select2.js" | asset }}"></script>
    <script src="{{ "js/devoops.js" | asset }}"></script>
    <script type="text/javascript">
        host = {{ .host }};
    </script>
{{ end }}

{{ define "css" }}
    <link href="{{ "plugins/select2/select2.css" | asset }}" rel="stylesheet"/>
{{ end }}

{{ define "content" }}
    <div id="inner-content" class="col-xs-10 col-sm-10">
        <br/>

        {{ if isset .flash "registerComplete" }}
        <p id="flash-register" class="txt-success" style="text-align: center;">{{ .flash.registerComplete }}</p>
        <script type="application/javascript">
            setTimeout(function() {
                $('#flash-register').fadeOut();
            }, 5000);

        </script>
        {{ end }}
        <!-- <span class="btn btn-success btn-add" id="edit-button"><a href="#" id="edit-server"><i class="fa fa-pencil"></i><span> Edit Servers</span></a></span>
        <span id="edit-control">
        <span class="btn btn-success btn-add"><a href="#" id="add-server"><i class="fa fa-plus"></i><span> Add Server</span></a></span>
        <span class="placeholder-left"></span>
        <span class="btn btn-danger btn-add"><a href="#" id="remove-server"><i class="fa fa-minus"></i><span> Remove Server</span></a></span>
            </span> -->

        <div id="page-selection"></div>
        <div id="server-table" class="table-responsive">
            <table class="table table-hover">

            </table>
        </div>
        <div id="slider"></div>
        <div id="chartContainer" style="height: 800px; width: 100%;"></div>
    </div>
{{ end }}
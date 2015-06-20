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

        {{ if isset .flash "success" }}
        <p id="flash-success" class="txt-success" style="text-align: center;">{{ .flash.success }}</p>
        <script type="application/javascript">
            setTimeout(function() {
                $('#flash-success').fadeOut();
            }, 5000);
        </script>
        {{ end }}

        {{ if isset .flash "notice" }}
        <p id="flash-notice" class="txt-info" style="text-align: center;">{{ .flash.notice }}</p>
        <script type="application/javascript">
            setTimeout(function() {
                $('#flash-notice').fadeOut();
            }, 5000);
        </script>
        {{ end }}

        {{ if isset .flash "warning" }}
        <p id="flash-warning" class="txt-warning" style="text-align: center;">{{ .flash.warning }}</p>
        <script type="application/javascript">
            setTimeout(function() {
                $('#flash-warning').fadeOut();
            }, 5000);
        </script>
        {{ end }}

        {{ if isset .flash "error" }}
        <p id="flash-error" class="txt-error" style="text-align: center;">{{ .flash.error }}</p>
        <script type="application/javascript">
            setTimeout(function() {
                $('#flash-error').fadeOut();
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
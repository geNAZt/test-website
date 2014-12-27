<!DOCTYPE html>

<html>
<head>
    <title>Test</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">

    <script src="//ajax.googleapis.com/ajax/libs/jquery/2.0.3/jquery.min.js"></script>
    <script type="text/javascript">
        $(function() {
            var conn;

            if (window["WebSocket"]) {
                conn = new WebSocket("ws://{{.host}}/ws");

                conn.onclose = function(evt) {
                    console.log("Connection closed")
                };

                conn.onmessage = function(evt) {
                    console.log(evt.data)
                };

                conn.onopen = function() {
                    conn.send("time:test");
                };
            } else {
                console.log($("<div><b>Your browser does not support WebSockets.</b></div>"))
            }
        });
    </script>
</head>

<body>
{{range $key, $val := .user}}
{{$key}}
{{$val.Name}}
{{end}}
{{.visitCounter}}
</body>
</html>

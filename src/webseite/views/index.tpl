<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>MineTracker</title>
    <meta name="description" content="description">
    <meta name="author" content="geNAZt">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="{{ "plugins/bootstrap/bootstrap.css" | asset }}" rel="stylesheet">
    <link href="{{ "plugins/jquery-ui/jquery-ui.min.css" | asset }}" rel="stylesheet">
    <link href="http://netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.css" rel="stylesheet">
    <link href='http://fonts.googleapis.com/css?family=Righteous' rel='stylesheet' type='text/css'>
    <link href="{{ "css/style.css" | asset }}" rel="stylesheet">
    <!-- HTML5 shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!--[if lt IE 9]>
    <script src="http://getbootstrap.com/docs-assets/js/html5shiv.js"></script>
    <script src="http://getbootstrap.com/docs-assets/js/respond.min.js"></script>
    <![endif]-->
    <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
    <!--<script src="http://code.jquery.com/jquery.js"></script>-->
    <script src="{{ "plugins/jquery/jquery-2.1.0.min.js" | asset }}"></script>
    <script src="{{ "plugins/jquery-ui/jquery-ui.min.js" | asset }}"></script>
    <!-- Include all compiled plugins (below), or include individual files as needed -->
    <script src="{{ "plugins/bootstrap/bootstrap.min.js" | asset }}"></script>
    <!-- All functions for this theme + document.ready processing -->
    <script src="{{ "js/devoops.js" | asset }}"></script>
    <script src="{{ "js/jquery.canvasjs.min.js" | asset }}"></script>
    <script src="{{ "js/jquery.bootpag.min.js" | asset }}"></script>
    <script type="text/javascript">
        host = {{ .host }};
    </script>
</head>
<body>
<!--Start Header-->
<header class="navbar">
    <div class="container-fluid expanded-panel">
        <div class="row">
            <div id="logo" class="col-xs-12 col-sm-12">
                <a href="/">MineTracker</a>
            </div>
        </div>
    </div>
</header>
<!--End Header-->
<!--Start Container-->
<div id="main" class="container-fluid">
    <div class="row">
        <!--Start Content-->
        <div id="content" class="col-xs-12 col-sm-12">
            <div id="inner-content" class="col-xs-10 col-sm-10">
                <br/>

                <h2 class="table-header">Tracked Servers</h2>

                <div id="page-selection"></div>
                <div id="server-table" class="table-responsive">
                    <table class="table table-hover">

                    </table>
                </div>
                <div id="chartContainer" style="height: 800px; width: 100%;"></div>
            </div>
        </div>
        <!--End Content-->
    </div>
</div>
<!--End Container-->
</body>
</html>

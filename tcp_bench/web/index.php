<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Tcp Latency Bench</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="">
    <meta name="author" content="">

    <!-- Le styles -->
    <link href="css/bootstrap.min.css" rel="stylesheet" media="screen">
    <style>
        body {
        padding-top: 60px; /* 60px to make the container go all the way to the bottom of the topbar */
        }
    </style>

    <!-- Fav and touch icons -->
    <link rel="apple-touch-icon-precomposed" sizes="144x144" href="../assets/ico/apple-touch-icon-144-precomposed.png">
    <link rel="apple-touch-icon-precomposed" sizes="114x114" href="../assets/ico/apple-touch-icon-114-precomposed.png">
    <link rel="apple-touch-icon-precomposed" sizes="72x72" href="../assets/ico/apple-touch-icon-72-precomposed.png">
    <link rel="apple-touch-icon-precomposed" href="../assets/ico/apple-touch-icon-57-precomposed.png">
    <link rel="shortcut icon" href="../assets/ico/favicon.png">
</head>

<body>

<div class="navbar navbar-inverse navbar-fixed-top">
    <div class="navbar-inner">
        <div class="container">
            <button type="button" class="btn btn-navbar" data-toggle="collapse" data-target=".nav-collapse">
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <a class="brand" href="#">Latency Bench</a>
            <div class="nav-collapse collapse">
                <ul class="nav">
                    <li class="active"><a href="#">Home</a></li>
                </ul>
            </div><!--/.nav-collapse -->
        </div>
    </div>
</div>

<div class="container">

    <h1>Test runs</h1>

    <?php
        echo '<div class="btn-group" role="group" aria-label="First group">';
        foreach (glob("*.json") as $file) {
          echo '<button id="' . $file . '" type="button" class="btn btn-default result-button">' . $file . '</button>';
        }
        echo '</div>';
    ?>

    <div id="ex0"></div>

</div> <!-- /container -->

<!-- Le javascript
================================================== -->
<!-- Placed at the end of the document so the pages load faster -->
<script src="http://code.jquery.com/jquery.js"></script>
<script src="js/bootstrap.min.js"></script>
<script type="text/javascript" src="https://www.google.com/jsapi"></script>
<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/underscore.js/1.7.0/underscore-min.js"></script>
<script type="text/javascript">

    google.load('visualization', '1', {packages: ['corechart']});
    google.setOnLoadCallback(drawChart);
    var results = [];

    $('.result-button').on('click', function() {
       var id = $(this).attr('id');
       $.getJSON('/' + id).done(function(r) {
         results = r.Results;
         drawChart();
       });
    });

    function drawChart() {
      var data = new google.visualization.DataTable();
      data.addColumn('date', 'X');
      data.addColumn('number', 'Latency');
      var formatter_long = new google.visualization.DateFormat({ pattern: 'hh:mm:ss' });

      data.addRows(_(results).chain().map(function(point) {
        return [new Date(point.Timesent), parseInt(point.Elapsed, 10) || 0];
      }).sortBy(function(r) {
        return r[0];
      }).value());

      formatter_long.format(data, 0);

      var options = {
        width: 1000,
        height: 563,
        hAxis: {
          title: 'Time'
        },
        vAxis: {
          title: 'Latency (ms)'
        }
      };

      var chart = new google.visualization.LineChart(
        document.getElementById('ex0'));

      chart.draw(data, options);

    }

</script>

</body>
</html>

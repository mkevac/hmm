package main

import (
	"log"
	"text/template"
)

var indexHtmlTemplate *template.Template
var mainJsTemplate *template.Template

var indexHtmlTemplateString = `
<!DOCTYPE html>
<html lang="en">
<head>
    <title></title>
    <meta charset="utf-8" />
    <script src="https://code.jquery.com/jquery-3.2.1.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/moment.js/2.18.1/moment.min.js"></script>
    <script src="https://cdn.plot.ly/plotly-latest.min.js"></script>
</head>
<body>
<div id="container1"></div>
<script src="main.js"></script>
</body>
</html>

`

var mainJsTemplateString = `
var chart1;
var headers = [];

$(function () {

    function wsurl() {
        var l = window.location;
        return ((l.protocol === "https:") ? "wss://" : "ws://") + l.hostname + (((l.port != 80) && (l.port != 443)) ? ":" + l.port : "") + "/ws";
    }

    ws = new WebSocket(wsurl());
    ws.onopen = function () {
        ws.onmessage = function (evt) {

            var data = JSON.parse(evt.data);
            console.log(data);
            var splitted = data.split(/\s*[\s,]\s*/);

            if (headers.length == 0) {
                headers = splitted;

                var chartData = [];
                var chartLayout = {
                    autosize: true,
                    yaxis: {
                        tickformat: ".5s"
                    }
                };

                for (i = 1; i < headers.length; i++) {
                    chartData.push({x: [], y: [], type: "scatter", name: headers[i]})
                }

                chart1 = Plotly.newPlot('container1', chartData, chartLayout);

                return;
            }

            //var d = moment(splitted[0]).format('HH:mm:ss');

            var xses = [];
            var yses = [];
            var numbers = [];

            for (i = 1; i < splitted.length; i++) {
                xses.push([splitted[0]]);
                yses.push([parseFloat(splitted[i])]);
                numbers.push(i - 1);
            }

            Plotly.extendTraces('container1', {x: xses, y: yses}, numbers, 86400);
        }
    };
});
`

func init() {
	var err error

	indexHtmlTemplate, err = template.New("").Parse(indexHtmlTemplateString)
	if err != nil {
		log.Fatal("Error while parsing index.html template")
	}

	mainJsTemplate, err = template.New("").Parse(mainJsTemplateString)
	if err != nil {
		log.Fatal("Error while parsing main.js template")
	}
}

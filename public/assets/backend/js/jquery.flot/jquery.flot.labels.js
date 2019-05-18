/*
 * The MIT License

Copyright (c) 2012 by Matt Burland

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
/*
This plugin will draw labels next to the plotted points on your graph. Tested on 
a scatter graph, may or may not work with other graph types. Best suited to 
situations involving a smaller number of points.

usage -
    <style type="text/css">
    .myCSSClass
    {
        font-size: 9px;
        color: #AD8200;
        padding: 2px;
        opacity: 0.80;
    }
    </style>

    <script type="text/javascript">
    var names = [
        "foo",
        "bar"
    ];
    var data = { data: [[1,1],[2,2]], showLabels: true, labels: names, labelPlacement: "left", labelClass: "myCSSClass" };
    $.plot($("#placeholder), [data], options);
    </script>

For each series that you want labeled you need to set showLabels to true, set labels to an array of label names (strings),
set the labelClass to a css class you have defined for your labels and optionally set labelPlacement to one of "left", "right", 
"above" or "below" (below by default if not specified). Placement can be fine tuned by setting the margins in your label class.
Note: if the labelClass is not explicitly supplied in the development version of flot (> v0.7), the plugin will auto generate
a label class as "seriesLabelx" where x is the 1-based index of the data series. I.e. the first dataseries will be seriesLabel1,
the second seriesLabel2, etc.
For the names, the array should be the same length as the data. If any are missing (null) then the label for that point will
be skipped. For example, to label only the 1st and 3rd points:
    
    var names = ["foo", null, "bar"];

Update: Version 0.2

Added support for drawing labels using canvas.fillText. The advantages are that, in theory, drawing to the canvas should be 
faster, but the primary reason is that in some browsers, the labels added as absolutely positioned div elements won't show up
if you print the page. So if you want to print your graphs, you should probably use canvasRender.
The disadvantage is that you lose the flexibility of defining the label with a CSS class.

Options added to series (with defaults):

            canvasRender: false,                // false will add divs to the DOM rather than use canvas.fillText
            cColor: "#000",                     // color for the text if using canvasRender
            cFont: "9px, san-serif",            // font for the text if using canvasRender
            cPadding: 4                         // Padding to add when using canvasRender (where padding is added depends on
                                                // labelPlacement)

Also, version 0.2 takes into account the radius of the data points when placing the labels.
*/

(function ($) {

    function init(plot) {
        plot.hooks.drawSeries.push(drawSeries);
        plot.hooks.shutdown.push(shutdown);
        if (plot.hooks.processOffset) {         // skip if we're using 0.7 - just add the labelClass explicitly.
            plot.hooks.processOffset.push(processOffset);
        }
    }

    function processOffset(plot, offset) {
        // Check to see if each series has a labelClass defined. If not, add a default one.
        // processOptions gets called before the data is loaded, so we can't do this there.
        var series = plot.getData();
        for (var i = 0; i < series.length; i++) {
            if (!series[i].canvasRender && series[i].showLabels && !series[i].labelClass) {
                series[i].labelClass = "seriesLabel" + (i + 1);
            }
        }
    }

    function drawSeries(plot, ctx, series) {
        if (!series.showLabels || !(series.labelClass || series.canvasRender) || !series.labels || series.labels.length == 0) {
            return;
        }
        ctx.save();
        if (series.canvasRender) {
            ctx.fillStyle = series.cColor;
            ctx.font = series.cFont;
        }

        for (i = 0; i < series.data.length; i++) {
            if (series.labels[i]) {
                var loc = plot.pointOffset({ x: series.data[i][0], y: series.data[i][1] });
                var offset = plot.getPlotOffset();
                if (loc.left > 0 && loc.left < plot.width() && loc.top > 0 && loc.top < plot.height())
                    drawLabel(series.labels[i], loc.left, loc.top);
            }
        }
        ctx.restore();

        function drawLabel(contents, x, y) {
            var radius = series.points.radius;
            if (!series.canvasRender) {
                var elem = $('<div class="' + series.labelClass + '">' + contents + '</div>').css({ position: 'absolute' }).appendTo(plot.getPlaceholder());
                switch (series.labelPlacement) {
                    case "above":
                        elem.css({
                            top: y - (elem.height() + radius),
                            left: x - elem.width() / 2
                        });
                        break;
                    case "left":
                        elem.css({
                            top: y - elem.height() / 2,
                            left: x - (elem.width() + radius)
                        });
                        break;
                    case "right":
                        elem.css({
                            top: y - elem.height() / 2,
                            left: x + radius /*+ 15 */
                        });
                        break;
                    default:
                        elem.css({
                            top: y + radius/*+ 10*/,
                            left: x - elem.width() / 2
                        });
                }
            }
            else {
                //TODO: check boundaries
                var tWidth = ctx.measureText(contents).width;
                switch (series.labelPlacement) {
                    case "above":
                        x = x - tWidth / 2;
                        y -= (series.cPadding + radius);
                        ctx.textBaseline = "bottom";
                        break;
                    case "left":
                        x -= tWidth + series.cPadding + radius;
                        ctx.textBaseline = "middle";
                        break;
                    case "right":
                        x += series.cPadding + radius;
                        ctx.textBaseline = "middle";
                        break;
                    default:
                        ctx.textBaseline = "top";
                        y += series.cPadding + radius;
                        x = x - tWidth / 2;

                }
                ctx.fillText(contents, x, y);
            }
        }

    }

    function shutdown(plot, eventHolder) {
        var series = plot.getData();
        for (var i = 0; i < series.length; i++) {
            if (!series[i].canvasRender && series[i].labelClass) {
                $("." + series[i].labelClass).remove();
            }
        }
    }

    // labelPlacement options: below, above, left, right
    var options = {
        series: {
            showLabels: false,
            labels: [],
            labelClass: null,
            labelPlacement: "below",
            canvasRender: false,
            cColor: "#000",
            cFont: "9px, san-serif",
            cPadding: 4
        }
    };

    $.plot.plugins.push({
        init: init,
        options: options,
        name: "seriesLabels",
        version: "0.2"
    });
})(jQuery);
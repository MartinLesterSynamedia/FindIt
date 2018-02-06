// Some javascript/JQuery to make the FindIt work

var circleSize = 40;
var step = 4;
var points = new Array();

function fixCanvas() {
	var img = $("#image");
	// TODO: Use JQuery selector
	var canvas = document.getElementById("draw");
	canvas.width = img.width();
	canvas.height = img.height();
}

function recordClick(event) {
	var num_keys = $("#keys > div").length;
	if (points.length == num_keys) {
		return;
	}

	var ex = event.offsetX;
	var ey = event.offsetY;
	points.push({x:ex, y:ey});
	var out = "Clicked at: ";
	for (i=0; i<points.length; i++) {
		out += "<br/>" + points[i].x + ", " + points[i].y;
	}	
	$("#click").html(out);

	drawTarget(ex, ey);
	
	if (points.length == num_keys) {
		$("#move").html("Verifying the selected points");
	}
}

function debugMousePos(event) {
	$("#move").html( "Mouse at " + event.offsetX + ", " + event.offsetY );
}

function drawTarget(x, y) {
	// TODO: Use JQuery selector
	var ctx = document.getElementById("draw").getContext('2d')
	for (i=circleSize; i>0; i-=step) {
		ctx.beginPath();
		ctx.arc(x, y, i, 0, 2 * Math.PI, false);
		var rgba = "rgba(255,255,0," + (1.01 - (i/circleSize)) + ")";
		ctx.fillStyle = rgba;
		ctx.fill();
	}
}

function revealPoints() {
	// TODO: Show the hidden help text

	$.get( "reveal/", function(data) {
		// TODO: Use JQuery selector
		var ctx = document.getElementById("draw").getContext('2d')
		var meta = $.parseJSON(data);
		meta["Key_rects"].forEach( function(item, index) {
			var px = (item.Max.X - item.Min.X) / 2;
			var py = (item.Max.Y - item.Min.Y) / 2;
			var radius = Math.max(px, py);
			px += item.Min.X;
			py += item.Min.Y;

			ctx.beginPath();
			ctx.arc(px, py, radius, 0, 2 * Math.PI, false);
			ctx.strokeStyle = "red";
			ctx.lineWidth = 3;
			ctx.stroke();
		});
	});	
}
const canvas_black = document.getElementById("canvas-black");
const ctx_black = canvas_black.getContext("2d");
ctx_black.globalCompositeOperation = "source-over";

const canvas_red = document.getElementById("canvas-red");
const ctx_red = canvas_red.getContext("2d");
ctx_red.globalCompositeOperation = "destination-out";

const black_button = document.getElementById("color-black");
const red_button = document.getElementById("color-red");
const white_button = document.getElementById("color-white");
const send_button = document.getElementById("send");

let coordinates = { x: 0, y: 0 };
let is_painting = false;
let color = "black";

document.addEventListener("mousedown", startDrawing);
document.addEventListener("mouseup", stopDrawing);
document.addEventListener("mousemove", draw);

black_button.addEventListener("click", () => {
	color = "black";
	ctx_black.globalCompositeOperation = "source-over";
	ctx_red.globalCompositeOperation = "destination-out";
});
red_button.addEventListener("click", () => {
	color = "red";
	ctx_black.globalCompositeOperation = "destination-out";
	ctx_red.globalCompositeOperation = "source-over";
});
white_button.addEventListener("click", () => {
	color = "white";
	ctx_black.globalCompositeOperation = "destination-out";
	ctx_red.globalCompositeOperation = "destination-out";
});
send_button.addEventListener("click", sendImage);

ctx_black.fillStyle = "white";
ctx_black.rect(0, 0, 264, 176);
ctx_black.fill();

function setPosition(event) {
	let rect = canvas_red.getBoundingClientRect();
	let scaleX = canvas_red.width / rect.width;
	let scaleY = canvas_red.height / rect.height;

	let x = (event.clientX - rect.left) * scaleX;
	let y = (event.clientY - rect.top) * scaleY;

	coordinates.x = x;
	coordinates.y = y;
}

function startDrawing(event) {
	is_painting = true;
	setPosition(event);
}

function stopDrawing() {
	is_painting = false;
}

function draw(event) {
	if (!is_painting) return;

	ctx_black.beginPath();
	ctx_black.lineWidth = 5;
	ctx_black.lineCap = "round";
	ctx_black.strokeStyle = "#000000";
	ctx_black.moveTo(coordinates.x, coordinates.y);

	ctx_red.beginPath();
	ctx_red.lineWidth = 5;
	ctx_red.lineCap = "round";
	ctx_red.strokeStyle = "#ff0000";
	ctx_red.moveTo(coordinates.x, coordinates.y);

	setPosition(event);

	ctx_black.lineTo(coordinates.x, coordinates.y);
	ctx_black.stroke();

	ctx_red.lineTo(coordinates.x, coordinates.y);
	ctx_red.stroke();
}

function sendImage() {
	let data = new FormData();
	data.append("black", canvas_black.toDataURL());
	data.append("red", canvas_red.toDataURL());

	let xhr = new XMLHttpRequest();
	xhr.open("POST", "/process", true);
	xhr.send(data);

	alert("You just sent a drawing to my e-ink screen!");
}

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
const faq_button = document.getElementById("faq");

let canvas_rect = canvas_red.getBoundingClientRect();
const canvas_scale_x = canvas_red.width / canvas_rect.width;
const canvas_scale_y = canvas_red.height / canvas_rect.height;

let coordinates = { x: 0, y: 0 };
let touches = [];
let is_painting = false;
let color = "black";

document.addEventListener("mousedown", (event) => {
	is_painting = true;
	setPosition(event);
});
document.addEventListener("mouseup", () => {
	is_painting = false;
});
document.addEventListener("mousemove", drawCursor);

document.addEventListener("touchstart", (event) => {
	event.preventDefault();
	is_painting = true;
	setTouches(event.touches);
});
document.addEventListener("touchend", (event) => {
	event.preventDefault();
	is_painting = false;
});
document.addEventListener(
	"touchmove",
	(event) => {
		event.preventDefault();
		drawTouches(event);
	},
	false,
);

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
faq_button.addEventListener("click", () => {
	alert(
		"1. Can you add more colors?\nNo the screen itself only supports black red and white so adding more colors is not possible.\n2. What does this do?\nThis drawing canvas allows you to send an image to an e-ink display in my room thats always on.\n3. How long will my image be visible?\nIt will show for a minute then it waits for a new one to come in to replace it.",
	);
});

black_button.addEventListener("touchend", () => {
	color = "black";
	ctx_black.globalCompositeOperation = "source-over";
	ctx_red.globalCompositeOperation = "destination-out";
});
red_button.addEventListener("touchend", () => {
	color = "red";
	ctx_black.globalCompositeOperation = "destination-out";
	ctx_red.globalCompositeOperation = "source-over";
});
white_button.addEventListener("touchend", () => {
	color = "white";
	ctx_black.globalCompositeOperation = "destination-out";
	ctx_red.globalCompositeOperation = "destination-out";
});
send_button.addEventListener("touchend", sendImage);

ctx_black.fillStyle = "white";
ctx_black.rect(0, 0, 264, 176);
ctx_black.fill();

function setPosition(event) {
	canvas_rect = canvas_red.getBoundingClientRect();

	let x = (event.clientX - canvas_rect.left) * canvas_scale_x;
	let y = (event.clientY - canvas_rect.top) * canvas_scale_y;

	coordinates.x = x;
	coordinates.y = y;
}

function setTouches(event_touches) {
	canvas_rect = canvas_red.getBoundingClientRect();
	touches = [];

	for (let i = 0; i < event_touches.length; i++) {
		let touch = event_touches[i];

		let x = (touch.clientX - canvas_rect.left) * canvas_scale_x;
		let y = (touch.clientY - canvas_rect.top) * canvas_scale_y;

		touches.push({ x: x, y: y });
	}
}

function drawCursor(event) {
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

function drawTouches(event) {
	if (!is_painting) return;

	let old_touches = touches;
	setTouches(event.touches);

	for (let i = 0; i < old_touches.length; i++) {
		ctx_black.beginPath();
		ctx_black.lineWidth = 5;
		ctx_black.lineCap = "round";
		ctx_black.strokeStyle = "#000000";
		ctx_black.moveTo(old_touches[i].x, old_touches[i].y);

		ctx_red.beginPath();
		ctx_red.lineWidth = 5;
		ctx_red.lineCap = "round";
		ctx_red.strokeStyle = "#ff0000";
		ctx_red.moveTo(old_touches[i].x, old_touches[i].y);

		ctx_black.lineTo(touches[i].x, touches[i].y);
		ctx_black.stroke();

		ctx_red.lineTo(touches[i].x, touches[i].y);
		ctx_red.stroke();
	}
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


function status(name) {
	var status = document.getElementById("status");
	status.className = name;
	status.innerHTML = name.toUpperCase();

	var favicon = document.getElementById("favicon")
	favicon.href = "/assets/"+name+".png";

	document.getElementById("loading").classList.add("hidden");
}

function loading() {
	document.getElementById("loading").classList.remove("hidden");
}
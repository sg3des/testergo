'use strict';

function status(name) {
	var status = document.getElementById("status");
	status.className = name;
	status.innerHTML = name.toUpperCase();

	var favicon = document.getElementById("favicon");
	favicon.href = "/assets/"+name+".png";

	var loading = document.getElementById("loading");
	loading.classList.add("hidden");
}

function loading() {
	var loading = document.getElementById("loading");
	loading.classList.remove("hidden");
}
const storageUrl = window.location.protocol + "//" + window.location.hostname + ":" + window.location.port + "/storage";
let pathHistory = [];

$(document).ready(function () {
	$("#currentPath").text("Current path: " + "/");
	show(false, "/");
});

function getCurrentPath() {
	return pathHistory[pathHistory.length - 1];
}

function getFile(path, filename) {
	let filePath = storageUrl + "/" + path;
	axios({
		url: filePath,
		method: "GET",
		responseType: 'blob'
	}).then((response) => {
		const url = window.URL.createObjectURL(new Blob([response.data]));
		const link = document.createElement('a');
		link.href = url;
		link.setAttribute('download', filename);
		document.body.append(link);
		link.click();
		document.body.removeChild(link);
	})
}

function goBack() {
	if (pathHistory.length < 2) {
		return;
	}
	lastPath = pathHistory.pop();
	lastPath = pathHistory.pop();
	show(false, lastPath);
}

function show(reload, path) {
	if (!reload) {
		pathHistory.push(path);
	}
	const buildList = (folderData) => {
		let result = "<table>";
		for (let i = 0; i < folderData.length; i++) {
			let folderElement = folderData[i];
			result += "<tr>"
			result += "<td>" + folderElement.name + "</td>"
			if (folderElement.isfolder) {
				result += "<td><button onclick=\"show(false,'" + folderElement.path + "')\">OPEN</button></td>"
			} else {
				result += "<td><button onclick=\"getFile('" + folderElement.path + "', '" + folderElement.name + "')\">DOWNLOAD</button></td>"
			}
			result += "<td><button onclick=\"remove('" + folderElement.path + "')\">REMOVE</button></td>";
			result += "</tr>"
		}
		result += "</table>"
		return result;
	};

	$.ajax({
		url: storageUrl,
		data: {
			operation: "show",
			path: path
		},
		type: "GET",
		success: (result) => {
			$("#currentPath").text("Current path: " + path);
			$("#content").empty();
			$("#content").append(
				buildList(result)
			)
		},
		error: (error) => {
			console.log(error);
		}
	});
}

function createDirectory(folderPath) {
	$.ajax({
		url: storageUrl,
		data: {
			operation: "mkdir",
			operand1: folderPath,
		},
		success: (result) => {

		},
		error: (error) => {
			console.log(error);
		}
	})

}

function remove(path) {
	$.ajax({
		url: storageUrl,
		data: {
			operation: "remove",
			path: path
		},
		success: (result) => {
			show(true, getCurrentPath());
		},
		error: (error) => {
			console.log(error);
		}
	})
}

function newFolder() {
	let folderName = document.getElementById("folderName").value;
	document.getElementById("folderName").value = "";
	if (folderName == "") {
		return;
	}
	let fullPath = getCurrentPath() + "/" + folderName;

	$.ajax({
		url: storageUrl,
		data: {
			operation: "mkdir",
			path: fullPath
		},
		success: (result) => {
			show(true, getCurrentPath());
		},
		error: (error) => {
			console.log(error);
		}
	})

}
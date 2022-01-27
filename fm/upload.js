let dropArea = document.getElementById('drop-area');

['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
	dropArea.addEventListener(eventName, (event) => {
		event.preventDefault();
		event.stopPropagation();
	}, false)
});

['dragenter', 'dragover'].forEach(eventName => {
	dropArea.addEventListener(eventName, (event) => {
		dropArea.classList.add('highlight')
	}, false)
});

['dragleave', 'drop'].forEach(eventName => {
	dropArea.addEventListener(eventName, (event) => {
		dropArea.classList.remove('highlight');
	}, false)
});


dropArea.addEventListener('drop', (event) => {
	let dt = e.dataTransfer;
	let files = dt.files;
	files.forEach(uploadFile);
}, false);

function handleFiles(files) {
	([...files]).forEach(uploadFile)
}

function uploadFile(file) {
	const url = window.location.protocol + "//" + window.location.hostname + ":" + window.location.port + "/storage";
	let formData = new FormData();
	formData.append('operation', 'add');
	formData.append('path', getCurrentPath());
	formData.append('file', file);

	fetch(url, {
		method: 'POST',
		body: formData
	}).then(() => {
		show(true, getCurrentPath());
		console.log("DONE");
	}).catch(() => {
		console.log("ERROR");
	})
}
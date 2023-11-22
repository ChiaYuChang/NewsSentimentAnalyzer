function ShowInfoToast(message, x = 50, y = 10, duration = 3000, destination = "") {
    Toastify({
        text: message,
        className: "info-toast",
        style: {
            "background": "#19A7CE",
            "color": "#F6F1F1",
            "font-weight": "bold",
        },
        offset: {
            x: x, // horizontal axis - can be a number or a string indicating unity. eg: '2em'
            y: y, // vertical axis - can be a number or a string indicating unity. eg: '2em'
        },
        duration: duration,
        destination: destination,
    }).showToast();
}

function ShowAlertToast(message, x = 50, y = 10, duration = 3000, destination = "", onClick) {
    Toastify({
        text: message,
        className: "alert-toast",
        style: {
            "background": "#ca3c3c",
            "color": "#F6F1F1",
            "font-weight": "bold",
        },
        offset: {
            x: x, // horizontal axis - can be a number or a string indicating unity. eg: '2em'
            y: y, // vertical axis - can be a number or a string indicating unity. eg: '2em'
        },
        stopOnFocus: true,
        close: true,
        duration: duration,
        destination: destination,
        onClick: onClick,
    }).showToast();
}
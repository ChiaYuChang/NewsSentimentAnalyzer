function showPassword(prefix) {
    const password = document.getElementById(prefix + "password");
    const hide = document.getElementById(prefix + "show");
    const show = document.getElementById(prefix + "hide");
    if (password.type === "password") {
        password.type = "text";
        show.className = "hide"
        hide.className = "btn"
    } else {
        password.type = "password";
        show.className = "btn"
        hide.className = "hide"
    }
}
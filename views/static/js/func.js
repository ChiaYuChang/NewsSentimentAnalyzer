function showPassword(prefix) {
    const password = document.getElementById(prefix + "password");
    const hide = document.getElementById(prefix + "show");
    const show = document.getElementById(prefix + "hide");
    if (password.type === "password") {
        password.type = "text";
        show.classList.remove("btn")
        show.classList.add("hide")
        hide.classList.remove("hide")
        hide.classList.add("btn")
    } else {
        password.type = "password";
        show.classList.add("btn")
        show.classList.remove("hide")
        hide.classList.add("hide")
        hide.classList.remove("btn")
    }
}
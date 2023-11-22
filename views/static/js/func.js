function showPassword(prefix) {
    if (prefix !== '') {
        prefix += "-";
    }

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

function isValidatePassword(criterion, password) {
    const matches = password.match(criterion.regx);

    return (
        matches &&
        criterion.min <= matches.length &&
        (criterion.max < 0 || matches.length <= criterion.max)
    )
}

function updateValidationStatus(criteria, password) {
    let all_valid = true;
    criteria.forEach((criterion) => {
        let is_valid = isValidatePassword(criterion, password)
        all_valid = all_valid && is_valid;
        if (is_valid) {
            criterion.element.classList.remove("invalid");
            criterion.element.classList.add("valid");
        } else {
            criterion.element.classList.remove("valid");
            criterion.element.classList.add("invalid");
        }
    });
    return all_valid;
}

class Switch {
    constructor(element) {
        this.element = element;
        this.status = false;
        this.on_func = () => { }
        this.off_func = () => { }
    }
    switch(ok, msg) {
        if (ok) {
            this.on(msg);
        } else {
            this.off();
        }
        return ok;
    }
    with_on_func(func) {
        this.on_func = func;
        return this;
    }
    with_off_func(func) {
        this.off_func = func;
        return this;
    }
    on(msg) {
        if (this.status === true && this.element.innerText === msg) {
            console.log("stay on");
            return;
        }
        console.log("turning on");
        this.element.classList.remove("hide");
        this.element.classList.add("alert");
        this.element.innerText = msg;
        this.status = true;
        this.on_func();
    }
    off() {
        if (this.status === false) {
            console.log("stay off");
            return;
        }
        console.log("turning off");
        this.element.classList.remove("alert");
        this.element.classList.add("hide");
        this.status = false;
        this.off_func();
    }
}


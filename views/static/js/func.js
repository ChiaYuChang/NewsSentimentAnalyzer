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

// var new_pwd_alert_switch = undefined;
// var old_pwd_alert_switch = undefined;

// function changPasswordFormSubmit() {
//     const old_pwd_el = document.querySelector('input[name="old-password"]')
//     const new_pwd_el = document.querySelector('input[name="new-password"]')

//     if (new_pwd_alert_switch === undefined) {
//         new_pwd_alert_switch = new Switch(document.getElementById("new-alert"))
//         new_pwd_alert_switch.with_on_func(() => {
//             new_pwd_el.classList.add("highlight-alert")
//         })
//         new_pwd_alert_switch.with_off_func(() => {
//             new_pwd_el.classList.remove("highlight-alert")
//         })
//     }

//     if (old_pwd_alert_switch === undefined) {
//         old_pwd_alert_switch = new Switch(document.getElementById("old-alert"))
//         old_pwd_alert_switch.with_on_func(() => {
//             old_pwd_el.classList.add("highlight-alert")
//         })
//         old_pwd_alert_switch.with_off_func(() => {
//             old_pwd_el.classList.remove("highlight-alert")
//         })
//     }

//     obj = {
//         "old": old_pwd_el.value,
//         "new": new_pwd_el.value,
//     }

//     let required_field_missing = false
//     required_field_missing |= old_pwd_alert_switch.switch(
//         obj["old"] === "",
//         "Your current password is missing")
//     required_field_missing |= new_pwd_alert_switch.switch(
//         obj["new"] === "",
//         "Your new password cannot be blank")
//     if (required_field_missing) {
//         return
//     }

//     if (obj["new"] === obj["old"]) {
//         new_pwd_alert_switch.on("Password not changed")
//         return
//     }
//     new_pwd_alert_switch.off();

//     if (password_is_valid !== true) {
//         new_pwd_alert_switch.on("Your new password is invalid");
//         return
//     }
//     new_pwd_alert_switch.off();

//     fetch('/v1/change-password', {
//         method: 'PATCH',
//         headers: {
//             "Content-Type": "application/json",
//             "Accept": "application/json"
//         },
//         credentials: "same-origin",
//         body: JSON.stringify(obj),
//     }).then(resp => {
//         resp.json().then(data => {
//             console.log(data)
//             if (data.status === 200) {
//                 console.log("success")
//             } else {
//                 old_pwd_alert_switch.switch(
//                     data.password_not_matched,
//                     "Password did not match")
//                 new_pwd_alert_switch.switch(data.password_not_changed,
//                     "Password did not change")
//             }
//         });
//     }).catch(err => {
//         console.log(err)
//     });
// }
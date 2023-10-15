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
    let allValid = true;
    criteria.forEach((criterion) => {
        let isValid = isValidatePassword(criterion, password)
        allValid = allValid && isValid;
        if (isValid) {
            criterion.element.classList.remove("invalid");
            criterion.element.classList.add("valid");
        } else {
            criterion.element.classList.remove("valid");
            criterion.element.classList.add("invalid");
        }
    });
    return allValid;
}


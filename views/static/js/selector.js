function insertSelectElement(where, index, name, classes, opts) {
    let newSelect = document.createElement("select");
    newSelect.name = name
    newSelect.setAttribute("aria-label", where.id)
    newSelect.classList = classes

    opts.forEach(opt => {
        if (index > 0 && opt.value === "") {
            return
        }
        newOpt = document.createElement("option")
        newOpt.value = opt.value;
        newOpt.textContent = opt.txt;
        newSelect.appendChild(newOpt);
    });

    where.appendChild(newSelect);
    return newSelect;
}

function deleteSelectElement(where) {
    let lastChild = where.lastElementChild;
    if (lastChild) {
        where.removeChild(lastChild);
    }
}

function addListenerToBtn(iPosId, iBtnId, dBtnId, maxDiv, opts, alertMsg) {
    const position = document.getElementById(iPosId);
    if (!!position) {
        const insertButton = document.getElementById(iBtnId);
        const deleteButton = document.getElementById(dBtnId);
        const divLimit = maxDiv;
        let counter = 0;

        let fstEle = insertSelectElement(position, counter, `${iPosId}[${counter}]`, "form-input", opts);
        counter++
        deleteButton.classList.add("pure-button-disabled")

        insertButton.addEventListener("click", () => {
            if (fstEle.value === "") {
                alert(`Please exclude 'All' when selecting multiple ${iPosId}.`);
                return
            }

            if (counter < divLimit) {
                insertSelectElement(position, counter, `${iPosId}[${counter}]`, "form-input", opts);
                counter++;
                fstEle.classList.add("pure-button-disabled");
                deleteButton.classList.remove("pure-button-disabled");
            } else {
                alert(alertMsg);
            }

            if (counter >= divLimit) {
                insertButton.classList.add("pure-button-disabled")
            }
        });

        deleteButton.addEventListener("click", () => {
            if (counter > 1) {
                deleteSelectElement(position)
                counter--;
                insertButton.classList.remove("pure-button-disabled");
            }

            if (counter <= 1) {
                fstEle.classList.remove("pure-button-disabled");
                deleteButton.classList.add("pure-button-disabled");
            }
        });
    }
}

function getTimeZone() {
    const fTimeTZ = document.getElementById("from-time-tz");
    const tTimeTZ = document.getElementById("to-time-tz");
    const hiddenTZFormField = document.getElementById("timezone");
    const tz = Intl.DateTimeFormat().resolvedOptions().timeZone;

    if (!!hiddenTZFormField) {
        hiddenTZFormField.value = tz;
    }
    if (!!fTimeTZ) {
        fTimeTZ.innerText = tz;
    }
    if (!!tTimeTZ) {
        tTimeTZ.innerText = tz;
    }
}
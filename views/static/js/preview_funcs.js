var pcid = location.pathname.split('/')[3];

var list;
var masterCheckbox;
var itemCheckboxes = [];

document.addEventListener("DOMContentLoaded", (event) => {
    masterCheckbox = document.getElementById('select-all');

    initList();
    // Add an event listener to the master checkbox
    masterCheckbox.addEventListener('change', () => {
        // Set the state of all other checkboxes to match the master checkbox
        itemCheckboxes.forEach((obj) => {
            obj.checkbox.checked = masterCheckbox.checked;
        });
    });

    getPreviewItems(pcid);
})

function initList() {
    function item_func(values) {
        return `<tr>
                    <td><input type='checkbox' class='pure-checkbox' id='${values["id"]}'></td>
                    <td>
                        <a href=${values["link"]}><h5>${values["title"]}</h5></a>
                        <p>${values["description"]}</p>
                    </td>
                    <td>${values["publication_date"]}</td>
                </tr>`;
    }

    let options = {
        valueNames: ["title", "description", "category", "pubDate"],
        item: item_func,
    };
    list = new List("item-table", options);
}

async function getPreviewItems(pcid) {
    const response = await fetch(`/v1/preview/fetch-next-page/${pcid}`, {
        method: 'GET',
    });

    response.json().then(data => {
        console.log(data);
        if (data["has_next"] === false) {
            let el = document.getElementById("more");
            el.classList.add("pure-button-disabled")
        }

        if ("error" in data) {
            console.log(data["error"]);
            if (data["error"]["url"] !== "") {
                location.href = data["error"]["url"];
            }
            return
        }

        if (!("items" in data)) {
            ShowAlertToast(
                "Nothing found"
            )
            window.history.go(-1);
            return
        }
        data["items"].forEach(element => {
            tmp = new Date(element["publication_date"]);
            element["publication_date"] = tmp.toLocaleString();
            list.add(element)
        });

        data["items"].forEach(element => {
            cb = document.getElementById(element["id"]);
            itemCheckboxes.push({ id: element["id"], checkbox: cb });
            cb.addEventListener('change', function () {
                // If any checkbox is unchecked, uncheck the master checkbox
                masterCheckbox.checked = itemCheckboxes.every(function (obj) {
                    return obj.checkbox.checked;
                });
            });
        })
        masterCheckbox.checked = false;
    }).catch(err => {
        console.log("Error", err);
    })

}

async function submit(pcid) {
    const fdata = new URLSearchParams();
    if (masterCheckbox.checked) {
        fdata.append("select_all", true);
    } else {
        fdata.append("select_all", false);
        let selectedDataId = [];
        itemCheckboxes.forEach((obj) => {
            if (obj.checkbox.checked === true) {
                selectedDataId.push(obj.id);
            }
        })

        if (selectedDataId.length === 0) {
            ShowAlertToast(
                message = "Please select at least 1 item",
            )
            return
        }
        selectedDataId.forEach((dataId, i) => {
            console.log(`add item[${i}] ${dataId}`);
            fdata.append(`item[${i}]`, dataId);
        });
    }
    const response = await fetch(`/v1/preview/${pcid}`, {
        method: 'POST',
        body: fdata
    });

    response.json().then((data) => {
        ShowInfoToast(message = "Done", destination = data["url"])
        setTimeout(() => {
            window.location.href = data["url"];
        }, 3000);
    }).catch(err => {
        ShowAlertToast(
            message = `Code ${err.code}: ${err.message}`
        )
        console.log("Error", err);
    })
}

function goToPreviousPage() {
    window.history.go(-1);
}

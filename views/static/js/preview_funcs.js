var pcid = location.pathname.split('/')[3];

var list
var data = [];

function init_list() {
    function item_func(values) {
        return `<tr>
                    <td><input type='checkbox' class='pure-checkbox' name='item[${values["id"]}]' value='${values["id"]}'></td>
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
        } else {
            data["items"].forEach(element => {
                tmp = new Date(element["publication_date"]);
                element["publication_date"] = tmp.toLocaleString();
                list.add(element)
            });
        }
    }).catch(err => {
        console.log("Error", err);
    })
}

document.addEventListener("DOMContentLoaded", (event) => {
    init_list();
    getPreviewItems(pcid);
})
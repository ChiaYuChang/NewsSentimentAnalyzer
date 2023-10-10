var pagerCache = new Map();
var detailCache = new Map();

addEventListener("DOMContentLoaded", (event) => {
    showLoadingAnimation()
    getJobs();
});

function showLoadingAnimation(colspan = 5) {
    let el = document.getElementById("query-result-tbody")
    el.replaceChildren();

    let tr = document.createElement("tr")
    tr.setAttribute("id", "loading-animation")

    let td = document.createElement("td")
    td.setAttribute("colspan", colspan)

    let container = document.createElement("div")
    container.setAttribute("class", "animation-container")

    let outer = document.createElement("div")
    outer.setAttribute("class", "base-box rotation-box-outer")

    let inner = document.createElement("div")
    inner.setAttribute("class", "base-box rotation-box-inner")

    let midbox = document.createElement("div")
    midbox.setAttribute("class", "base-box fixed-box")

    let text = document.createElement("span")
    text.setAttribute("class", "loading-text")
    text.textContent = "Loading..."

    midbox.appendChild(text)
    es = [outer, inner, midbox];
    es.forEach(e => {
        container.appendChild(e)
    })

    td.appendChild(container)
    tr.appendChild(td)
    el.appendChild(tr)
}

function showNoJob(jstatus, colspan = 5) {
    let el = document.getElementById("query-result-tbody")
    pager["jstatus"] = jstatus;

    let tr = document.createElement("tr")
    tr.setAttribute("id", "loading-animation")

    let td = document.createElement("td")
    td.setAttribute("colspan", colspan)
    td.setAttribute("align", "center")

    let h3 = document.createElement("h3")
    h3.textContent = "No Jobs"

    td.appendChild(h3)
    tr.appendChild(td)
    el.appendChild(tr)
}

function updatePageButton() {
    let pbtn = document.getElementById("prev-page-q")
    if (pager.page === 0) {
        pbtn.classList.add("pure-button-disabled")
        pbtn.setAttribute("onclick", "#")
    } else {
        pbtn.classList.remove("pure-button-disabled")
        pbtn.setAttribute("onclick", "queryPage(qPrev)")
    }

    let nbtn = document.getElementById("next-page-q")
    if (pager.page === last_page_of_each_jstatus[pager.jstatus] - 1) {
        nbtn.classList.add("pure-button-disabled")
        nbtn.setAttribute("onclick", "#")
    } else {
        nbtn.classList.remove("pure-button-disabled")
        nbtn.setAttribute("onclick", "queryPage(qNext)")
    }
}

function updateQuery(jstatus, n) {
    if (pager.jstatus === jstatus) {
        // debug only
        // console.log(`current status "${jstatus}" not change`)
        return
    }

    let el = document.getElementById("query-result-tbody");
    el.replaceChildren();

    pager.jstatus = jstatus
    if (n === 0) {
        showNoJob(jstatus, 5)
        return
    }
    pager.fjid = maxInt32
    pager.tjid = maxInt32
    pager.page = 0

    updatePageButton();
    showLoadingAnimation();
    getJobs();
    return
}

function queryPage(direction) {
    if (direction) {
        pager.page += 1
        pager.tjid = pager.fjid
    } else {
        pager.page -= 1
        pager.tjid = pager.fjid
    }

    pager.direction = direction
    updatePageButton()
    showLoadingAnimation();
    getJobs();
}

function getKey(jstatus, page) {
    return `${jstatus}-${page}`;
}

var jobList
var item_func = function (values) {
    return `
        <tr class='job-data' onclick=${values["onclick"]}>",
        <th class='job-id mono'>${values["job-id"]}</th>",
        <td><div class='job-status' style='width:90%;text-align:center' status=${values["job-status"]}>${values["job-status"]}</div></td>",
        <td class='job-news_src'>${values["job-news_src"]}</td>",
        <td class='job-analyzer'>${values["job-analyzer"]}</td>",
        <td class='job-updated_at mono'>${values["job-updated_at"]}</td>",
        </tr>`
}

function newList(data) {
    for (let i = 0; i < data.length; i++) {
        data[i]["onclick"] = `getJobDetails(${data[i]["job-id"]})`
    }

    var options = {
        valueNames: [
            "job-id", "job-status", "job-news_src", "job-analyzer", "job-created_at", "job-updated_at",
            { attr: "onclick", name: "job-details" },
        ], item: item_func,
    };

    let el = document.getElementById("query-result-tbody");
    el.replaceChildren()

    jobList = new List('qurey-result', options, data);
    jobList.sort('job-id', { order: "desc" })
}

async function getJobs() {
    const fdata = new URLSearchParams();

    let ckey = getKey(pager.jstatus, pager.page)
    if (pagerCache.has(ckey)) {
        // debug only
        // console.log("from cache")

        let data = pagerCache.get(ckey)["data"];
        pager.fjid = data[data.length - 1]["job-id"];
        pager.tjid = data[0]["job-id"];
        newList(data)
    } else {
        // debug only
        // console.log("from query")

        for (const key in pager) {
            fdata.append(key, pager[key]);
        }

        const response = await fetch("/v1/job", {
            method: "POST",
            body: fdata,
        });

        response.json()
            .then(data => {
                if (data.length === 0) {
                    return
                }
                pager.fjid = data[data.length - 1]["job-id"];
                pager.tjid = data[0]["job-id"];
                pagerCache.set(ckey, {
                    "from": pager.fjid,
                    "to": pager.tjid,
                    "data": data,
                });
                newList(data)
            })
            .catch(err => { console.error("Error:", err) });
    }
}

async function getJobDetails(id) {
    let detailEl = document.getElementById("detail");
    if (detailEl.getAttribute("job-id") === ('' + id)) { return }

    var data
    if (detailCache.has(id)) {
        data = detailCache.get(id)

        // debug only
        // console.log("read from cache")
    } else {
        const response = await fetch(`/v1/job/${id}`);
        if (response.status != 200) {
            detailEl.setAttribute("hidden", "");
            detailEl.removeAttribute("job-id");
            return
        }
        data = await response.json();
        detailCache.set(id, data)

        // debug only
        // console.log("read from query")
    }
    const dtbodyEl = document.getElementById("detail-table-body")
    dtbodyEl.replaceChildren()
    detailEl.removeAttribute("hidden");
    detailEl.setAttribute("job-id", data["job-id"]);

    tr = document.createElement("tr")
    var detailsfields = [
        {
            "row_header": "Job ID",
            "field_name": "job-id",
            "is_mono": false,
        },
        {
            "row_header": "Owner",
            "field_name": "job-owner",
            "is_mono": false,
        },
        {
            "row_header": "Status",
            "field_name": "job-status",
            "is_mono": false,
        },
        {
            "row_header": "News API",
            "field_name": "job-news_api",
            "is_mono": false,
        },
        {
            "row_header": "News API Query",
            "field_name": "job-news_api_query",
            "is_mono": true,
        },
        {
            "row_header": "Analyzer",
            "field_name": "job-analyzer",
            "is_mono": false,
        },
        {
            "row_header": "Analyzer Query",
            "field_name": "job-analyzer_query",
            "is_mono": true,
        },
        {
            "row_header": "Created At",
            "field_name": "job-created_at",
            "is_mono": true,
        },
        {
            "row_header": "Updated At",
            "field_name": "job-updated_at",
            "is_mono": true,
        },
    ]

    detailsfields.forEach((f) => {
        let tr = document.createElement("tr")
        let th = document.createElement("th")
        th.textContent = f.row_header
        th.setAttribute("scope", "row")
        th.setAttribute("style", "width:20%;min-width:8.5rem;text-align:left")

        let td = document.createElement("td")
        if (f.field_name === "job-analyzer_query") {
            let pre = document.createElement("pre")
            let jsn = JSON.stringify(data[f.field_name], null, "\t");
            pre.textContent = jsn
            td.classList.add("mono")
            td.appendChild(pre)
        } else if (f.field_name === "job-status") {
            let div = document.createElement("div")
            div.setAttribute("class", "job-status")
            div.setAttribute("status", data[f.field_name])
            div.textContent = data[f.field_name];
            div.setAttribute("style", "width:6rem")
            td.appendChild(div)
        } else {
            td.textContent = data[f.field_name]
            if (f.is_mono) {
                td.classList.add("mono")
            }
        }

        tr.appendChild(th)
        tr.appendChild(td)
        dtbodyEl.appendChild(tr)
    })
}

var pagerCache = new Map();
var detailCache = new Map();

const urlParams = new URLSearchParams(window.location.search);
var selectedJobId = parseInt(urlParams.get('jid'));

addEventListener("DOMContentLoaded", (event) => {
    showLoadingAnimation();
    updatePageButton();
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
    // status not change
    if (pager.jstatus === jstatus) {
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
    if (values["job-id"] == selectedJobId) {
        return `
        <tr class='job-data highlight-alert' id='job-${values["job-id"]}' onclick=${values["onclick"]}>",
        <th class='job-id mono'>${values["job-id"]}</th>",
        <td><div class='job-status' status=${values["job-status"]}>${values["job-status"]}</div></td>",
        <td class='job-news_src'>${values["job-news_src"]}</td>",
        <td class='job-analyzer'>${values["job-analyzer"]}</td>",
        <td class='job-updated_at mono'>${values["job-updated_at"]}</td>",
        </tr>`
    }
    return `
        <tr class='job-data' id='job-${values["job-id"]}' onclick=${values["onclick"]}>",
        <th class='job-id mono'>${values["job-id"]}</th>",
        <td><div class='job-status' status=${values["job-status"]}>${values["job-status"]}</div></td>",
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
            "job-id", "job-status", "job-news_src", "job-analyzer",
            "job-created_at", "job-updated_at",
            { attr: "onclick", name: "job-details" },
        ], item: item_func,
    };

    let el = document.getElementById("query-result-tbody");
    el.replaceChildren()

    jobList = new List('qurey-result', options, data);
    // jobList.sort('job-id', { order: "desc" })
}

async function getJobs() {
    const fdata = new URLSearchParams();

    let ckey = getKey(pager.jstatus, pager.page)
    if (pagerCache.has(ckey)) {
        // use cache
        let data = pagerCache.get(ckey)["data"];
        pager.fjid = data[data.length - 1]["job-id"];
        pager.tjid = data[0]["job-id"];
        newList(data)
    } else {
        // fetch new data
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
            .catch(err => { 
                console.error("Error:", err)
            });
    }
}

var detailsFields = [
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

async function getJobDetails(id) {
    if (selectedJobId !== null) {
        let el = document.getElementById(`job-${selectedJobId}`);
        if (el !== null ) {
            el.classList.remove("highlight-alert");
        }        
    }
    selectedJobId= id;
    let el = document.getElementById(`job-${selectedJobId}`);
    if (el !== null ) {
        el.classList.add("highlight-alert");
    }     

    const detailEl = document.getElementById("detail");
    if (detailEl.getAttribute("job-id") === ('' + id)) { return }

    var data
    if (detailCache.has(id)) {
        data = detailCache.get(id)
    } else {
        const response = await fetch(`/v1/job/${id}`);
        if (response.status != 200) {
            detailEl.setAttribute("hidden", "");
            detailEl.removeAttribute("job-id");
            return
        }
        data = await response.json();
        detailCache.set(id, data)
    }
    const dtbodyEl = document.getElementById("detail-table-body")
    dtbodyEl.replaceChildren()
    detailEl.removeAttribute("hidden");
    detailEl.setAttribute("job-id", data["job-id"]);

    tr = document.createElement("tr")
    detailsFields.forEach((f) => {
        let tr = document.createElement("tr")
        let th = document.createElement("th")
        th.textContent = f.row_header
        th.setAttribute("scope", "row")

        let td = document.createElement("td")
        switch (f.field_name) {
            case "job-analyzer_query":
                let pre = document.createElement("pre");
                let jsn = JSON.parse(data["job-analyzer_query"]);
                pre.textContent = JSON.stringify(jsn, null, 2);
                td.classList.add("mono");
                td.appendChild(pre);
                break;
            case "job-status":
                let div = document.createElement("div")
                div.setAttribute("class", "job-status")
                div.setAttribute("status", data["job-status"])
                div.textContent = data["job-status"];
                td.appendChild(div)
                break;
            case "job-news_api_query":
                td.textContent = decodeURIComponent(data["job-news_api_query"]);
                td.classList.add("mono")
                break;
            default:
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

var pagerCache = new Map();
var detailCache = new Map();

const jStatusAll = "all";
const jStatusCreated = "created";
const jStatusRunning = "running";
const jStatusDone = "done";
const jStatusFailure = "failure";
const jStatusCanceled = "canceled";

const maxInt32 = 2147483647

const qNext = true
const qPrev = false

var pager = {
    jstatus: jStatusAll,
    fjid: maxInt32,
    tjid: maxInt32,
    page: 0,
    direction: qNext,
}

var lastP = {
    all: 0,
    created: 0,
    running: 0,
    done: 0,
    failure: 0,
    canceled: 0,
}

addEventListener("DOMContentLoaded", (event) => {
    showLoadingAnimation()
    console.log(pageSize)
    console.log(nStatus)

    var jss = [jStatusAll, jStatusCreated, jStatusRunning, jStatusDone, jStatusFailure, jStatusCanceled]
    jss.forEach((js, index) => {
        lastP[js] = Math.ceil(nStatus[index] / pageSize)
    })
    console.log(lastP)
    getJobs();
});

function showLoadingAnimation() {
    let el = document.getElementById("query-result-tbody")
    el.replaceChildren();

    let tr = document.createElement("tr")
    tr.setAttribute("id", "loading-animation")

    let td = document.createElement("td")
    td.setAttribute("colspan", "6")

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

function showNoJob(jstatus) {
    let el = document.getElementById("query-result-tbody")
    pager["jstatus"] = jstatus;

    let tr = document.createElement("tr")
    tr.setAttribute("id", "loading-animation")

    let td = document.createElement("td")
    td.setAttribute("colspan", "6")
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
    if (pager.page === lastP[pager.jstatus] - 1) {
        nbtn.classList.add("pure-button-disabled")
        nbtn.setAttribute("onclick", "#")
    } else {
        nbtn.classList.remove("pure-button-disabled")
        nbtn.setAttribute("onclick", "queryPage(qNext)")
    }
}

function updateQuery(jstatus) {
    if (pager.jstatus === jstatus) {
        console.log(`current status "${jstatus}" not change`)
        return
    }

    let el = document.getElementById("query-result-tbody");
    el.replaceChildren();

    pager.jstatus = jstatus
    pager.fjid = maxInt32
    pager.tjid = maxInt32
    pager.page = 0

    updatePageButton()
    showLoadingAnimation();
    getJobs();
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

async function getJobs() {
    const fdata = new URLSearchParams();

    let ckey = getKey(pager.jstatus, pager.page)
    if (pagerCache.has(ckey)) {
        console.log("from cache")
        let data = pagerCache.get(ckey)["data"];
        pager.fjid = data[data.length - 1].id;
        pager.tjid = data[0].id;
        appendRowsToTable(data);
    } else {
        console.log("from query")
        for (const key in pager) {
            fdata.append(key, pager[key]);
        }

        const response = await fetch("/v1/job", {
            method: "POST",
            body: fdata,
        });

        response.json()
            .then(data => {
                pager.fjid = data[data.length - 1].id;
                pager.tjid = data[0].id;
                pagerCache.set(ckey, {
                    "from": pager.fjid,
                    "to": pager.tjid,
                    "data": data,
                });
                appendRowsToTable(data);
            })
            .catch(err => { console.error("Error:", err) });
    }
}

var fields = ["id", "status", "news_src", "analyzer", "created_at", "updated_at"]

function appendRowsToTable(data) {
    let el = document.getElementById("query-result-tbody");
    el.replaceChildren()

    data.forEach(job => {
        tr = document.createElement("tr")
        tr.setAttribute("onclick", `getJobDetails(${job["id"]})`)
        tr.setAttribute("class", "job-row")

        fields.forEach(f => {
            td = document.createElement("td")
            td.classList.add("job-" + f)
            if (f == "status") {
                div = document.createElement("div")
                div.setAttribute("class", job[f]["class"])
                div.setAttribute("style", "width:6rem;text-align:center")
                div.textContent = job[f]["text"]
                td.appendChild(div)
            } else {
                td.textContent = job[f]
            }
            tr.appendChild(td)
        });
        tr.appendChild(td)
        el.appendChild(tr)
    });
}

async function getJobDetails(id) {
    let detailEl = document.getElementById("detail");
    if (detailEl.getAttribute("jid") === ('' + id)) { return }

    var data
    if (detailCache.has(id)) {
        data = detailCache.get(id)
        console.log("read from cache")
    } else {
        const response = await fetch(`/v1/job/${id}`);
        if (response.status != 200) {
            detailEl.setAttribute("hidden", "");
            detailEl.removeAttribute("jid");
            return
        }
        data = await response.json();
        detailCache.set(id, data)
        console.log("read from query")
    }
    const dtbodyEl = document.getElementById("detail-table-body")
    dtbodyEl.replaceChildren()
    detailEl.removeAttribute("hidden");
    detailEl.setAttribute("jid", data.jid);

    let jsn = JSON.stringify(data.analyzer_query, null, "\t");
    tr = document.createElement("tr")

    let fields = ["jid", "owner", "status", "news_api", "news_api_query", "analyzer", "analyzer_query", "created_at", "updated_at"]
    display_fields = ["Job ID", "Owner", "Status", "News API", "News API Query", "Analyzer", "Analyzer Query", "Created At", "Updated At"]
    fields.forEach((field, index) => {
        let tr = document.createElement("tr")
        let th = document.createElement("th")
        th.textContent = display_fields[index]
        th.setAttribute("scope", "row")
        th.setAttribute("style", "width:7rem;text-align:left")

        let td = document.createElement("td")
        if (field === "analyzer_query") {
            let pre = document.createElement("pre")
            pre.textContent = jsn
            td.appendChild(pre)
        } else if (field === "status") {
            let div = document.createElement("div")
            div.setAttribute("class", `job - status ${data[field].class}`)
            div.textContent = data[field].text;
            div.setAttribute("style", "width:6rem;text-align:center")
            td.appendChild(div)
        } else {
            td.textContent = data[field]
        }

        tr.appendChild(th)
        tr.appendChild(td)
        dtbodyEl.appendChild(tr)
    })
}

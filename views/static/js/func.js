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

async function getJob(id) {
    const detailEl = document.getElementById("detail");
    if (detailEl.getAttribute("jid") === ('' + id)) { return }
    const response = await fetch(`/v1/job/${id}`);
    const dtbodyEl = document.getElementById("detail-table-body")

    dtbodyEl.replaceChildren()
    if (response.status != 200) {
        detailEl.setAttribute("hidden", "");
        detailEl.removeAttribute("jid");
        return
    }
    const data = await response.json();
    detailEl.removeAttribute("hidden");
    detailEl.setAttribute("jid", data["JobID"]);

    let jsn = JSON.stringify(data, null, "\t");
    tr = document.createElement("tr")
    let fields = ["JobID", "Owner", "Status", "NewsAPI", "NewsAPIQuery", "Analyzer", "AnalyzerQuery", "CreatedAt", "UpdatedAt"]
    fields.forEach(field => {
        let tr = document.createElement("tr")

        let th = document.createElement("th")
        th.textContent = field
        th.setAttribute("scope", "row")
        th.setAttribute("style", "width:5rem;text-align:left")

        let td = document.createElement("td")
        if (field === "AnalyzerQuery") {
            let pre = document.createElement("pre")
            pre.textContent = jsn
            td.appendChild(pre)
        } else if (field === "Status") {
            let div = document.createElement("div")
            div.classList.add(data[field]["Class"])
            div.textContent = data[field]["Text"]
            div.setAttribute("style", "width:5rem;text-align:center")
            td.appendChild(div)
        } else {
            td.textContent = data[field]
        }

        tr.appendChild(th)
        tr.appendChild(td)
        dtbodyEl.appendChild(tr)
    })
}
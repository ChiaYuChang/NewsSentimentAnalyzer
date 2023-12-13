const pcid = location.pathname.split('/')[3];

const urlParams = new URLSearchParams(window.location.search);
const aid = urlParams.get('aid')
const eid = urlParams.get('eid')

var api_name_to_id = {
    "openai": 5,
    "cohere": 6
}

var llm_api = "openai";
var has_sent = false;
const useOpenAI = document.getElementById('openai');
const useCohere = document.getElementById('cohere');
const doEmbedding = document.getElementById('do-embedding');
const doSentiment = document.getElementById('do-sentiment');
const openaiOpt = document.getElementById('openai-option');
const cohereOpt = document.getElementById('cohere-option');
const embeddingOpt = document.querySelectorAll('div.option[type="embedding"]');
const sentimentOpt = document.querySelectorAll('div.option[type="sentiment"]');
const submitEl = document.querySelector('input[type="submit"]');

document.addEventListener('DOMContentLoaded', () => {
    useCohere.classList.add("btn-inactive");
    cohereOpt.style.display = 'none';
    
    useOpenAI.addEventListener('click', () => {
        openaiOpt.style.display = 'block';
        useOpenAI.classList.remove("btn-inactive");
        cohereOpt.style.display = 'none';
        useCohere.classList.add("btn-inactive");
        llm_api = "openai";
    });

    useCohere.addEventListener('click', () => {
        openaiOpt.style.display = 'none';
        useOpenAI.classList.add("btn-inactive");
        cohereOpt.style.display = 'block';
        useCohere.classList.remove("btn-inactive");
        llm_api = "cohere";
    });

    doEmbedding.addEventListener('change', () => {
        if (doEmbedding.checked) {
            embeddingOpt.forEach(el => {
                el.style.display = 'block';
            });
            submitEl.classList.remove("pure-button-disabled")
        } else {
            embeddingOpt.forEach(el => {
                el.style.display = 'none';
            });
            if (!doSentiment.checked) {
                submitEl.classList.add("pure-button-disabled")
            }
        }
    });

    doSentiment.addEventListener('change', () => {
        if (doSentiment.checked) {
            sentimentOpt.forEach(el => {
                el.style.display = 'block';
            });
            submitEl.classList.remove("pure-button-disabled")
        } else {
            sentimentOpt.forEach(el => {
                el.style.display = 'none';
            });
            if (!doEmbedding.checked) {
                submitEl.classList.add("pure-button-disabled")
            }
        }
    });
});

async function submitForm() {
    if (has_sent) {
        ShowAlertToast(
            message = "You have already submitted this request",
            x = 50, y = 10, duration = 3000,
            destination = "",
        )
        return;
    }
    const api_options = document.getElementById(`${llm_api}-option`);
    const fdata = new URLSearchParams(new FormData(api_options));
    let llm_api_id = api_name_to_id[llm_api];

    fdata.append("do-embedding", doEmbedding.checked);
    fdata.append("do-sentiment", doSentiment.checked);
    fdata.append("llm-api-id", llm_api_id);

    console.log(fdata);

    const response = await fetch("/v1/analyzer/" + `${pcid}?aid=${aid}&eid=${eid}`, {
        method: 'POST',
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: fdata
    });
    
    response.json().then((data) => {
        console.log(data);
        if ('error' in data) {
            err = data['error'];
            switch (err['code']) {
            case 410:
                ShowAlertToast(
                    message = "Preview has been expired, please create a new one",
                    x = 50, y = 10, duration = 3000,
                    destination = "",
                )
                has_sent = true;
                break;
            case 500:
                if ('pgx_code' in err) {
                    if (err['pgx_code'] === '23505') {
                        ShowAlertToast(
                            message = "You have already submitted this request",
                            x = 50, y = 10, duration = 3000,
                            destination = "",
                        )
                        break;
                    }
                } 
                ShowAlertToast(
                    message = "unknown error",
                    x = 50, y = 10, duration = 3000,
                    destination = "",
                )
                break;
            }
            // setTimeout(function() {
            //     location.href = err['url'];
            // }, 3000);
            return
        }
        ShowInfoToast(
            message = "Done!",
        );
        has_sent = true;
        setTimeout(function() {
            window.location.href = data.url;
        }, 3000);

    }).catch(err => {
        ShowAlertToast(
            message = err,
            x = 50, y = 10, duration = 3000,
            destination = "",
        )
    })
}
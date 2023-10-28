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

document.addEventListener("DOMContentLoaded", function () {
    getTimeZone();

    let categoryOpts = [
        { value: "", txt: "All" },

        { value: "business", txt: "Business" },
        { value: "entertainment", txt: "Entertainment" },
        { value: "environment", txt: "Environment" },
        { value: "food", txt: "Food" },
        { value: "health", txt: "Health" },
        { value: "politics", txt: "Politics" },
        { value: "science", txt: "Science" },
        { value: "sports", txt: "Sports" },
        { value: "technology", txt: "Technology" },
        { value: "top", txt: "Top" },
        { value: "tourism", txt: "Tourism" },
        { value: "world", txt: "World" },
    ];

    let countryOpts = [
        { value: "", txt: "All" },

        { value: "af", txt: "Afghanistan" },
        { value: "al", txt: "Albania" },
        { value: "dz", txt: "Algeria" },
        { value: "ao", txt: "Angola" },
        { value: "ar", txt: "Argentina" },
        { value: "au", txt: "Australia" },
        { value: "at", txt: "Austria" },
        { value: "az", txt: "Azerbaijan" },
        { value: "bh", txt: "Bahrain" },
        { value: "bd", txt: "Bangladesh" },
        { value: "bb", txt: "Barbados" },
        { value: "by", txt: "Belarus" },
        { value: "be", txt: "Belgium" },
        { value: "bm", txt: "Bermuda" },
        { value: "bt", txt: "Bhutan" },
        { value: "bo", txt: "Bolivia" },
        { value: "ba", txt: "Bosnia And Herzegovina" },
        { value: "br", txt: "Brazil" },
        { value: "bn", txt: "Brunei" },
        { value: "bg", txt: "Bulgaria" },
        { value: "bf", txt: "Burkinafasco" },
        { value: "kh", txt: "Cambodia" },
        { value: "cm", txt: "Cameroon" },
        { value: "ca", txt: "Canada" },
        { value: "cv", txt: "CapeVerde" },
        { value: "ky", txt: "Cayman Islands" },
        { value: "cl", txt: "Chile" },
        { value: "cn", txt: "China" },
        { value: "co", txt: "Colombia" },
        { value: "km", txt: "Comoros" },
        { value: "cr", txt: "Costa Rica" },
        { value: "hr", txt: "Croatia" },
        { value: "cu", txt: "Cuba" },
        { value: "cy", txt: "Cyprus" },
        { value: "cz", txt: "Czech Republic" },
        { value: "ci", txt: "Côte d&#39;Ivoire" },
        { value: "cd", txt: "Democratic Republic of the Congo" },
        { value: "dk", txt: "Denmark" },
        { value: "dj", txt: "Djibouti" },
        { value: "dm", txt: "Dominica" },
        { value: "do", txt: "Dominican Republic" },
        { value: "ec", txt: "Ecuador" },
        { value: "eg", txt: "Egypt" },
        { value: "sv", txt: "ElSalvador" },
        { value: "ee", txt: "Estonia" },
        { value: "et", txt: "Ethiopia" },
        { value: "fj", txt: "Fiji" },
        { value: "fi", txt: "Finland" },
        { value: "fr", txt: "France" },
        { value: "pf", txt: "French Polynesia" },
        { value: "ga", txt: "Gabon" },
        { value: "ge", txt: "Georgia" },
        { value: "de", txt: "Germany" },
        { value: "gh", txt: "Ghana" },
        { value: "gr", txt: "Greece" },
        { value: "gt", txt: "Guatemala" },
        { value: "gn", txt: "Guinea" },
        { value: "ht", txt: "Haiti" },
        { value: "hn", txt: "Honduras" },
        { value: "hk", txt: "Hong Kong" },
        { value: "hu", txt: "Hungary" },
        { value: "is", txt: "Iceland" },
        { value: "in", txt: "India" },
        { value: "id", txt: "Indonesia" },
        { value: "iq", txt: "Iraq" },
        { value: "ie", txt: "Ireland" },
        { value: "il", txt: "Israel" },
        { value: "it", txt: "Italy" },
        { value: "jm", txt: "Jamaica" },
        { value: "jp", txt: "Japan" },
        { value: "jo", txt: "Jordan" },
        { value: "kz", txt: "Kazakhstan" },
        { value: "ke", txt: "Kenya" },
        { value: "kw", txt: "Kuwait" },
        { value: "kg", txt: "Kyrgyzstan" },
        { value: "lv", txt: "Latvia" },
        { value: "lb", txt: "Lebanon" },
        { value: "ly", txt: "Libya" },
        { value: "lt", txt: "Lithuania" },
        { value: "lu", txt: "Luxembourg" },
        { value: "mo", txt: "Macau" },
        { value: "mk", txt: "Macedonia" },
        { value: "mg", txt: "Madagascar" },
        { value: "mw", txt: "Malawi" },
        { value: "my", txt: "Malaysia" },
        { value: "mv", txt: "Maldives" },
        { value: "ml", txt: "Mali" },
        { value: "mt", txt: "Malta" },
        { value: "mr", txt: "Mauritania" },
        { value: "mx", txt: "Mexico" },
        { value: "md", txt: "Moldova" },
        { value: "mn", txt: "Mongolia" },
        { value: "me", txt: "Montenegro" },
        { value: "ma", txt: "Morocco" },
        { value: "mz", txt: "Mozambique" },
        { value: "mm", txt: "Myanmar" },
        { value: "na", txt: "Namibia" },
        { value: "np", txt: "Nepal" },
        { value: "nl", txt: "Netherland" },
        { value: "nz", txt: "Newzealand" },
        { value: "ne", txt: "Niger" },
        { value: "ng", txt: "Nigeria" },
        { value: "kp", txt: "North Korea" },
        { value: "no", txt: "Norway" },
        { value: "om", txt: "Oman" },
        { value: "pk", txt: "Pakistan" },
        { value: "pa", txt: "Panama" },
        { value: "py", txt: "Paraguay" },
        { value: "pe", txt: "Peru" },
        { value: "ph", txt: "Philippines" },
        { value: "pl", txt: "Poland" },
        { value: "pt", txt: "Portugal" },
        { value: "pr", txt: "Puertorico" },
        { value: "ro", txt: "Romania" },
        { value: "ru", txt: "Russia" },
        { value: "rw", txt: "Rwanda" },
        { value: "ws", txt: "Samoa" },
        { value: "sm", txt: "SanMarino" },
        { value: "sa", txt: "Saudiarabia" },
        { value: "sn", txt: "Senegal" },
        { value: "rs", txt: "Serbia" },
        { value: "sg", txt: "Singapore" },
        { value: "sk", txt: "Slovakia" },
        { value: "si", txt: "Slovenia" },
        { value: "sb", txt: "Solomon Islands" },
        { value: "so", txt: "Somalia" },
        { value: "za", txt: "South Africa" },
        { value: "kr", txt: "South Korea" },
        { value: "es", txt: "Spain" },
        { value: "lk", txt: "Sri Lanka" },
        { value: "sd", txt: "Sudan" },
        { value: "se", txt: "Sweden" },
        { value: "ch", txt: "Switzerland" },
        { value: "sy", txt: "Syria" },
        { value: "tw", txt: "Taiwan" },
        { value: "tj", txt: "Tajikistan" },
        { value: "tz", txt: "Tanzania" },
        { value: "th", txt: "Thailand" },
        { value: "to", txt: "Tonga" },
        { value: "tn", txt: "Tunisia" },
        { value: "tr", txt: "Turkey" },
        { value: "tm", txt: "Turkmenistan" },
        { value: "ug", txt: "Uganda" },
        { value: "ua", txt: "Ukraine" },
        { value: "ae", txt: "United Arabemirates" },
        { value: "gb", txt: "United Kingdom" },
        { value: "us", txt: "United States" },
        { value: "uy", txt: "Uruguay" },
        { value: "uz", txt: "Uzbekistan" },
        { value: "ve", txt: "Venezuela" },
        { value: "vi", txt: "Vietnam" },
        { value: "ye", txt: "Yemen" },
        { value: "zm", txt: "Zambia" },
        { value: "zw", txt: "Zimbabwe" },
    ];

    let languageOpts = [
        { value: "", txt: "All" },

        { value: "af", txt: "Afrikaans" },
        { value: "sq", txt: "Albanian" },
        { value: "am", txt: "Amharic" },
        { value: "ar", txt: "Arabic" },
        { value: "as", txt: "Assamese" },
        { value: "az", txt: "Azerbaijani" },
        { value: "be", txt: "Belarusian" },
        { value: "bn", txt: "Bengali" },
        { value: "bs", txt: "Bosnian" },
        { value: "bg", txt: "Bulgarian" },
        { value: "my", txt: "Burmese" },
        { value: "ca", txt: "Catalan" },
        { value: "ckb", txt: "CentralKurdish" },
        { value: "zh", txt: "Chinese" },
        { value: "hr", txt: "Croatian" },
        { value: "cs", txt: "Czech" },
        { value: "da", txt: "Danish" },
        { value: "nl", txt: "Dutch" },
        { value: "en", txt: "English" },
        { value: "et", txt: "Estonian" },
        { value: "pi", txt: "Filipino" },
        { value: "fi", txt: "Finnish" },
        { value: "fr", txt: "French" },
        { value: "ka", txt: "Georgian" },
        { value: "de", txt: "German" },
        { value: "el", txt: "Greek" },
        { value: "gu", txt: "Gujarati" },
        { value: "he", txt: "Hebrew" },
        { value: "hi", txt: "Hindi" },
        { value: "hu", txt: "Hungarian" },
        { value: "is", txt: "Icelandic" },
        { value: "id", txt: "Indonesian" },
        { value: "it", txt: "Italian" },
        { value: "jp", txt: "Japanese" },
        { value: "kh", txt: "Khmer" },
        { value: "rw", txt: "Kinyarwanda" },
        { value: "ko", txt: "Korean" },
        { value: "lv", txt: "Latvian" },
        { value: "lt", txt: "Lithuanian" },
        { value: "lb", txt: "Luxembourgish" },
        { value: "mk", txt: "Macedonian" },
        { value: "ms", txt: "Malay" },
        { value: "ml", txt: "Malayalam" },
        { value: "mt", txt: "Maltese" },
        { value: "mi", txt: "Maori" },
        { value: "mr", txt: "Marathi" },
        { value: "mn", txt: "Mongolian" },
        { value: "ne", txt: "Nepali" },
        { value: "no", txt: "Norwegian" },
        { value: "or", txt: "Oriya" },
        { value: "ps", txt: "Pashto" },
        { value: "fa", txt: "Persian" },
        { value: "pl", txt: "Polish" },
        { value: "pt", txt: "Portuguese" },
        { value: "pa", txt: "Punjabi" },
        { value: "ro", txt: "Romanian" },
        { value: "ru", txt: "Russian" },
        { value: "sm", txt: "Samoan" },
        { value: "sr", txt: "Serbian" },
        { value: "sn", txt: "Shona" },
        { value: "si", txt: "Sinhala" },
        { value: "sk", txt: "Slovak" },
        { value: "sl", txt: "Slovenian" },
        { value: "so", txt: "Somali" },
        { value: "es", txt: "Spanish" },
        { value: "sw", txt: "Swahili" },
        { value: "sv", txt: "Swedish" },
        { value: "tg", txt: "Tajik" },
        { value: "ta", txt: "Tamil" },
        { value: "te", txt: "Telugu" },
        { value: "th", txt: "Thai" },
        { value: "tr", txt: "Turkish" },
        { value: "tk", txt: "Turkmen" },
        { value: "uk", txt: "Ukrainian" },
        { value: "ur", txt: "Urdu" },
        { value: "uz", txt: "Uzbek" },
        { value: "vi", txt: "Vietnamese" },
    ];

    addListenerToBtn("category", "insert-category-btn", "delete-category-btn", 5, categoryOpts, "You can add a maximum of 5 categories.");
    addListenerToBtn("country", "insert-country-btn", "delete-country-btn", 5, countryOpts, "You can add a maximum of 5 countries.");
    addListenerToBtn("language", "insert-lang-btn", "delete-lang-btn", 5, languageOpts, "You can add a maximum of 5 languages.");

})
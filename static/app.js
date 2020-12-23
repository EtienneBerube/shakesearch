const WORK_IDS = {
    0: "THE SONNETS",
    1: "ALL’S WELL THAT ENDS WELL",
    2: "ANTONY AND CLEOPATRA",
    3: "AS YOU LIKE IT",
    4: "THE COMEDY OF ERRORS",
    5: "THE TRAGEDY OF CORIOLANUS",
    6: "CYMBELINE",
    7: "THE TRAGEDY OF HAMLET, PRINCE OF DENMARK",
    8: "THE FIRST PART OF KING HENRY THE FOURTH",
    9: "THE SECOND PART OF KING HENRY THE FOURTH",
    10: "THE LIFE OF KING HENRY V",
    11: "THE FIRST PART OF HENRY THE SIXTH",
    12: "THE SECOND PART OF KING HENRY THE SIXTH",
    13: "THE THIRD PART OF KING HENRY THE SIXTH",
    14: "KING HENRY THE EIGHTH",
    15: "KING JOHN",
    16: "THE TRAGEDY OF JULIUS CAESAR",
    17: "THE TRAGEDY OF KING LEAR",
    18: "LOVE’S LABOUR’S LOST",
    19: "MACBETH",
    20: "MEASURE FOR MEASURE",
    21: "THE MERCHANT OF VENICE",
    22: "THE MERRY WIVES OF WINDSOR",
    23: "A MIDSUMMER NIGHT’S DREAM",
    24: "MUCH ADO ABOUT NOTHING",
    25: "OTHELLO",
    26: "PERICLES, PRINCE OF TYRE",
    27: "KING RICHARD THE SECOND",
    28: "KING RICHARD THE THIRD",
    29: "THE TRAGEDY OF ROMEO AND JULIET",
    30: "THE TAMING OF THE SHREW",
    31: "THE TEMPEST",
    32: "THE LIFE OF TIMON OF ATHENS",
    33: "THE TRAGEDY OF TITUS ANDRONICUS",
    34: "THE HISTORY OF TROILUS AND CRESSIDA",
    35: "TWELFTH NIGHT",
    36: "THE TWO GENTLEMEN OF VERONA",
    37: "THE TWO NOBLE KINSMEN",
    38: "THE WINTER’S TALE",
    39: "A LOVER’S COMPLAINT",
    40: "THE PASSIONATE PILGRIM",
    41: "THE PHOENIX AND THE TURTLE",
    42: "THE RAPE OF LUCRECE",
    43: "VENUS AND ADONIS"
}

const populateSelect = () => {
    let ele = document.getElementById('sel');
    for (let id in WORK_IDS) {
        // POPULATE SELECT ELEMENT WITH JSON.
        ele.innerHTML = ele.innerHTML +
            `<option value="${id}">${WORK_IDS[id]}</option>`;
    }
}

const Controller = {
    search: (ev) => {
        ev.preventDefault();
        const form = document.getElementById("form");
        const data = Object.fromEntries(new FormData(form));
        const selectedWork = document.getElementById("sel").value
        let url =`/search?q=${data.query}`
        if(selectedWork !== "all"){
            url+=`&w=${selectedWork}`
        }
        const response = fetch(url).then((response) => {
            response.json().then((results) => {
                Controller.updateTable(results);
            });
        });
    },

    updateTable: (results) => {
        const table = document.getElementById("table-body");
        const numResults = document.getElementById("num_result");
        const rows = [];
        for (let result of results.results) {
            rows.push(`<tr>${result}<tr/>`);
        }
        table.innerHTML = rows;
        numResults.innerText = `Server found ${results.results.length} results in ${results.time}`
    },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
populateSelect()

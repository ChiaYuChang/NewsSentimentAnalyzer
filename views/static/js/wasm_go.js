const go = new Go();

WebAssembly.instantiateStreaming(
    fetch("/static/wasm/app.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
    });

function sortTable(which) {
    let err = sort_table(which)
    if (err != null) {
        console.log(err)
    }
}
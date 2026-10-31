import { add } from "./math";

let c = 0;
render(c);

const btn = document.getElementById("add");
btn.onclick = () => {
    c = add(c);
    render(c);
};

function render(c) {
    const count = document.getElementById("count");
    count.textContent = `COUNT: ${c}`;
}


const modal = document.getElementById("myModal");

const openModalBtn = document.getElementById("openModalBtn");

const closeIcon = document.querySelector(".close-icon");

openModalBtn.addEventListener("click", function () {
  modal.style.display = "block";
  onSelectLanguage();
});

function closeModal() {
  modal.style.display = "none";
}

function onSelectLanguage() {
  var selectElement = document.getElementById("language");
  var selectedValue =
    selectElement.options[selectElement.selectedIndex].value || "html";
  let result_element = document.querySelector("#highlighting-content");
  result_element.className = "language-" + selectedValue;
  update(document.getElementById("editing").value);
}

closeIcon.addEventListener("click", function () {
  closeModal();
});

window.addEventListener("click", function (event) {
  if (event.target === modal) {
    closeModal();
  }
});

function update(text) {
  let result_element = document.querySelector("#highlighting-content");
  if (text[text.length - 1] == "\n") {
    text += " ";
  }
  result_element.innerHTML = text
    .replace(new RegExp("&", "g"), "&amp;")
    .replace(new RegExp("<", "g"), "&lt;");
  hljs.highlightElement(result_element);
}

function sync_scroll(element) {
  let result_element = document.querySelector("#highlighting");
  result_element.scrollTop = element.scrollTop;
  result_element.scrollLeft = element.scrollLeft;
}

function check_tab(element, event) {
  let code = element.value;
  if (event.key == "Tab") {
    event.preventDefault();
    let before_tab = code.slice(0, element.selectionStart);
    let after_tab = code.slice(
      element.selectionEnd,
      element.value.length
    );
    let cursor_pos = element.selectionStart + 1;
    element.value = before_tab + "\t" + after_tab;
    element.selectionStart = cursor_pos;
    element.selectionEnd = cursor_pos;
    update(element.value);
  }
}
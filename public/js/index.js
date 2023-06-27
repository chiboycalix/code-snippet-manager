const createSnippetModal = document.getElementById("createSnippetModal");
const updateSnippetModal = document.getElementById("updateSnippetModal");

const openCreateSnippetModal = document.getElementById(
  "openCreateSnippetModalBtn"
);
const openUpdateSnippetModal = document.getElementsByClassName(
  "fa-pen-to-square"
);
const createSnippetCloseModalIcon = document.querySelector(
  ".create-snippet-close-icon"
);
const updateSnippetCloseModalIcon = document.querySelector(
  ".update-snippet-close-icon"
);

for (let i = 0; i < openUpdateSnippetModal.length; i++) {
  openUpdateSnippetModal[i].addEventListener("click", function (event) {
    console.log(event.currentTarget, "hii")
    updateSnippetModal.style.display = "block";
    onSelectLanguage();
  });
}

createSnippetCloseModalIcon.addEventListener("click", function () {
  closeCreateSnippetModal();
});

updateSnippetCloseModalIcon.addEventListener("click", function () {
  closeUpdateSnippetModal();
});

function closeCreateSnippetModal() {
  createSnippetModal.style.display = "none";
}

function closeUpdateSnippetModal() {
  updateSnippetModal.style.display = "none";
}

window.addEventListener("click", function (event) {
  if (event.target === createSnippetModal) {
    closeCreateSnippetModal();
  }
});

window.addEventListener("click", function (event) {
  if (event.target === updateSnippetModal) {
    closeUpdateSnippetModal();
  }
});

function onSelectLanguage() {
  var selectElement = document.getElementById("language");
  var selectedValue = selectElement.options[selectElement.selectedIndex].value || "html";
  let result_element = document.querySelector("#highlighting-content");
  result_element.className = "language-" + selectedValue;
  update(document.getElementById("editing").value);
}

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
    let after_tab = code.slice(element.selectionEnd, element.value.length);
    let cursor_pos = element.selectionStart + 1;
    element.value = before_tab + "\t" + after_tab;
    element.selectionStart = cursor_pos;
    element.selectionEnd = cursor_pos;
    update(element.value);
  }
}

function handleCopy(event) {
  const text =
    event.target.parentNode.parentNode.parentNode.querySelector(
      ".custom-snippet"
    ).innerText;

  if (typeof navigator !== "undefined" && navigator.clipboard) {
    navigator.clipboard
      .writeText(text)
      .then(() => {
        // show toast
        const toast = document.createElement("div");
        toast.classList.add("toast");
        toast.innerText = "Copied to clipboard!!!";
        document.body.appendChild(toast);
        setTimeout(() => {
          toast.remove();
        }, 1000);
      })
      .catch((error) => {
        console.log(error);
        console.log(navigator, "navigator");
      });
  } else {
    // Clipboard API not supported, implement alternative approach
    console.log("Clipboard API not supported");
    console.log(navigator.clipboard, "navigator");
  }
}

function deleteSnippet(event) {
  const swalWithCustomButtons = Swal.mixin({
    customClass: {
      confirmButton: "confirmButton",
      cancelButton: "cancelButton",
    },
    buttonsStyling: false,
  });

  swalWithCustomButtons
    .fire({
      title: "Are you sure?",
      showCancelButton: true,
      confirmButtonText: "Delete",
      background: "#1a1a1a",
      color: "#fff",
      padding: "5rem 1rem",
    })
    .then((result) => {
      if (result.isConfirmed) {
        handleDelete(event);
      } else {
        // do nothing
      }
    });
}

function handleDelete(event) {
  const snippetId = event.target.dataset.snippetid;
  const snippet = event.target.parentNode.parentNode.parentNode;
  const url = `/snippets/${snippetId}`;

  fetch(url, {
    method: "POST",
  })
    .then((response) => {
      if (response.status === 200) {
        snippet.remove();
      }
    })
    .then((data) => {});
}

function handleEdit(event) {
  const snippetId = event.target.dataset.snippetid;
  const url = `/snippets/${snippetId}`;

  fetch(url, {
    method: "PUT",
  })
    .then((response) => {
      if (response.status === 200) {
        return response.json();
      }
    })
    .then((data) => {
      const { title, language, code } = data;
      const titleInput = document.querySelector("#title");
      const languageInput = document.querySelector("#language");
      const codeInput = document.querySelector("#editing");
      const submitButton = document.querySelector("#submit-button");

      titleInput.value = title;
      languageInput.value = language;
      codeInput.value = code;
      submitButton.innerText = "Update";
      submitButton.dataset.snippetid = snippetId;
      modal.style.display = "block";
      onSelectLanguage();
    });
}

function handleUpdate(event) {
  const snippetId = event.target.dataset.snippetid;
  const url = `/snippets/${snippetId}`;

  const title = document.querySelector("#title").value;
  const language = document.querySelector("#language").value;
  const code = document.querySelector("#editing").value;

  const data = {
    title,
    language,
    code,
  };

  fetch(url, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((response) => {
      if (response.status === 200) {
        window.location.reload();
      }
    })
    .then((data) => {});
}

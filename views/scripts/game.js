function isMobile() {
  return window.matchMedia("(max-width: 600px)").matches;
}

function addFields(digit) {
  let parent = document.createElement("div");
  parent.id = "input-container";
  parent.className = "input-container";
  let form = document.getElementById("form-container");
  let container = form.insertBefore(parent, form.childNodes[4]);

  // dynamically generate input fields for game.html
  for (i = 0; i < digit; i++) {
    let input = document.createElement("input");
    input.type = "number";
    input.min = "0";
    input.max = "9";
    input.name = "digit" + i;
    input.id = "digit" + i;
    input.className = "digit-input";
    input.tabIndex = -1;
    container.appendChild(input);

    if (isMobile()) {
      input.readOnly = true;
    }
  }

  if (isMobile()) {
    configureMobileKeyboard();
    const hintsTable = document.getElementById("hints-table");
    document.getElementById("popup-hints").appendChild(hintsTable);
  } else {
    setDigitInputEventListeners();
  }
}

// on large devices, use system keyboard
function setDigitInputEventListeners() {
  let boxes = document.getElementsByClassName("digit-input");
  let secondPress = false;
  document.getElementById("digit0").focus();

  // allow only single digit inputs
  Array.from(boxes).forEach((box, index, array) => {
    box.addEventListener("input", (event) => {
      if (box.value.length == 1 && index != array.length - 1) {
        // focus on the next sibling
        box.nextElementSibling.focus();
      }
    });
    box.addEventListener("keyup", (event) => {
      if (event.key == "Backspace" && !box.value && index > 0 && secondPress) {
        // focus on the previous sibling
        box.previousElementSibling.focus();
        secondPress = false;
      } else if (event.key == "Backspace" && !box.value && index > 0) {
        secondPress = true;
      } else if (
        event.key == "Enter" &&
        box.value &&
        index == array.length - 1
      ) {
        document.getElementById("submit-btn").focus();
      }
    });
  });
}

// concat strings from input boxes
function calcInputs() {
  let boxes = document.getElementsByClassName("digit-input");
  let val = "";
  Array.from(boxes).forEach((box) => {
    if (box.value) {
      val += box.value;
    }
  });
  return val;
}

function clearInputsAndReFocus() {
  const inputs = document.getElementsByClassName("digit-input");
  Array.from(inputs).forEach((i) => {
    i.value = "";
  });
  const box = document.getElementById("digit0");
  box.focus();
}

// mobile: enter digits via custom keyboard
function configureMobileKeyboard() {
  let n = document.getElementsByClassName("digit-input").length;
  let cursor = 0;
  let keys = Array.from(document.getElementsByClassName("key"));

  document.getElementById("digit0").focus();

  keys.forEach((key) => {
    if (key.id === "key-delete") {
      key.addEventListener("click", () => {
        if (cursor == 0) {
          return;
        }
        let box = document.getElementById("digit" + (cursor - 1));
        box.value = "";
        cursor--;
        if (box.previousElementSibling !== null) {
          box.previousElementSibling.focus();
        }
      });
    } else if (key.id === "key-enter") {
      key.addEventListener("click", () => {
        cursor = 0;
      });
    } else {
      key.addEventListener("click", () => {
        if (cursor == n) {
          return;
        }
        let box = document.getElementById("digit" + cursor);
        box.value = key.textContent;
        if (cursor != n) {
          cursor++;
          if (cursor == n) {
            box.blur();
          } else {
            box.nextElementSibling.focus();
          }
        }
      });
    }
  });
}

function addFields(digit) {
  let parent = document.createElement("div");
  parent.id = "input-container";
  parent.className = "input-container";
  let form = document.getElementById("form-container");
  let container = form.insertBefore(parent, form.childNodes[4]);

  // dynamically generate input fields for game.html
  for (i = 0; i < digit; i++) {
    let input = document.createElement("div");
    input.type = "number";
    input.min = "0";
    input.max = "9";
    input.name = "digit" + (i + 1);
    input.id = "digit" + (i + 1);
    input.className = "digit-input";
    input.tabIndex = -1;
    container.appendChild(input);
  }

  let boxes = document.getElementsByClassName("digit-input");
  let secondPress = false;
  document.getElementById("digit1").focus();

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

//
// mobile custom keyboard
//

const cursor = 1; // global variable to keep track of which the current input box

const keys = document.querySelectorAll(".key");
const deleteKey = document.getElementById("delete");
const enterKey = document.getElementById("enter");

// mobile: enter digits via custom keyboard
keys.forEach((key) => {
  const value = key.textContent;
  if (key.id === "delete") {
    key.addEventListener("click", () => {
      let box = document.getElementById("digit" + cursor);
      box.textContent = "";
      if (cursor != 1) {
        cursor--;
        box.previousElementSibling.focus();
      }
    });
  } else if (key.id === "enter") {
    key.addEventListener("click", () => {
      input.textContent = "";
    });
  } else {
    key.addEventListener("click", () => {});
  }
});

// TODO: so the problem is i want to show animations on the input fields using focus()
// but i dont want to show the built-in keyboard on mobile

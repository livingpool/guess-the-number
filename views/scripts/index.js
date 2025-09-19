// htmx.logAll();
document.addEventListener("DOMContentLoaded", (event) => {
  document.addEventListener("htmx:beforeSwap", (evt) => {
    if (evt.detail.xhr.status == 422) {
      evt.detail.shouldSwap = true;
      evt.detail.isError = false;
      if (evt.detail.target.id == "tbody") {
        Swal.fire({
          title: "u entered the wrong digits",
          icon: "warning",
          confirmButtonText: "OK",
        });
      }
    } else if (evt.detail.xhr.status == 404) {
      // handle player id not found
      evt.detail.isError = false;
      Swal.fire({
        title: "i couldnt find your id sorry",
        icon: "error",
        confirmButtonText: "Return Home",
        allowOutsideClick: false,
        allowEscapeKey: false,
      }).then(() => {
        const btn = document.getElementById("return-btn");
        htmx.trigger(btn, "click");
      });
    }
  });
  document.addEventListener("htmx:afterSwap", (evt) => {
    if (evt.detail.elt.id == "form-container") {
      const digit = parseInt(
        document.getElementById("form-container").getAttribute("data-digit")
      );
      addFields(digit);
    } else if (evt.detail.elt.id == "leaderboard-tbody") {
      openPopup();
      insertIndex();
    }
  });

  // TODO: there is a werid problem where this request's elt & target are both tbody. See if i can reproduce it!
  // handle winning condition
  document.addEventListener("htmx:afterSettle", (evt) => {
    if (evt.detail.target.id == "tbody" && evt.detail.xhr.status != 422) {
      const digit = document
        .getElementById("form-container")
        .getAttribute("data-digit");
      const matchDigit = document.getElementById("tbody").rows.length
        ? document.getElementById("tbody").rows[0].cells[2].textContent[0]
        : null;
      const ans = document.getElementById("tbody").rows.length
        ? document
            .getElementById("tbody")
            .rows[0].cells[1].textContent.split(" ")[1]
        : null;
      if (matchDigit && matchDigit == digit) {
        notifyResult(ans);
      }
    }
  });

  document.addEventListener("htmx:configRequest", (evt) => {
    // add param to request GET /check
    if (evt.detail.elt.id == "submit-btn") {
      evt.detail.parameters["guess"] = calcInputs(); // add a new param to request
      var inputs = document.getElementsByClassName("digit-input"); // clear input
      Array.from(inputs).forEach((i) => {
        i.value = "";
      });
      const box = document.getElementById("digit1");
      box.focus();
    }

    // add param to request GET /show-board
    if (evt.detail.elt.id == "leaderboard-btn") {
      const digit = parseInt(
        document.getElementById("form-container").getAttribute("data-digit")
      );
      const name = document
        .getElementById("form-container")
        .getAttribute("name");
      evt.detail.parameters["digit"] = digit;
      evt.detail.parameters["name"] = name;
    }
  });
});

// sweetalert2 notify game result and save record
async function notifyResult(ans) {
  const { value: name } = await Swal.fire({
    title: "you won the game!",
    icon: "success",
    input: "text",
    inputLabel: "enter your name to save the result",
    confirmButtonText: "insert",
    showCancelButton: true,
    text: `the answer is ${ans}`,
    inputValidator: (value) => {
      if (!value) {
        return "just gimme a name bro";
      }
    },
  });

  const existingName = document
    .getElementById("form-container")
    .getAttribute("name");

  if (name && !existingName) {
    // save record in db
    const digit = parseInt(
      document.getElementById("form-container").getAttribute("data-digit")
    );
    const attempts = document.getElementById("tbody").rows.length;

    fetch(window.location.origin + "/save-record", {
      method: "POST",
      body: JSON.stringify({
        digits: digit,
        name: name,
        attempts: attempts,
      }),
      headers: {
        "Content-Type": "application/json",
      },
    }).catch((err) => console.error(err));
    document.getElementById("form-container").setAttribute("name", name);
  } else if (name && existingName && name != existingName) {
    // degenerates spamming my ass
    await Swal.fire({
      title: "dont spam names okay?",
      icon: "warning",
      confirmButtonText: "im a good cat",
    }).then((result) => {
      if (result.isConfirmed) {
        const btn = document.getElementById("return-btn");
        htmx.trigger(btn, "click");
      }
    });
  } else {
    // enter same name, same (or worse) result so no api call
  }
}

// htmx.logAll();
document.addEventListener("DOMContentLoaded", (event) => {
  document.addEventListener("htmx:beforeSwap", (evt) => {
    if (evt.detail.xhr.status == 422) {
      evt.detail.shouldSwap = true;
      evt.detail.isError = false;
      if (evt.detail.target == htmx.find("#hints-table tbody")) {
        Swal.fire({
          title: "u entered the wrong digits",
          theme: htmlEl.dataset.theme,
          icon: "warning",
          confirmButtonText: "OK",
        });
      }
    } else if (evt.detail.xhr.status == 404) {
      // handle player id not found
      evt.detail.isError = false;
      Swal.fire({
        title: "i couldnt find your id sorry",
        theme: htmlEl.dataset.theme,
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
      if (digit) {
        addFields(digit);
      }
    } else if (evt.detail.elt == htmx.find("#popup-leaderboard tbody")) {
      openLeaderboardPopup();
      insertIndex();
    }
  });

  // TODO: there is a werid problem where this request's elt & target are both tbody. See if i can reproduce it!
  // handle winning condition
  document.addEventListener("htmx:afterSettle", (evt) => {
    const tbody = htmx.find("#hints-table tbody");
    if (evt.detail.target == tbody && evt.detail.xhr.status != 422) {
      const digit = document
        .getElementById("form-container")
        .getAttribute("data-digit");
      const matchDigit = tbody.rows.length
        ? tbody.rows[0].cells[2].textContent[0]
        : null;
      const ans = tbody.rows.length
        ? tbody.rows[0].cells[1].textContent.split(" ")[1]
        : null;
      if (matchDigit && matchDigit == digit) {
        notifyResult(ans);
      }
    }
  });

  document.addEventListener("htmx:configRequest", (evt) => {
    // add param to request GET /check
    if (evt.detail.elt.id == "submit-btn" || evt.detail.elt.id == "key-enter") {
      evt.detail.parameters["guess"] = calcInputs(); // add a new param to request
      clearInputsAndReFocus();
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
  let turnstileToken;
  const { value: name } = await Swal.fire({
    title: "you won the game!",
    theme: htmlEl.dataset.theme,
    icon: "success",
    input: "text",
    confirmButtonText: "insert",
    showCancelButton: true,
    html: `<p>the answer is ${ans}</p>
            <p>enter your name to save the result</p>`,
    didOpen: () => {
      Swal.disableButtons();
      Swal.getInput().insertAdjacentHTML(
        "afterend",
        `<div id="turnstile-widget">`
      );

      const captchaSiteKey = document
        .getElementById("outer")
        .getAttribute("captcha-site-key");
      console.log(captchaSiteKey);

      turnstile.render("#turnstile-widget", {
        sitekey: captchaSiteKey,
        size: "compact",
        callback: (token) => {
          console.debug("challenge success", token);
          turnstileToken = token;
          Swal.enableButtons();
        },
        "error-callback": (errorCode) => {
          console.error("challenge error", errorCode);
          Swal.fire({
            title: "captcha internal error",
            theme: htmlEl.dataset.theme,
            icon: "error",
            text: "if this issue persists, report to tim pls",
            confirmButtonText: "close",
          });
        },
        "expired-callback": () => {
          console.info("token expired");
        },
        "timeout-callback": () => {
          console.info("challenge timed out");
        },
      });
    },
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
    const attempts = htmx.find("#hints-table tbody").rows.length;

    fetch(window.location.origin + `/save-record?token=${turnstileToken}`, {
      method: "POST",
      body: JSON.stringify({
        digits: digit,
        name: name,
        attempts: attempts,
      }),
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then(async (resp) => {
        const msg = await resp.text();
        if (!resp.ok) {
          throw new Error(`status: ${resp.status}, msg: ${msg}`);
        }
        console.info("record inserted");
        document.getElementById("form-container").setAttribute("name", name);
      })
      .catch((err) => {
        console.error(err);
        Swal.fire({
          title: "failed to save record",
          theme: htmlEl.dataset.theme,
          icon: "error",
          html: `<p>${err}</p><p>if this issue persists, report to tim pls</p>`,
          confirmButtonText: "close",
        });
      });
  } else if (name && existingName && name != existingName) {
    // degenerates spamming my ass
    await Swal.fire({
      title: "dont spam names okay?",
      theme: htmlEl.dataset.theme,
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
